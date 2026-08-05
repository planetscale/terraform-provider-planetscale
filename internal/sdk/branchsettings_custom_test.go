package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

func TestDatabaseBranchesSafeMigrationsSettings(t *testing.T) {
	t.Parallel()

	var safeMigrations bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/organizations/org/databases/db/branches/br", r.URL.Path[:len("/organizations/org/databases/db/branches/br")])
		require.NotEmpty(t, r.Header.Get("Authorization"))

		switch {
		case r.URL.Path == "/organizations/org/databases/db/branches/br" && r.Method == http.MethodGet:
		case r.URL.Path == "/organizations/org/databases/db/branches/br/safe-migrations" && r.Method == http.MethodPost:
			safeMigrations = true
		case r.URL.Path == "/organizations/org/databases/db/branches/br/safe-migrations" && r.Method == http.MethodDelete:
			safeMigrations = false
		default:
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
			"id":                            "branch-id",
			"safe_migrations":               safeMigrations,
			"vtgate_size":                   "vg.c1.large",
			"vtgate_name":                   "VTG_320",
			"vtgate_count":                  1,
			"vtgate_max_count":              2,
			"vtgate_autoscaling":            false,
			"vtgate_target_cpu_utilization": 50,
		}))
	}))
	t.Cleanup(server.Close)

	client := New(
		WithServerURL(server.URL),
		WithClient(server.Client()),
		WithSecurity(shared.Security{ServiceTokenID: "token-id", ServiceToken: "token"}),
	)

	settings, err := client.DatabaseBranches.GetVitessBranchSettings(context.Background(), "org", "db", "br")
	require.NoError(t, err)
	require.False(t, settings.SafeMigrations)
	require.Equal(t, "VTG_320", settings.VTGateName)
	require.Equal(t, int64(2), *settings.VTGateMaxCount)
	require.Equal(t, int64(50), *settings.VTGateTargetCPUUtilization)

	settings, err = client.DatabaseBranches.SetVitessBranchSafeMigrations(context.Background(), "org", "db", "br", true)
	require.NoError(t, err)
	require.True(t, settings.SafeMigrations)

	settings, err = client.DatabaseBranches.SetVitessBranchSafeMigrations(context.Background(), "org", "db", "br", false)
	require.NoError(t, err)
	require.False(t, settings.SafeMigrations)
}

func TestDatabaseBranchesVTGateConfiguration(t *testing.T) {
	t.Parallel()

	resize := VitessBranchResizeRequest{
		ID:                         "resize-id",
		State:                      "completed",
		VTGateSize:                 "vg.c1.large",
		VTGateName:                 "VTG_320",
		VTGateCount:                1,
		VTGateMaxCount:             Int64(2),
		VTGateAutoscaling:          true,
		VTGateTargetCPUUtilization: Int64(50),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/organizations/org/databases/db/branches/br/resizes", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPut:
			var request UpdateVitessBranchVTGateConfigurationRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&request))
			require.NotNil(t, request.VTGateAutoscaling)
			require.True(t, *request.VTGateAutoscaling)
			require.Equal(t, int64(2), *request.VTGateMaxCount)
			require.NoError(t, json.NewEncoder(w).Encode(resize))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := New(
		WithServerURL(server.URL),
		WithClient(server.Client()),
		WithSecurity(shared.Security{ServiceTokenID: "token-id", ServiceToken: "token"}),
	)
	enabled := true
	updated, err := client.DatabaseBranches.UpdateVitessBranchVTGateConfiguration(
		context.Background(),
		"org",
		"db",
		"br",
		UpdateVitessBranchVTGateConfigurationRequest{
			VTGateAutoscaling: &enabled,
			VTGateMaxCount:    Int64(2),
		},
	)
	require.NoError(t, err)
	require.True(t, updated.VTGateAutoscaling)
}
