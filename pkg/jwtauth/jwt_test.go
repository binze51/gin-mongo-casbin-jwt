package jwtauth

import (
	"context"
	"fmt"
	"testing"
)

var JWT_RSAPRIVATE string = `-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJBAKPYT/wbFM4hF4cExiWLVg61CRpJeiHvanM3k/k219Rv8meG0X/l
7oi+T8Uaq9eWpXy0iA/S0uwgk6DhPX72q80CAwEAAQJAKvsLwG66Pnif22N9N0wd
/y2ufZ7Y0x4jJqZlwvKIG8n2sCxR4c6M3baLR748YFGZWhGA4uNESANgzliz4oRX
AQIhAM69A8CPDcS0FS335anY1Fb4woVI23V1enTQ7akY+XO1AiEAyuLQOosa5Kyz
9Se4OO7yDQ+i/KLtunOOhXj3t0EOlrkCIFTuuhflrVZeVUUpTqTUe4evctqm7+H1
fXV4T+rkY7bxAiBX7B0TEc9wxAsktbbXLW3GDT2zwCPHxmZAH2EykEXzOQIgZW0z
EPpITA0X1xnGYyZLJSuXzejlX2S5vVZ8rBfgc7s=
-----END RSA PRIVATE KEY-----`

var JWT_PUBLIC string = `-----BEGIN PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAKPYT/wbFM4hF4cExiWLVg61CRpJeiHv
anM3k/k219Rv8meG0X/l7oi+T8Uaq9eWpXy0iA/S0uwgk6DhPX72q80CAwEAAQ==
-----END PUBLIC KEY-----`

func TestSiginjwt(t *testing.T) {
	conf := JwtConfig{}
	// 只提供pem私钥，该私钥对应的JKW会给网关使用，JWKURL
	conf.RSAPrivateSecret = JWT_RSAPRIVATE
	conf.RSAPublicSecret = JWT_PUBLIC
	conf.Expired = 86400
	authTest := New(&conf)
	token, err := authTest.IssueToken(context.Background(), false, "tian", "heming-0", "", "", "", "admin", false)
	fmt.Printf("%+v\n", token)
	fmt.Printf("%+v\n", err)
	fmt.Print("==========================================\n")
	userid, role, err := authTest.ParseUser(context.Background(), token.AccessToken)
	fmt.Printf("%+v\n", userid)
	fmt.Printf("%+v\n", role)
	fmt.Printf("%+v\n", err)

}
