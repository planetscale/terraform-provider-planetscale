package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

func TestNekiConfigurationProfileDataSourceReadsEnabledExtensions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		w.Header().Set("Content-Type", "application/json")
		switch req.URL.Path {
		case "/v1/organizations/org/databases/db/branches/main/configuration-profiles/default":
			_, _ = w.Write([]byte(`{"name":"default","state":"ready","storage":{}}`))
		case "/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/parameters":
			_, _ = w.Write([]byte(`[]`))
		case "/v1/organizations/org/databases/db/branches/main/configuration-profiles/default/extensions":
			require.Empty(t, req.URL.RawQuery)
			_, _ = w.Write([]byte(`[
				{"name":"hll","enabled":true,"can_enable":true},
				{"name":"auto_explain","enabled":false,"can_enable":true},
				{"name":"pgextwlist","enabled":true,"can_enable":false}
			]`))
		default:
			t.Errorf("unexpected request: %s", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	source := &NekiConfigurationProfileDataSource{client: sdk.New(
		sdk.WithServerURL(server.URL+"/v1"),
		sdk.WithClient(server.Client()),
		sdk.WithSecurity(shared.Security{ServiceToken: "token", ServiceTokenID: "id"}),
	)}
	var schemaResponse datasource.SchemaResponse
	source.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)
	require.True(t, schemaResponse.Schema.Attributes["extensions"].IsComputed())

	config, diags := types.ObjectValueFrom(ctx, schemaResponse.Schema.Type().(types.ObjectType).AttrTypes, &NekiConfigurationProfileDataSourceModel{
		Organization: types.StringValue("org"),
		Database:     types.StringValue("db"),
		Branch:       types.StringValue("main"),
		Name:         types.StringValue("default"),
	})
	require.False(t, diags.HasError(), diags.Errors())
	raw, err := config.ToTerraformValue(ctx)
	require.NoError(t, err)
	response := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResponse.Schema}}

	source.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Raw: raw, Schema: schemaResponse.Schema}}, &response)

	require.False(t, response.Diagnostics.HasError(), response.Diagnostics.Errors())
	var extensions []string
	diags = response.State.GetAttribute(ctx, path.Root("extensions"), &extensions)
	require.False(t, diags.HasError(), diags.Errors())
	require.Equal(t, []string{"hll"}, extensions)
}
