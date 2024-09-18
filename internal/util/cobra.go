package util

import (
	"slices"

	"github.com/spf13/cobra"
)

func CmdsContains(cmds []*cobra.Command, target *cobra.Command) bool {
	return slices.ContainsFunc(cmds, func(cmd *cobra.Command) bool {
		return cmd.Name() == target.Name()
	})
}
