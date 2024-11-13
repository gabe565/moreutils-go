package cmdutil

import (
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

const DisableTTYAnnotation = "disable-tty"

func DisableTTY() cobrax.Option {
	return func(cmd *cobra.Command) {
		if cmd.Annotations == nil {
			cmd.Annotations = make(map[string]string, 1)
		}
		cmd.Annotations[DisableTTYAnnotation] = "true"
	}
}
