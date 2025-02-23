# Run all tests
test:
	@go test ./...

# Run parser tests specifically
test-parser:
	@go test ./pkg/... -run TestParse

# Run tests with coverage
coverage:
	@go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# Clean up generated files
clean:
	@rm -f coverage.out

.PHONY: test test-parser coverage clean
