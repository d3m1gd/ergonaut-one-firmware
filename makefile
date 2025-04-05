all:
	@go generate .
	@go run ./util/gen/behavior/... behavior.go
	@go run .

