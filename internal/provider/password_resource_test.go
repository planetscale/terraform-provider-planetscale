package provider

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccPasswordResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	branchName := acctest.RandomWithPrefix("branch")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")

	ipv4Cidr := []string{"10.0.0.0/24"}
	// ipv6Cidr := []string{"2001:db8::/64"}
	ipv4Addr := []string{"10.0.0.1"}
	// ipv6Addr := []string{"2001:db8::1"}
	multiIPv4Cidrs := []string{"10.0.0.0/8", "10.1.0.0/8"}
	// multiIPv6Cidrs := []string{"2001:db8::/64", "2001:db9::/64"}
	// mixedCidrs := []string{"10.0.0.0/8", "2001:db8::/64"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName, ipv4Cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "name", passwdName),
					resource.TestCheckResourceAttr("planetscale_password.test", "role", "admin"),
					resource.TestCheckResourceAttr("planetscale_password.test", "branch", branchName),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", ipv4Cidr[0]),
				),
			},

			// ImportState testing
			{
				ResourceName:      "planetscale_password.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Import requires: 'organization,database,branch,id' but 'id' of the password
				// is only known after creation. Use a func to retrieve the ID from the state:
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					id, err := getPasswordIDFromState(s, "planetscale_password.test")
					if err != nil {
						return "", err
					}
					// Import requires: 'organization,database,branch,id'
					return fmt.Sprintf("%s,%s,%s,%s", testAccOrg, dbName, branchName, id), nil
				},
				// The actual password is not returned by the API, so we can't verify it here:
				ImportStateVerifyIgnore: []string{"plaintext"},
			},

			// Update tests:

			// Update 'name' attribute:
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", ipv4Cidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "name", passwdName+"-new"),
				),
			},
			// Update 'cidrs' attribute with single ipv6 range:
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", ipv6Cidr),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", ipv6Cidr[0]),
			// 	),
			// },
			// Update `cidrs` with multiple ipv4 ranges
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", multiIPv4Cidrs),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", multiIPv4Cidrs[0]),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.1", multiIPv4Cidrs[1]),
				),
			},
			// Update `cidrs` with multiple ipv6 ranges
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", multiIPv6Cidrs),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", multiIPv6Cidrs[0]),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.1", multiIPv6Cidrs[1]),
			// 	),
			// },
			// Update `cidrs` with ipv4 + ipv6 mixed ranges
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", mixedCidrs),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", mixedCidrs[0]),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.1", mixedCidrs[1]),
			// 	),
			// },
			// Update `cidrs` with a single ipv4 address without /32
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", ipv4Addr),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", ipv4Addr[0]),
				),
			},
			// Update `cidrs` with a single ipv6 address without /128
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", ipv6Addr),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.0", ipv6Addr[0]),
			// 	),
			// },
		},
	})
}

// TestAccPasswordResource_OutOfBandDelete tests the out-of-band deletion of a branch password.
// In this test we simulate the password has been deleted out of band, perhaps by
// a user on the console or using pscale CLI.
// https://github.com/planetscale/terraform-provider-planetscale/issues/53
func TestAccPasswordResource_OutOfBandDelete(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	branchName := acctest.RandomWithPrefix("branch")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")

	passId := ""

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "role", "admin"),
					resource.TestCheckResourceAttr("planetscale_password.test", "branch", branchName),
				),
			},
			// ImportState testing
			{
				ResourceName:      "planetscale_password.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					id, err := getPasswordIDFromState(s, "planetscale_password.test")
					if err != nil {
						return "", err
					}
					passId = id // save the ID for use in later steps' PreConfig func:
					// Import requires: 'organization,database,branch,id'
					return fmt.Sprintf("%s,%s,%s,%s", testAccOrg, dbName, branchName, id), nil
				},
				// The actual password is not returned by the API, so we can't verify it here:
				ImportStateVerifyIgnore: []string{"plaintext"},
			},
			// Test out-of-bands deletion of the database should produce a plan to recreate, not error.
			{
				ResourceName:       "planetscale_password.test",
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				PreConfig: func() {
					ctx := context.Background()
					if _, err := testAccAPIClient.DeletePassword(ctx, testAccOrg, dbName, branchName, passId); err != nil {
						t.Fatalf("PreConfig: failed to delete password: %s", err)
					}
				},
			},
		},
	})
}

func testAccPasswordResourceConfig(dbName, branchName, passwdName string, cidrs []string) string {
	const tmpl = `
resource "planetscale_database" "test" {
  name           = "{{.DBName}}"
  organization   = "{{.Org}}"
  cluster_size   = "PS-10"
  default_branch = "main"
}

resource "planetscale_branch" "test" {
  name          = "{{.BranchName}}"
  organization  = "{{.Org}}"
  database      = planetscale_database.test.name
  parent_branch = planetscale_database.test.default_branch
}

resource "planetscale_password" "test" {
  name         = "{{.PasswdName}}"
  organization = "{{.Org}}"
  database     = planetscale_database.test.name
  branch       = planetscale_branch.test.name
	{{if .CIDRs}}cidrs = {{.CIDRs}}{{end}}
}
`
	t := template.Must(template.New("config").Parse(tmpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, map[string]string{
		"Org":        testAccOrg,
		"DBName":     dbName,
		"BranchName": branchName,
		"PasswdName": passwdName,
		"CIDRs":      terraformStringList(cidrs),
	})
	if err != nil {
		return ""
	}
	return buf.String()
}

func terraformStringList(items []string) string {
	if len(items) == 0 {
		return `"null"`
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("%q", item)
	}
	return fmt.Sprintf("[%s]", strings.Join(quoted, ", "))
}

func getPasswordIDFromState(state *terraform.State, resourceName string) (string, error) {
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	if rawState == nil {
		return "", fmt.Errorf("resource %s not found in state", resourceName)
	}
	return rawState["id"], nil
}
