package subcommands

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChoose(t *testing.T) {
	t.Run("moreutils", func(t *testing.T) {
		cmd, err := Choose("moreutils")
		require.Error(t, err)
		assert.Nil(t, cmd)
	})

	t.Run("other", func(t *testing.T) {
		cmd, err := Choose("other")
		require.Error(t, err)
		assert.Nil(t, cmd)
	})

	for _, sub := range All() {
		t.Run(sub.Name(), func(t *testing.T) {
			cmd, err := Choose(sub.Name())
			require.NoError(t, err)
			assert.Equal(t, sub.Name(), cmd.Name())
		})
	}

	t.Run("combine alias", func(t *testing.T) {
		cmd, err := Choose("combine")
		require.NoError(t, err)
		assert.Equal(t, "combine", cmd.Name())
	})

	t.Run("zrun prefix", func(t *testing.T) {
		cmd, err := Choose("zcat")
		require.NoError(t, err)
		assert.Equal(t, "zrun", cmd.Name())
	})
}

func TestWithout(t *testing.T) {
	assert.Len(t, slices.Collect(Without(nil)), len(All())-len(DefaultExcludes()))
}
