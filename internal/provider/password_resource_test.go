package provider

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestValidateNonOverlappingCIDRs(t *testing.T) {
	tests := []struct {
		name    string
		cidrs   []string
		wantErr bool
	}{
		{
			name:    "empty list is valid",
			cidrs:   []string{},
			wantErr: false,
		},
		{
			name:    "single CIDR is valid",
			cidrs:   []string{"10.0.0.0/24"},
			wantErr: false,
		},
		{
			name:    "non-overlapping IPv4 CIDRs are valid",
			cidrs:   []string{"10.0.0.0/24", "10.0.1.0/24", "192.168.1.0/24"},
			wantErr: false,
		},
		{
			name:    "non-overlapping IPv6 CIDRs are valid",
			cidrs:   []string{"2001:db8::/32", "2001:db9::/32"},
			wantErr: false,
		},
		{
			name:    "mixed non-overlapping IPv4/IPv6 CIDRs are valid",
			cidrs:   []string{"10.0.0.0/24", "2001:db8::/32"},
			wantErr: false,
		},
		{
			name:    "duplicated IPv4 CIDRs are invalid",
			cidrs:   []string{"10.0.0.0/24", "10.0.0.0/24"},
			wantErr: true,
		},
		{
			name:    "overlapping IPv4 CIDRs are invalidd",
			cidrs:   []string{"10.0.0.0/8", "10.1.0.0/8"},
			wantErr: true,
		},
		{
			name:    "containing IPv4 CIDRs",
			cidrs:   []string{"10.0.0.0/16", "10.0.1.0/24"},
			wantErr: true,
		},
		{
			name:    "duplicated IPv6 CIDRs are invalid",
			cidrs:   []string{"2001:db8::/32", "2001:db8::/32"},
			wantErr: true,
		},
		{
			name:    "overlapping IPv6 CIDRs are invalid",
			cidrs:   []string{"2001:db8::/24", "2001:db8::/32"},
			wantErr: true,
		},
		{
			name:    "containing IPv6 CIDRs",
			cidrs:   []string{"2001:db8::/32", "2001:db8:1::/48"},
			wantErr: true,
		},
		{
			name:    "invalid CIDR format",
			cidrs:   []string{"10.0.0.0/24", "invalid"},
			wantErr: true,
		},
		{
			name:    "invalid CIDR prefix length",
			cidrs:   []string{"10.0.0.0/24", "10.0.0.0/33"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNonOverlappingCIDRs(tt.cidrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateNonOverlappingCIDRs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccPasswordResource(t *testing.T) {
	dbName := acctest.RandomWithPrefix("testacc-passwd-db")
	branchName := acctest.RandomWithPrefix("branch")
	passwdName := acctest.RandomWithPrefix("testacc-passwd")

	ipv4Cidr := []string{"10.0.0.0/24"}
	// ipv6Cidr := []string{"2001:db8::/64"}
	multiIPv4Cidrs := []string{"10.0.0.0/16", "10.1.0.0/16"}
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
					resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", ipv4Cidr[0]),
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

			// TODO: when the API supports ipv6:
			// Update 'cidrs' attribute with single ipv6 range:
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", ipv6Cidr),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
			// 		resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", ipv6Cidr[0]),
			// 	),
			// },

			// Update `cidrs` with multiple ipv4 ranges
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", multiIPv4Cidrs),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
					resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", multiIPv4Cidrs[0]),
					resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", multiIPv4Cidrs[1]),
				),
			},

			// TODO: when the API supports ipv6:
			// Update `cidrs` with multiple ipv6 ranges
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", multiIPv6Cidrs),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
			// 		resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", multiIPv6Cidrs[0]),
			// 		resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", multiIPv6Cidrs[1]),
			// 	),
			// },

			// TODO: when the API supports ipv6:
			// Update `cidrs` with ipv4 + ipv6 mixed ranges
			// {
			// 	Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", mixedCidrs),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "2"),
			// 		resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", mixedCidrs[0]),
			// 		resource.TestCheckTypeSetElemAttr("planetscale_password.test", "cidrs.*", mixedCidrs[1]),
			// 	),
			// },

			// Test removal of `cidrs`
			{
				Config: testAccPasswordResourceConfig(dbName, branchName, passwdName+"-new", nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("planetscale_password.test", "cidrs"),
					resource.TestCheckResourceAttr("planetscale_password.test", "cidrs.#", "1"),
				),
			},
		},
	})
}

func TestAccPasswordResource_ValidationFailures(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
                    resource "planetscale_password" "test" {
                        organization = "test-org"
                        database    = "test-db"
                        branch     = "main"
                        cidrs      = ["192.168.1.1"]  # Missing prefix
                    }`,
				ExpectError: regexp.MustCompile("CIDR notation required"),
			},
			{
				Config: `
                    resource "planetscale_password" "test" {
                        organization = "test-org"
                        database    = "test-db"
                        branch     = "main"
                        cidrs      = ["2001:db8::1"]  # Missing prefix
                    }`,
				ExpectError: regexp.MustCompile("CIDR notation required"),
			},
			{
				Config: `
                    resource "planetscale_password" "test" {
                        organization = "test-org"
                        database    = "test-db"
                        branch     = "main"
                        cidrs      = ["10.0.0.0/8", "10.1.0.0/8"]  # overlapping
                    }`,
				ExpectError: regexp.MustCompile("CIDR.+overlaps"),
			},
			{
				Config: `
                    resource "planetscale_password" "test" {
                        organization = "test-org"
                        database    = "test-db"
                        branch     = "main"
                        cidrs      = ["999.0.0.1/32"]  # invalid
                    }`,
				ExpectError: regexp.MustCompile("invalid CIDR"),
			},
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
		return `null`
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
