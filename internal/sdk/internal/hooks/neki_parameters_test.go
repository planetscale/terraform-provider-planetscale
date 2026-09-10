package hooks

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReconcileNekiParameters(t *testing.T) {
	t.Parallel()

	details := []nekiParameterDetail{
		{
			Namespace:    "pgconf",
			Name:         "managed_default",
			Value:        "25",
			DefaultValue: "25",
		},
		{
			Namespace:    "pgconf",
			Name:         "unmanaged_default",
			Value:        "5min",
			DefaultValue: "5min",
		},
		{
			Namespace:    "pgconf",
			Name:         "managed_changed",
			Value:        "30",
			DefaultValue: "25",
		},
		{
			Namespace:    "router",
			Name:         "remote_non_default",
			Value:        "on",
			DefaultValue: "off",
		},
	}
	managed := map[string]map[string]string{
		"pgconf": {
			"managed_default": "025",
			"managed_changed": "25",
		},
	}

	require.Equal(t, map[string]map[string]string{
		"pgconf": {
			"managed_default": "25",
			"managed_changed": "30",
		},
		"router": {
			"remote_non_default": "on",
		},
	}, reconcileNekiParameters(details, managed))
}

func TestNekiConfigurationProfileParametersHookReconcilesAndStripsClientState(t *testing.T) {
	t.Parallel()

	var calls int
	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, "false", req.URL.Query().Get("internal"))
		require.False(t, req.URL.Query().Has(terraformManagedParametersQuery))

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`[
				{"namespace":"pgconf","name":"max_connections","value":"25","default_value":"25"},
				{"namespace":"pgconf","name":"archive_timeout","value":"1min","default_value":"5min"},
				{"namespace":"router","name":"pool_size","value":"10","default_value":"10"}
			]`)),
			Request: req,
		}, nil
	}))

	managedJSON, err := json.Marshal(map[string]map[string]string{
		"pgconf": {"max_connections": "025"},
	})
	require.NoError(t, err)
	query := url.Values{
		"internal":                      {"false"},
		terraformManagedParametersQuery: {string(managedJSON)},
	}
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/parameters?"+query.Encode(),
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, calls)

	var payload struct {
		Parameters map[string]map[string]string `json:"parameters"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	require.Equal(t, map[string]map[string]string{
		"pgconf": {
			"archive_timeout": "1min",
			"max_connections": "25",
		},
	}, payload.Parameters)
}

func TestNekiAdminParametersHookReconcilesAndStripsClientState(t *testing.T) {
	t.Parallel()

	var calls int
	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, "false", req.URL.Query().Get("internal"))
		require.False(t, req.URL.Query().Has(terraformManagedParametersQuery))

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`[
				{"namespace":"admin","name":"recovery-poll-interval","value":"10s","default_value":"10s"},
				{"namespace":"admin","name":"acceptable-replication-lag","value":"20s","default_value":"10s"}
			]`)),
			Request: req,
		}, nil
	}))

	managedJSON, err := json.Marshal(map[string]map[string]string{
		"admin": {"recovery-poll-interval": "10s"},
	})
	require.NoError(t, err)
	query := url.Values{
		"internal":                      {"false"},
		terraformManagedParametersQuery: {string(managedJSON)},
	}
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/admin/parameters?"+query.Encode(),
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, calls)

	var payload struct {
		Parameters map[string]map[string]string `json:"parameters"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	require.Equal(t, map[string]map[string]string{
		"admin": {
			"acceptable-replication-lag": "20s",
			"recovery-poll-interval":     "10s",
		},
	}, payload.Parameters)
}

func TestNekiSidecarParametersHookReconcilesAndStripsClientState(t *testing.T) {
	t.Parallel()

	var calls int
	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		require.Equal(t, "false", req.URL.Query().Get("internal"))
		require.False(t, req.URL.Query().Has(terraformManagedParametersQuery))

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`[
				{"namespace":"sidecar","name":"pool-capacity","value":"100","default_value":"100"},
				{"namespace":"sidecar","name":"pool-idle-timeout","value":"5m","default_value":"10m"}
			]`)),
			Request: req,
		}, nil
	}))

	managedJSON, err := json.Marshal(map[string]map[string]string{
		"sidecar": {"pool-capacity": "100"},
	})
	require.NoError(t, err)
	query := url.Values{
		"internal":                      {"false"},
		terraformManagedParametersQuery: {string(managedJSON)},
	}
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/sidecars/default/parameters?"+query.Encode(),
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, calls)

	var payload struct {
		Parameters map[string]map[string]string `json:"parameters"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	require.Equal(t, map[string]map[string]string{
		"sidecar": {
			"pool-capacity":     "100",
			"pool-idle-timeout": "5m",
		},
	}, payload.Parameters)
}

func TestNekiConfigurationProfileParametersHookReturnsEmptyMapWithoutRelevantParameters(t *testing.T) {
	t.Parallel()

	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`[
				{"namespace":"pgconf","name":"max_connections","value":"25","default_value":"25"}
			]`)),
			Request: req,
		}, nil
	}))

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/parameters?internal=false",
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.NoError(t, err)
	require.JSONEq(t, `{"parameters":{}}`, readResponseBody(t, res))
}

func TestNekiConfigurationProfileParametersHookRejectsNotFound(t *testing.T) {
	t.Parallel()

	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"error":"not found"}`)),
			Request:    req,
		}, nil
	}))

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/parameters?internal=false",
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.Nil(t, res)
	require.ErrorContains(t, err, "refusing to treat its parent resource as missing")
}

func TestNekiConfigurationProfileParametersHookRejectsInvalidManagedState(t *testing.T) {
	t.Parallel()

	var calls int
	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		return nil, nil
	}))

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/parameters?parameters=not-json",
		nil,
	)
	require.NoError(t, err)

	_, err = wrappedClient.Do(req)
	require.ErrorContains(t, err, "decode prior Terraform Neki parameters")
	require.Zero(t, calls)
}

func TestNekiConfigurationProfileParametersHookPassesThroughOtherRequests(t *testing.T) {
	t.Parallel()

	var calls int
	hook := NewNekiParametersHook()
	_, wrappedClient := hook.SDKInit("https://api.planetscale.com", testHTTPClient(func(req *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"state":"ready"}`)),
			Request:    req,
		}, nil
	}))

	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.planetscale.com/v1/organizations/org/databases/db/branches/main/configuration-profiles/default",
		nil,
	)
	require.NoError(t, err)

	res, err := wrappedClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, 1, calls)
	require.JSONEq(t, `{"state":"ready"}`, readResponseBody(t, res))
}

func readResponseBody(t *testing.T, res *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(body)
}
