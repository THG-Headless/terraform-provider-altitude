terraform {
  required_providers {
    altitude = {
      version = "0.0.2"
      source  = "thg-headless/altitude"
    }
  }
}

provider "altitude" {
  api_key = "{{APIKEY}}"
}

resource "altitude_mte_config" "item_1" {
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


