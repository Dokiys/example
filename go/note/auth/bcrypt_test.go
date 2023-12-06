package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	var pwd = []byte("123456")
	bcrypt, err := GenBcrypt(pwd)
	assert.NoError(t, err)

	ok, err := ValidateBcrypt(pwd, bcrypt)
	assert.NoError(t, err)
	assert.True(t, ok)
}
