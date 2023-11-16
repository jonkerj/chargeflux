package smartevse

import (
	"log/slog"
	"reflect"
)

func VisitNumericFields(data interface{}, visit func(string, int64)) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)
		typ := t.Field(i)

		if tag := typ.Tag.Get("chargeflux"); tag != "" {
			switch val.Kind() {
			case reflect.Bool:
				if val.Bool() {
					visit(tag, 1)
				} else {
					visit(tag, 0)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				visit(tag, val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				visit(tag, int64(val.Uint()))
			default:
				slog.Warn("unknown data type", "tag", tag, "kind", val.Kind())
			}
		}

		if val.Kind() == reflect.Struct {
			VisitNumericFields(val.Interface(), visit)
		}
	}
}
