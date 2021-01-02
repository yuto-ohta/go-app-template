# go-app-template

```
├── docker
│   ├── app
│   └── db
└── src
    ├── api
    ├── config
    ├── domain
    ├── errors
    ├── infrastructure
    ├── usecase
    └── utils
```
## ■What is

* golangでapiをつくりたい
    * 将来的には、このリポジトリをテンプレートにして、いろいろなアプリをつくりたい
* WIPです

## ■Setup

### 1. build & run
`docker-compose build` (first time only)

`docker-compose up`

* ※たまにdbの起動が間に合わずに失敗するので、そのときは何度か施行してくださいmm

### 2. create data

`cd ${リポジトリのルート}`

`go run ./src/config/db/initialize/initialize_local_data.go`

## ■Routing

* [See Me](https://github.com/yuto-ohta/go-app-template/blob/master/src/config/routes/router.go)
* `http://localhost:1323/user/:id`
    * ユーザー取得
* `http://localhost:1323/user/new?name=:userName`
    * ユーザー登録
