package from_request

import (
	"fmt"
	"reflect"
	"strconv"
)

func getBitSize(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Int, reflect.Uint:
		return 32
	case reflect.Int8, reflect.Uint8:
		return 8
	case reflect.Int16, reflect.Uint16:
		return 16
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64
	default:
		panic("only int, uint and float can support getBitSize")
	}
}

func ConvertFromStr(strValue string, v reflect.Value) error {
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.String:
		v.SetString(strValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intV, err := strconv.ParseInt(strValue, 10, getBitSize(v))
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(intV).Convert(v.Type()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintV, err := strconv.ParseUint(strValue, 10, getBitSize(v))
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(uintV).Convert(v.Type()))
	case reflect.Float32, reflect.Float64:
		bitSize := getBitSize(v)
		floatV, err := strconv.ParseFloat(strValue, bitSize)
		if err != nil {
			return err
		}
		v.Set(reflect.ValueOf(floatV).Convert(v.Type()))
	case reflect.Bool:
		boolV, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}
		v.SetBool(boolV)
	default:
		return fmt.Errorf("un support")
	}
	return nil
}
