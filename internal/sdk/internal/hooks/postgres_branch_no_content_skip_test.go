package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testHTTPClient func(*http.Request) (*http.Response, error)

func (t testHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return t(req)
}

func TestPostgresBranchNoContentSkipHookSkipsNextGetAfterPatch204(t *testing.T) {
	var getCalls int

	hook := NewPostgresBranchNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		switch req.Method {
		case http.MethodPatch:
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		case http.MethodGet:
			getCalls++
			require.FailNow(t, "wrapped client should skip outbound GET after PATCH 204")
		}
		return nil, nil
	}))

	patchReq, err := http.NewRequest(http.MethodPatch, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/changes", strings.NewReader(`{"cluster_size":"PS_20_AWS_X86"}`))
	require.NoError(t, err, "failed to build patch request")

	patchRes, err := wrappedClient.Do(patchReq)
	require.NoError(t, err, "patch request failed")
	require.Equal(t, http.StatusOK, patchRes.StatusCode, "expected transformed patch response status 200")

	getReq, err := http.NewRequest(http.MethodGet, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/changes/oq4hzhavm3um", nil)
	require.NoError(t, err, "failed to build get request")

	getRes, err := wrappedClient.Do(getReq)
	require.NoError(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getRes.StatusCode, "expected synthetic get response status 200")
	require.Equal(t, 0, getCalls, "expected wrapped client to short-circuit get")

	bodyBytes, err := io.ReadAll(getRes.Body)
	require.NoError(t, err, "failed to read synthetic body")

	var payload map[string]any
	require.NoError(t, json.Unmarshal(bodyBytes, &payload), "failed to decode synthetic payload")

	require.Equal(t, "completed", payload["state"], "expected synthetic state completed")
}

func TestPostgresBranchNoContentSkipHookPassesThroughGetWithoutMarker(t *testing.T) {
	var getCalls int

	hook := NewPostgresBranchNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet {
			getCalls++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
				Body:    io.NopCloser(strings.NewReader(`{"id":"oq4hzhavm3um","state":"completed"}`)),
				Request: req,
			}, nil
		}

		require.FailNowf(t, "unexpected method in test client", "method=%s", req.Method)
		return nil, nil
	}))

	getReq, err := http.NewRequest(http.MethodGet, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/changes/oq4hzhavm3um", nil)
	require.NoError(t, err, "failed to build get request")

	getRes, err := wrappedClient.Do(getReq)
	require.NoError(t, err, "get request failed")
	require.Equal(t, http.StatusOK, getRes.StatusCode, "expected passthrough status 200")
	require.Equal(t, 1, getCalls, "expected base client to be called once")
}
