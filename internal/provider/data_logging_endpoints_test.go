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
				Config: `data "logging_endpoints" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.logging_endpoints.test", "logging_endpoints.type", TEST_TYPE),
					resource.TestCheckResourceAttr("data.logging_endpoints.test", "logging_endpoints.environmentid", TEST_ENVIRONMENTID),
					resource.TestCheckResourceAttr("data.logging_endpoints.test", "logging_endpoints.config", TEST_CONFIG),
					resource.TestCheckResourceAttr("data.hashicups_coffees.test", "id", "placeholder"),
				),
			},
		},
	})
}
