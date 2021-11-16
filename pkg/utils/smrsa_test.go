package utils

import (
	"encoding/hex"
	"log"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/tjfoc/gmsm/sm4"
)

func TestSm2(t *testing.T) {
	source := "112121"
	log.Println([]byte(source))
	gm := NewGM1("/Users/fire/worksapce/gosvcpublic")
	result, err := gm.SM2EN([]byte(source))
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Logf("加密后内容: %v", result)
	result1 := result[2:]
	t.Logf("截取后的内容: %v", result1)
	data, _ := hex.DecodeString("04" + result1)
	resultDe, err := gm.SM2DE(data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("解密后内容: %v", resultDe)
}
func TestSM4(t *testing.T) {
	key := uuid.NewV4()
	t.Logf("字节长度: %v", len(key.Bytes()))
	t.Logf("十六进制长度: %v", len(key.String()))
	t.Logf("替换-后长度: %v", len(strings.ReplaceAll(key.String(), "-", "")))

	iv, _ := hex.DecodeString("31313131313131313131313131313131")
	data := "这个是中国"
	hexData := hex.EncodeToString([]byte(data))
	t.Logf("key: %v\n", key)
	t.Logf("替换后key: %v\n", strings.ReplaceAll(key.String(), "-", ""))
	t.Logf("data: %v\n", data)
	t.Logf("hex data: %v\n", hexData)
	_ = sm4.SetIV(iv)
	keyContent, _ := hex.DecodeString(strings.ReplaceAll(key.String(), "-", ""))
	ecbMsg, _ := sm4.Sm4Ecb(keyContent, []byte(hexData), true)
	t.Logf("加密后的值: %v", hex.EncodeToString(ecbMsg))
	ecbDec, err := sm4.Sm4Ecb(keyContent, ecbMsg, false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("解密后十六进制: %v", string(ecbDec))

}
