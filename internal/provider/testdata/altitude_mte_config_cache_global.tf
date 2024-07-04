resource "altitude_mte_config" "cache-field-test" {
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
        append_path_prefix   = "foo"
      }
    ]
    basic_auth = {
      username = "foobar",
      password = "barfoo"
    }
    cache = [
      {
        keys = {
          headers = ["foo"]
          cookies = ["bar"]
        }
        ttl_seconds = 100
      }
    ]
  }
  environment_id = "%s"
}
