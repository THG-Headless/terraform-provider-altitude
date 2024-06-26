package provider

import (
	"fmt"
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
				Config: testAccKVResourceConfigWithBasicAuth(TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.basic_auth.username", "foobar"),
				),
			},
			{
				Config: testAccKVResourceConfigWithBasicAuth(TEST_ENVIRONMENT_ID, SECONDARY_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", SECONDARY_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.basic_auth.username", "foobar"),
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
				Config: testAccKVResourceConfigWithoutBasicAuth(TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
			{
				Config: testAccKVResourceConfigWithoutBasicAuth(TEST_ENVIRONMENT_ID, SECONDARY_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", SECONDARY_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", TEST_ENVIRONMENT_ID),
				),
			},
		},
	})
}

func TestAccConfigWithCacheMaxAgeResource(t *testing.T) {
	var TEST_ENVIRONMENT_ID = randomString(10)
	var INITIAL_HOST = "www.thgaltitude.com"
	var CACHE_MAX_AGE = 360
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResourceConfigNoCacheMaxAge(TEST_ENVIRONMENT_ID, INITIAL_HOST),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckNoResourceAttr("altitude_mte_config.cache-field-test", "config.routes.0.cache_max_age"),
				),
			},
			{
				Config: testAccKVResourceConfigCacheMaxAge(TEST_ENVIRONMENT_ID, INITIAL_HOST, int64(CACHE_MAX_AGE)),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.routes.0.host", INITIAL_HOST),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "environment_id", TEST_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.routes.0.cache_max_age", fmt.Sprintf("%d", CACHE_MAX_AGE)),
				),
			},
		},
	})
}

func TestAccConfigWithConditionalHeadersCreateUpdateDelete(t *testing.T) {
	var matching_header = "testHeader"
	var pattern = "*123*"
	var new_header = "testNewHeader"
	var match_value = "foo"
	var no_match_value = "bar"
	var updated_matching_header = "header2"
	var env_id = randomString(9)
	var host = "www.thgaltitude.com"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKVResourceConfigWithConditionalHeaders(matching_header, pattern, new_header, match_value, no_match_value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.matching_header", matching_header),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.pattern", pattern),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.new_header", new_header),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.match_value", match_value),
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.no_match_value", no_match_value),
				),
			},
			{
				Config: testAccKVResourceConfigWithConditionalHeaders(updated_matching_header, pattern, new_header, match_value, no_match_value),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers.0.matching_header", updated_matching_header),
				),
			},
			{
				Config: testAccKVResourceConfigWithoutBasicAuth(env_id, host),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("altitude_mte_config.cache-field-test", "config.conditional_headers"),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "config.routes.0.host", host),
					resource.TestCheckResourceAttr("altitude_mte_config.tester", "environment_id", env_id),
				),
			},
		},
	})
}

func testAccKVResourceConfigWithBasicAuth(environmentId string, host string) string {
	return fmt.Sprintf(`
resource "altitude_mte_config" "tester" {
  config = {
    routes = [
      {
        host                 = "%s"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location      = "London"
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
`, host, environmentId)
}

func testAccKVResourceConfigWithoutBasicAuth(environmentId string, host string) string {
	return fmt.Sprintf(`
resource "altitude_mte_config" "tester" {
  config = {
    routes = [
      {
        host                 = "%s"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location		 = "London"
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
`, host, environmentId)
}

func testAccKVResourceConfigNoCacheMaxAge(environmentId string, host string) string {
	return fmt.Sprintf(`
resource "altitude_mte_config" "cache-field-test" {
  config = {
    routes = [
      {
        host                 = "%s"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location		 = "London"
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
`, host, environmentId)
}

func testAccKVResourceConfigCacheMaxAge(environmentId string, host string, cacheMaxAge int64) string {
	return fmt.Sprintf(`
	resource "altitude_mte_config" "cache-field-test" {
	  config = {
		routes = [
		  {
			host                 = "%s"
			path                 = "/test"
			enable_ssl           = true
			preserve_path_prefix = true
			shield_location		 = "London"
			cache_max_age  		 = %d
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
	`, host, cacheMaxAge, environmentId)
}

func testAccKVResourceConfigWithConditionalHeaders(matchHeader string, pattern string, newHeader string, matchValue string, noMatchValue string) string {
	return fmt.Sprintf(`
	resource "altitude_mte_config" "cache-field-test" {
	  config = {
		routes = [
		  {
			host                 = "testhost"
			path                 = "/test"
			enable_ssl           = true
			preserve_path_prefix = true
			shield_location		 = "London"
		  },
		]
		conditional_headers = [
			{
				matching_header = "%s"
				pattern         = "%s"
				new_header      = "%s"
				match_value     = "%s"
				no_match_value  = "%s"
			}
		]
	  }
	  environment_id = "testenvid"
	}
	`, matchHeader, pattern, newHeader, matchValue, noMatchValue)
}