package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
)

const skippedKeyspaceResizeRequestID = "__speakeasy_skipped_keyspace_resize_request__"

var keyspaceResizePathPattern = regexp.MustCompile(`^(/v1/organizations/[^/]+/databases/[^/]+/branches/[^/]+/keyspaces/[^/]+)/resizes$`)

// VitessKeyspaceResizeSkipHook skips the keyspace resize that every
// planetscale_vitess_keyspace update sends when the requested size and
// replicas already match the keyspace, such as when only parameters change.
// The API rejects a resize that changes nothing.
type VitessKeyspaceResizeSkipHook struct{}

var _ sdkInitHook = (*VitessKeyspaceResizeSkipHook)(nil)

func NewVitessKeyspaceResizeSkipHook() *VitessKeyspaceResizeSkipHook {
	return &VitessKeyspaceResizeSkipHook{}
}

func (h *VitessKeyspaceResizeSkipHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	if client == nil {
		return baseURL, client
	}

	return baseURL, &vitessKeyspaceResizeSkipClient{client: client}
}

type vitessKeyspaceResizeSkipClient struct {
	client HTTPClient
}

func (c *vitessKeyspaceResizeSkipClient) Do(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil || req.Method != http.MethodPut || req.Body == nil {
		return c.client.Do(req)
	}
	matches := keyspaceResizePathPattern.FindStringSubmatch(req.URL.Path)
	if matches == nil {
		return c.client.Do(req)
	}

	payload, err := io.ReadAll(req.Body)
	_ = req.Body.Close()
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(payload))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(payload)), nil
	}

	if c.resizeChangesNothing(req, matches[1], payload) {
		return syntheticJSONResponse(req, http.StatusOK, map[string]string{
			"id":    skippedKeyspaceResizeRequestID,
			"state": "completed",
		})
	}
	return c.client.Do(req)
}

// resizeChangesNothing reports whether the keyspace already has the requested
// size and replicas. Any doubt, such as a failed lookup, sends the resize.
func (c *vitessKeyspaceResizeSkipClient) resizeChangesNothing(req *http.Request, keyspacePath string, payload []byte) bool {
	var requested struct {
		ClusterSize   *string `json:"cluster_size"`
		ExtraReplicas *int64  `json:"extra_replicas"`
	}
	if err := json.Unmarshal(payload, &requested); err != nil {
		return false
	}

	res, err := sendDerivedRequest(c.client, req, http.MethodGet, keyspacePath, nil)
	if err != nil {
		return false
	}
	if res.StatusCode != http.StatusOK {
		drainAndClose(res.Body)
		return false
	}

	var current struct {
		ClusterName   *string `json:"cluster_name"`
		ExtraReplicas *int64  `json:"extra_replicas"`
	}
	if err := decodeAndClose(res.Body, &current); err != nil {
		return false
	}

	sizeMatches := requested.ClusterSize == nil ||
		(current.ClusterName != nil && *current.ClusterName == *requested.ClusterSize)
	replicasMatch := requested.ExtraReplicas == nil ||
		(current.ExtraReplicas != nil && *current.ExtraReplicas == *requested.ExtraReplicas)
	return sizeMatches && replicasMatch
}
