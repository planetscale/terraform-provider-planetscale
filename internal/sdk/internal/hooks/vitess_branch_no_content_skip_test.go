package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVitessBranchNoContentSkipHookCompletesNoOpResizeLocally(t *testing.T) {
	var getCalls int

	hook := NewVitessBranchNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		switch req.Method {
		case http.MethodPut:
			return &http.Response{
				StatusCode: http.StatusNoContent,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			}, nil
		case http.MethodGet:
			getCalls++
			require.FailNow(t, "wrapped client should complete the sentinel resize locally")
		}
		return nil, nil
	}))

	putReq, err := http.NewRequest(http.MethodPut, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/resizes", strings.NewReader(`{}`))
	require.NoError(t, err)

	putRes, err := wrappedClient.Do(putReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, putRes.StatusCode)

	var putPayload map[string]any
	require.NoError(t, json.NewDecoder(putRes.Body).Decode(&putPayload))
	require.Equal(t, skippedResizeRequestID, putPayload["id"])

	getReq, err := http.NewRequest(http.MethodGet, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/resizes/"+skippedResizeRequestID, nil)
	require.NoError(t, err)

	getRes, err := wrappedClient.Do(getReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	require.Equal(t, 0, getCalls)

	var getPayload map[string]any
	require.NoError(t, json.NewDecoder(getRes.Body).Decode(&getPayload))
	require.Equal(t, "completed", getPayload["state"])
}

func TestVitessBranchNoContentSkipHookPassesThroughRealResizeRequests(t *testing.T) {
	var calls int

	hook := NewVitessBranchNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"id":"oq4hzhavm3um","state":"resizing"}`)),
			Request:    req,
		}, nil
	}))

	getReq, err := http.NewRequest(http.MethodGet, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/resizes/oq4hzhavm3um", nil)
	require.NoError(t, err)

	getRes, err := wrappedClient.Do(getReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getRes.StatusCode)
	require.Equal(t, 1, calls)
}

func TestVitessBranchNoContentSkipHookPassesThrough204OnOtherPaths(t *testing.T) {
	hook := NewVitessBranchNoContentSkipHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
			Request:    req,
		}, nil
	}))

	putReq, err := http.NewRequest(http.MethodPut, "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br/safe-migrations", strings.NewReader(`{"safe_migrations":true}`))
	require.NoError(t, err)

	putRes, err := wrappedClient.Do(putReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, putRes.StatusCode)
}
