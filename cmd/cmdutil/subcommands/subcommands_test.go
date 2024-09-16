package subcommands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChoose(t *testing.T) {
	t.Run("correct command is chosen", func(t *testing.T) {
		cmd, err := Choose("moreutils")
		assert.Error(t, err)
		assert.Nil(t, cmd)

		cmd, err = Choose("other")
		assert.Error(t, err)
		assert.Nil(t, cmd)

		for _, sub := range All() {
			cmd, err := Choose(sub.Name())
			assert.NoError(t, err)
			assert.Equal(t, sub.Name(), cmd.Name())
		}
	})
}
