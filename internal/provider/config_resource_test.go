package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigWithBasicAuthResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(11)
	var INITIAL_HOST = "www.thgaltitude.com"
	var SECONDARY_HOST = "www.altitude.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResource("testdata/altitude_mte_config_basic_auth_included.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.basic_auth.username", "foobar"),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_basic_auth_included.tf", TEST_ENVIRONMENT_ID, SECONDARY_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.routes.0.host", SECONDARY_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.basic_auth.username", "foobar"),
				),
			},
		},
	})
}

func TestAccConfigWithoutBasicAuthResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	var INITIAL_HOST = "www.thgaltitude.com"
	var SECONDARY_HOST = "www.altitude.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResource("testdata/altitude_mte_config_basic_auth_excluded.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_basic_auth_excluded.tf", TEST_ENVIRONMENT_ID, SECONDARY_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "config.routes.0.host", SECONDARY_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.basic-auth-test", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
		},
	})
}

func TestAccConfigWithCacheResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	var INITIAL_HOST = "www.thgaltitude.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResource("testdata/altitude_mte_config_cache_excluded.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_cache_key.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.path_rules.any_match.0", "/test**"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.keys.headers.0", "foo"),
					resource.TestCheckNoResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.ttl_seconds"),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_cache_max_age.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.path_rules.any_match.0", "/test**"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.ttl_seconds", "100"),
					resource.TestCheckNoResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.keys"),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_cache_full.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.path_rules.any_match.0", "/test**"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.ttl_seconds", "100"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.keys.headers.0", "foo"),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_cache_global.tf", TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.path_rules"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.ttl_seconds", "100"),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.cache.0.keys.headers.0", "foo"),
				),
			},
		},
	})
}

func TestAccConfigWithConditionalHeadersCreateUpdateDelete(t *testing.T) {
	var matching_header = "testHeader"
	var pattern = ".*123.*"
	var new_header = "testNewHeader"
	var match_value = "foo"
	var no_match_value = "bar"
	var updated_matching_header = "header2"
	var env_id = randomString(10)
	var host = "www.thgaltitude.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResource("testdata/altitude_mte_config_conditional_headers_included_1.tf", env_id, host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.matching_header", matching_header),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.pattern", pattern),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.new_header", new_header),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.match_value", match_value),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.no_match_value", no_match_value),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_config_conditional_headers_included_2.tf", env_id, host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers.0.matching_header", updated_matching_header),
				),
			},
			{
				Config: testAccKVResource("testdata/altitude_mte_configconditional_headers_excluded.tf", env_id, host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("altitude_mte_config.cond-header-test", "config.conditional_headers"),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "config.routes.0.host", host),
					resource.TestCheckResourceAttr("altitude_mte_config.cond-header-test", "environment_id", env_id),
				),
			},
		},
	})
}

func testAccKVResource(fileResource string, environmentId string, host string) string {
	b, err := os.ReadFile(fileResource)
	if err != nil {
		fmt.Println(err)
	}
	str := string(b)
	return fmt.Sprintf(str, host, environmentId)
}


