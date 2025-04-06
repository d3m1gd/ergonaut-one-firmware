all:
	@go generate .
	@go test ./...
	@go run ./util/gen/reference/... reference.go
	@go run .

