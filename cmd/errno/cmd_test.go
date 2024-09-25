package errno

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrno(t *testing.T) {
	if !Supported {
		t.SkipNow()
	}

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"number", []string{"32"}, "EPIPE 32 broken pipe\n", require.NoError},
		{"name", []string{"EPIPE"}, "EPIPE 32 broken pipe\n", require.NoError},
		{"search", []string{"-s", "broken"}, "EPIPE 32 broken pipe\n", require.NoError},
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

	t.Run("list", func(t *testing.T) {
		cmd := New()
		cmd.SetArgs([]string{"-ls"})
		var buf strings.Builder
		cmd.SetOut(&buf)
		require.NoError(t, cmd.Execute())
		assert.NotEmpty(t, buf.String())
	})
}

func TestErrnoUnsupported(t *testing.T) {
	if Supported {
		t.SkipNow()
	}

	cmd := New()
	require.Error(t, cmd.Execute())
}
