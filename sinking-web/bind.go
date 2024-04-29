package sinking_web

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// BindAll 绑定所有参数
func (c *Context) BindAll(obj interface{}) error {
	_ = c.bind(c.AllParam(), obj)
	_ = c.bind(c.AllQuery(), obj)
	_ = c.bind(c.AllForm(), obj)
	_ = c.BindJson(obj)
	if obj == nil {
		return errors.New("the param bind error")
	}
	return nil
}

// BindForm 绑定post参数
func (c *Context) BindForm(obj interface{}) error {
	return c.bind(c.AllForm(), obj)
}

// BindQuery 绑定get参数
func (c *Context) BindQuery(obj interface{}) error {
	return c.bind(c.AllQuery(), obj)
}

// BindParam 绑定路径参数
func (c *Context) BindParam(obj interface{}) error {
	return c.bind(c.AllParam(), obj)
}

// BindJson 绑定json
func (c *Context) BindJson(obj interface{}) error {
	body := c.Body()
	err := json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return err
	}
	return nil
}

// bind 通用绑定
func (c *Context) bind(params map[string]string, obj interface{}) error {
	keys := reflect.TypeOf(obj).Elem()
	values := reflect.ValueOf(obj).Elem()
	for i := 0; i < keys.NumField(); i++ {
		name := keys.Field(i).Tag.Get(BindFormTagName)
		if name == "" {
			name = keys.Field(i).Name
		}
		if params[name] == "" {
			var isNull bool
			var defaultValue string
			value := values.Field(i)
			switch value.Kind() {
			case reflect.String:
				isNull = value.String() == ""
				defaultValue = value.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				isNull = value.Int() == 0
				defaultValue = fmt.Sprintf("%d", value.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				isNull = value.Uint() == 0
				defaultValue = fmt.Sprintf("%d", value.Uint())
			case reflect.Float32, reflect.Float64:
				isNull = value.Float() == 0
				defaultValue = fmt.Sprintf("%f", value.Float())
			case reflect.Bool:
				isNull = value.Bool() == false
				if value.Bool() {
					defaultValue = "true"
				} else {
					defaultValue = "false"
				}
			default:
				isNull = false
				defaultValue = value.String()
			}
			if !isNull {
				params[name] = defaultValue
			} else {
				defaultValue = keys.Field(i).Tag.Get(BindDefaultValueTagName)
				if defaultValue != "" {
					params[name] = defaultValue
				}
			}
		}
		err := setWithProperType(params[name], values.Field(i), keys.Field(i))
		if err != nil {
			continue
		}
	}
	return nil
}

// setWithProperType 类型转换
func setWithProperType(val string, value reflect.Value, field reflect.StructField) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, 0, value)
	case reflect.Int8:
		return setIntField(val, 8, value)
	case reflect.Int16:
		return setIntField(val, 16, value)
	case reflect.Int32:
		return setIntField(val, 32, value)
	case reflect.Int64:
		switch value.Interface().(type) {
		case time.Duration:
			return setTimeDuration(val, value)
		}
		return setIntField(val, 64, value)
	case reflect.Uint:
		return setUintField(val, 0, value)
	case reflect.Uint8:
		return setUintField(val, 8, value)
	case reflect.Uint16:
		return setUintField(val, 16, value)
	case reflect.Uint32:
		return setUintField(val, 32, value)
	case reflect.Uint64:
		return setUintField(val, 64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, 32, value)
	case reflect.Float64:
		return setFloatField(val, 64, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		switch value.Interface().(type) {
		case time.Time:
			return setTimeField(val, field, value)
		}
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(val), value.Addr().Interface())
	default:
		return nil
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}

	switch tf := strings.ToLower(timeFormat); tf {
	case "unix", "unixnano":
		tv, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}

		d := time.Duration(1)
		if tf == "unixnano" {
			d = time.Second
		}

		t := time.Unix(tv/int64(d), tv%int64(d))
		value.Set(reflect.ValueOf(t))
		return nil
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}

func setArray(vals []string, value reflect.Value, field reflect.StructField) error {
	for i, s := range vals {
		err := setWithProperType(s, value.Index(i), field)
		if err != nil {
			return err
		}
	}
	return nil
}

func setTimeDuration(val string, value reflect.Value) error {
	d, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(d))
	return nil
}
