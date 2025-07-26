package bootstrap

import (
	"fmt"
	"os"
	"regexp"
	"server/app/constant"
	"server/app/util"
	"server/app/util/database"
	"server/app/util/file"
	"server/public"
	"strings"
)

// LoadDatabase 初始化数据库
func LoadDatabase() {
	util.Database = database.NewSqlite(getDbFile())
	if util.Database.DbError != nil {
		panic(util.Database.DbError)
		return
	}
	err := checkDatabaseInit()
	if err != nil {
		panic(err)
		return
	}
}

// checkDatabaseInit 判断数据库是否初始化
func checkDatabaseInit() error {
	var tables []string
	err := util.Database.Db.Raw(`SELECT name FROM sqlite_master WHERE type = 'table' and name like 'cloud_%'`).Scan(&tables).Error
	if err != nil {
		return err
	}
	sql := public.Sql
	if len(tables) >= strings.Count(sql, "create table") {
		if checkContainsAllElements(tables, getSqlCreateTables(sql)) {
			return nil
		}
	}
	//初始化dbFile
	deleteDbFile()
	util.Database = database.NewSqlite(getDbFile())
	if util.Database.DbError != nil {
		return util.Database.DbError
	}
	//新建数据表
	lines := strings.Split(public.Sql, "\n")
	sqlStmt := ""
	successCount := 0
	errorCount := 0
	for _, line := range lines {
		if len(line) > 0 && !(strings.HasPrefix(line, "--") || line == "" || strings.HasPrefix(line, "/*")) {
			sqlStmt += line
			if strings.HasSuffix(line, ";") && line != "COMMIT;" {
				sqlStmt = strings.ReplaceAll(sqlStmt, "INSERT INTO ", "INSERT IGNORE INTO ")
				ctx := util.Database.Db.Exec(sqlStmt)
				if ctx.Error != nil {
					err = ctx.Error
					errorCount++
				} else {
					successCount++
				}
				sqlStmt = ""
			}
		}
	}
	if errorCount > 0 {
		return fmt.Errorf("创建数据表失败")
	}
	return nil
}

// 获取并初始化db文件
func getDbFile() string {
	path := constant.DBPath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	f := file.NewDisk(path)
	_ = f.AutoCreate(constant.DBFile)
	return strings.ReplaceAll(path+constant.DBFile, "//", "")
}

// deleteDbFile 删除db文件
func deleteDbFile() {
	f := getDbFile()
	_ = os.Remove(f)
}

// getSqlCreateTables 获取建表语句中的表
func getSqlCreateTables(sql string) []string {
	createTableRegex := regexp.MustCompile("(?i)CREATE TABLE (\\w+)\\s")
	matches := createTableRegex.FindAllStringSubmatch(sql, -1)
	var tableNames []string
	for _, match := range matches {
		tableNames = append(tableNames, match[1])
	}
	return tableNames
}

// checkContainsAllElements 判断A是否全部含有B的元素
func checkContainsAllElements(A, B []string) bool {
	for _, b := range B {
		found := false
		for _, a := range A {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
