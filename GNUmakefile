default: testacc

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -trimpath .

.PHONY: lint
lint:
	golangci-lint run -v ./...

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test -parallel=2 ./... -v $(TESTARGS) -timeout 120m

.PHONY: generate-docs
generate-docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

.PHONY: generate
generate:
	bash ./script/update_openapi_spec
	bash ./script/generate
	go generate ./...

.PHONY: sweep
sweep:
	bash ./script/sweep