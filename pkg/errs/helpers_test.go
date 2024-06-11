package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBubbleOriErr(t *testing.T) {
	type args struct {
		err *Errs
	}
	tests := []struct {
		name string
		args args
		want *Errs
	}{
		{
			"1",
			args{NewErrContext().SetError(errors.New("level 1"))},
			NewErrContext().SetError(errors.New("level 1")),
		},
		{
			"2",
			args{NewErrContext().SetError(NewErrContext().SetError(errors.New("level 2")))},
			NewErrContext().SetError(errors.New("level 2")),
		},
		{
			"3",
			args{NewErrContext().SetError(NewErrContext().SetError(NewErrContext().SetError(errors.New("level 3"))))},
			NewErrContext().SetError(errors.New("level 3")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BubbleOriErr(tt.args.err); got.Error() != tt.want.Error() {
				t.Errorf("BubbleOriErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaskingError(t *testing.T) {
	MaskingErrorMessage = "1"

	err2 := errors.New("")

	err1 := NewErrContext().SetCode(GeneralErrorHandler).SetError(err2)
	assert.Error(t, MaskingError(err1))
	assert.Error(t, MaskingError(err2))
	assert.Equal(t, MaskingError(err1), MaskingError(err2))
	assert.Nil(t, MaskingError(nil))
}
