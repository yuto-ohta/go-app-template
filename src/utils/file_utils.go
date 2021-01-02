package utils

import (
	"log"
	"path/filepath"
	"runtime"
)

/*
	指定したファイルの絶対パスを取得する
	ディレクトリのパスは呼び出し元のものとなる

	ex) "/projectRoot/hoge/huga/main.go"にて、fileName "piyopiyo.txt"で実行したとき
		→  return "/projectRoot/hoge/huga/piyopiyo.txt"
*/
func GetFilePathWithCurrentDir(fileName string) string {
	_, callerAbsPath, _, ok := runtime.Caller(1)
	if !ok {
		log.Fatal("ファイル読込エラー, 呼び出し元のパスの取得に失敗しました")
	}

	callerDir := filepath.Dir(callerAbsPath)
	resFilePath := filepath.Join(callerDir, fileName)
	return resFilePath
}
