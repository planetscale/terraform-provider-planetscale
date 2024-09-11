package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationsDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.planetscale_organizations.test", "organizations.#", "1"),
					resource.TestCheckResourceAttr("data.planetscale_organizations.test", "organizations.0.name", testAccOrg),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.admin_only_production_access"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.billing_email"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.can_create_databases"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.created_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.database_count"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.features.insights"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.free_databases_remaining"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.has_past_due_invoices"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.id"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.idp_managed_roles"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.plan"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.single_tenancy"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.sleeping_database_count"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.sso"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.sso_directory"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.planetscale_organizations.test", "organizations.0.valid_billing_info"),
				),
			},
		},
	})
}

const testAccOrganizationsDataSourceConfig = `
data "planetscale_organizations" "test" {

}
`
