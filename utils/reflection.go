package utils

import (
	"reflect"
)

func SetReflctionTag(model any, tag string, tagValue string, value any) bool {
	v := reflect.ValueOf(model)

	if v.Kind() != reflect.Ptr {
		return false
	}

	v = v.Elem()
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)

		if fieldType.Tag.Get(tag) == tagValue {
			fieldValue := v.Field(i)

			if !fieldValue.CanSet() {
				return false
			}

			newVal := reflect.ValueOf(value)

			if newVal.Type().AssignableTo(fieldValue.Type()) {
				fieldValue.Set(newVal)
				return true
			}

			if newVal.Type().ConvertibleTo(fieldValue.Type()) {
				fieldValue.Set(newVal.Convert(fieldValue.Type()))
				return true
			}

			return false

		}

	}

	return true
}

func StructToMap(model any, tag string) map[string]any {
	result := make(map[string]any)

	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get(tag)

		if tag != "" && tag != "-" {
			result[tag] = reflect.ValueOf(model).Elem().Field(i).Interface()
		}
	}
	return result
}

func ReturnMetadataTable(model any, tag string) []map[string]string {
	result := []map[string]string{}

	t := reflect.TypeOf(model)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := range t.NumField() {
		el := make(map[string]string)
		field := t.Field(i)
		tagField := field.Tag.Get(tag)

		if tagField == "" || tagField == "-" {
			continue
		}

		el["Field"] = tagField
		el["Type"] = ParserTypesByDatabases(field.Type.String())

		result = append(result, el)
	}

	return result
}
