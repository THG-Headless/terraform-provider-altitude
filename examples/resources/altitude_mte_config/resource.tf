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
    conditional_headers = [
      {
        matching_header = "foo"
        pattern         = "*.pattern.*"
        new_header      = "bar"
        match_value     = "match"
        no_match_value  = "no match"
      }
    ]
  }
  environment_id = "test"
}