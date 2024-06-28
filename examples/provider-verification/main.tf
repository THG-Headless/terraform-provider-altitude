terraform {
  required_providers {
    altitude = {
      source = "thg-headless/altitude"
    }
  }
}

provider "altitude" {
  client_id     = "{{CLIENT ID}}"
  client_secret = "{{CLIENT SECRET}}"
  mode          = "Local"
}

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
    basic_auth = {
      username = "joe"
      password = "test"
    }
    conditional_headers = [{
      matching_header = "foo"
      pattern         = "*.123.*"
      new_header      = "bar"
      match_value     = "woop"
      no_match_value  = "sadge"
    }]
  }
  environment_id = "test"
}