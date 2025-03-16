package environ

import (
	"os"
	"reflect"
)

func Load[T any]() T {
	var conf T

	val := reflect.Indirect(reflect.ValueOf(&conf))
	for index := 0; index < val.NumField(); index++ {
		fname := val.Type().Field(index).Name
		field := val.FieldByName(fname)

		tag := val.Type().Field(index).Tag.Get("env")

		value, found := os.LookupEnv(tag)
		if !found {
			value = val.Type().Field(index).Tag.Get("default")
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		}
	}

	return conf
}
