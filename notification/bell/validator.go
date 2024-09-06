package bell

import (
	"fmt"
	"reflect"
)

func validatePayload(payload NotificationPayload) error {
	v := reflect.ValueOf(payload)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			return fmt.Errorf("missing required field: %s", v.Type().Field(i).Name)
		}
		if field.Kind() == reflect.Interface && field.IsNil() {
			return fmt.Errorf("missing required field: %s", v.Type().Field(i).Name)
		}
	}
	return nil
}
