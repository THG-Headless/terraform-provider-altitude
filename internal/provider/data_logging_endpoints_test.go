package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLoggingEndpointsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "altitude_mte_logging_endpoints" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.type"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.environment_id"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.config.dataset"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.config.project_id"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.config.table"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.config.email"),
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.config.secret_key"),
				),
			},
		},
	})
}
