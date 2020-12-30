# go-app-template

```
├── docker
│   ├── app
│   └── db
│       ├── data
│       │   ├── development
│       │   ├── mysql
│       │   ├── performance_schema
│       │   └── sys
│       └── init
└── src
    ├── api
    │   └── controller
    │       └── test
    ├── config
    │   ├── db
    │   └── routes
    ├── domain
    │   └── repository
    ├── errors
    │   ├── messages
    │   └── test
    │       └── mock
    ├── infrastructure
    │   └── model
    ├── tmp
    └── usecase
        └── impl
            └── test
```
## ■What is

* golangでapiをつくりたい
    * 将来的には、このリポジトリをテンプレートにして、いろいろなアプリをつくりたい
* WIPです

## ■Setup

`docker-compose build` (first time only)

`docker-compose up`

## ■Routing

* [See Me](https://github.com/yuto-ohta/go-app-template/blob/051b1270883b7ee1b472902812d149bba9180387/src/config/routes/router.go#L24)
* `http://localhost:1323/user/:id`
    * ユーザー取得
