package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-altitude/internal/provider/client"
)

//Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &loggingEndpointsDataSource{}
	_ datasource.DataSourceWithConfigure = &loggingEndpointsDataSource{}
)

//NewLoggingEndpointsDataSource is a helper function to simplify the provider implementation.
func NewLoggingEndpointsDataSource() datasource.DataSource {
	return &loggingEndpointsDataSource{}
}

//coffeesDataSource is the data source implementation.
type loggingEndpointsDataSource struct {
	client *client.Client
}

//loggingEndpointsDataSourceModel maps the data source schema data.
type loggingEndpointsDataSourceModel struct {
	Type			types.String								`tfsdk:"type"`
	EnvironmentId	types.String								`tfsdk:"environmentId:`
	Config			getAbstractAccessLoggingConfigModel	`tfsdk:"config"`
}

type getAbstractAccessLoggingConfigModel struct {
}

//Configure adds the provider configured client to the data source.
func (d *loggingEndpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *altitudeProvider.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData), 
		)

		return
	}

	d.client = client
}
 

// Metadata returns the data source type name.
func (d *loggingEndpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_logging_endpoints"
}


//Schema defines the schema for the data source.
func(d *loggingEndpointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	
}


//Read refreshes the terraform state with the latest data.
func (d *loggingEndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state loggingEndpointsDataSourceModel

	loggingEndpoints, err := d.client.ReadMTELoggingEndpoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get logging endpoints from Altitude provider",
			err.Error(),
		)
		return
	}

	state.Type = types.StringValue(loggingEndpoints.Type)
	state.EnvironmentId = types.StringValue(loggingEndpoints.EnvironmentId)
	state.Config = AgetAbstractAccessLoggingConfigModel(loggingEndpoints.Config)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
