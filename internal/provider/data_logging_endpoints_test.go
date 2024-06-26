package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type GetAbstractAccessLoggingConfigModel struct {
	Dataset   types.String           `tfsdk:"dataset"`
	ProjectId types.String           `tfsdk:"projectid"`
	Table     types.String           `tfsdk:"table"`
	Email     types.String           `tfsdk:"email"`
	Headers   []BqLoggingHeaderModel `tfsdk:"headers"`
	SecretKey types.String           `tfsdk:"secretkey"`
}

type BqLoggingHeaderModel struct {
	ColumnName   types.String `tfsdk:"columnname"`
	HeaderName   types.String `tfsdk:"headername"`
	DefaultValue types.String `tfsdk:"defaultvalue"`
}

func TestAccLoggingEndpointsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "altitude_mte_logging_endpoints" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.altitude_mte_logging_endpoints.test", "endpoints.0.type"),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.environmentid"),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.config"),
					resource.TestCheckResourceAttr("data.altitude_mte_logging_endpoints.test", "endpoints.0.id"),
				),
			},
		},
	})
}
