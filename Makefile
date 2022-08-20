run:first second proxy

first:
	go run cmd/node/main.go
second:
	go run cmd/node/main.go ":8081"
proxy:
	go run cmd/proxy/main.go