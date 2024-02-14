package changelog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_builder2B(t *testing.T) {
	a := func(s string) strings.Builder {
		var b strings.Builder
		b.WriteString(s)

		return b
	}
	tests := []struct {
		name string
		args strings.Builder
		want []byte
	}{
		{
			name: "builder2B",
			args: a("test"),
			want: []byte("test"),
		},
		{
			name: "builder2B",
			args: a("test\n  \n\ntest"),
			want: []byte("test\n\ntest"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, builder2B(tt.args), "builder2B(%v)", tt.args)
		})
	}
}
