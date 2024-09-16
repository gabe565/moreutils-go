package editor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		want    []string
		wantErr require.ErrorAssertionFunc
	}{
		{"default", "", []string{"vim"}, require.NoError},
		{"vim", "/usr/bin/vim", []string{"/usr/bin/vim"}, require.NoError},
		{"vscode", "code --wait --new-window", []string{"code", "--wait", "--new-window"}, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("EDITOR", tt.env)
			got, err := Get()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
