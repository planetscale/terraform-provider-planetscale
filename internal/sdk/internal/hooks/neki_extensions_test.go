package hooks

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNekiExtensionsReadsEnabledCustomerExtensions(t *testing.T) {
	t.Parallel()

	hook := &NekiExtensionsHook{}
	_, client := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		require.Empty(t, req.URL.RawQuery)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`[
				{"name":"hll","enabled":true,"can_enable":true},
				{"name":"pg_cron","enabled":false,"can_enable":true},
				{"name":"pgextwlist","enabled":true,"can_enable":false},
				{"name":"bloom","can_enable":false}
			]`)),
		}, nil
	}))
	req, err := http.NewRequest(http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/extensions?extensions=%5B%22pg_cron%22%5D", nil)
	require.NoError(t, err)

	res, err := client.Do(req)

	require.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"extensions":["hll"]}`, string(body))
}

func TestNekiExtensionsPreservesEmptySelection(t *testing.T) {
	t.Parallel()

	hook := &NekiExtensionsHook{}
	_, client := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"name":"pgextwlist","enabled":true,"can_enable":false}]`)),
		}, nil
	}))
	req, err := http.NewRequest(http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/extensions?extensions=%5B%5D", nil)
	require.NoError(t, err)

	res, err := client.Do(req)

	require.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"extensions":[]}`, string(body))
}

func TestNekiExtensionsLeavesOmittedSelectionUnmanaged(t *testing.T) {
	t.Parallel()

	hook := &NekiExtensionsHook{}
	_, client := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"name":"hll","enabled":true,"can_enable":true}]`)),
		}, nil
	}))
	req, err := http.NewRequest(http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/extensions", nil)
	require.NoError(t, err)

	res, err := client.Do(req)

	require.NoError(t, err)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.JSONEq(t, `{"extensions":null}`, string(body))
}
