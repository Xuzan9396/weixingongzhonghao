package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type HttpConfig struct {
	Port string
}

type WxConfig struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
}

type RedisConfig struct {
	Host     string
	Password string
	Database int
}

type UploadConfig struct {
	Path string
}

type ConfigJson struct {
	Http   HttpConfig   `json:"http"`
	Wx     WxConfig     `json:"wx"`
	Redis  RedisConfig  `json:"redis"`
	Upload UploadConfig `json:"upload"`
}

var G_JsonConfig *ConfigJson

func InitConfigJson() (err error) {

	if G_JsonConfig == nil {

		G_JsonConfig = &ConfigJson{}
		var (
			content  []byte
			filename string
		)

		if filename == "" {

			filename = "./config.json"
		} else {
			filename = strings.TrimSuffix(filename, "yaml") + "json"
		}
		// 1, 把配置文件读进来
		if content, err = ioutil.ReadFile(filename); err != nil {

			return
		}
		// 2, 做JSON反序列化
		if err = json.Unmarshal(content, &G_JsonConfig); err != nil {
			return
		}

	}

	return
}
