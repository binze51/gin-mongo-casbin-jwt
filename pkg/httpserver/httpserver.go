package httpserver

import (
	"fmt"
	"net/http"

	"gitee.com/binze/binzekeji/pkg/router"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HttpServConfig struct {
	Port        string `yaml:"port"`
	ServiceName string `yaml:"servicename"`
	InitRoot    bool   `yaml:"initroot"`
	RootPwd     string `yaml:"rootpwd"`
	Http2       bool   `yaml:"http2"`
}

//New 构建httpserver
func New(conf *HttpServConfig, gr *router.Router) *http.Server {
	if conf.Http2 {
		if gin.IsDebugging() {
			fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] "+fmt.Sprintf("Listening and serving H2C on %s\n", conf.Port))
		}
		h2cSer := &http2.Server{}

		return &http.Server{
			Addr:    ":" + conf.Port,
			Handler: h2c.NewHandler(gr, h2cSer),
		}
	}
	if gin.IsDebugging() {
		fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] "+fmt.Sprintf("Listening and serving HTTP on %s\n", conf.Port))
	}
	return &http.Server{
		Addr:    ":" + conf.Port,
		Handler: gr,
	}
}
