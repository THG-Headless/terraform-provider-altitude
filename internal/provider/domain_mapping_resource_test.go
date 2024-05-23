package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainMappingResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	var TEST_DOMAIN = randomString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainMapping(TEST_DOMAIN, TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain", TEST_DOMAIN),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain_mapping", TEST_ENVIRONMENT_ID),
				),
			},
			{
				Config: testAccDomainMapping(TEST_DOMAIN, TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain", TEST_DOMAIN),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain_mapping", TEST_ENVIRONMENT_ID),
				),
			},
		},
	})
}

func testAccDomainMapping(domain string, environmentId string) string {
	return fmt.Sprintf(`
resource "altitude_mte_domain_mapping" "tester" {
  domain         = "%s"
  environment_id = "%s"
}
`, domain, environmentId)
}
