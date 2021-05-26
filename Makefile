TEST_FUNCTION?=
server:
	source .env; \
	go run server.go
check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)
docs: check-swagger
	swagger generate spec -o ./swagger/swagger.yaml --scan-models
serve-docs:check-swagger
	swagger serve -F=swagger ./swagger/swagger.yaml
test:
	go test ./...
test-verbose:
	go test ./... -v
test-specific:
	go test ./... -v -run $(TEST_FUNCTION)