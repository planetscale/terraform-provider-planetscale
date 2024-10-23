default: testacc

.PHONY: lint
lint:
	golangci-lint run -v ./...

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: generate
generate:
	bash ./script/update_openapi_spec
	bash ./script/generate
	go generate ./...