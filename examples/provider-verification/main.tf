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
