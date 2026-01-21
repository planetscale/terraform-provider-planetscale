package provider

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const testAccOrg = "planetscale-terraform-testing"

// Returns a mapping of provider type names to provider server implementations,
// suitable for acceptance testing via the
// resource.TestCase.ProtoV6ProtocolFactories field.
func testAccProviders() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"planetscale": providerserver.NewProtocol6WithError(New("test")()),
	}
}

// Immediately fails testing if the PLANETSCALE_SERVICE_TOKEN and
// PLANETSCALE_SERVICE_TOKEN_ID environment variables are not set.
func testAccPreCheck(t *testing.T) {
	if os.Getenv("PLANETSCALE_SERVICE_TOKEN") != "" && os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID") != "" {
		return
	}

	t.Fatal("Both PLANETSCALE_SERVICE_TOKEN and PLANETSCALE_SERVICE_TOKEN_ID must be set for acceptance tests")
}

// randomWithPrefix generates a random string with the given prefix.
func randomWithPrefix(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, rand.Intn(1000000))
}
