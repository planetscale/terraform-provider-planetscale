package hooks

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"slices"
)

var nekiExtensionsPathPattern = regexp.MustCompile(
	`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/configuration-profiles/[^/]+/extensions$`,
)
var postgresExtensionsPathPattern = regexp.MustCompile(
	`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/extensions$`,
)

type NekiExtensionsHook struct{}

var _ sdkInitHook = (*NekiExtensionsHook)(nil)

func (h *NekiExtensionsHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}
	return baseURL, &nekiExtensionsClient{client: client}
}

type nekiExtensionsClient struct {
	client HTTPClient
}

func (c *nekiExtensionsClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil || req.Method != http.MethodGet ||
		(!nekiExtensionsPathPattern.MatchString(req.URL.Path) && !postgresExtensionsPathPattern.MatchString(req.URL.Path)) {
		return c.client.Do(req)
	}

	query := req.URL.Query()
	terraformManaged := query.Get("terraform_managed") == "true"
	query.Del("terraform_managed")
	req.URL.RawQuery = query.Encode()

	managed, err := takeTerraformManagedExtensions(req)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil || res == nil {
		return res, err
	}
	if res.StatusCode == http.StatusNotFound {
		if res.Body != nil {
			_ = res.Body.Close()
		}
		if nekiExtensionsPathPattern.MatchString(req.URL.Path) {
			return nil, fmt.Errorf("neki extensions endpoint returned HTTP 404; refusing to treat its configuration profile as missing")
		}
		return nil, fmt.Errorf("extensions endpoint returned HTTP 404; refusing to treat its parent resource as missing")
	}
	if res.StatusCode != http.StatusOK {
		return res, nil
	}

	var extensions []struct {
		Name      string `json:"name"`
		Enabled   bool   `json:"enabled"`
		CanEnable bool   `json:"can_enable"`
	}
	if err := decodeAndClose(res.Body, &extensions); err != nil {
		return nil, fmt.Errorf("decode extensions response: %w", err)
	}
	if terraformManaged && managed == nil {
		return res, replaceResponseBody(res, map[string]any{"extensions": nil})
	}

	enabled := []string{}
	remaining := make(map[string]bool)
	for _, extension := range extensions {
		if extension.Enabled && extension.CanEnable {
			remaining[extension.Name] = true
		}
	}
	for _, name := range managed {
		if remaining[name] {
			enabled = append(enabled, name)
			delete(remaining, name)
		}
	}
	enabled = append(enabled, slices.Sorted(maps.Keys(remaining))...)
	return res, replaceResponseBody(res, map[string]any{"extensions": enabled})
}

func takeTerraformManagedExtensions(req *http.Request) ([]string, error) {
	query := req.URL.Query()
	encoded := query.Get("extensions")
	query.Del("extensions")
	req.URL.RawQuery = query.Encode()

	var managed []string
	if encoded != "" {
		if err := json.Unmarshal([]byte(encoded), &managed); err != nil {
			return nil, fmt.Errorf("decode prior Terraform extensions: %w", err)
		}
	}
	return managed, nil
}
