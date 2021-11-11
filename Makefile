NAME=main
COVERAGE_MIN=50.0

## coverage: Show coverage report in browser
coverage: test
	go tool cover -html=cp.out

## test: Run test and enforce go coverage
test:
	go test ./... -coverprofile cp.out

	$(eval COVERAGE_CURRENT = $(shell go tool cover -func=cp.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}' ))
	$(eval COVERAGE_PASSED = $(shell echo "$(COVERAGE_CURRENT) >= $(COVERAGE_MIN)" | bc -l ))

	@if [ $(COVERAGE_PASSED) == 0 ] ; then \
		echo "coverage is $(COVERAGE_CURRENT) below required threshold $(COVERAGE_MIN)"; \
		exit 2; \
    fi

## bench: Run benchmark test
bench:
	go test -bench=.

## watch: development with air
watch:
	air -c .air.default.toml

## build: Build binary applications
build:
	@go generate ./...
	@echo building binary to ./dist/${NAME}
	@go build -o ./dist/${NAME} .

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run with parameter options: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
