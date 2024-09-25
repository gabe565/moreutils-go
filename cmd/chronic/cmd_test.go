package chronic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChronic(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStdout string
		wantStderr string
		wantErr    require.ErrorAssertionFunc
	}{
		{
			"success",
			[]string{"sh", "-c", "echo stdout; echo stderr;"},
			"",
			"",
			require.NoError,
		},
		{
			"fail",
			[]string{"sh", "-c", "echo stdout; echo stderr >&2; exit 1;"},
			"stdout\n",
			"stderr\nError: exit status 1\n",
			require.Error,
		},
		{
			"flag stderr",
			[]string{"--stderr", "sh", "-c", "echo stdout; echo stderr >&2;"},
			"stdout\n",
			"stderr\n",
			require.NoError,
		},
		{
			"flag verbose",
			[]string{"--verbose", "sh", "-c", "echo stdout; echo stderr >&2; exit 1;"},
			"STDOUT:\nstdout\n\nSTDERR:\n\nRETVAL: 1\n",
			"stderr\nError: exit status 1\n",
			require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			var stderr strings.Builder
			cmd.SetErr(&stderr)
			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.wantStdout, stdout.String())
			assert.Equal(t, tt.wantStderr, stderr.String())
		})
	}
}
