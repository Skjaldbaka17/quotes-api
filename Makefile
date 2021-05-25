TEST_TO_RUN?=
server:
	source .env; \
	go run server.go
test:
	go test ./...
test-verbose:
	go test ./... -v
test-specific:
	go test ./... -v -run $(TEST_TO_RUN)