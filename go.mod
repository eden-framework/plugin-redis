module github.com/eden-framework/plugin-redis

go 1.14

replace k8s.io/client-go => k8s.io/client-go v0.18.8

require (
	github.com/eden-framework/eden-framework v1.0.13
	github.com/go-redis/redis/v8 v8.2.1
	github.com/profzone/envconfig v1.4.6
)
