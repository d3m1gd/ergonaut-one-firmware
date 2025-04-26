all:
	@go install golang.org/x/tools/cmd/stringer@latest
	@go generate .
	@go run ./util/gen/reference/... reference.go
	@go test ./...
	@go run .

