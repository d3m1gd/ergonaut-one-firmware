all:
	@go generate .
	@go test ./...
	@go run .

deps:
	@go install golang.org/x/tools/cmd/stringer@latest
