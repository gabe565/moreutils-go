package vipe

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVipe(t *testing.T) {
	var sed string
	if runtime.GOOS == "darwin" {
		sed = "sed -i '' 's/test/testing/'"
	} else {
		sed = "sed -i 's/test/testing/'"
	}

	tests := []struct {
		name     string
		editor   string
		args     []string
		stdin    string
		want     string
		wantIsRe bool
		wantErr  require.ErrorAssertionFunc
	}{
		{"run", sed, nil, "test\n", "testing\n", false, require.NoError},
		{"suffix", `sh -c 'echo "$0" >"$0"'`, []string{"--suffix=bin"}, "", ".*.bin\n", true, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", tt.editor)
			cmd := New()
			cmd.SetArgs(tt.args)
			cmd.SetIn(strings.NewReader(tt.stdin))
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			if tt.wantIsRe {
				assert.Regexp(t, tt.want, &stdout)
			} else {
				assert.Equal(t, tt.want, stdout.String())
			}
		})
	}
}
