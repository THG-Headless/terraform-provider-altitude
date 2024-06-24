package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLoggingEndpointsDataSource(t *testing.T) {
	var TEST_TYPE = randomString(10)
	var TEST_ENVIRONMENTID = randomString(10)
	var TEST_CONFIG = randomString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "altitude_mte_logging_endpoints" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.type", TEST_TYPE),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.environmentid", TEST_ENVIRONMENTID),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.config", TEST_CONFIG),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.id", "placeholder"),
				),
			},
		},
	})
}
