package helpers

import "testing"

func TestIsEmailValid(t *testing.T) {
	type args struct {
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Check Email Format :: ", args{e: "xltouchpoint@gmail.com"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmailValid(tt.args.e); got != tt.want {
				t.Errorf(" IsEmailValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
