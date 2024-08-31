package provider

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		serviceTokenName  = os.Getenv("PLANETSCALE_SERVICE_TOKEN_NAME")
		serviceTokenValue = os.Getenv("PLANETSCALE_SERVICE_TOKEN")
	)
	switch {
	case accessToken != "":
	case serviceTokenName != "" && serviceTokenValue != "":
	default:
		t.Fatalf("must have either PLANETSCALE_ACCESS_TOKEN or both of (PLANETSCALE_SERVICE_TOKEN_NAME, PLANETSCALE_SERVICE_TOKEN)")
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

		serviceTokenName := os.Getenv("PLANETSCALE_SERVICE_TOKEN_NAME")
		serviceTokenValue := os.Getenv("PLANETSCALE_SERVICE_TOKEN")
		if serviceTokenName != "" && serviceTokenValue != "" {
			rt := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				r.Header.Set("Authorization", serviceTokenName+":"+serviceTokenValue)
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
