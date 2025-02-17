package isutf8

import (
	"compress/gzip"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, compress bool, content string) string {
	path := "isutf8-test-*.txt"
	if compress {
		path += ".gz"
	}
	temp, err := os.CreateTemp(t.TempDir(), path)
	require.NoError(t, err)

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
	var bad strings.Builder
	gzw := gzip.NewWriter(&bad)
	_, err := gzw.Write([]byte("test"))
	require.NoError(t, err)
	require.NoError(t, gzw.Close())

	tests := []struct {
		name       string
		args       []string
		stdin      io.Reader
		wantStdout string
		wantErr    require.ErrorAssertionFunc
	}{
		{
			"all plain",
			[]string{
				tempFile(t, false, "abc"),
				tempFile(t, false, "def"),
			},
			nil,
			"",
			require.NoError,
		},
		{
			"all compressed",
			[]string{
				tempFile(t, true, "abc"),
				tempFile(t, true, "def"),
			},
			nil,
			`.+\/isutf8-test-\d+\.txt\.gz: line 1, char 1, byte 1
.+\/isutf8-test-\d+\.txt\.gz: line 1, char 1, byte 1
`,
			require.Error,
		},
		{
			"not latin",
			[]string{tempFile(t, false, "世界")},
			nil,
			"",
			require.NoError,
		},
		{
			"stdin plain",
			nil,
			strings.NewReader("abc"),
			"",
			require.NoError,
		},
		{
			"stdin compressed",
			nil,
			strings.NewReader(bad.String()),
			`\(standard input\): line 1, char 1, byte 1
`,
			require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			if tt.stdin != nil {
				cmd.SetIn(tt.stdin)
			}
			var stdout strings.Builder
			cmd.SetOut(&stdout)
			tt.wantErr(t, cmd.Execute())
			assert.Regexp(t, "^"+tt.wantStdout+"$", stdout.String())
		})
	}
}
