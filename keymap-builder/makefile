all:
	@go generate .
	@go run ./util/gen/reference/... reference.go
	@go test ./...
	@go run .

deps:
	@go install golang.org/x/tools/cmd/stringer@latest
