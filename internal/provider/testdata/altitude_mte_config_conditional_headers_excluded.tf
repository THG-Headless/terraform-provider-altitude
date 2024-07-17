resource "altitude_mte_config" "cond-header-test" {
	  config = {
		routes = [
		  {
			host                 = "%s"
			path                 = "/test"
			enable_ssl           = true
			preserve_path_prefix = true
			shield_location		 = "London"
		  }
		]
	  }
	  environment_id = "%s"
	}