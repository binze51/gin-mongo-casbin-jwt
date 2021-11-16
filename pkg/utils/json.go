package utils

import (
	jsoniter "github.com/json-iterator/go"
)

// 定义JSON操作
var (
	json              = jsoniter.ConfigCompatibleWithStandardLibrary
	JsonMarshal       = json.Marshal
	JsonUnmarshal     = json.Unmarshal
	JsonMarshalIndent = json.MarshalIndent
	JsonNewDecoder    = json.NewDecoder
	JsonNewEncoder    = json.NewEncoder
)

// MarshalToString JSON编码为字符串
func JsonMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}
