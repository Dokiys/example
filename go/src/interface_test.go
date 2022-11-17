package src

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterfaceEqual(t *testing.T) {
	var i interface{}
	i = "-"
	assert.Equal(t, "-", i)
	i = 0
	assert.Equal(t, 0, i)
}
