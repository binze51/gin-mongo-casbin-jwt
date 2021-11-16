package config

import (
	"fmt"
	"testing"
)

func TestFromYamlFile(t *testing.T) {
	envflag := "dev"
	cc := FromYamlFile(envflag)
	fmt.Printf("jwt:%+v", cc.JwtConfig)
	fmt.Printf("engineconf:%+v", cc.EngineConfig)
	fmt.Printf("h2servconf:%+v", cc.HttpServConfig)
	fmt.Printf("mongoadapter:%+v", cc.MongoAdapterConf)
	fmt.Printf("mongorsconf:%+v", cc.MongoRSConfig)
	fmt.Printf("ZapConfig:%+v", cc.ZapConfig)

	t.Log("ok")

}
