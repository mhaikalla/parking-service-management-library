package condutils

import "testing"

func TestIDRLayoutFormatter(t *testing.T) {
	type args struct {
		amount int64
		sep    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test < 1000 1st",
			args{amount: int64(1), sep: "."},
			"1",
		},
		{
			"test < 1000 2nd",
			args{amount: int64(123), sep: "."},
			"123",
		},
		{
			"test < 1000 3rd",
			args{amount: int64(999), sep: "."},
			"999",
		},
		{
			"test < 1000 3rd",
			args{amount: int64(0000), sep: "."},
			"0",
		},
		{
			"test > 999 1st",
			args{amount: int64(1000), sep: "."},
			"1.000",
		},
		{
			"test > 999 2nd",
			args{amount: int64(12345), sep: "."},
			"12.345",
		},
		{
			"test > 999 3rd",
			args{amount: int64(12345), sep: ","},
			"12,345",
		},
		{
			"test > 999 4th",
			args{amount: int64(1234567890), sep: "."},
			"1.234.567.890",
		},
		{
			"test > 999 4th",
			args{amount: int64(1234567890), sep: ","},
			"1,234,567,890",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IDRLayoutFormatter(tt.args.amount, tt.args.sep); got != tt.want {
				t.Errorf("IDRLayoutFormatter() = %v, want %v", got, tt.want)
			}
		})
	}
}
