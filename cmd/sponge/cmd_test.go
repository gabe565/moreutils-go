package sponge

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, content string) string {
	temp, err := os.CreateTemp("", "sponge-test-*.txt")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = temp.Close()
		_ = os.Remove(temp.Name())
	})

	_, err = temp.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, temp.Close())

	return temp.Name()
}

func TestSponge(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"replace", nil, "test\n", "test\n", require.NoError},
		{"append", []string{"--append"}, "test\n", "previous\ntest\n", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := tempFile(t, "previous\n")
			cmd := New()
			cmd.SetArgs(append(tt.args, temp))
			cmd.SetIn(strings.NewReader(tt.stdin))
			tt.wantErr(t, cmd.Execute())
			b, err := os.ReadFile(temp)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(b))
		})
	}
}
