package sdk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

func TestGetPostgresBouncerTerraformState(t *testing.T) {
	t.Parallel()

	type seenRequest struct {
		method string
		path   string
	}
	seenRequests := make(chan seenRequest, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenRequests <- seenRequest{
			method: r.Method,
			path:   r.URL.EscapedPath(),
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"type": "PostgresBouncer",
			"id": "bouncer-id",
			"name": "my-bouncer",
			"target": "primary",
			"bouncer_size": "PGB_5",
			"replicas_per_cell": 1,
			"actor": {
				"id": "actor-id",
				"display_name": "Actor Name",
				"avatar_url": "https://example.com/avatar.png"
			},
			"parameters": {
				"pgbouncer": {
					"default_pool_size": "100"
				}
			}
		}`))
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
	res, err := client.Bouncers.GetPostgresBouncerTerraformState(
		context.Background(),
		operations.GetPostgresBouncerTerraformStateRequest{
			Organization: "org",
			Database:     "db",
			Branch:       "br",
			Bouncer:      "my-bouncer",
		},
	)

	require.NoError(t, err)
	got := <-seenRequests
	require.Equal(t, http.MethodGet, got.method)
	require.Equal(t, "/organizations/org/databases/db/branches/br/bouncers/my-bouncer/terraform-state", got.path)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotNil(t, res.Object)
	require.Equal(t, "100", res.Object.Parameters["pgbouncer"]["default_pool_size"])
	require.Equal(t, "actor-id", res.Object.Actor.ID)
}
