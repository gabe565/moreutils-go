package parallel

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParallel(t *testing.T) {
	temp := t.TempDir()

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(temp))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(wd))
	})

	cmd := New()
	cmd.SetArgs([]string{"sh", "-c", `echo "$0" > "$0"`, "--", "a", "b", "c", "d"})
	var stdout strings.Builder
	cmd.SetOut(&stdout)
	require.NoError(t, cmd.Execute())
	assert.Equal(t, "", stdout.String())

	entries, err := os.ReadDir(".")
	require.NoError(t, err)
	assert.Len(t, entries, 4)
}
