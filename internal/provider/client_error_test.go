package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
	"github.com/stretchr/testify/require"
)

// The API rejects invalid input with a 422 whose body is {"code": ..., "message": ...}.
// The SDK must surface that message as the error Terraform reports.
func TestSDKReportsValidationErrorMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"code":"unprocessable","message":"Minimum storage bytes must be greater than or equal to 10 GB"}`))
	}))
	defer srv.Close()

	client := sdk.New(
		sdk.WithServerURL(srv.URL+"/v1"),
		sdk.WithSecurity(shared.Security{ServiceToken: "token", ServiceTokenID: "id"}),
	)

	res, err := client.APINekiShardConfigurationProfiles.CreateNekiConfigurationProfile(context.Background(), operations.CreateNekiConfigurationProfileRequest{
		Organization: "big-bang",
		Database:     "neki",
		Branch:       "main",
	})
	require.Nil(t, res)
	require.EqualError(t, err, "PlanetScale API returned HTTP 422 (unprocessable): Minimum storage bytes must be greater than or equal to 10 GB")
}

// 404 must keep flowing through as a response so Read can drop deleted resources from state.
func TestSDKStillReturnsNotFoundAsResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"not_found","message":"Not Found"}`))
	}))
	defer srv.Close()

	client := sdk.New(
		sdk.WithServerURL(srv.URL+"/v1"),
		sdk.WithSecurity(shared.Security{ServiceToken: "token", ServiceTokenID: "id"}),
	)

	res, err := client.APINekiShardConfigurationProfiles.GetNekiConfigurationProfile(context.Background(), operations.GetNekiConfigurationProfileRequest{
		Organization:         "big-bang",
		Database:             "neki",
		Branch:               "main",
		ConfigurationProfile: "myprof",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
