server:
	source .env; \
	go run server.go
test:
	go test ./...
test-verbose:
	go test ./... -v