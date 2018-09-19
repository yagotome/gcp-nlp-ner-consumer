package nlp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntitiesMatch(t *testing.T) {
	entities := make(Entities, 3)
	entities[0] = &Entity{Name: "Brasil"}
	entities[1] = &Entity{Name: "Rio de Janeiro"}
	entities[2] = &Entity{Name: "Universidade do Estado do Rio de Janeiro"}
	words := strings.Split("O Rio do Rio de Janeiro que fica no Brasil continua lindo", " ")

	assert.Equal(t, entities[0], entities.Match(words, 9))
	assert.Nil(t, entities.Match(words, 0))
	assert.Equal(t, entities[1], entities.Match(words, 3))
	assert.Nil(t, entities.Match(words, 1))
}
