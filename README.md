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
    * ログイン・ログアウト
    * ユーザーテーブルのCRUD
        * 全件取得→ ソート, ページング
* その他
    * アーキテクチャはDDD(風)
    * 基本的なエラー機構, テストコード

## ■Setup

### 1. build & run
`docker-compose build` (first time only)

`docker-compose up`

※たまにdbの起動が間に合わずに失敗します

### 2. create data

`cd ${リポジトリのルート}`

`docker-compose exec app go run ./config/db/localdata/script/initialize_local_data.go`

## ■Routing

[See Me](https://github.com/yuto-ohta/go-app-template/blob/master/src/config/route/router.go)

### ▼ログイン　POST `http://localhost:1323/login`


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

### ▼ログアウト　GET `http://localhost:1323/logout`

### ▼ユーザー登録　POST `http://localhost:1323/users/new`

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

### ▼ユーザー取得　GET `http://localhost:1323/users/:id`
    
### ▼ユーザー全取得　GET `http://localhost:1323/users?orderBy=name&order=ASC&limit=1&offset=1`

オプション: orderBy, order, limit, offset

### ▼ユーザー更新　PUT `http://localhost:1323/users/:id/update`

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

* ログイン必須
* パターン: ユーザー名のみ, パスワードのみ, 両方
    

### ▼ユーザー削除　DELETE `http://localhost:1323/users/:id`

* ログイン必須
