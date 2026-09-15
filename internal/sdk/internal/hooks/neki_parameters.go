package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const terraformManagedParametersQuery = "parameters"

var nekiParametersPathPattern = regexp.MustCompile(
	`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/(?:admin|configuration-profiles/[^/]+|sidecars/[^/]+)/parameters$`,
)

// NekiParametersHook reconciles parameter-detail responses with the keys
// already managed in Terraform state. The generated SDK carries prior state in
// a client-only query parameter, which this hook removes before calling the
// API.
type NekiParametersHook struct{}

var _ sdkInitHook = (*NekiParametersHook)(nil)

func NewNekiParametersHook() *NekiParametersHook {
	return &NekiParametersHook{}
}

func (h *NekiParametersHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &nekiParametersClient{client: client}
}

type nekiParametersClient struct {
	client HTTPClient
}

func (c *nekiParametersClient) Do(req *http.Request) (*http.Response, error) {
	if !isNekiParametersRequest(req) {
		return c.client.Do(req)
	}

	managed, err := takeTerraformManagedParameters(req)
	if err != nil {
		return nil, err
	}
	extensions, err := takeTerraformManagedExtensions(req)
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
		return nil, fmt.Errorf(
			"neki parameters endpoint returned HTTP %d; refusing to treat its parent resource as missing",
			res.StatusCode,
		)
	}
	if res.StatusCode != http.StatusOK {
		return res, nil
	}

	return reconcileNekiParametersResponse(res, managed, extensions != nil)
}

func isNekiParametersRequest(req *http.Request) bool {
	return req != nil &&
		req.URL != nil &&
		req.Method == http.MethodGet &&
		nekiParametersPathPattern.MatchString(req.URL.Path)
}

func takeTerraformManagedParameters(req *http.Request) (map[string]map[string]string, error) {
	query := req.URL.Query()
	encoded := query.Get(terraformManagedParametersQuery)
	query.Del(terraformManagedParametersQuery)
	req.URL.RawQuery = query.Encode()

	managed := map[string]map[string]string{}
	if encoded == "" {
		return managed, nil
	}

	if err := json.Unmarshal([]byte(encoded), &managed); err != nil {
		return nil, fmt.Errorf("decode prior Terraform Neki parameters: %w", err)
	}

	return managed, nil
}

type nekiParameterDetail struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	DefaultValue string `json:"default_value"`
}

func reconcileNekiParameters(
	details []nekiParameterDetail,
	managed map[string]map[string]string,
) map[string]map[string]string {
	parameters := map[string]map[string]string{}

	for _, detail := range details {
		_, wasManaged := managed[detail.Namespace][detail.Name]
		if detail.Value == detail.DefaultValue && !wasManaged {
			continue
		}

		if parameters[detail.Namespace] == nil {
			parameters[detail.Namespace] = map[string]string{}
		}
		parameters[detail.Namespace][detail.Name] = detail.Value
	}

	return parameters
}

func reconcileNekiParametersResponse(
	res *http.Response,
	managed map[string]map[string]string,
	managesExtensions bool,
) (*http.Response, error) {
	defer func() {
		_ = res.Body.Close()
	}()

	var details []nekiParameterDetail
	if err := json.NewDecoder(res.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("decode Neki configuration profile parameters response: %w", err)
	}

	parameters := reconcileNekiParameters(details, managed)
	if managesExtensions {
		delete(parameters["pgconf"], "shared_preload_libraries")
		delete(parameters["pgconf"], "session_preload_libraries")
		if len(parameters["pgconf"]) == 0 {
			delete(parameters, "pgconf")
		}
	}
	payload, err := json.Marshal(map[string]any{"parameters": parameters})
	if err != nil {
		return nil, fmt.Errorf("encode reconciled Neki configuration profile parameters response: %w", err)
	}

	res.Body = io.NopCloser(bytes.NewReader(payload))
	res.ContentLength = int64(len(payload))
	if res.Header == nil {
		res.Header = make(http.Header)
	}
	res.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	res.Header.Set("Content-Type", "application/json")

	return res, nil
}
