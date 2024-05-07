package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRulesMappingResource(t *testing.T) {
	var TEST_RULES_ID = randomString(10)
	var TEST_DOMAIN = randomString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRulesMapping(TEST_DOMAIN, TEST_RULES_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain", TEST_DOMAIN),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "rules_id", TEST_RULES_ID),
				),
			},
			{
				Config: testAccRulesMapping(TEST_DOMAIN, TEST_RULES_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "domain", TEST_DOMAIN),
					resource.TestCheckResourceAttr("altitude_mte_domain_mapping.tester", "rules_id", TEST_RULES_ID),
				),
			},
		},
	})
}

func testAccRulesMapping(domain string, rulesId string) string {
	return fmt.Sprintf(`
resource "altitude_mte_rules_mapping" "tester" {
  domain         = "%s"
  rules_id = "%s"
}
`, domain, rulesId)
}
