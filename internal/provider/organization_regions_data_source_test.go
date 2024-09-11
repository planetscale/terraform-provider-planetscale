package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var regions = []string{
	"aws-us-east-2",
	"us-east",
	"us-west",
	"eu-west",
	"ap-south",
	"ap-southeast",
	"ap-northeast",
	"eu-central",
	"aws-ap-southeast-2",
	"aws-sa-east-1",
	"gcp-us-central1",
	"aws-eu-west-2",
	"gcp-us-east4",
	"gcp-northamerica-northeast1",
	"gcp-asia-northeast3",
}

func TestAccOrganizationRegionsDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccOrganizationRegionsDataSourceConfig(testAccOrg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith("data.planetscale_organization_regions.test", "regions.#", checkIntegerMin(1)),
					resource.TestCheckResourceAttrWith("data.planetscale_organization_regions.test", "regions.0.slug", checkOneOf(regions...)),
				),
			},
		},
	})
}

func testAccOrganizationRegionsDataSourceConfig(org string) string {
	return fmt.Sprintf(`
	data "planetscale_organization_regions" "test" {
		organization = %[1]q
	}`,
		org,
	)
}
