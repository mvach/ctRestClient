.PHONY: test coverage coverage-html coverage-func generate lint build clean help

test:
	go test -coverprofile=coverage.out $$(go list ./... | grep -v fakes)

coverage-report: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "HTML coverage report generated: coverage.html"

coverage-func: test
	go tool cover -func=coverage.out

# Generate a clean coverage summary
coverage-summary:
	./scripts/coverage_summary

lint:
	golangci-lint run

build:
	./scripts/build_binaries

generate:
	./scripts/generate_fakes

clean:
	rm -f coverage.out coverage.html
	rm -rf dist/
