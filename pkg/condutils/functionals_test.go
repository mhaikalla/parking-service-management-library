package condutils

import (
	"errors"
	"reflect"
	"testing"
)

func TestWhen(t *testing.T) {
	type args struct {
		condition   bool
		expressions []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"test 1",
			args{
				condition: true,
				expressions: []interface{}{
					func(arg string) string {
						return arg
					},
					"foobar",
				},
			},
			"foobar",
		},
		{
			"test 2",
			args{
				condition: true,
				expressions: []interface{}{
					func(arg string) (string, error) {
						return arg, nil
					},
					"foobar",
				},
			},
			[]interface{}{"foobar", nil},
		},
		{
			"test 3",
			args{
				condition: true,
				expressions: []interface{}{
					func(arg string) (string, error) {
						return arg, errors.New("test")
					},
					"foobar",
				},
			},
			[]interface{}{"foobar", errors.New("test")},
		},
		{
			"test 4",
			args{
				condition: true,
				expressions: []interface{}{
					"foo",
				},
			},
			errors.New("func is first entry to expressions"),
		},
		{
			"test 5",
			args{
				condition: true,
				expressions: []interface{}{
					func() string { return "foo" },
				},
			},
			"foo",
		},
		{
			"test 6",
			args{
				condition:   true,
				expressions: []interface{}{},
			},
			errors.New("atleast 1 value provided on expressions"),
		},
		{
			"test 7",
			args{
				condition: false,
				expressions: []interface{}{
					func() int { return 1 },
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := When(tt.args.condition, tt.args.expressions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("When() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEval(t *testing.T) {
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"test 1",
			args{arg: func() string { return "foo" }},
			"foo",
		},
		{
			"test 2",
			args{arg: "foo"},
			"foo",
		},
		{
			"test 3",
			args{arg: 1},
			1,
		},
		{
			"test 4",
			args{arg: &struct{ name string }{"foo"}},
			&struct{ name string }{"foo"},
		},
		{
			"test 5",
			args{arg: errors.New("test")},
			errors.New("test"),
		},
		{
			"test 6",
			args{arg: func() (int, int) { return 1, 1 }},
			[]interface{}{1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Eval(tt.args.arg); !reflect.DeepEqual(got, tt.want) {
				t.Log(tt.name, reflect.DeepEqual(got, tt.want))
				t.Log(tt.name, reflect.ValueOf(got).Kind())
				t.Log(tt.name, reflect.ValueOf(tt.want).Kind())
				t.Log(Eval(tt.args.arg))
				t.Errorf("Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOr(t *testing.T) {
	type args struct {
		testValue    interface{}
		defaultValue interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"test 1", args{testValue: "", defaultValue: "foo"}, "foo"},
		{"test 2", args{testValue: 0, defaultValue: 10}, 10},
		{"test 3", args{testValue: struct{ name string }{}, defaultValue: "foo"}, "foo"},
		{"test 4", args{testValue: nil, defaultValue: 1}, 1},
		{"test 5", args{testValue: 0, defaultValue: 1}, 1},
		{"test 6", args{testValue: 10, defaultValue: 1}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Or(tt.args.testValue, tt.args.defaultValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrs(t *testing.T) {
	type args struct {
		val1   interface{}
		val2   interface{}
		values []interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{"1", args{values: []interface{}{1, 0}}, 1},
		{"2", args{values: []interface{}{1, 0, 10}}, 1},
		{"3", args{values: []interface{}{0, false, nil, true}}, true},
		{"3", args{values: []interface{}{0, false, errors.New("test"), true}}, errors.New("test")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ors(tt.args.val1, tt.args.val2, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstantly(t *testing.T) {
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name         string
		args         args
		suplliedArgs interface{}
		want         interface{}
	}{
		{"test 1", args{arg: 1}, "foo", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Constantly(tt.args.arg); !reflect.DeepEqual(got(tt.suplliedArgs), tt.want) {
				t.Errorf("Constantly() = %v, want %v", got(tt.suplliedArgs), tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args struct {
		fn    func(arg interface{}) interface{}
		colls []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			"test 1",
			args{
				fn:    func(s interface{}) interface{} { return s },
				colls: []interface{}{"foo", "bar"},
			},
			[]interface{}{"foo", "bar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.fn, tt.args.colls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Constantly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeNth(t *testing.T) {
	type args struct {
		n    int
		seqs []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{"nil seqs", args{n: 1, seqs: nil}, []interface{}{}},
		{"empty seqs", args{n: 10, seqs: []interface{}{}}, []interface{}{}},
		{"1 item", args{n: 1, seqs: []interface{}{1}}, []interface{}{1}},
		{"nth out of bound", args{n: 10, seqs: []interface{}{1}}, []interface{}{}},
		{"nth 0", args{n: 0, seqs: []interface{}{1, 2, 3, 4}}, []interface{}{1}},
		{"nth -1", args{n: 0, seqs: []interface{}{1, 2, 3, 4}}, []interface{}{1}},
		{"nth 2", args{n: 2, seqs: []interface{}{1, 2, 3, 4}}, []interface{}{2, 4}},
		{"nth 2, even seqs", args{n: 2, seqs: []interface{}{1, 2, 3}}, []interface{}{2}},
		{"nth 1, return seqs", args{n: 1, seqs: []interface{}{1, 2, 3}}, []interface{}{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TakeNth(tt.args.n, tt.args.seqs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeNth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFnName(t *testing.T) {
	type args struct {
		fn interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test 1",
			args{fn: getFnName},
			"github.com/mhaikalla/parking-service-management-library/pkg/condutils.getFnName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFnName(tt.args.fn); got != tt.want {
				t.Errorf("getFnName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvalSliceReflectValue(t *testing.T) {
	type args struct {
		values []reflect.Value
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"test 1",
			args{values: []reflect.Value{reflect.ValueOf(1), reflect.ValueOf(1)}},
			[]interface{}{1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EvalSliceReflectValue(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvalSliceReflectValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvalExpr(t *testing.T) {
	type args struct {
		fnAndArgs []interface{}
	}
	tests := []struct {
		name         string
		args         args
		wantReturned []interface{}
	}{
		{
			"test 1",
			args{fnAndArgs: []interface{}{
				func(a, b int) int { return a + b },
				1,
				10,
			}},
			[]interface{}{11},
		},
		{
			"test 2",
			args{fnAndArgs: []interface{}{
				func(a, b int) (int, int) { return a, b },
				1,
				10,
			}},
			[]interface{}{1, 10},
		},
		{
			"test 3",
			args{fnAndArgs: []interface{}{
				func() {},
				1,
				2,
			}},
			[]interface{}{errors.New("call github.com/mhaikalla/parking-service-management-library/pkg/condutils.TestEvalExpr.func3, panicked: reflect: Call with too many input arguments")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotReturned := EvalExpr(tt.args.fnAndArgs...); !reflect.DeepEqual(EvalSliceReflectValue(gotReturned), tt.wantReturned) {
				t.Errorf("EvalExpr() = %v, want %v", EvalSliceReflectValue(gotReturned), tt.wantReturned)
			}
		})
	}
}

func TestSliceOfInterface(t *testing.T) {
	type args struct {
		maybeSlice interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{"test 1", args{maybeSlice: 1}, []interface{}{1}},
		{"test 2", args{maybeSlice: "foo"}, []interface{}{"foo"}},
		{"test 3", args{maybeSlice: []string{"foo", "bar"}}, []interface{}{"foo", "bar"}},
		{"test 4", args{maybeSlice: []int{1, 2}}, []interface{}{1, 2}},
		{
			"test 5",
			args{
				maybeSlice: []reflect.Value{
					reflect.ValueOf(1),
					reflect.ValueOf("foobar"),
					reflect.ValueOf(errors.New("test")),
				},
			},
			[]interface{}{1, "foobar", errors.New("test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceOfInterface(tt.args.maybeSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceOfInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberFormat(t *testing.T) {
	type args struct {
		number       float64
		decimals     uint
		decPoint     string
		thousandsSep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{number: float64(12), decimals: 1, decPoint: ".", thousandsSep: "."}, "12.0"},
		{"2", args{number: float64(-12), decimals: 1, decPoint: ".", thousandsSep: "."}, "-12.0"},
		{"3", args{number: float64(100.99999), decimals: 0, decPoint: ".", thousandsSep: "."}, "101"},
		{"4", args{number: float64(1000000.99999), decimals: 0, decPoint: ".", thousandsSep: "."}, "1.000.001"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberFormat(tt.args.number, tt.args.decimals, tt.args.decPoint, tt.args.thousandsSep); got != tt.want {
				t.Errorf("NumberFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
