package mispipe

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMispipe(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"success fail", []string{"echo success", "cat; exit 1"}, "success\n", require.NoError},
		{"fail success", []string{"echo fail; exit 1", "cat"}, "fail\n", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.want, stdout.String())
		})
	}
}
