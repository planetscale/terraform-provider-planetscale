package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const bouncerTerraformChangesURL = "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/bouncers/my-bouncer/terraform-changes"

func TestPostgresBouncerNoContentSkipHookTransformsPut204(t *testing.T) {
	hook := NewPostgresBouncerNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	}))

	putReq, err := http.NewRequest(http.MethodPut, bouncerTerraformChangesURL, strings.NewReader(`{"bouncer_size":"PGB_5"}`))
	require.NoError(t, err, "failed to build put request")

	putRes, err := wrappedClient.Do(putReq)
	require.NoError(t, err, "put request failed")
	require.Equal(t, http.StatusOK, putRes.StatusCode, "expected 204 to be transformed into a synthetic 200")

	bodyBytes, err := io.ReadAll(putRes.Body)
	require.NoError(t, err, "failed to read synthetic body")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(bodyBytes, &payload), "failed to decode synthetic payload")
	require.Equal(t, "completed", payload["state"], "expected synthetic state completed")
}

func TestPostgresBouncerNoContentSkipHookPassesThrough200(t *testing.T) {
	responseBody := `{"id":"dh4i3vg2mv3r","state":"resizing"}`

	hook := NewPostgresBouncerNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header: map[string][]string{
				"Content-Type": {"application/json"},
			},
			Body:    io.NopCloser(strings.NewReader(responseBody)),
			Request: req,
		}, nil
	}))

	putReq, err := http.NewRequest(http.MethodPut, bouncerTerraformChangesURL, strings.NewReader(`{"bouncer_size":"PGB_10"}`))
	require.NoError(t, err, "failed to build put request")

	putRes, err := wrappedClient.Do(putReq)
	require.NoError(t, err, "put request failed")
	require.Equal(t, http.StatusOK, putRes.StatusCode, "expected 200 to pass through")

	bodyBytes, err := io.ReadAll(putRes.Body)
	require.NoError(t, err, "failed to read passthrough body")
	require.Equal(t, responseBody, string(bodyBytes), "expected passthrough body to be preserved")
}

func TestPostgresBouncerNoContentSkipHookPassesThrough204OnOtherPaths(t *testing.T) {
	hook := NewPostgresBouncerNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	}))

	putReq, err := http.NewRequest(http.MethodPut, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/bouncers/my-bouncer/resizes", strings.NewReader(`{}`))
	require.NoError(t, err, "failed to build put request")

	putRes, err := wrappedClient.Do(putReq)
	require.NoError(t, err, "put request failed")
	require.Equal(t, http.StatusNoContent, putRes.StatusCode, "expected 204 on a non-terraform-changes path to pass through unmodified")
}
