package cmdutil

import "github.com/spf13/cobra"

type Option func(cmd *cobra.Command)

const DisableTTYAnnotation = "without"

func DisableTTY() Option {
	return func(cmd *cobra.Command) {
		if cmd.Annotations == nil {
			cmd.Annotations = make(map[string]string, 1)
		}
		cmd.Annotations[DisableTTYAnnotation] = "true"
	}
}
