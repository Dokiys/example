package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {
	var key = "123456"
	var text = "永远年轻，永远热泪盈眶"
	token, err := GenJwtToken([]byte(key), text, time.Minute)
	assert.NoError(t, err)

	var result string
	claims, err := ParseJwtToken[string](func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	}, token, result)
	assert.NoError(t, err)
	assert.Equal(t, text, claims.Data)
}
