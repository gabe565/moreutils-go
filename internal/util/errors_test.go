package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinErrors(t *testing.T) {
	type args struct {
		errs []error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no error", args{nil}, false},
		{"one error", args{[]error{errors.New("test")}}, true},
		{"multiple", args{[]error{errors.New("test"), errors.New("test2")}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JoinErrors(tt.args.errs...)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
