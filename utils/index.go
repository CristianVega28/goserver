package utils

import (
	"errors"
	"reflect"
)

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func CheckTypesForResponse(value any) ([]map[string]any, error) {
	var response = make([]map[string]any, 0)
	if arr, ok := value.([]any); ok {
		for _, v := range arr {
			if m, ok := v.(map[string]any); ok {
				response = append(response, m)
			}
		}
	} else {
		return nil, errors.New("invalid type for response, expected []map[string]any")
	}

	return response, nil
}

func GetAt[T any](arr []T, index int) (T, bool) {
	var zero T
	if index < 0 || index >= len(arr) {
		return zero, false
	}
	return arr[index], true
}
