package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"github.com/tjfoc/gmsm/x509"
	"go.uber.org/zap"
)

var zlog *zap.Logger
var (
	publicStr1  = "04BB34D657EPIE8490E66EF577E6B3CEA28B739511E787WEF71B7F38F241D87F18A5A93DF74E90FF94F4EB907F271A36B295B851F971DA5418F4915E2C1A23D6E"
	privateStr1 = "0B1CE43928BC21B8E82B5C065EDB534CB86532B1900A49D49FCFT3762D2997FA"
)

type GM struct {
	PublicKey  *sm2.PublicKey
	PrivateKey *sm2.PrivateKey
}

func NewGM(str string) *GM {
	publicKey, err := x509.ReadPublicKeyFromHex(str)
	if err != nil {
		zlog.Fatal(err.Error())
		return nil
	}
	return &GM{
		PublicKey: publicKey,
	}
}

//国密实例 构造函数
func NewGMInstance(privateHexStr, publicHexStr string) *GM {
	publicKey, err := x509.ReadPublicKeyFromHex(publicHexStr)
	if err != nil {
		panic(err)
	}
	privateKey, err := x509.ReadPrivateKeyFromHex(privateHexStr)
	if err != nil {
		panic(err)
	}
	return &GM{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
func NewGM1(dir string) *GM {
	//pubPem,err:=ioutil.ReadFile(dir+ "./pk.pem")
	//if err!=nil{
	//	zlog.Fatal(err.Error())
	//}
	//priPem,err:=ioutil.ReadFile(dir+"./sk.pem")
	//if err!=nil{
	//	zlog.Fatal(err.Error())
	//}
	publicKey, err := x509.ReadPublicKeyFromHex(publicStr1)
	if err != nil {
		zlog.Fatal(err.Error())
	}
	privateKey, err := x509.ReadPrivateKeyFromHex(privateStr1)
	if err != nil {
		zlog.Fatal(err.Error())
	}
	return &GM{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}
func (g *GM) GenerateUUIDStr() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
func (g *GM) SM2EN(data []byte) (string, error) {
	do, err := sm2.Encrypt(g.PublicKey, data, rand.Reader, sm2.C1C3C2)
	if err != nil {
		zlog.Error(err.Error())
		return "", err
	}
	return hex.EncodeToString(do), nil
}
func (g *GM) SM2DE(data []byte) (string, error) {
	do, err := sm2.Decrypt(g.PrivateKey, data, sm2.C1C3C2)
	if err != nil {
		zlog.Error(err.Error())
		return "", err
	}
	return string(do), err
}
func (g *GM) SM3Sign(data []byte) (string, error) {
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum), nil
}
func (g *GM) SM3Verify(sign string, data []byte) bool {
	h := sm3.New()
	h.Write(data)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum) == sign
}
func (g *GM) SM4EcbEn(key, data []byte) (string, error) {
	iv := []byte("")
	sm4.SetIV(iv)
	ecbMsg, err := sm4.Sm4Ecb([]byte(key), data, true)
	if err != nil {
		zlog.Error(err.Error())
		return "", err
	}
	zlog.Info(fmt.Sprintf("ecbMsg: %x\n", ecbMsg))
	return hex.EncodeToString(ecbMsg), err
}
func (g *GM) SM4EcbDe(key, data []byte) ([]byte, error) {
	iv := []byte("")
	sm4.SetIV(iv)
	ecbDec, err := sm4.Sm4Ecb(key, data, false)
	if err != nil {
		return nil, err
	}
	return ecbDec, nil
}

//无填充堆成加密
func (g *GM) SM4NoPaddingDe(key []byte, in []byte) (out []byte, err error) {
	if len(key) != sm4.BlockSize {
		return nil, errors.New("SM4: invalid key size " + strconv.Itoa(len(key)))
	}
	inData := in
	out = make([]byte, len(inData))
	c, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(inData)/16; i++ {
		in_tmp := inData[i*16 : i*16+16]
		out_tmp := make([]byte, 16)
		c.Decrypt(out_tmp, in_tmp)
		copy(out[i*16:i*16+16], out_tmp)
	}
	return out, nil
}
