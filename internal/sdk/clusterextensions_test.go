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

func TestPostgresManagedExtensionsPreserveOmittedSelection(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/v1/organizations/org/databases/db/branches/main/extensions", req.URL.Path)
		require.Empty(t, req.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"name":"hll","enabled":true,"can_enable":true}]`))
	}))
	defer server.Close()

	client := New(
		WithServerURL(server.URL+"/v1"),
		WithClient(server.Client()),
		WithSecurity(shared.Security{ServiceToken: "token", ServiceTokenID: "id"}),
	)
	request := operations.GetPostgresBranchManagedExtensionsRequest{
		Organization: "org",
		Database:     "db",
		Branch:       "main",
	}

	response, err := client.ClusterExtensions.GetPostgresBranchManagedExtensions(context.Background(), request)
	require.NoError(t, err)
	require.Nil(t, response.Object.Extensions)

	request.Extensions = []string{}
	response, err = client.ClusterExtensions.GetPostgresBranchManagedExtensions(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, []string{"hll"}, response.Object.Extensions)
}
