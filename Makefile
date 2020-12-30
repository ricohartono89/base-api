PROJECT_NAME := "base-api"
PKG := "./"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all dep build clean test coverage coverhtml lint api

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -race . ./routes ./env -v -coverprofile=coverage.out
	@go tool cover -func=coverage.out

testlocal: ## Run unittests
	@go test -failfast . ./routes ./env

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	sh ./tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

dep: ## Get the dependencies
	@go get -v -d ./...

build: api ## Build the binary file

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

create-migration: ## Create new migration file. It takes parameter `file` as filename. Usage: `make create-migration file=add_column_time`
	ls -t db/migrations/*.sql | head -1 | awk -F"migrations/" '{print $$2}' | awk -F"_" '{print $$1}' | { read cur_v; expr $$cur_v + 1; } | { read new_v; printf "%06d" $$new_v; } | { read v; touch db/migrations/$$v"_$(file)".up.sql; touch db/migrations/$$v"_$(file)".down.sql; }
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

api: dep ## Build the binary file
	@go build -i -v $(PKG)


