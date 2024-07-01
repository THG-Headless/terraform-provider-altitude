package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-altitude/internal/provider/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LoggingEndpointsDataSource{}
	_ datasource.DataSourceWithConfigure = &LoggingEndpointsDataSource{}
)

// NewLoggingEndpointsDataSource is a helper function to simplify the provider implementation.
func NewLoggingEndpointsDataSource() datasource.DataSource {
	return &LoggingEndpointsDataSource{}
}

// loggingEndpointsDataSource is the data source implementation.
type LoggingEndpointsDataSource struct {
	client *client.Client
}

// loggingEndpointsDataSourceModel maps the data source schema data.

type LoggingEndpointsDataSourceModel struct {
	Endpoints []LoggingEndpointDataSourceModel `tfsdk:"endpoints"`
}
type LoggingEndpointDataSourceModel struct {
	ID            types.String                        `tfsdk:"id"`
	Type          types.String                        `tfsdk:"type"`
	EnvironmentId types.String                        `tfsdk:"environmentid"`
	Config        GetAbstractAccessLoggingConfigModel `tfsdk:"config"`
}

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

// Configure adds the provider configured client to the data source.
func (d *LoggingEndpointsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	resourceData, ok := req.ProviderData.(*ConfiguredData)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ConfiguredData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = resourceData.client
}

// Metadata returns the data source type name.
func (d *LoggingEndpointsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mte_logging_endpoints"
}

// Schema defines the schema for the data source.
func (d *LoggingEndpointsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoints": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"environmentid": schema.StringAttribute{
							Computed: true,
						},
						"config": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"dataset": schema.StringAttribute{
									Computed: true,
								},
								"projectid": schema.StringAttribute{
									Computed: true,
								},
								"table": schema.StringAttribute{
									Computed: true,
								},
								"email": schema.StringAttribute{
									Computed: true,
								},
								"headers": schema.ListNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"columnname": schema.StringAttribute{
												Computed: true,
											},
											"headername": schema.StringAttribute{
												Computed: true,
											},
											"defaultvalue": schema.StringAttribute{
												Computed: true,
											},
										},
									},
								},
								"secretkey": schema.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the terraform state with the latest data.
func (d *LoggingEndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LoggingEndpointsDataSourceModel

	loggingEndpoints, err := d.client.ReadMTELoggingEndpoints()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get logging endpoints from Altitude provider",
			err.Error(),
		)
		return
	}

	for _, endpoint := range loggingEndpoints.Endpoints {
		endpointState := LoggingEndpointDataSourceModel{
			Type:          types.StringValue(endpoint.Type),
			EnvironmentId: types.StringValue(endpoint.EnvironmentId),
			Config:        transformToConfigResourceModel(endpoint.Config),
			ID:            types.StringValue("placeholder"),
		}

		state.Endpoints = append(state.Endpoints, endpointState)
	}

	for i, endpoint := range loggingEndpoints.Endpoints {
		endpointState := LoggingEndpointDataSourceModel{
			Type:          types.StringValue(endpoint.Type),
			EnvironmentId: types.StringValue(endpoint.EnvironmentId),
			Config:        transformToConfigResourceModel(endpoint.Config),
			ID:            types.StringValue("placeholder"),
		}

		state.Endpoints[i] = endpointState
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func transformToConfigResourceModel(d client.MTELoggingEndpointsConfig) GetAbstractAccessLoggingConfigModel {
	var headerModels = make([]BqLoggingHeaderModel, len(d.Headers))

	for i, r := range d.Headers {
		var headersBody = BqLoggingHeaderModel{
			ColumnName:   types.StringValue(r.ColumnName),
			HeaderName:   types.StringValue(r.HeaderName),
			DefaultValue: types.StringValue(r.DefaultValue),
		}

		headerModels[i] = headersBody
	}

	model := GetAbstractAccessLoggingConfigModel{
		Dataset:   types.StringValue(d.Dataset),
		ProjectId: types.StringValue(d.ProjectId),
		Table:     types.StringValue(d.Table),
		Email:     types.StringValue(d.Email),
		Headers:   headerModels,
		SecretKey: types.StringValue(d.SecretKey),
	}

	return model
}
