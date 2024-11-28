default: run

.PHONY: run
run: # Run the app using go run
	go run cmd/dumpj/main.go

.PHONY: test
test: # Test the app using go test
	go test -v ./tests/...

.PHONY: clean
clean: ## Remove _files
	rm -rf _files/zips/*.zip
	rm -rf _files/unzipped/**
	rm _files/dumpj.db

sqlgen: ## Run sqlc
	docker run --rm -v $(shell pwd):/src -w /src sqlc/sqlc generate

migration: ## Create a new migration
	cd ./sql/migrations && goose -s create $(m) sql && cd ../..

migrate_up: ## Run migrate up
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./_files/dumpj.db goose -dir=./sql/migrations up

migrate_down: ## Run migrate up
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./_files/dumpj.db goose -dir=./sql/migrations down

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
