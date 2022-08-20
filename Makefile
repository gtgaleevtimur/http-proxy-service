run:
	go run cmd/node/main.go
	go run cmd/node/main.go ':8081'
	go run proxy/main.go