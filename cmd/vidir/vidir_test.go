package vidir

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gabe565/moreutils/internal/cmdutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempFile(t *testing.T, dir string, name string) {
	temp, err := os.Create(filepath.Join(dir, name))
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = temp.Close()
		_ = os.Remove(temp.Name())
	})

	_, err = temp.WriteString(name)
	require.NoError(t, err)
	require.NoError(t, temp.Close())
}

func TestRun(t *testing.T) {
	temp := t.TempDir()

	wd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})
	require.NoError(t, os.Chdir(temp))

	tempFile(t, temp, "a")
	tempFile(t, temp, "b")
	tempFile(t, temp, "c")
	tempFile(t, temp, "d")

	// Swap a and b, rename c to newname, remove d
	t.Setenv("EDITOR", `sh -c 'cat > "$0" <<EOT
0002	a
0001	b
0003	newname
EOT'`)

	cmd := New(cmdutil.DisableTTY())
	var buf strings.Builder
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--verbose"})
	require.NoError(t, cmd.Execute())

	// Check verbose logs
	assert.Equal(t, `"a" => "a~"
"b" => "a"
"a~" => "b"
"c" => "newname"
`, buf.String())

	// Check dir contents
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	assert.Len(t, entries, 3)
	for _, entry := range entries {
		b, err := os.ReadFile(entry.Name())
		require.NoError(t, err)

		switch entry.Name() {
		case "a":
			assert.Equal(t, "b", string(b))
		case "b":
			assert.Equal(t, "a", string(b))
		case "newname":
			assert.Equal(t, "c", string(b))
		default:
			t.Error("unexpected entry:", entry.Name())
		}
	}
}
