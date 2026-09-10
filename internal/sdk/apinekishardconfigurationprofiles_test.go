package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

func TestUpdateNekiConfigurationProfileUsesTerraformProviderUserAgent(t *testing.T) {
	t.Parallel()

	type seenRequest struct {
		userAgent string
		body      map[string]json.RawMessage
	}
	requests := make(chan seenRequest, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPatch, req.Method)
		require.Equal(t, "/organizations/org/databases/db/branches/main/configuration-profiles/default", req.URL.EscapedPath())
		var body map[string]json.RawMessage
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		requests <- seenRequest{
			userAgent: req.UserAgent(),
			body:      body,
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := New(
		WithServerURL(server.URL),
		WithClient(server.Client()),
		WithSecurity(shared.Security{
			ServiceToken:   "token",
			ServiceTokenID: "token-id",
		}),
	)
	_, err := client.APINekiShardConfigurationProfiles.UpdateNekiConfigurationProfile(
		context.Background(),
		operations.UpdateNekiConfigurationProfileRequest{
			Organization:         "org",
			Database:             "db",
			Branch:               "main",
			ConfigurationProfile: "default",
			Body: &operations.UpdateNekiConfigurationProfileRequestBody{
				Replicas: Int64(2),
			},
		},
	)
	require.NoError(t, err)
	got := <-requests
	require.Contains(t, got.userAgent, "terraform-provider-planetscale")
	require.NotContains(t, got.body, "parameters")
}

func TestUpdateNekiConfigurationProfilePreservesExplicitEmptyParameters(t *testing.T) {
	t.Parallel()

	requestBodies := make(chan map[string]json.RawMessage, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var body map[string]json.RawMessage
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		requestBodies <- body

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := New(
		WithServerURL(server.URL),
		WithClient(server.Client()),
		WithSecurity(shared.Security{
			ServiceToken:   "token",
			ServiceTokenID: "token-id",
		}),
	)
	_, err := client.APINekiShardConfigurationProfiles.UpdateNekiConfigurationProfile(
		context.Background(),
		operations.UpdateNekiConfigurationProfileRequest{
			Organization:         "org",
			Database:             "db",
			Branch:               "main",
			ConfigurationProfile: "default",
			Body: &operations.UpdateNekiConfigurationProfileRequestBody{
				Parameters: map[string]map[string]string{},
			},
		},
	)
	require.NoError(t, err)

	body := <-requestBodies
	require.Contains(t, body, "parameters")
	require.JSONEq(t, `{}`, string(body["parameters"]))
}
