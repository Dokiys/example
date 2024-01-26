package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AESCBCPKCS5PaddingEncrypt(key, iv, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)

	text = PKCS7Padding(text, blockMode.BlockSize())

	var crypted = make([]byte, len(text))
	blockMode.CryptBlocks(crypted, text)

	return crypted, nil
}

func AESCBCPKCS5PaddingDecrypt(key, iv, crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	var text = make([]byte, len(crypted))
	blockMode.CryptBlocks(text, crypted)

	return PKCS7UnPadding(text)
}

func TestAESCBCPKCS5Padding(t *testing.T) {
	var key = []byte("I5Y424l2vx7/+v1N5gmvB3DSos/OdCUh")
	var iv = []byte("gUX4fNgBpxGUTHas")
	var originText = []byte("Hello Work!")

	crypted, err := AESCBCPKCS5PaddingEncrypt(key, iv, originText)
	assert.NoError(t, err)
	text, err := AESCBCPKCS5PaddingDecrypt(key, iv, crypted)
	assert.NoError(t, err)
	assert.Equal(t, originText, text)
}

// PKCS7Padding 实现PKCS7填充，这也可以用于PKCS5填充
func PKCS7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS7UnPadding 实现PKCS7去填充
func PKCS7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("unpadding error: input data is empty")
	}

	padding := int(src[length-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, fmt.Errorf("unpadding error: invalid padding value %v", padding)
	}

	return src[:length-padding], nil
}
