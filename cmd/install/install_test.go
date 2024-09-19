package install

import (
	"os"
	"path/filepath"
	"slices"
	"syscall"
	"testing"

	"github.com/gabe565/moreutils/cmd/cmdutil/subcommands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	cmd := New()
	temp := t.TempDir()
	cmd.SetArgs([]string{temp})
	require.NoError(t, cmd.Execute())

	cmdCount := len(slices.Collect(subcommands.Without(subcommands.DefaultExcludes())))
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	require.Len(t, entries, cmdCount)

	exec, err := os.Executable()
	require.NoError(t, err)
	execStat, err := os.Stat(exec)
	require.NoError(t, err)
	execSys, ok := execStat.Sys().(*syscall.Stat_t)
	require.True(t, ok)

	for _, entry := range entries {
		stat, err := os.Lstat(filepath.Join(temp, entry.Name()))
		require.NoError(t, err)
		sys, ok := stat.Sys().(*syscall.Stat_t)
		require.True(t, ok)
		assert.Equal(t, execSys.Ino, sys.Ino)
	}
}

func TestSymlink(t *testing.T) {
	cmd := New()
	temp := t.TempDir()
	cmd.SetArgs([]string{"--symbolic", temp})
	require.NoError(t, cmd.Execute())

	cmdCount := len(slices.Collect(subcommands.Without(subcommands.DefaultExcludes())))
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	require.Len(t, entries, cmdCount)

	exec, err := os.Executable()
	require.NoError(t, err)

	for _, entry := range entries {
		link, err := os.Readlink(filepath.Join(temp, entry.Name()))
		require.NoError(t, err)
		assert.Equal(t, exec, link)
	}
}

func TestRelative(t *testing.T) {
	cmd := New()
	temp := t.TempDir()
	cmd.SetArgs([]string{"--symbolic", "--relative", temp})
	require.NoError(t, cmd.Execute())

	cmdCount := len(slices.Collect(subcommands.Without(subcommands.DefaultExcludes())))
	entries, err := os.ReadDir(temp)
	require.NoError(t, err)
	require.Len(t, entries, cmdCount)

	exec, err := os.Executable()
	require.NoError(t, err)

	execRel, err := filepath.Rel(temp, exec)
	require.NoError(t, err)

	for _, entry := range entries {
		link, err := os.Readlink(filepath.Join(temp, entry.Name()))
		require.NoError(t, err)
		assert.Equal(t, execRel, link)
	}
}
