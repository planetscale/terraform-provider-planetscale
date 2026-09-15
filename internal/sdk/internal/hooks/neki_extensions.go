package hooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
)

var nekiExtensionsPathPattern = regexp.MustCompile(
	`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/configuration-profiles/[^/]+/extensions$`,
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
	if req == nil || req.URL == nil || req.Method != http.MethodGet || !nekiExtensionsPathPattern.MatchString(req.URL.Path) {
		return c.client.Do(req)
	}

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
		return nil, fmt.Errorf("neki extensions endpoint returned HTTP 404; refusing to treat its configuration profile as missing")
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
		return nil, fmt.Errorf("decode Neki extensions response: %w", err)
	}

	var enabled []string
	if managed != nil {
		enabled = []string{}
		for _, extension := range extensions {
			if extension.Enabled && extension.CanEnable {
				enabled = append(enabled, extension.Name)
			}
		}
		sort.Strings(enabled)
	}
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
			return nil, fmt.Errorf("decode prior Terraform Neki extensions: %w", err)
		}
	}
	return managed, nil
}
