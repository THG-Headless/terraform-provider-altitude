package provider

import (
	"context"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
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
	client  *http.Client
	baseUrl string
	apiKey  string
}

// ProviderModel describes the provider data model.
type ProviderModel struct {
	ApiKey  types.String `tfsdk:"api_key"`
	BaseUrl types.String `tfsdk:"base_url"`
}

func (p *altitudeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "altitude"
	resp.Version = p.version
}

func (p *altitudeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Altitude API Key",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Altitude API URL",
				Optional:    true,
			},
		},
	}
}

func (p *altitudeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	apiToken := os.Getenv("ALTITUDE_API_KEY")
	baseUrl := "https://api.platform.thgaltitude.com"

	var data ProviderModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiKey.ValueString() != "" {
		apiToken = data.ApiKey.ValueString()
	}
	if data.BaseUrl.ValueString() != "" {
		apiToken = data.BaseUrl.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"While configuring the provider, the API token was not found in "+
				"the ALTITUDE_API_KEY environment variable or provider "+
				"configuration block api_token attribute.",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	client := http.DefaultClient
	var downstreamData = ConfiguredData{
		client:  client,
		baseUrl: baseUrl,
		apiKey:  apiToken,
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
