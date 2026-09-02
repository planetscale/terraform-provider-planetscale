package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const postgresBranchGetURL = "https://api.planetscale.com/v1/organizations/org/databases/db/branches/br"

func postgresBranchGetClient(status int, body string) HTTPClient {
	hook := NewPostgresBranchRegionSlugHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	}))
	return wrappedClient
}

func TestPostgresBranchRegionSlugHookHoistsSlug(t *testing.T) {
	client := postgresBranchGetClient(http.StatusOK, `{"id":"br","region":{"slug":"us-east","display_name":"US East"}}`)

	req, err := http.NewRequest(http.MethodGet, postgresBranchGetURL, nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.JSONEq(t, `{
		"id": "br",
		"region": {"slug": "us-east", "display_name": "US East"},
		"region_slug": "us-east"
	}`, string(bodyBytes))
	require.Equal(t, int64(len(bodyBytes)), res.ContentLength)
}

func TestPostgresBranchRegionSlugHookPassesThroughUnexpectedBody(t *testing.T) {
	responseBody := `{"id":"br"}`
	client := postgresBranchGetClient(http.StatusOK, responseBody)

	req, err := http.NewRequest(http.MethodGet, postgresBranchGetURL, nil)
	require.NoError(t, err)
	res, err := client.Do(req)
	require.NoError(t, err)
	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, responseBody, string(bodyBytes))
}

func TestPostgresBranchRegionSlugHookIgnoresOtherPaths(t *testing.T) {
	responseBody := `{"region":{"slug":"us-east"}}`
	client := postgresBranchGetClient(http.StatusOK, responseBody)

	req, err := http.NewRequest(http.MethodGet, postgresBranchGetURL+"/backups", nil)
	require.NoError(t, err)
	res, err := client.Do(req)
	require.NoError(t, err)
	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, responseBody, string(bodyBytes))
}

func TestPostgresBranchRegionSlugHookIgnoresNon200(t *testing.T) {
	responseBody := `{"code":"not_found"}`
	client := postgresBranchGetClient(http.StatusNotFound, responseBody)

	req, err := http.NewRequest(http.MethodGet, postgresBranchGetURL, nil)
	require.NoError(t, err)
	res, err := client.Do(req)
	require.NoError(t, err)
	bodyBytes, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Equal(t, responseBody, string(bodyBytes))
}
