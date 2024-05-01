package provider

import (
	"context"
	"os"
	"terraform-provider-altitude/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	Audience     types.String `tfsdk:"audience"`
}

func (p *altitudeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "altitude"
	resp.Version = p.version
}

func (p *altitudeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				MarkdownDescription: "Altitude API URL. Defaults to `https://api.platform.thgaltitude.com`.",
				Optional:            true,
			},
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
			"audience": schema.StringAttribute{
				Description: "The Audience for an issued token, usually varies between test and prod environments.",
				Optional:    true,
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
			"Unknown HashiCups API Host",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the default variable.",
		)
	}

	if config.ClientId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown HashiCups API Username",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALTITUDE_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown HashiCups API Password",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALTITUDE_CLIENT_SECRET environment variable.",
		)
	}

	if config.Audience.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("audience"),
			"Unknown HashiCups API Audience",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API audience. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALTIUDE_AUDIENCE environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	clientId := os.Getenv("ALTITUDE_CLIENT_ID")
	clientSecret := os.Getenv("ALTITUDE_CLIENT_SECRET")
	audience := os.Getenv("ALTIUDE_AUDIENCE")
	baseUrl := "https://api.platform.thgaltitude.com"

	if !config.Audience.IsNull() {
		audience = config.Audience.ValueString()
	}

	if !config.ClientId.IsNull() {
		clientId = config.ClientId.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		clientSecret = config.ClientSecret.ValueString()
	}

	if !config.BaseUrl.IsNull() {
		baseUrl = config.BaseUrl.ValueString()
	} else {
		resp.Diagnostics.AddWarning(
			"Using Default Base URL",
			"The default base URL of "+baseUrl+" is being used. Please set the base_url configuration value if you do not want to use this default.",
		)

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

	if audience == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Altitude Audience",
			"The provider cannot create the Altitude API client as there is a missing or empty value for the Altitude Audience. "+
				"Set the audience value in the configuration or use the ALTIUDE_AUDIENCE environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	client, err := client.New(
		baseUrl,
		clientId,
		clientSecret,
		audience,
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
