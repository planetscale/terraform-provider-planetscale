default: lint test

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -trimpath .

.PHONY: lint
lint:
	golangci-lint run -v ./...

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests only ..."
	go test -v -cover ./...

# Run acceptance tests. These create real resources.
.PHONY: testacc
testacc:
	TF_ACC=1 go test -parallel=2 ./... -v $(TESTARGS) -timeout 120m

.PHONY: generate-docs
generate-docs:
	go tool github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: generate
generate:
	bash ./script/update_openapi_spec
	bash ./script/generate
	go generate ./...

.PHONY: sweep
sweep:
	bash ./script/sweep