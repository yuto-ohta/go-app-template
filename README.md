# go-app-template

```
├── docker
│   ├── app
│   └── db
└── src
    ├── api
    ├── apperror
    ├── apputil
    ├── config
    ├── domain
    ├── infrastructure
    └── usecase
```
## ■What is

* golangで基本的なAPIをつくる
* 機能
    * ユーザーテーブルのCRUD
    * Rの全件取得はソート、ページングもする
* その他
    * アーキテクチャはDDD(風)
    * 基本的なエラー機構, テストコード
* WIP
    * ログイン・ログアウト
    * TODO
        * ログアウト機能
        * ユーザー更新、ユーザー削除を未ログイン時にUnauthorizedにする

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
* POST `http://localhost:1323/login`
    * ログイン

```
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "id": 1,
  "password": "Test1111"
}' \
 'http://localhost:1323/login'
```

* POST `http://localhost:1323/users/new`
    * ユーザー登録

```
curl -i -X POST \
   -H "Content-Type:application/json" \
   -d \
'{
  "name": "ハルキゲニア",
  "password": "Test1111"
}' \
 'http://localhost:1323/users/new'
```

* GET `http://localhost:1323/users/:id`
    * ユーザー取得
    
* GET `http://localhost:1323/users?orderBy=name&order=ASC&limit=1&offset=1`
    * ユーザー全取得
    * オプション: orderBy, order, limit, offset

* PUT `http://localhost:1323/users/:id/update`
    * ユーザー更新
    * パターン: ユーザー名のみ, パスワードのみ, 両方
    
```
curl -i -X PUT \
   -H "Content-Type:application/json" \
   -d \
'{
  "name": "ハルキゲニア",
  "password": "Test1111"
}' \
 'http://localhost:1323/users/1/update'
```

* DELETE `http://localhost:1323/users/:id`
    * ユーザー削除
