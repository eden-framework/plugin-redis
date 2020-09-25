generator: generator.go
	go build -buildmode=plugin -o=./generator.so ./generator.go
