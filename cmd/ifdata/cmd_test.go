package ifdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfdataUsage(t *testing.T) {
	cmd := New()
	assert.NotPanics(t, func() {
		assert.NoError(t, cmd.Usage())
	})
}
