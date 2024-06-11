package condutils

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

// Or check if testValue is empty, if empty return defaultValue.
func Or(testValue, defaultValue interface{}) interface{} {
	if !IsEmpty(testValue) {
		return testValue
	}
	return defaultValue
}

// Ors like Or but accept multi args.
func Ors(val1, val2 interface{}, values ...interface{}) interface{} {
	if len(values) == 0 {
		return Or(val1, val2)
	}
	return Ors(Or(val1, val2), values[0], values[1:]...)
}

// Map iterate colls and apply fn to each item on colls, return colls of fn result.
func Map(fn func(arg interface{}) interface{}, colls []interface{}) []interface{} {
	result := []interface{}{}
	for _, i := range colls {
		evaluated := fn(i)
		result = append(result, evaluated)
	}
	return result
}

// EvalSliceReflectValue eval slice of reflect value to an interface.
func EvalSliceReflectValue(values []reflect.Value) interface{} {
	multiReturn := []interface{}{}
	for _, v := range values {
		multiReturn = append(multiReturn, v.Interface())
	}

	return multiReturn
}

// Eval if arg is func with return value eval arg and return the result
// if multi value returned will be []interface{}
// if not func, return arg. func must pure functions
func Eval(arg interface{}) interface{} {
	if IsFunc(arg) {
		res := reflect.ValueOf(arg).Call([]reflect.Value{})
		if len(res) == 1 {
			return res[0].Interface()
		}
		return EvalSliceReflectValue(res)
	}
	return arg
}

// Constantly give arg and return func that accept single argument
// subsequent call to returned function wil return the same value (arg).
func Constantly(arg interface{}) func(args ...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		return arg
	}
}

// getFnName get function name
func getFnName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// EvalExpr eval fnAndArgs as expression
// expressions must contains at least 1 func at first, the rest will be added as arguments to func
// if multi value returned will be []reflect.Value{}
// on panic, will append error to return value.
func EvalExpr(fnAndArgs ...interface{}) (returned []reflect.Value) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("call %v, panicked: %v", getFnName(fnAndArgs[0]), r)
			returned = append(returned, reflect.ValueOf(err))
		}
	}()
	fn := fnAndArgs[0]
	args := []reflect.Value{}
	for _, arg := range fnAndArgs[1:] {
		args = append(args, reflect.ValueOf(arg))
	}

	return reflect.ValueOf(fn).Call(args)
}

// When if condition true eval expressions
// expressions must contains at least 1 func at first, the rest will be added as arguments to func
// return nil or result of func. func must pure functions.
func When(condition bool, expression ...interface{}) interface{} {

	if condition {
		if IsEmpty(expression) {
			return errors.New("atleast 1 value provided on expressions")
		}

		if !IsFunc(expression[0]) {
			return errors.New("func is first entry to expressions")
		}

		if len(expression) < 2 {
			return Eval(expression[len(expression)-1])
		}

		res := EvalExpr(expression...)

		if len(res) == 1 {
			return res[0].Interface()
		}

		return EvalSliceReflectValue(res)
	}

	return nil
}

// TakeNth take an nth entry from slice
func TakeNth(n int, seqs []interface{}) []interface{} {
	if seqs == nil {
		return []interface{}{}
	}

	if n < 1 {
		return seqs[:1]
	}

	if n > len(seqs) {
		return []interface{}{}
	}

	result := []interface{}{}
	l := len(seqs)

	for i := n; i <= l; i += n {
		result = append(result, seqs[i-1])
	}

	return result
}

// SliceOfInterface coerce maybeSlice as []interface{}
// if maybeSlice isn't slice wrap it inside []interface{}.
func SliceOfInterface(maybeSlice interface{}) []interface{} {
	if !IsSlice(maybeSlice) {
		return SliceOfInterface([]interface{}{maybeSlice})
	}

	if sliceIface, ok := maybeSlice.([]interface{}); ok {
		return sliceIface
	}

	if sliceRF, ok := maybeSlice.([]reflect.Value); ok {
		return EvalSliceReflectValue(sliceRF).([]interface{})
	}

	s := reflect.ValueOf(maybeSlice)
	n := s.Len()
	r := []interface{}{}
	for i := 0; i < n; i++ {
		r = append(r, s.Index(i).Interface())
	}
	return r
}

// NumberFormat result under thousands not expected like thousands.
func NumberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	var prefix, suffix string
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}
