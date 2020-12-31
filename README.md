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
    └── usecase
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

* 初期データ作成用のSQLは ["/docker/db/manual/initialize_local_data.sql"](https://github.com/yuto-ohta/go-app-template/blob/master/docker/db/manual/initialize_local_data.sql) にあります。
* MySQLクライアントなどを用いて、手動で適用してください

## ■Routing

* [See Me](https://github.com/yuto-ohta/go-app-template/blob/051b1270883b7ee1b472902812d149bba9180387/src/config/routes/router.go#L24)
* `http://localhost:1323/user/:id`
    * ユーザー取得
