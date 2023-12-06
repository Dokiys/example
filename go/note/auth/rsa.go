package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func RSAEncrypt(msg, pemPubKey []byte) (string, error) {
	pk, err := parsePublicKey(pemPubKey)
	if err != nil {
		return "", err
	}
	// 添加随机值以确保相同的 msg 不会产生同样的 ciphertext
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pk, msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RSADecrypt(encrypted string, pemPriKey []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	pk, err := parsePrivateKey(pemPriKey)
	if err != nil {
		return "", err
	}

	msg, err := rsa.DecryptPKCS1v15(nil, pk, ciphertext)
	if err != nil {
		return "", err
	}
	return string(msg), err
}

func RSASign(data []byte, pemPriKey []byte, hash crypto.Hash) (string, error) {
	h := hash.New()
	h.Write(data)
	var hashed = h.Sum(nil)

	pk, err := parsePrivateKey(pemPriKey)
	if err != nil {
		return "", err
	}

	// 先对原始消息进行哈希计算，然后再对哈希值进行签名
	bs, err := rsa.SignPKCS1v15(nil, pk, hash, hashed)

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bs), nil
}

func RSAVerify(src []byte, sign string, pemPubKey []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	var hashed = h.Sum(nil)

	signData, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	pk, err := parsePublicKey(pemPubKey)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(pk, hash, hashed, signData)
	if err != nil {
		return err
	}
	return nil
}

func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("PublicKey format error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if pub, ok := pubInterface.(*rsa.PublicKey); ok {
		return pub, nil
	}

	return nil, errors.New("PublicKey type error")
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("PrivateKey format error")
	}

	switch block.Type {
	case "RSA PRIVATE KEY", "PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, errors.New("PrivateKey type error")
	}
}
