package smartevse

import (
	"log/slog"
	"reflect"

	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func addInfluxFields(data interface{}, p *write.Point) {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	for i := 0; i < t.NumField(); i++ {
		val := v.Field(i)
		typ := t.Field(i)

		if tag := typ.Tag.Get("influx"); tag != "" {
			switch val.Kind() {
			case reflect.Bool:
				if val.Bool() {
					addField(p, tag, 1)
				} else {
					addField(p, tag, 0)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i := val.Int()
				addField(p, tag, i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i := val.Uint()
				addField(p, tag, i)
			default:
				slog.Warn("unknown data type", "tag", tag, "kind", val.Kind())
			}
		}

		if val.Kind() == reflect.Struct {
			addInfluxFields(val.Interface(), p)
		}
	}

}

func addField(p *write.Point, field string, v interface{}) {
	slog.Debug("adding field", "field", field, "value", v)
	p = p.AddField(field, v)
}

func (s *SmartEVSESettings) AddFieldsToPoint(p *write.Point) {
	addInfluxFields(*s, p)
}
