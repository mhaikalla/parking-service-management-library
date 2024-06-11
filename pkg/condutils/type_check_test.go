package condutils

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mhaikalla/parking-service-management-library/pkg/errs"
)

func TestIsEmpty(t *testing.T) {
	var emptyTime time.Time
	rfEmptyStr := reflect.ValueOf("")
	rfNil := reflect.ValueOf(nil)
	rfZeroNumber := reflect.ValueOf(0)
	rfStr := reflect.ValueOf("foo")
	rfNumber := reflect.ValueOf(1234)
	intVar := 1
	stringVar := "foo"
	zeroVar := 0
	emptyString := ""
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{arg: nil}, true},
		{"test 2", args{arg: false}, true},
		{"test 3", args{arg: true}, false},
		{"test 4", args{arg: 0}, true},
		{"test 5", args{arg: 1}, false},
		{"test 6", args{arg: &struct{ name string }{"foo"}}, false},
		{"test 7", args{arg: &struct{ name string }{}}, true},
		{"test 8", args{arg: struct{ name string }{"foo"}}, false},
		{"test 9", args{arg: struct{ name string }{}}, true},
		{"test 10", args{arg: ""}, true},
		{"test 11", args{arg: "foo"}, false},
		{"test 12", args{arg: []string{"foo"}}, false},
		{"test 13", args{arg: []string{}}, true},
		{"test 14", args{arg: &intVar}, false},
		{"test 15", args{arg: &stringVar}, false},
		{"test 16", args{arg: &zeroVar}, false},
		{"test 17", args{arg: &emptyString}, false},
		{"test 18", args{arg: struct{ Name struct{ First string } }{}}, true},
		{"test 19", args{arg: &struct{ Name struct{ First string } }{}}, true},
		{"test 20", args{arg: struct{ Name struct{ First string } }{struct{ First string }{"foo"}}}, false},
		{"test 21", args{arg: &struct{ Name struct{ First string } }{struct{ First string }{"foo"}}}, false},
		{"test 22", args{arg: rfEmptyStr}, true},
		{"test 23", args{arg: rfNil}, true},
		{"test 24", args{arg: rfZeroNumber}, true},
		{"test 25", args{arg: rfStr}, false},
		{"test 26", args{arg: rfNumber}, false},
		{"test 27", args{arg: emptyTime}, true},
		{"test 28", args{arg: time.Time{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.arg); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsError(t *testing.T) {
	type args struct {
		maybeErr interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 interface{}
	}{
		{
			"error",
			args{maybeErr: errors.New("test")},
			true,
			errors.New("test"),
		},
		{
			"integer",
			args{maybeErr: 64},
			false,
			64,
		},
		{
			"string",
			args{maybeErr: "foobar"},
			false,
			"foobar",
		},
		{
			"map structure",
			args{maybeErr: map[string]interface{}{"Error": "test"}},
			false,
			map[string]interface{}{"Error": "test"},
		},
		{
			"slice",
			args{maybeErr: []interface{}{"Error", "test"}},
			false,
			[]interface{}{"Error", "test"},
		},
		{
			"reflect.Value",
			args{maybeErr: reflect.ValueOf(errors.New("test"))},
			true,
			errors.New("test"),
		},
		{
			"nil",
			args{maybeErr: nil},
			false,
			nil,
		},
		{
			"integer",
			args{maybeErr: 1},
			false,
			1,
		},
		{
			"pointer",
			args{maybeErr: &struct{ Name string }{}},
			false,
			&struct{ Name string }{},
		},
		{
			"pointer with field Error",
			args{maybeErr: &struct{ Error error }{}},
			false,
			&struct{ Error error }{},
		},
		{
			"package github.com/mhaikalla/parking-service-management-library/pkg/errs",
			args{maybeErr: errs.NewErrContext().SetError(errors.New("test"))},
			true,
			errs.NewErrContext().SetError(errors.New("test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsError(tt.args.maybeErr)
			if got != tt.want {
				t.Errorf("IsError() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsError() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestIsFunc(t *testing.T) {
	type args struct {
		maybeFunc interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybeFunc: nil}, false},
		{"test 2", args{maybeFunc: "foo"}, false},
		{"test 3", args{maybeFunc: 1}, false},
		{"test 4", args{maybeFunc: func() {}}, true},
		{"test 5", args{maybeFunc: func() string { return "foo" }}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFunc(tt.args.maybeFunc); got != tt.want {
				t.Errorf("IsFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPtr(t *testing.T) {
	type args struct {
		maybePtr interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybePtr: nil}, false},
		{"test 2", args{maybePtr: "foo"}, false},
		{"test 3", args{maybePtr: 1}, false},
		{"test 4", args{maybePtr: func() {}}, false},
		{"test 5", args{maybePtr: func() string { return "foo" }}, false},
		{"test 6", args{maybePtr: &struct{ Name string }{}}, true},
		{"test 7", args{maybePtr: &map[string]int{"a": 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPtr(tt.args.maybePtr); got != tt.want {
				t.Errorf("IsPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsComplete(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"test 1",
			args{target: struct{ Name string }{}},
			false,
		},
		{
			"test 2",
			args{target: 1},
			false,
		},
		{
			"test 3",
			args{target: struct{ Name string }{"foo"}},
			true,
		},
		{
			"test 4",
			args{target: struct {
				Name string
				age  int
			}{"foo", 0}},
			false,
		},
		{
			"test 5",
			args{target: &struct {
				Name string
				age  int
			}{"foo", 0}},
			false,
		},
		{
			"test 6",
			args{target: &struct {
				Name string
				age  int
			}{"foo", 1}},
			true,
		},
		{
			"test 7",
			args{target: &struct {
				Name string
				age  int
			}{"", 1}},
			false,
		},
		{
			"test 8",
			args{target: nil},
			false,
		},
		{
			"test 9",
			args{target: &struct {
				Name string
				do   func()
			}{"bar", nil}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsComplete(tt.args.target); got != tt.want {
				t.Errorf("IsComplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
	type args struct {
		maybeMap interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybeMap: nil}, false},
		{"test 2", args{maybeMap: "foo"}, false},
		{"test 3", args{maybeMap: 1}, false},
		{"test 4", args{maybeMap: func() {}}, false},
		{"test 5", args{maybeMap: func() string { return "foo" }}, false},
		{"test 6", args{maybeMap: map[string]bool{}}, true},
		{"test 7", args{maybeMap: []string{}}, false},
		{"test 8", args{maybeMap: &map[string]bool{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMap(tt.args.maybeMap); got != tt.want {
				t.Errorf("IsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSlice(t *testing.T) {
	type args struct {
		maybeSlice interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybeSlice: nil}, false},
		{"test 2", args{maybeSlice: "foo"}, false},
		{"test 3", args{maybeSlice: 1}, false},
		{"test 4", args{maybeSlice: func() {}}, false},
		{"test 5", args{maybeSlice: func() string { return "foo" }}, false},
		{"test 6", args{maybeSlice: map[string]bool{}}, false},
		{"test 7", args{maybeSlice: []string{}}, true},
		{"test 8", args{maybeSlice: &map[string]bool{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSlice(tt.args.maybeSlice); got != tt.want {
				t.Errorf("IsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCompound(t *testing.T) {
	type args struct {
		maybeCompound interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybeCompound: nil}, false},
		{"test 2", args{maybeCompound: "foo"}, false},
		{"test 3", args{maybeCompound: 1}, false},
		{"test 4", args{maybeCompound: func() {}}, false},
		{"test 5", args{maybeCompound: func() string { return "foo" }}, false},
		{"test 6", args{maybeCompound: map[string]bool{}}, true},
		{"test 7", args{maybeCompound: []string{}}, true},
		{"test 8", args{maybeCompound: &map[string]bool{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCompound(tt.args.maybeCompound); got != tt.want {
				t.Errorf("IsCompound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescEmptyField(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"test 1",
			args{target: struct{ Name string }{}},
			[]string{`field: Name; type of: string; has value: ""`},
		},
		{
			"test 2",
			args{target: 1},
			[]string{},
		},
		{
			"test 3",
			args{target: struct{ Name string }{"foo"}},
			[]string{},
		},
		{
			"test 4",
			args{target: struct {
				Name string
				age  int
			}{"foo", 0}},
			[]string{`field: age; type of: int; has value: 0`},
		},
		{
			"test 5",
			args{target: &struct {
				Name string
				age  int
			}{"foo", 0}},
			[]string{`field: age; type of: int; has value: 0`},
		},
		{
			"test 6",
			args{target: &struct {
				Name string
				age  int
			}{"foo", 1}},
			[]string{},
		},
		{
			"test 7",
			args{target: &struct {
				Name string
				age  int
			}{"", 1}},
			[]string{`field: Name; type of: string; has value: ""`},
		},
		{
			"test 8",
			args{target: nil},
			[]string{},
		},
		{
			"test 9",
			args{target: &struct {
				Name string
				do   func()
			}{"bar", nil}},
			[]string{`field: do; type of: func(); has value: (func())(nil)`},
		},
		{
			"test 10",
			args{target: &struct {
				Name string
				Age  int
			}{}},
			[]string{
				`field: Name; type of: string; has value: ""`,
				`field: Age; type of: int; has value: 0`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DescEmptyField(tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DescEmptyField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStruct(t *testing.T) {
	var emptyTime time.Time
	rfEmptyStr := reflect.ValueOf("")
	rfNil := reflect.ValueOf(nil)
	rfZeroNumber := reflect.ValueOf(0)
	rfStr := reflect.ValueOf("foo")
	rfNumber := reflect.ValueOf(1234)
	intVar := 1
	stringVar := "foo"
	zeroVar := 0
	emptyString := ""
	type args struct {
		maybeStruct interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test 1", args{maybeStruct: nil}, false},
		{"test 2", args{maybeStruct: false}, false},
		{"test 3", args{maybeStruct: true}, false},
		{"test 4", args{maybeStruct: 0}, false},
		{"test 5", args{maybeStruct: 1}, false},
		{"test 6", args{maybeStruct: &struct{ name string }{"foo"}}, true},
		{"test 7", args{maybeStruct: &struct{ name string }{}}, true},
		{"test 8", args{maybeStruct: struct{ name string }{"foo"}}, true},
		{"test 9", args{maybeStruct: struct{ name string }{}}, true},
		{"test 10", args{maybeStruct: ""}, false},
		{"test 11", args{maybeStruct: "foo"}, false},
		{"test 12", args{maybeStruct: []string{"foo"}}, false},
		{"test 13", args{maybeStruct: []string{}}, false},
		{"test 14", args{maybeStruct: &intVar}, false},
		{"test 15", args{maybeStruct: &stringVar}, false},
		{"test 16", args{maybeStruct: &zeroVar}, false},
		{"test 17", args{maybeStruct: &emptyString}, false},
		{"test 18", args{maybeStruct: struct{ Name struct{ First string } }{}}, true},
		{"test 19", args{maybeStruct: &struct{ Name struct{ First string } }{}}, true},
		{"test 20", args{maybeStruct: struct{ Name struct{ First string } }{struct{ First string }{"foo"}}}, true},
		{"test 21", args{maybeStruct: &struct{ Name struct{ First string } }{struct{ First string }{"foo"}}}, true},
		{"test 22", args{maybeStruct: rfEmptyStr}, true},
		{"test 23", args{maybeStruct: rfNil}, true},
		{"test 24", args{maybeStruct: rfZeroNumber}, true},
		{"test 25", args{maybeStruct: rfStr}, true},
		{"test 26", args{maybeStruct: rfNumber}, true},
		{"test 27", args{maybeStruct: emptyTime}, true},
		{"test 28", args{maybeStruct: time.Time{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStruct(tt.args.maybeStruct); got != tt.want {
				t.Errorf("IsStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
