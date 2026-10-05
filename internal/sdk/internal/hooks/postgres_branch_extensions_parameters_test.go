package hooks

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgresBranchExtensionsParametersRemovesManagedLoaders(t *testing.T) {
	t.Parallel()
	hook := &PostgresBranchExtensionsParametersHook{}
	_, client := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		require.Empty(t, req.URL.RawQuery)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"id":"branch-id",
				"parameters":{"pgconf":{"shared_preload_libraries":"hll","session_preload_libraries":"","max_connections":"50"}}
			}`)),
		}, nil
	}))
	req, err := http.NewRequest(http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main?extensions="+url.QueryEscape(`["hll"]`), nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"id":"branch-id","parameters":{"pgconf":{"max_connections":"50"}}}`, string(body))
}
