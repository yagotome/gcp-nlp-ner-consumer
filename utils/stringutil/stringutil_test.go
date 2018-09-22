package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitByPunctuation(t *testing.T) {
	expected := []string{"Foobar", "É um cara bem legal", "Pena que", "não pode ver mulher", "ela ta dançando"}
	actual := SplitByPunctuation("Foobar. É um cara bem legal!!! Pena que...\t não pode ver mulher!?\nela ta dançando.")
	assert.ElementsMatch(t, expected, actual)
}

func TestEscapeGroup(t *testing.T) {
	type args struct {
		group string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Single slash", args{"P/D"}, "P-D"},
		{"Single inverted slash", args{"P\\D"}, "P-D"},
		{"Multiple slashes", args{"A / / B/ D"}, "A - - B- D"},
		{"Multiple inverted slashes", args{"A\\B \\ C \\"}, "A-B - C -"},
		{"Mixed", args{"\\A/B \\ C/D"}, "-A-B - C-D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeGroup(tt.args.group); got != tt.want {
				t.Errorf("EscapeGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
