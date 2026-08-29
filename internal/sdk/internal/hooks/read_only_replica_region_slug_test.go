package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const readOnlyReplicaGetURL = "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/read-only-replicas/my-replica"

func replicaGetClient(status int, body string) HTTPClient {
	hook := NewReadOnlyReplicaRegionSlugHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body:    io.NopCloser(strings.NewReader(body)),
			Request: req,
		}, nil
	}))

	return wrappedClient
}

func TestReadOnlyReplicaRegionSlugHookHoistsSlug(t *testing.T) {
	client := replicaGetClient(http.StatusOK, `{"name":"my-replica","region":{"slug":"us-east","display_name":"US East"}}`)

	req, err := http.NewRequest(http.MethodGet, readOnlyReplicaGetURL, nil)
	require.NoError(t, err, "failed to build get request")

	res, err := client.Do(req)
	require.NoError(t, err, "get request failed")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err, "failed to read rewritten body")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(bodyBytes, &payload), "failed to decode rewritten payload")
	require.Equal(t, "us-east", payload["region_slug"], "expected region slug hoisted to region_slug")
	require.Equal(t, "my-replica", payload["name"], "expected original fields preserved")
	require.Equal(t, int64(len(bodyBytes)), res.ContentLength, "expected content length to match rewritten body")
}

func TestReadOnlyReplicaRegionSlugHookPassesThroughUnexpectedBody(t *testing.T) {
	responseBody := `{"name":"my-replica"}`
	client := replicaGetClient(http.StatusOK, responseBody)

	req, err := http.NewRequest(http.MethodGet, readOnlyReplicaGetURL, nil)
	require.NoError(t, err, "failed to build get request")

	res, err := client.Do(req)
	require.NoError(t, err, "get request failed")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err, "failed to read passthrough body")
	require.Equal(t, responseBody, string(bodyBytes), "expected body without region object to pass through unchanged")
}

func TestReadOnlyReplicaRegionSlugHookIgnoresOtherPaths(t *testing.T) {
	responseBody := `[{"name":"my-replica","region":{"slug":"us-east"}}]`
	client := replicaGetClient(http.StatusOK, responseBody)

	req, err := http.NewRequest(http.MethodGet, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/read-only-replicas", nil)
	require.NoError(t, err, "failed to build get request")

	res, err := client.Do(req)
	require.NoError(t, err, "get request failed")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err, "failed to read passthrough body")
	require.Equal(t, responseBody, string(bodyBytes), "expected list response to pass through unchanged")
}

func TestReadOnlyReplicaRegionSlugHookIgnoresNon200(t *testing.T) {
	responseBody := `{"code":"not_found"}`
	client := replicaGetClient(http.StatusNotFound, responseBody)

	req, err := http.NewRequest(http.MethodGet, readOnlyReplicaGetURL, nil)
	require.NoError(t, err, "failed to build get request")

	res, err := client.Do(req)
	require.NoError(t, err, "get request failed")

	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err, "failed to read passthrough body")
	require.Equal(t, responseBody, string(bodyBytes), "expected non-200 response to pass through unchanged")
}
