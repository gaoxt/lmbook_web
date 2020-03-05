module app

go 1.14

require (
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/gorilla/mux v1.7.4
	github.com/app v0.0.0-00010101000000-000000000000
)

replace github.com/app => ../app
