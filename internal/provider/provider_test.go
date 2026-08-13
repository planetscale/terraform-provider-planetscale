package provider

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/operations"
	"github.com/planetscale/terraform-provider-planetscale/internal/sdk/models/shared"
)

// Defaults to planetscale-terraform-testing; override with PLANETSCALE_TEST_ORG.
var testAccOrg = func() string {
	if org := os.Getenv("PLANETSCALE_TEST_ORG"); org != "" {
		return org
	}
	return "planetscale-terraform-testing"
}()

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

func testAccBackupBranch() string {
	if branch := os.Getenv("PLANETSCALE_TEST_BRANCH"); branch != "" {
		return branch
	}
	return "main"
}

const (
	backupWaitTimeout  = 5 * time.Minute
	backupWaitInterval = 10 * time.Second
)

func testAccSDKClient() *sdk.PlanetScale {
	serverURL := os.Getenv("PLANETSCALE_SERVER_URL")
	if serverURL == "" {
		serverURL = "https://api.planetscale.com/v1"
	}
	return sdk.New(
		sdk.WithServerURL(serverURL),
		sdk.WithSecurity(shared.Security{
			ServiceToken:   os.Getenv("PLANETSCALE_SERVICE_TOKEN"),
			ServiceTokenID: os.Getenv("PLANETSCALE_SERVICE_TOKEN_ID"),
		}),
	)
}

func waitForNoInFlightPostgresBackups(t *testing.T, database, branch string) {
	t.Helper()

	client := testAccSDKClient()
	waitForNoInFlightBackups(t, branch, func(ctx context.Context) (int, error) {
		inFlight := 0
		for _, state := range []operations.ListPostgresBranchBackupsQueryParamState{
			operations.ListPostgresBranchBackupsQueryParamStatePending,
			operations.ListPostgresBranchBackupsQueryParamStateRunning,
		} {
			resp, err := client.Backups.ListPostgresBranchBackups(ctx, operations.ListPostgresBranchBackupsRequest{
				Organization: testAccOrg,
				Database:     database,
				Branch:       branch,
				State:        state.ToPointer(),
			})
			if err != nil {
				return 0, err
			}
			if resp.Object != nil {
				inFlight += len(resp.Object.Data)
			}
		}
		return inFlight, nil
	})
}

func waitForNoInFlightVitessBackups(t *testing.T, database, branch string) {
	t.Helper()

	client := testAccSDKClient()
	waitForNoInFlightBackups(t, branch, func(ctx context.Context) (int, error) {
		inFlight := 0
		for _, state := range []operations.ListVitessBranchBackupsQueryParamState{
			operations.ListVitessBranchBackupsQueryParamStatePending,
			operations.ListVitessBranchBackupsQueryParamStateRunning,
		} {
			resp, err := client.Backups.ListVitessBranchBackups(ctx, operations.ListVitessBranchBackupsRequest{
				Organization: testAccOrg,
				Database:     database,
				Branch:       branch,
				State:        state.ToPointer(),
			})
			if err != nil {
				return 0, err
			}
			if resp.Object != nil {
				inFlight += len(resp.Object.Data)
			}
		}
		return inFlight, nil
	})
}

func waitForNoInFlightBackups(t *testing.T, branch string, countInFlight func(context.Context) (int, error)) {
	t.Helper()

	ctx := context.Background()
	deadline := time.Now().Add(backupWaitTimeout)
	for {
		inFlight, err := countInFlight(ctx)
		if err != nil {
			t.Fatalf("listing in-flight backups on branch %q: %s", branch, err)
		}
		if inFlight == 0 {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("%d backup(s) still in flight on branch %q after waiting %s", inFlight, branch, backupWaitTimeout)
		}
		t.Logf("waiting for %d in-flight backup(s) on branch %q to complete", inFlight, branch)
		time.Sleep(backupWaitInterval)
	}
}
