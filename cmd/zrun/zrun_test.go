package zrun

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
	path := "zrun-test-*.txt"
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

func setExecutable(t *testing.T, name string) {
	oldExec := os.Args[0]
	os.Args[0] = name
	t.Cleanup(func() {
		os.Args[0] = oldExec
	})
}

func TestZrun(t *testing.T) {
	for _, name := range []string{"zrun", "moreutils"} {
		t.Run(name, func(t *testing.T) {
			setExecutable(t, name)

			cmd := New()
			cmd.SetArgs([]string{
				"cat",
				tempFile(t, true, "compressed\n"),
				tempFile(t, false, "plain\n"),
			})
			var buf strings.Builder
			cmd.SetOut(&buf)
			require.NoError(t, cmd.Execute())
			assert.Equal(t, "compressed\nplain\n", buf.String())
		})
	}
}

func TestAlias(t *testing.T) {
	setExecutable(t, "zcat")

	cmd := New()
	cmd.SetArgs([]string{
		tempFile(t, true, "compressed\n"),
		tempFile(t, false, "plain\n"),
	})
	var buf strings.Builder
	cmd.SetOut(&buf)
	require.NoError(t, cmd.Execute())
	assert.Equal(t, "compressed\nplain\n", buf.String())
}
