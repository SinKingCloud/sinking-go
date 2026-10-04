package repository

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"server/app/util/page"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// SelectPage 查询页码或游标分页数据，定位字段需要共同唯一标识一条记录。
// 未指定定位字段时使用模型主键；主排序字段之外的单个定位值直接写入 cursor_last_id，多个定位值编码为 JSON 字符串数组。
func (r *Repository[T]) SelectPage(query *gorm.DB, queryPage *page.Query, customCursorLastField ...string) (*page.Result[T], error) {
	if r == nil || r.Database == nil || r.Database.Db == nil || query == nil || queryPage == nil || queryPage.PageSize <= 0 || queryPage.PageSize > 1000 {
		return nil, gorm.ErrInvalidData
	}
	if query.Statement == nil || query.Statement.Model == nil {
		return nil, errors.New("分页查询模型无效")
	}
	queryPage.HasMore = false
	queryPage.NextCursorId = ""
	queryPage.NextCursorLastId = ""
	stmt := &gorm.Statement{DB: r.Database.Db}
	if err := stmt.Parse(query.Statement.Model); err != nil {
		return nil, err
	}
	orderField := stmt.Schema.LookUpField(queryPage.OrderByField)
	if orderField == nil || orderField.DBName == "" {
		return nil, errors.New("分页排序字段无效")
	}
	locatorFields := stmt.Schema.PrimaryFields
	if len(customCursorLastField) > 0 {
		locatorFields = make([]*schema.Field, 0, len(customCursorLastField))
		for _, name := range customCursorLastField {
			locatorFields = append(locatorFields, stmt.Schema.LookUpField(name))
		}
	}
	if len(locatorFields) == 0 {
		return nil, errors.New("分页定位字段无效")
	}
	fields := []*schema.Field{orderField}
	seen := map[string]bool{orderField.DBName: true}
	for _, field := range locatorFields {
		if field == nil || field.DBName == "" {
			return nil, errors.New("分页定位字段无效")
		}
		if !seen[field.DBName] {
			fields = append(fields, field)
			seen[field.DBName] = true
		}
	}
	desc := strings.ToLower(queryPage.OrderByType) != "asc"
	cursorMode := queryPage.IsCursor()
	cursorSet := queryPage.CursorId != "" || queryPage.CursorLastId != ""
	if cursorMode && cursorSet {
		cursorValues := []string{queryPage.CursorId}
		switch len(fields) {
		case 1:
			if queryPage.CursorId == "" || queryPage.CursorLastId != "" {
				return nil, errors.New("分页游标无效")
			}
		case 2:
			if queryPage.CursorLastId == "" {
				return nil, errors.New("分页定位值无效")
			}
			cursorValues = append(cursorValues, queryPage.CursorLastId)
		default:
			var locatorValues []string
			if err := json.Unmarshal([]byte(queryPage.CursorLastId), &locatorValues); err != nil || len(locatorValues) != len(fields)-1 {
				return nil, errors.New("分页定位值无效")
			}
			cursorValues = append(cursorValues, locatorValues...)
		}
		modelValue := reflect.New(stmt.Schema.ModelType).Elem()
		ctx := context.Background()
		branches := make([]clause.Expression, 0, len(fields))
		prefix := make([]clause.Expression, 0, len(fields))
		for i, field := range fields {
			if err := field.Set(ctx, modelValue, cursorValues[i]); err != nil {
				return nil, errors.Join(errors.New("分页游标无效"), err)
			}
			value, _ := field.ValueOf(ctx, modelValue)
			column := clause.Column{Name: field.DBName}
			var comparison clause.Expression = clause.Gt{Column: column, Value: value}
			if desc {
				comparison = clause.Lt{Column: column, Value: value}
			}
			conditions := append([]clause.Expression{}, prefix...)
			branches = append(branches, clause.And(append(conditions, comparison)...))
			prefix = append(prefix, clause.Eq{Column: column, Value: value})
		}
		query = query.Where(clause.Or(branches...))
	}
	limit := queryPage.PageSize
	var total int64
	if cursorMode {
		limit++
	} else {
		const maxPageOffset int64 = 500000
		pageIndex := int64(queryPage.Page - 1)
		if pageIndex > maxPageOffset/int64(queryPage.PageSize) {
			return nil, errors.New("分页位置过深，请使用游标分页")
		}
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}
		query = query.Offset(int(pageIndex * int64(queryPage.PageSize)))
	}
	columns := make([]clause.OrderByColumn, 0, len(fields))
	for i, field := range fields {
		columns = append(columns, clause.OrderByColumn{Column: clause.Column{Name: field.DBName}, Desc: desc, Reorder: i == 0})
	}
	query = query.Clauses(clause.OrderBy{Columns: columns})
	list := make([]T, 0)
	if err := query.Limit(limit).Find(&list).Error; err != nil {
		return nil, err
	}
	if !cursorMode || len(list) <= queryPage.PageSize {
		return page.New(queryPage, total, list), nil
	}
	list = list[:queryPage.PageSize]
	last := reflect.ValueOf(list[len(list)-1])
	if !last.IsValid() {
		return nil, errors.New("分页数据无效")
	}
	for last.Kind() == reflect.Interface || last.Kind() == reflect.Ptr {
		if last.IsNil() {
			return nil, errors.New("分页数据无效")
		}
		last = last.Elem()
	}
	resultStmt := &gorm.Statement{DB: r.Database.Db}
	if err := resultStmt.Parse(last.Interface()); err != nil {
		return nil, err
	}
	cursorValues := make([]string, 0, len(fields))
	for _, field := range fields {
		resultField := resultStmt.Schema.LookUpField(field.DBName)
		if resultField == nil {
			return nil, errors.New("分页结果缺少排序字段")
		}
		value, _ := resultField.ValueOf(context.Background(), last)
		cursorValue, err := formatCursorValue(value)
		if err != nil {
			return nil, err
		}
		cursorValues = append(cursorValues, cursorValue)
	}
	queryPage.HasMore = true
	queryPage.NextCursorId = cursorValues[0]
	if len(cursorValues) == 2 {
		queryPage.NextCursorLastId = cursorValues[1]
	} else if len(cursorValues) > 2 {
		locatorValues, err := json.Marshal(cursorValues[1:])
		if err != nil {
			return nil, err
		}
		queryPage.NextCursorLastId = string(locatorValues)
	}
	return page.New(queryPage, total, list), nil
}

func formatCursorValue(fieldValue interface{}) (string, error) {
	for fieldValue != nil {
		value := reflect.ValueOf(fieldValue)
		if (value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr) && value.IsNil() {
			fieldValue = nil
			break
		}
		if valuer, ok := fieldValue.(driver.Valuer); ok {
			var err error
			fieldValue, err = valuer.Value()
			if err != nil {
				return "", err
			}
			break
		}
		if value.Kind() != reflect.Interface && value.Kind() != reflect.Ptr {
			break
		}
		fieldValue = value.Elem().Interface()
	}
	if fieldValue == nil {
		return "", errors.New("分页结果游标无效")
	}
	switch value := fieldValue.(type) {
	case time.Time:
		return value.Format(time.RFC3339Nano), nil
	case []byte:
		return string(value), nil
	default:
		return fmt.Sprint(value), nil
	}
}
