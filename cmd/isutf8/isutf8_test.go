package isutf8

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, compress bool, content string) string {
	path := "isutf8-test-*.txt"
	if compress {
		path += ".gz"
	}
	temp, err := os.CreateTemp("", path)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = temp.Close()
		_ = os.Remove(temp.Name())
	})

	w := io.WriteCloser(temp)
	if compress {
		w = gzip.NewWriter(temp)
	}

	_, err = w.Write([]byte(content))
	require.NoError(t, err)
	if compress {
		require.NoError(t, w.Close())
	}
	require.NoError(t, temp.Close())

	return temp.Name()
}

func TestIsUTF8(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr require.ErrorAssertionFunc
	}{
		{
			"all plain",
			[]string{
				tempFile(t, false, "abc"),
				tempFile(t, false, "def"),
			},
			require.NoError,
		},
		{
			"all compressed",
			[]string{
				tempFile(t, true, "abc"),
				tempFile(t, true, "def"),
			},
			require.Error,
		},
		{
			"not latin",
			[]string{tempFile(t, false, "世界")},
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
		})
	}
}
