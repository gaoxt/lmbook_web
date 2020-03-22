module app

go 1.14

require (
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/gorilla/mux v1.7.4
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.3 // indirect
)

replace github.com/app => ../app
