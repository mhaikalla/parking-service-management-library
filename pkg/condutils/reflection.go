package condutils

import (
	"fmt"
	"reflect"
)

// BindFunc binding func via  `reflect.MakeFunc`.
// `i` is pointer of the func want to bind.
// `f` is layout needed by `reflect.MakeFunc`.
func BindFunc(i interface{}, f func(args []reflect.Value) []reflect.Value) {
	impl := reflect.ValueOf(i).Elem()
	impl.Set(reflect.MakeFunc(impl.Type(), f))
}

// AsyncCall call a func with it's args and return chan.
// func must be single return value.
// Second args must be a func and the rest args is applied to the func.
// Return closed channel of type first return value of func.
func AsyncCall(args []reflect.Value) []reflect.Value {

	fn, fnArgs := args[0].Interface(), args[1].Interface().([]interface{})

	if checkErr := isCallValid(fn, fnArgs); checkErr != nil {
		panic(fmt.Sprintf("cannot evaluate expr with error: %v", checkErr))
	}

	ctype := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(fn).Out(0))
	c := reflect.MakeChan(ctype, 0)

	rfSlice := []reflect.Value{}

	for _, iface := range fnArgs {
		rfSlice = append(rfSlice, reflect.ValueOf(iface))
	}

	go func(c reflect.Value, f reflect.Value, a []reflect.Value) {
		defer c.Close()
		res := f.Call(a)[0]

		c.Send(res)
	}(c, reflect.ValueOf(fn), rfSlice)

	return []reflect.Value{c}
}

// isCallValid check if call to func and supplied args is match.
// we expect, func is single return value.
func isCallValid(fn interface{}, args []interface{}) error {
	if !IsFunc(fn) {
		return fmt.Errorf("expect func got %T", fn)
	}

	fnType := reflect.TypeOf(fn)
	numOut := fnType.NumOut()
	numIn := fnType.NumIn()
	lenArgs := len(args)

	if numOut > 1 && numOut > 0 {
		return fmt.Errorf("expect func return single value, got %d", numOut)
	}

	if numIn != lenArgs {
		return fmt.Errorf("func need %d args, but got %d args", numIn, lenArgs)
	}

	// for i, a := range args {
	// 	t := reflect.TypeOf(a)
	// 	if !fnType.In(i).ConvertibleTo(t) {
	// 		return false
	// 	}
	// }
	return nil
}
