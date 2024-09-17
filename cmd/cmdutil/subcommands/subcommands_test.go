package subcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChoose(t *testing.T) {
	t.Run("correct command is chosen", func(t *testing.T) {
		cmd, err := Choose("moreutils")
		require.Error(t, err)
		assert.Nil(t, cmd)

		cmd, err = Choose("other")
		require.Error(t, err)
		assert.Nil(t, cmd)

		for _, sub := range All() {
			cmd, err := Choose(sub.Name())
			require.NoError(t, err)
			assert.Equal(t, sub.Name(), cmd.Name())
		}
	})
}
