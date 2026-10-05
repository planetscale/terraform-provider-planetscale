package hooks

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PostgresBranchExtensionsParametersHook struct{}

var _ sdkInitHook = (*PostgresBranchExtensionsParametersHook)(nil)

func (h *PostgresBranchExtensionsParametersHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}
	return baseURL, &postgresBranchExtensionsParametersClient{client: client}
}

type postgresBranchExtensionsParametersClient struct {
	client HTTPClient
}

func (c *postgresBranchExtensionsParametersClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil || req.Method != http.MethodGet || !postgresBranchGetPathPattern.MatchString(req.URL.Path) {
		return c.client.Do(req)
	}

	managed, err := takeTerraformManagedExtensions(req)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil || res == nil || managed == nil || res.StatusCode != http.StatusOK {
		return res, err
	}

	var branch map[string]json.RawMessage
	if err := decodeAndClose(res.Body, &branch); err != nil {
		return nil, fmt.Errorf("decode Postgres branch response: %w", err)
	}
	var parameters map[string]map[string]json.RawMessage
	if err := json.Unmarshal(branch["parameters"], &parameters); err == nil {
		delete(parameters["pgconf"], "shared_preload_libraries")
		delete(parameters["pgconf"], "session_preload_libraries")
		if len(parameters["pgconf"]) == 0 {
			delete(parameters, "pgconf")
		}
		encoded, err := json.Marshal(parameters)
		if err != nil {
			return nil, fmt.Errorf("encode Postgres branch parameters: %w", err)
		}
		branch["parameters"] = encoded
	}
	return res, replaceResponseBody(res, branch)
}
