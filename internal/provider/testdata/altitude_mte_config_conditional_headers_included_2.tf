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
		conditional_headers = [
			{
				matching_header = "header2"
				pattern         = ".*123.*"
				new_header      = "testNewHeader"
				match_value     = "foo"
				no_match_value  = "bar"
			}
		]
	  }
	  environment_id = "%s"
	}