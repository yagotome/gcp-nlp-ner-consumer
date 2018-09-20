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
