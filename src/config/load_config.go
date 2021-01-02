package config

import (
	"encoding/json"
	appUtils "go-app-template/src/utils"
	"io/ioutil"
	"log"
	"path/filepath"
)

type configJSON struct {
	ProjectName string
}

var (
	Properties configJSON
)

const configFileName = "config.json"

func LoadConfig() {
	Properties = getConfigJSON()
}

func getConfigJSON() configJSON {
	configFilePath := appUtils.GetFilePathWithCurrentDir(configFileName)
	b, err := ioutil.ReadFile(filepath.Clean(configFilePath))
	if err != nil {
		log.Fatalf("設定ファイル読込エラー, Error: %v", err.Error())
	}

	var configJSON configJSON
	if err := json.Unmarshal(b, &configJSON); err != nil {
		log.Fatalf("設定ファイル読込エラー, 設定ファイルの構造がconfigJSONと異なっている可能性があります\nError: %v", err.Error())
	}

	return configJSON
}
