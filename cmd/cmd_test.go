package cmd

import (
	"testing"

	"gabe565.com/moreutils/internal/cmdutil/subcommands"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("correct command is chosen", func(t *testing.T) {
		cmd := New("moreutils")
		assert.Equal(t, "moreutils", cmd.Name())

		cmd = New("other")
		assert.Equal(t, "moreutils", cmd.Name())

		for _, sub := range subcommands.All() {
			cmd := New(sub.Name())
			assert.Equal(t, sub.Name(), cmd.Name())
		}
	})
}
