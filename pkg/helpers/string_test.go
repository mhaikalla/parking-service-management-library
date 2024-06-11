package helpers

import (
	"fmt"
	"strings"
	"testing"
)

func Test_InsertSpace(t *testing.T) {
	type args struct {
		text   string
		length int
	}

	tests := []struct {
		name string
		args args
		isOk bool
	}{
		{"Normal Test", args{text: "542640XXXXXX9386", length: 4}, true},
	}

	for _, tt := range tests {
		got := InsertSpaceNth(tt.args.text, tt.args.length)
		got = strings.ReplaceAll(got, "X", "*")
		fmt.Println(got)
	}
}
