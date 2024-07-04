resource "altitude_mte_config" "basic-auth-test" {
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
  }
  environment_id = "%s"
}
