package main

import (
	"fmt"
	"go-app-template/src/config/db"
	appUtils "go-app-template/src/utils"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

const (
	ddlFileName = "ddl.sql"
	dmlFileName = "dml.sql"
	sqlSubStr   = ";"
)

func main() {
	fmt.Println("Initialize Local Data Start!!-----------------------------------")

	// SQLファイルのパスを取得する
	ddlFilePath := appUtils.GetFilePathWithCurrentDir(ddlFileName)
	dmlFileName := appUtils.GetFilePathWithCurrentDir(dmlFileName)

	// SQLファイルを読み込む
	ddlByteSlice, err := ioutil.ReadFile(filepath.Clean(ddlFilePath))
	if err != nil {
		log.Fatalf("初期データファイル読込エラー, Error: %v", err.Error())
	}
	dmlByteSlice, err := ioutil.ReadFile(filepath.Clean(dmlFileName))
	if err != nil {
		log.Fatalf("初期データファイル読込エラー, Error: %v", err.Error())
	}

	// ddl実行
	ddlQueries := strings.Split(string(ddlByteSlice), sqlSubStr)
	for _, q := range ddlQueries {
		q = strings.TrimSpace(q)
		// SQLファイルの末尾の空行がqueryに入ることがあるため,そのときはスキップする
		if q == "" {
			continue
		}

		fmt.Println(q)
		result := db.Conn.Exec(q)

		if result.Error != nil {
			fmt.Println("Initialize Local Data Failed!!-----------------------------------")
			log.Fatalf("DDLの実行中にエラーが発生しました, Error: %v", result.Error)
		}
	}

	// dml実行
	dmlQueries := strings.Split(string(dmlByteSlice), sqlSubStr)
	err = db.Conn.Transaction(func(tx *gorm.DB) error {
		for _, q := range dmlQueries {
			q = strings.TrimSpace(q)
			// SQLファイルの末尾の空行がqueryに入ることがあるため,そのときはスキップする
			if q == "" {
				continue
			}

			fmt.Println(q)
			result := tx.Exec(q)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Initialize Local Data Failed!!-----------------------------------")
		log.Fatalf("DML実行中にエラーが発生したため、ロールバックしました, Error: %v", err)
	} else {
		fmt.Println("Initialize Local Data All Finished!!-----------------------------------")
	}
}
