package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

var readOnlyReplicaGetPathPattern = regexp.MustCompile(`^/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/read-only-replicas/[^/]+$`)

// ReadOnlyReplicaRegionSlugHook wraps the SDK HTTP client and copies the
// nested region slug to a top-level region_slug field on read-only replica
// GET responses. The Terraform schema maps region_slug back to the region
// attribute the create request configures, so reads and imports refresh it.
type ReadOnlyReplicaRegionSlugHook struct{}

var _ sdkInitHook = (*ReadOnlyReplicaRegionSlugHook)(nil)

func NewReadOnlyReplicaRegionSlugHook() *ReadOnlyReplicaRegionSlugHook {
	return &ReadOnlyReplicaRegionSlugHook{}
}

func (h *ReadOnlyReplicaRegionSlugHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &readOnlyReplicaRegionSlugClient{client: client}
}

type readOnlyReplicaRegionSlugClient struct {
	client HTTPClient
}

func (c *readOnlyReplicaRegionSlugClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil || res == nil || req == nil || req.URL == nil {
		return res, err
	}

	if req.Method != http.MethodGet || res.StatusCode != http.StatusOK || !readOnlyReplicaGetPathPattern.MatchString(req.URL.Path) {
		return res, nil
	}

	body, readErr := io.ReadAll(res.Body)
	drainAndClose(res.Body)
	if readErr != nil {
		return res, readErr
	}

	if rewritten, ok := hoistRegionSlug(body); ok {
		body = rewritten
	}

	res.Body = io.NopCloser(bytes.NewReader(body))
	res.ContentLength = int64(len(body))
	if res.Header == nil {
		res.Header = make(http.Header)
	}
	res.Header.Set("Content-Length", strconv.Itoa(len(body)))

	return res, nil
}

func hoistRegionSlug(body []byte) ([]byte, bool) {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, false
	}

	region, ok := payload["region"].(map[string]any)
	if !ok {
		return nil, false
	}

	slug, ok := region["slug"].(string)
	if !ok {
		return nil, false
	}

	payload["region_slug"] = slug

	rewritten, err := json.Marshal(payload)
	if err != nil {
		return nil, false
	}

	return rewritten, true
}
