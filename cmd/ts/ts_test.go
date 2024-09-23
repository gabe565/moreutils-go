package ts

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"no format", nil, "test\n", time.Now().Format(time.DateTime) + " test\n", require.NoError},
		{"format stamp", []string{"%b %e %H:%M:%S"}, "test\n", time.Now().Format(time.Stamp) + " test\n", require.NoError},
		{"invalid format", []string{"%g"}, "test\n", "", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			cmd.SetIn(strings.NewReader(tt.stdin))
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.want, stdout.String())
		})
	}
}
