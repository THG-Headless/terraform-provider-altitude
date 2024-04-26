terraform {
  required_providers {
    altitude = {
      source  = "thg-headless/altitude"
    }
  }
}

provider "altitude" {
  api_key = "{{APIKEY}}"
}

resource "altitude_mte_config" "item_1" {
  config = {
    routes = [
      {
        host                 = "yo"
        path                 = "yo"
        enable_ssl           = true
        preserve_path_prefix = true
        shield_location      = "New York City"
      }
    ]

    basic_auth = {
      username = "yo"
      password = "yo"
    }
  }
  environment_id = "123"

}


