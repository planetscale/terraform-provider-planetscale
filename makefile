default: lint test

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -trimpath .

.PHONY: download-openapi
download-openapi:
	curl -sSL https://planetscale.com/docs/openapi.yaml -o schemas/openapi.yaml

.PHONY: generate
generate:
	speakeasy run --skip-versioning

.PHONY: update-speakeasy
update-speakeasy:
	speakeasy update
	@VERSION=$$(speakeasy --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1) && \
		sed -i '' "s/^speakeasyVersion:.*/speakeasyVersion: $$VERSION/" .speakeasy/workflow.yaml
	$(MAKE) generate

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

.PHONY: sweep
sweep:
	bash ./script/sweep

## misc helpers
.PHONY: run-renovate-dry
run-renovate-dry:
	docker run --rm -e LOG_LEVEL=debug -e "GITHUB_COM_TOKEN=$$(gh auth token)" -v $$(pwd -P):/usr/src/app --pull always renovate/renovate --platform=local --include-paths '$(INCLUDE_PATHS)'