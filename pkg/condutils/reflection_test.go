package condutils

import (
	"reflect"
	"testing"
)

func TestBindFunc(t *testing.T) {
	e := func(a, b int) int { return a + b }

	var f1 func(fn interface{}, args ...interface{}) chan int

	type args struct {
		i interface{}
		f func(args []reflect.Value) []reflect.Value
	}

	tests := []struct {
		name   string
		args   args
		fnTest func(fn interface{}, args ...interface{}) chan int
	}{
		{
			"1",
			args{
				f: AsyncCall,
			},
			f1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.i = &tt.fnTest
			BindFunc(tt.args.i, tt.args.f)
			if tt.args.i == nil {
				t.Errorf("BindAsyncCall(i, f)=> %#v", tt.args.i)
			}

			if res := <-tt.fnTest(e, 1, 2); res != 3 {
				t.Errorf("%#v", res)
			}
		})
	}
}

func TestAsyncCall(t *testing.T) {
	type args struct {
		args []reflect.Value
	}
	tests := []struct {
		name string
		args args
		want []reflect.Value
	}{
		{
			"1",
			args{args: []reflect.Value{
				reflect.ValueOf(func(a, b int) int { return a + b }),
				reflect.ValueOf([]interface{}{1, 2}),
			}},
			[]reflect.Value{reflect.ValueOf(make(chan int))},
		},
		{
			"2",
			args{args: []reflect.Value{
				reflect.ValueOf(func(a, b int) int { return a + b }),
				reflect.ValueOf([]interface{}{1, 2}),
			}},
			[]reflect.Value{reflect.ValueOf(make(chan int))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsyncCall(tt.args.args); !reflect.DeepEqual(got[0].Type().Name(), tt.want[0].Type().Name()) {
				t.Errorf("AsyncCall() = %v, want %v", got[0].Type().Name(), tt.want[0].Type().Name())
			}
		})
	}
}
