package combine

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, content string) string {
	temp, err := os.CreateTemp(t.TempDir(), "combine-test-*.txt")
	require.NoError(t, err)

	_, err = temp.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, temp.Close())

	return temp.Name()
}

func TestCombine(t *testing.T) {
	tests := []struct {
		name    string
		f1      string
		f2      string
		op      operator
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"and", "a\nb\nc\n", "c\nd\ne\n", operatorAnd, "c\n", require.NoError},
		{"not", "a\nb\nc\n", "c\nd\ne\n", operatorNot, "a\nb\n", require.NoError},
		{"or", "a\nb\nc\n", "c\nd\ne\n", operatorOr, "a\nb\nc\nc\nd\ne\n", require.NoError},
		{"xor", "a\nb\nc\n", "c\nd\ne\n", operatorXor, "a\nb\nd\ne\n", require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f1 := tempFile(t, tt.f1)
			f2 := tempFile(t, tt.f2)

			cmd := New()
			cmd.SetArgs([]string{f1, tt.op.String(), f2})
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.want, stdout.String())
		})
	}
}
