package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"
)

type configJSON struct {
	ProjectName string
}

var (
	Config configJSON
)

func LoadConfig() {
	Config = getConfigJSON()
}

func getConfigJSON() configJSON {
	configFilePath := getConfigFilePath()
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("設定ファイル読込エラー, Error: %v", err.Error())
	}

	var configJSON configJSON
	if err := json.Unmarshal(b, &configJSON); err != nil {
		log.Fatalf("設定ファイル読込エラー, 設定ファイルの構造がconfigJSONと異なっている可能性があります\nError: %v", err.Error())
	}

	return configJSON
}

func getConfigFilePath() string {
	const ConfigFileName = "config.json"
	_, currentAbsPath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("設定ファイル読込エラー, 現在のパスの取得に失敗しました")
	}

	currentDir := filepath.Dir(currentAbsPath)
	configFilePath := filepath.Join(currentDir, ConfigFileName)
	return configFilePath
}
