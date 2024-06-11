package condutils

import (
	"fmt"
	"reflect"
	"unsafe"
)

const strFormatDescField = "field: %s; type of: %s; has value: %#v"

// IsEmpty check if arg is empty either zero by it's type or no entries on structure
func IsEmpty(arg interface{}) bool {

	if rf, ok := arg.(reflect.Value); ok {
		if !rf.IsValid() {
			return true
		}
		return IsEmpty(rf.Interface())
	}

	if arg != nil {
		switch {
		case IsMap(arg), IsSlice(arg):
			return reflect.ValueOf(arg).Len() < 1
		case IsStruct(arg):
			return reflect.Indirect(reflect.ValueOf(arg)).IsZero()
		default:
			argType := reflect.TypeOf(arg)
			zeroVal := reflect.Zero(argType).Interface()
			return reflect.DeepEqual(arg, zeroVal)
		}
	}
	return true
}

// IsError check if interface is type of error by poking it's method
// return boolean and the arg it test
func IsError(maybeErr interface{}) (bool, interface{}) {

	if maybeErr == nil {
		return false, nil
	}

	if reflect.TypeOf(maybeErr).Name() == "Value" {
		return IsError(maybeErr.(reflect.Value).Interface())
	}

	if reflect.ValueOf(maybeErr).MethodByName("Error").Kind() != reflect.Func {
		return false, maybeErr
	}

	return true, maybeErr
}

// IsFunc check if args is func kind
func IsFunc(maybeFunc interface{}) bool {
	if maybeFunc != nil {
		return reflect.TypeOf(maybeFunc).Kind() == reflect.Func
	}
	return false
}

// IsStruct check if arg is struct kind
func IsStruct(maybeStruct interface{}) bool {
	if maybeStruct != nil {
		if IsPtr(maybeStruct) {
			refVal := reflect.ValueOf(maybeStruct)
			return reflect.Indirect(refVal).Kind() == reflect.Struct
		}
		return reflect.TypeOf(maybeStruct).Kind() == reflect.Struct
	}
	return false
}

// IsPtr check if arg is pointer kind
func IsPtr(maybePtr interface{}) bool {
	if maybePtr != nil {
		if rf, ok := maybePtr.(reflect.Value); ok {
			if !rf.IsValid() {
				return false
			}
			return IsPtr(rf.Interface())
		}
		return reflect.TypeOf(maybePtr).Kind() == reflect.Ptr
	}
	return false
}

// IsMap check if arg is map kind
func IsMap(maybeMap interface{}) bool {
	if maybeMap != nil {
		if rf, ok := maybeMap.(reflect.Value); ok {
			if !rf.IsValid() {
				return false
			}
			return IsMap(rf.Interface())
		}

		if IsPtr(maybeMap) {
			refVal := reflect.ValueOf(maybeMap)
			return reflect.Indirect(refVal).Kind() == reflect.Map
		}
		return reflect.TypeOf(maybeMap).Kind() == reflect.Map
	}
	return false
}

// IsSlice check if arg is slice kind
func IsSlice(maybeSlice interface{}) bool {
	if maybeSlice != nil {
		if rf, ok := maybeSlice.(reflect.Value); ok {
			if !rf.IsValid() {
				return false
			}
			return IsSlice(rf.Interface())
		}

		if IsPtr(maybeSlice) {
			refVal := reflect.ValueOf(maybeSlice)
			return reflect.Indirect(refVal).Kind() == reflect.Slice
		}
		return reflect.TypeOf(maybeSlice).Kind() == reflect.Slice
	}
	return false
}

// IsCompound ...
func IsCompound(maybeCompound interface{}) bool {
	return IsMap(maybeCompound) ||
		IsSlice(maybeCompound) ||
		IsStruct(maybeCompound)
}

func structToRV(strct interface{}) reflect.Value {
	v := reflect.ValueOf(strct)
	if !IsPtr(strct) {
		newAddr := reflect.New(v.Type())
		newAddr.Elem().Set(v)
		strct = newAddr.Interface()
	}

	return reflect.Indirect(reflect.ValueOf(strct))
}

// IsComplete check if all field on struct not zero value
// if not struct return false
// please be carefull using this function since we use unsafe pointer
// for accessing private field
func IsComplete(target interface{}) bool {
	if IsStruct(target) {
		el := structToRV(target)
		num := el.NumField()
		for i := 0; i < num; i++ {
			f := el.Field(i)
			if !f.CanSet() {
				/* #nosec G103 */
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			}

			if IsEmpty(f.Interface()) {
				return false
			}
		}
		return true
	}
	return false
}

// DescEmptyField describe zeroed/empty field of struct recursively,
// return []string{} with each entry described zeroed/empty field,
// return empty slice of string if not struct type.
func DescEmptyField(target interface{}) []string {
	result := []string{}
	if IsStruct(target) {
		el := structToRV(target)
		num := el.NumField()
		for i := 0; i < num; i++ {
			f := el.Field(i)
			if !f.CanSet() {
				/* #nosec G103 */
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			}

			if IsStruct(f.Interface()) {
				result = append(result, DescEmptyField(f.Interface())...)
				continue
			}

			if IsEmpty(f.Interface()) {
				result = append(result, fmt.Sprintf(strFormatDescField, el.Type().Field(i).Name, f.Type(), f))
			}
		}
	}
	return result
}
