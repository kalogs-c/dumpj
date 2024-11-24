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

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
