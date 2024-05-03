package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigWithBasicAuthResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResourceConfigWithBasicAuth(TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", "www.thgaltitude.com"),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.basic_auth.username", "foobar"),
				),
			},
			{
				Config: testAccKVResourceConfigWithBasicAuth(TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", "www.thgaltitude.com"),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.basic_auth.username", "foobar"),
				),
			},
		},
	})
}


func TestAccConfigWithoutBasicAuthResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResourceConfigWithoutBasicAuth(TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", "www.thgaltitude.com"),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
			{
				Config: testAccKVResourceConfigWithoutBasicAuth(TEST_ENVIRONMENT_ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", "www.thgaltitude.com"),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
		},
	})
}

func testAccKVResourceConfigWithBasicAuth(environmentId string) string {
	return fmt.Sprintf(`
resource "altitude_mte_config" "tester" {
  config = {
    routes = [
      {
        host                 = "www.thgaltitude.com"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
				shield_location			 = "London"
      },
      {
        host                 = "docs.thgaltitude.com"
        path                 = "/docs"
        enable_ssl           = false
        preserve_path_prefix = false
				append_path_prefix	 = "foo"
      }
    ]
		basic_auth = {
			username = "foobar",
			password = "barfoo"
		}
  }
  environment_id = "%s"
}
`, environmentId)
}

func testAccKVResourceConfigWithoutBasicAuth(environmentId string) string {
	return fmt.Sprintf(`
resource "altitude_mte_config" "tester" {
  config = {
    routes = [
      {
        host                 = "www.thgaltitude.com"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
				shield_location			 = "London"
      },
      {
        host                 = "docs.thgaltitude.com"
        path                 = "/docs"
        enable_ssl           = false
        preserve_path_prefix = false
				append_path_prefix	 = "foo"
      }
    ]
  }
  environment_id = "%s"
}
`, environmentId)
}
