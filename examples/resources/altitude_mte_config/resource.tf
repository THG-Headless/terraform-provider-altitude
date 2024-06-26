resource "altitude_mte_config" "config" {
  config = {
    routes = [
      {
        host                 = "www.thgaltitude.com"
        path                 = "/test"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location      = "London"
        cache_max_age        = 360
      }
    ]
  }
  environment_id = "test"
}