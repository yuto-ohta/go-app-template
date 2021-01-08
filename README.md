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

※たまにdbの起動が間に合わずに失敗します

### 2. create data

`cd ${リポジトリのルート}`

`docker-compose exec app go run ./config/db/localdata/script/initialize_local_data.go`

## ■Routing

* [See Me](https://github.com/yuto-ohta/go-app-template/blob/master/src/config/route/router.go)
* POST `http://localhost:1323/user/new`
    * ユーザー登録

```
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{"name":"ハルキゲニア"}' \
 'http://localhost:1323/user/new'
```

* GET `http://localhost:1323/user/:id`
    * ユーザー取得
    
* GET `http://localhost:1323/users`
    * ユーザー全取得

* PUT `http://localhost:1323/user/:id/update`
    * ユーザー更新
    
```
curl -i -X PUT \
   -H "Content-Type:application/json" \
   -d \
'{"name":"オパビニア"}' \
 'http://localhost:1323/user/1/update'
```

* DELETE `http://localhost:1323/user/:id`
    * ユーザー削除