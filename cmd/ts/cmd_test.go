package ts

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTs(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"no format", []string{}, "test\n", now.Format(time.DateTime) + " test\n", require.NoError},
		{"format stamp", []string{"%b %e %H:%M:%S"}, "test\n", `[A-Z][a-z]{2} [\d ]\d \d{2}:\d{2}:\d{2}` + " test\n", require.NoError},
		{"invalid format", []string{"%g"}, "test\n", "", require.Error},
		{"relative", []string{"-r"}, "INFO " + now.Add(-1*time.Hour).Format(time.RFC3339) + " abc", "INFO " + `1h0m\ds ago` + " abc\n", require.NoError},
		{"relative format", []string{"-r", "%a %b %e %T %Y"}, "INFO " + now.Add(-1*time.Hour).Format(time.RFC3339) + " abc", "INFO " + `[A-Z][a-z]{2} [A-Z][a-z]{2} [\d ]\d \d{2}:\d{2}:\d{2} \d{4}` + " abc\n", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			cmd.SetIn(strings.NewReader(tt.stdin))
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			assert.Regexp(t, "^"+tt.want+"$", stdout.String())
		})
	}
}

func Test_validArgs(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		cmd := New()
		got, directive := validArgs(cmd, nil, "")
		assert.NotEqual(t, cobra.ShellCompDirectiveError, directive)
		assert.NotEmpty(t, got)
	})

	t.Run("increment flag", func(t *testing.T) {
		cmd := New()
		require.NoError(t, cmd.Flags().Set(FlagIncrement, "true"))
		got, directive := validArgs(cmd, nil, "")
		assert.NotEqual(t, cobra.ShellCompDirectiveError, directive)
		assert.NotEmpty(t, got)
	})
}
