package utils

import (
	"fmt"
	"testing"
)

const (
	key = "MQZJlxv3vq0nV7PL"
	iv  = "MQZJlxv3vq0nV7PL"
)

func TestEncrypt(t *testing.T) {
	//123456字符串加密
	t.Log(NewAes(key, iv).Encrypt("123456"))
	fmt.Println(NewAes(key, iv).Encrypt("123456")) //iVu1bcIHgITHxcFiZysjtw==
}

func TestDecrypt(t *testing.T) {
	//解密
	t.Log(NewAes(key, iv).Decrypt("iVu1bcIHgITHxcFiZysjtw=="))
	fmt.Println(NewAes(key, iv).Decrypt("iVu1bcIHgITHxcFiZysjtw==")) //123456
}

func BenchmarkEncryptAndDecrypt(b *testing.B) {
	b.ResetTimer()
	aes := NewAes(key, iv)
	for i := 0; i < b.N; i++ {
		encryptString, _ := aes.Encrypt("123456")
		aes.Decrypt(encryptString)
	}
}
