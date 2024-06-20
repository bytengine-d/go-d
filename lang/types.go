package lang

import (
	"reflect"
)

func GetTargetTypeFromValue(v any) reflect.Type {
	return GetTargetType(reflect.TypeOf(v))
}

func GetTargetType(t reflect.Type) reflect.Type {
	var (
		targetType reflect.Type
	)
	if t.Kind() == reflect.Struct {
		targetType = t
	} else if t.Kind() == reflect.Pointer {
		targetType = t.Elem()
	}
	return targetType
}
