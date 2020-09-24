module github.com/eden-framework/redis

go 1.14

replace (
	github.com/eden-framework/eden-framework v1.0.12 => /Users/liyiwen/Documents/golang/src/github.com/eden-framework/eden-framework
	k8s.io/client-go => k8s.io/client-go v0.18.8
)

require (
	github.com/eden-framework/eden-framework v1.0.12
	github.com/go-redis/redis/v8 v8.2.1
	github.com/profzone/envconfig v1.4.5
)
