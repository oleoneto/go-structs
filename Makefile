clean-test:
	@go fmt ./...
	@go clean -testcache

tidy:
	@go mod tidy

test: clean-test
	go test -cover -coverprofile=coverage.out -p 1 ./... | tee test.log
	go tool cover -html=coverage.out
