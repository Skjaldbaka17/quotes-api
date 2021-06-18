TEST_FUNCTION?=
server:
	go run server.go
check-swagger:
	which swagger || (go get -u github.com/go-swagger/go-swagger/cmd/swagger)
docs: check-swagger
	swagger generate spec --input ./swagger/tags.yaml -o ./swagger/swagger.yaml --scan-models
serve-docs:check-swagger
	swagger serve -F=swagger ./swagger/swagger.yaml
test:
	go test ./...
test-verbose:
	go test ./... -v
test-specific:
	go test ./... -run $(TEST_FUNCTION) -v -failfast