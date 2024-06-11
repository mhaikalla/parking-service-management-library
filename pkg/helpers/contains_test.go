package helpers

import "testing"

func TestStringContains(t *testing.T) {
	type args struct {
		slices     []string
		comparison string
	}
	exampleString := []string{"BILL", "PAYBILL", "FTTH", "TOPUP"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test success StringContains", args{comparison: "PAYBILL", slices: exampleString}, true},
		{"test false StringContains", args{comparison: "PAYBILL2", slices: exampleString}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringContains(tt.args.slices, tt.args.comparison); got != tt.want {
				t.Errorf("StringContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntContains(t *testing.T) {
	type args struct {
		slices     []int
		comparison int
	}
	exampleSlicesInt := []int{20000, 25000, 30000, 50000}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test success IntContains", args{slices: exampleSlicesInt, comparison: 20000}, true},
		{"test false IntContains", args{slices: exampleSlicesInt, comparison: 12345}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntContains(tt.args.slices, tt.args.comparison); got != tt.want {
				t.Errorf("IntContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Contains(t *testing.T) {
	type args struct {
		slices     []float64
		comparison int
	}
	exampleSlicesFloat64 := []float64{20000.0, 25000.0, 30000.0, 50000.0}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test success Float64Contains", args{slices: exampleSlicesFloat64, comparison: 20000}, true},
		{"test false Float64Contains", args{slices: exampleSlicesFloat64, comparison: 12345}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64Contains(tt.args.slices, tt.args.comparison); got != tt.want {
				t.Errorf("Float64Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	type args struct {
		a    string
		list []string
	}
	exampleStringInSlice := []string{"BILL", "PAYBILL", "FTTH", "TOPUP"}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test success StringInSlice", args{a: "PAYBILL", list: exampleStringInSlice}, true},
		{"test false StringInSlice", args{a: "PAYBILL2", list: exampleStringInSlice}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("StringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURLContains(t *testing.T) {
	type args struct {
		comparison string
		slices     []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URLContains(tt.args.comparison, tt.args.slices); got != tt.want {
				t.Errorf("URLContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInt(tt.args.s); got != tt.want {
				t.Errorf("IsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
