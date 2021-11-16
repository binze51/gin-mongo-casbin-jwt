package config

import (
	"io/ioutil"

	"gitee.com/binze/binzekeji/pkg/casbinx"
	"gitee.com/binze/binzekeji/pkg/engine"
	"gitee.com/binze/binzekeji/pkg/httpserver"
	"gitee.com/binze/binzekeji/pkg/jwtauth"
	"gitee.com/binze/binzekeji/pkg/logger"
	"gitee.com/binze/binzekeji/pkg/mongors"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

//系统配置
type Config struct {
	engine.EngineConfig
	httpserver.HttpServConfig
	mongors.MongoRSConfig
	logger.ZapConfig
	jwtauth.JwtConfig
	casbinx.MongoAdapterConf
	Menuconf []Menu `yaml:"menuconf" json:"menus"`
	Auvconf  struct {
		ServeHost string `yaml:"serverhost"`
		ServePort string `yaml:"serverport"`
		BlockSize int    `yaml:"blocksize"`
	} `yaml:"auvconf"`
	AirPortConf struct {
		ServeHost string `yaml:"serverhost"`
		ServePort string `yaml:"serverport"`
		BlockSize int    `yaml:"blocksize"`
	} `yaml:"airportconf"`
}

//菜单
type Menu struct {
	Id          string `yaml:"id" json:"id"`
	Description string `yaml:"description" json:"description"`
	ActionId    int    `yaml:"actionId" json:"actionId"`
}

//FromYamlFile 解析yaml文件
func FromYamlFile(env string) *Config {
	var c Config
	switch env {
	case "dev":
		data, err := ioutil.ReadFile("./config/dev.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			panic(err)
		}
		gin.SetMode(gin.DebugMode)
	case "test":
		data, err := ioutil.ReadFile("./config/test.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			panic(err)
		}
		gin.SetMode(gin.TestMode)
	case "uat":
		data, err := ioutil.ReadFile("./config/uat.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			panic(err)
		}
		gin.SetMode(gin.ReleaseMode)
	case "prod":
		data, err := ioutil.ReadFile("./config/prod.yaml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			panic(err)
		}
		gin.SetMode(gin.ReleaseMode)
	}
	return &c
}
