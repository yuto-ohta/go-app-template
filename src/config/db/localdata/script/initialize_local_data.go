package main

import "go-app-template/src/config/db/localdata"

// 初期データ投入
func main() {
	localdata.InitializeLocalData()
}
