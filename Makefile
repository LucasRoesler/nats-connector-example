
.PROJECT=nats-connector-example
.REPO=theaxer
.TAG=latest

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the job image
	docker build --tag ${.REPO}/${.PROJECT}:${.TAG} .
