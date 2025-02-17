package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/planetscale/terraform-provider-planetscale/internal/client/planetscale"
	"golang.org/x/oauth2"
)

const testAccOrg = "planetscale-terraform-testing"

var debugProvider = os.Getenv("TF_PS_PROVIDER_DEBUG") != ""

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"planetscale": providerserver.NewProtocol6WithError(New("test", debugProvider)()),
}

var testAccAPIClient *planetscale.Client

func testAccPreCheck(t *testing.T) {
	var (
		accessToken       = os.Getenv("PLANETSCALE_ACCESS_TOKEN")
		serviceTokenID    = os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID")
		serviceTokenValue = os.Getenv("PLANETSCALE_SERVICE_TOKEN")
	)
	switch {
	case accessToken != "":
	case serviceTokenID != "" && serviceTokenValue != "":
	default:
		t.Fatalf("must have either PLANETSCALE_ACCESS_TOKEN or both of (PLANETSCALE_SERVICE_TOKEN_ID, PLANETSCALE_SERVICE_TOKEN)")
	}

	// TODO: factor client creation out of the provider.go Configure() func so we can
	//       more easily re-use it here and maintain the logic around access and service-token lookups
	if testAccAPIClient == nil {
		accessToken := os.Getenv("PLANETSCALE_ACCESS_TOKEN")
		if accessToken != "" {
			tok := &oauth2.Token{AccessToken: accessToken}
			rt := &oauth2.Transport{
				Base:   http.DefaultTransport,
				Source: oauth2.StaticTokenSource(tok),
			}
			testAccAPIClient = planetscale.NewClient(&http.Client{Transport: rt}, nil)
		}

		serviceTokenID := os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID")
		serviceTokenValue := os.Getenv("PLANETSCALE_SERVICE_TOKEN")
		if serviceTokenID != "" && serviceTokenValue != "" {
			rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				r.Header.Set("Authorization", serviceTokenID+":"+serviceTokenValue)
				return http.DefaultTransport.RoundTrip(r)
			})
			testAccAPIClient = planetscale.NewClient(&http.Client{Transport: rt}, nil)
		}
	}
}

func checkIntegerMin(minimum int) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if v < minimum {
			return fmt.Errorf("value %d is less than %d", v, minimum)
		}
		return nil
	}
}

func checkOneOf(values ...string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		for _, valid := range values {
			if value == valid {
				return nil
			}
		}
		return fmt.Errorf("value %q is not one of %s", value, strings.Join(values, ", "))
	}
}

// checkExpectUpdate is a helper function for resource.TestStep/ConfigPlanChecks
// to assert that the plan should updated the resource in place.
func checkExpectUpdate(resourceName string) resource.ConfigPlanChecks { //nolint:unparam
	return resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}
}

// checkExpectRecreate is a helper function for resource.TestStep/ConfigPlanChecks
// to assert that the plan should recreate the resource.
func checkExpectRecreate(resourceName string) resource.ConfigPlanChecks { //nolint:unused
	return resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionReplace),
		},
	}
}

func TestProviderConfigure(t *testing.T) {
	cases := map[string]struct {
		config      map[string]tftypes.Value
		expectWarn  bool
		expectError bool
	}{
		"service_token_id only": {
			config: map[string]tftypes.Value{
				"service_token_name": tftypes.NewValue(tftypes.String, nil), // deprecated
				"service_token_id":   tftypes.NewValue(tftypes.String, "token123"),
				"service_token":      tftypes.NewValue(tftypes.String, "secret"),
				"endpoint":           tftypes.NewValue(tftypes.String, nil),
				"access_token":       tftypes.NewValue(tftypes.String, nil),
			},
			expectWarn:  false,
			expectError: false,
		},

		"both access_token and service_token specified": {
			config: map[string]tftypes.Value{
				"service_token_name": tftypes.NewValue(tftypes.String, nil), // deprecated
				"service_token_id":   tftypes.NewValue(tftypes.String, "token123"),
				"service_token":      tftypes.NewValue(tftypes.String, "secret"),
				"endpoint":           tftypes.NewValue(tftypes.String, nil),
				"access_token":       tftypes.NewValue(tftypes.String, "acctoken123"),
			},
			expectWarn:  false,
			expectError: true,
		},

		"deprecated service_token_name used": {
			config: map[string]tftypes.Value{
				"service_token_name": tftypes.NewValue(tftypes.String, "token123"), // deprecated
				"service_token_id":   tftypes.NewValue(tftypes.String, nil),
				"service_token":      tftypes.NewValue(tftypes.String, "secret"),
				"endpoint":           tftypes.NewValue(tftypes.String, nil),
				"access_token":       tftypes.NewValue(tftypes.String, nil),
			},
			expectWarn:  true,
			expectError: false,
		},

		"deprecated service_token_name is used with service_token_id": {
			config: map[string]tftypes.Value{
				"service_token_name": tftypes.NewValue(tftypes.String, "token123"), // deprecated
				"service_token_id":   tftypes.NewValue(tftypes.String, nil),
				"service_token":      tftypes.NewValue(tftypes.String, "secret"),
				"endpoint":           tftypes.NewValue(tftypes.String, nil),
				"access_token":       tftypes.NewValue(tftypes.String, nil),
			},
			expectWarn:  true,
			expectError: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			p := New("test", false)()

			var schemaResp provider.SchemaResponse
			p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

			objectType := tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"endpoint":           tftypes.String,
					"access_token":       tftypes.String,
					"service_token_id":   tftypes.String,
					"service_token_name": tftypes.String,
					"service_token":      tftypes.String,
				},
			}

			var req provider.ConfigureRequest
			req.Config = tfsdk.Config{
				Raw:    tftypes.NewValue(objectType, tc.config),
				Schema: schemaResp.Schema,
			}
			var resp provider.ConfigureResponse

			p.Configure(ctx, req, &resp)

			// t.Logf("Diagnostics: %v", resp.Diagnostics)

			if tc.expectWarn && resp.Diagnostics.WarningsCount() == 0 {
				t.Error("expected warning but got none")
			}
			if !tc.expectWarn && resp.Diagnostics.WarningsCount() != 0 {
				t.Errorf("unexpected warning: %v", resp.Diagnostics)
			}

			if tc.expectError && !resp.Diagnostics.HasError() {
				t.Error("expected error but got none")
			}
			if !tc.expectError && resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}
}
