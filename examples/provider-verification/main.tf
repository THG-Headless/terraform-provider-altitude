terraform {
  required_providers {
    altitude = {
      source = "thg-headless/altitude"
    }
  }
}

provider "altitude" {
  client_id     = "sJ8kVG1yCFyD8qOTLAS5yEy2F28ZM7qO"
  client_secret = "w455Q1EFsVkOKxm4u05Reuiy84mp1Ka3ihHSfHBctHvxNtzZ61W-HWxxETHzsfZN"
  mode          = "UAT"
}

data "altitude_mte_logging_endpoints" "test" {}

output "output" {
  value = data.altitude_mte_logging_endpoints.test
}