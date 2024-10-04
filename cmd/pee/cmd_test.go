package pee

import (
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type lockedWriter struct {
	w  io.Writer
	mu sync.Mutex
}

func (l *lockedWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

func TestPee(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdin   string
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{
			"cat md5 base64",
			[]string{"cat", "md5sum"},
			"test\n",
			"test\nd8e8fca2dc0f896fd7cb4cb0031ba249  -\n",
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := New()
			cmd.SetArgs(tt.args)
			cmd.SetIn(strings.NewReader(tt.stdin))
			var stdout strings.Builder
			cmd.SetOut(&lockedWriter{w: &stdout})
			tt.wantErr(t, cmd.Execute())
			assert.Equal(t, tt.want, stdout.String())
		})
	}
}
