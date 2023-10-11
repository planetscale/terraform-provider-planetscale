package provider

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"planetscale": providerserver.NewProtocol6WithError(New("test", false)()),
}

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
}

func checkIntegerMin(min int) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if v < min {
			return fmt.Errorf("value %d is less than %d", v, min)
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
		return fmt.Errorf("valud %q is not one of %s", value, strings.Join(values, ", "))
	}
}
