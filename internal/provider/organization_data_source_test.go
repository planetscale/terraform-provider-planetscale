package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.planetscale_organization.test", "name", "planetscale-terraform-testing"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "billing_email"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "database_count"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "features.insights"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "has_past_due_invoices"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "id"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "idp_managed_roles"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "plan"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "single_tenancy"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "sleeping_database_count"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "sso"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "sso_directory"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "updated_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_organization.test", "valid_billing_info"),
				),
			},
		},
	})
}

const testAccOrganizationDataSourceConfig = `
data "planetscale_organization" "test" {
	name = "planetscale-terraform-testing"
}
`
