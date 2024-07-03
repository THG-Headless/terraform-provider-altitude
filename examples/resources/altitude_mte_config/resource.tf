resource "altitude_mte_config" "config" {
  config = {
    routes = [
      {
        host                 = "www.thgaltitude.com"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location      = "London"
      }
    ]
    cache = [
      {
        path_rules = {
          any_match = [
            "/foo**"
          ]
        }
        keys = {
          headers = ["X-Header"]
          cookies = ["X-Cookie"]
        }
        ttl_seconds = 100
      }
    ]
  }
  environment_id = "test"
}