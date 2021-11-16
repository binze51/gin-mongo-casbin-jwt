package utils

import (
	"bytes"
	cryptoAes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// var _ Aes = (*aes)(nil)

// type Aes interface {
// 	i()
// 	// Encrypt 加密
// 	Encrypt(encryptStr string) (string, error)

// 	// Decrypt 解密
// 	Decrypt(decryptStr string) (string, error)
// }

func NewAes(key, iv string) *Aes {
	return &Aes{
		key: key,
		iv:  iv,
	}
}

type Aes struct {
	key string
	iv  string
}

//Encrypt 字符串加密
func (a *Aes) Encrypt(encryptStr string) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = pkcs5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(a.iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

//Decrypt 密串解密成字符串
func (a *Aes) Decrypt(decryptStr string) (string, error) {
	decryptBytes, err := base64.URLEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(a.iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted = pkcs5UnPadding(decrypted)
	return string(decrypted), nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(decrypted []byte) []byte {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	return decrypted[:(length - unPadding)]
}
