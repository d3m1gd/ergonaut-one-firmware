all:
	@go generate .
	@go run ./util/gen/reference/... reference.go
	@go test ./...
	@go run .

