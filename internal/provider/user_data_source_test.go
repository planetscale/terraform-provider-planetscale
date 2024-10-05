package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	if os.Getenv("PLANETSCALE_ACCESS_TOKEN") == "" {
		t.Skip("PLANETSCALE_ACCESS_TOKEN is required for planetscale_user data source tests.")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccUserDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "id"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "name"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "avatar_url"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "created_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "default_organization.%"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "directory_managed"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "display_name"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "email"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "email_verified"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "managed"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "sso"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "two_factor_auth_configured"),
					resource.TestCheckResourceAttrSet("data.planetscale_user.current", "updated_at"),
				),
			},
		},
	})
}

const testAccUserDataSourceConfig = `
data "planetscale_user" "current" {}
`
