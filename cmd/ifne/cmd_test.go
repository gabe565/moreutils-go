package ifne

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIfne(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"not empty", []string{"sh", "-c", "cat; echo ran"}, "input\n", "input\nran\n", require.NoError},
		{"empty", []string{"echo", "ran"}, "", "", require.NoError},
		{"invert not empty", []string{"--invert", "sh", "-c", "cat; echo ran"}, "input\n", "", require.NoError},
		{"invert empty", []string{"--invert", "echo", "ran"}, "", "ran\n", require.NoError},
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
