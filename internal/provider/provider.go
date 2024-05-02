package provider

import (
	"context"
	"os"
	"terraform-provider-altitude/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &altitudeProvider{}
var _ provider.ProviderWithFunctions = &altitudeProvider{}

type altitudeProvider struct {
	version string
}

type ConfiguredData struct {
	client *client.Client
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	BaseUrl      types.String `tfsdk:"base_url"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
	Mode         types.String `tfsdk:"mode"`
}

func (p *altitudeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "altitude"
	resp.Version = p.version
}

func (p *altitudeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"client_id": schema.StringAttribute{
				Description: "The unique identifier for the Auth0 Application.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret for the Auth0 Application. Used to sign and validate the Client ID specified.",
				Optional:    true,
				Sensitive:   true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "The environment selected for development which in turn sets the base URL if not specified. This value can be either `Production`, `UAT` or `Local`. It defaults to Local.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{string(client.Production),
						string(client.UAT),
						string(client.Local),
					}...,
					),
				},
			},
		},
	}
}

func (p *altitudeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseUrl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_url"),
			"Unknown Altitude API Base URL",
			"The provider cannot create the Altitude API client as there is an unknown configuration value for the Altitude API base URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the default variable.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown Altitude API Client ID",
			"The provider cannot create the Altitude API client as there is an unknown configuration value for the Altitude API client ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALTITUDE_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown Altitude API Client Secret",
			"The provider cannot create the Altitude API client as there is an unknown configuration value for the Altitude API client secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALTITUDE_CLIENT_SECRET environment variable.",
		)
	}

	if config.Mode.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("mode"),
			"Unknown Altitude API Mode",
			"The provider cannot create the Altitude API client as there is an unknown configuration value for the Altitude API mode. "+
				"Either target apply the source of the value first or set the value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	clientId := os.Getenv("ALTITUDE_CLIENT_ID")
	clientSecret := os.Getenv("ALTITUDE_CLIENT_SECRET")
	mode := client.Local

	if !config.Mode.IsNull() {
		mode = client.Mode(config.Mode.ValueString())
	}

	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if clientId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Altitude Client ID",
			"The provider cannot create the Altitude API client as there is a missing or empty value for the Altitude Client ID. "+
				"Set the client_id value in the configuration or use the ALTITUDE_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if clientSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Altitude Client Secret",
			"The provider cannot create the Altitude API client as there is a missing or empty value for the Altitude Client Secret. "+
				"Set the client_secret value in the configuration or use the ALTITUDE_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	client, err := client.New(
		clientId,
		clientSecret,
		mode,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Incorrect Client Configuration",
			"While configuring the provider, the Client could not be created "+
				"successfully. The error returned from the initialisation was:\n"+err.Error(),
		)
	}
	var downstreamData = ConfiguredData{
		client: client,
	}
	resp.DataSourceData = &downstreamData
	resp.ResourceData = &downstreamData
}

func (p *altitudeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewMTEConfigResource,
		NewMTEDomainMappingResource,
	}
}

func (p *altitudeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *altitudeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &altitudeProvider{
			version: version,
		}
	}
}
