package config

import (
	"go-app-template/src/apputil"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const _configFileName = "config.yml"

func GetConfig() map[interface{}]interface{} {
	m := loadConfigFile()
	res := make(map[interface{}]interface{}, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func loadConfigFile() map[interface{}]interface{} {
	configFilePath := apputil.GetFilePathWithCurrentDir(_configFileName)
	b, err := ioutil.ReadFile(filepath.Clean(configFilePath))
	if err != nil {
		log.Fatalf("設定ファイル読込エラー, Error: %v", err.Error())
	}

	m := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(b, &m); err != nil {
		log.Fatalf("設定ファイル読込エラー, Error: %v", err.Error())
	}

	return m
}
