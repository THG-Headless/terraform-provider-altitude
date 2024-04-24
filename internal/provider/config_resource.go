package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MTEConfigResource{}
var _ resource.ResourceWithImportState = &MTEConfigResource{}

func NewMTEConfigResource() resource.Resource {
	return &MTEConfigResource{}
}

type MTEConfigResource struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

type MTEConfigResourceModel struct {
	Routes []RouteModel `tfsdk:"routes"`
	BasicAuth BasicAuthModel `tfsdk:"basic_auth"`
}

type RouteModel struct {
	host types.String `tfsdk:"host"`
	path types.String `tfsdk:"path"`
	enableSsl types.Bool `tfsdk:"enable_ssl"`
	preservePathPrefix types.Bool `tfsdk:"preserve_path_prefix"`
	cacheKey CacheKeyModel `tfsdk:"cache_key"`
	appendPathPrefix types.String `tfsdk:"append_path_prefix"`
	shieldLocation ShieldLocation `tfsdk:"shield_location"`
}

type ShieldLocation string

const (
  London string = "London"
  Manchester string = "Manchester"
  Frankfurt = "Frankfurt"
  Madrid = "Madrid"
  New_York_City = "New York City"
  Los_Angeles = "Los Angeles"
  Toronto = "Toronto"
  Johannesburg = "Johannesburg"
  Seoul = "Seoul"
  Sydney = "Sydney"
  Tokyo = "Tokyo"
  Hong_Kong = "Hong Kong"
  Mumbai = "Mumbai"
  Singapore = "Singapore"
)

type CacheKeyModel struct {
	headers []types.String `tfsdk:"headers"`
	cookies []types.String `tfsdk:"cookies"`
}

type BasicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata implements resource.Resource.
func (m *MTEConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mte_config"
}

// Schema implements resource.Resource.
func (m *MTEConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Config lock to enable mte",

		Attributes: map[string]schema.Attribute{
			"routes": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							Required: true,
							MarkdownDescription: "yo",
						},
					},
				},
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTEConfigResource) ImportState(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Create implements resource.Resource.
func (m *MTEConfigResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (m *MTEConfigResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements resource.Resource.
func (m *MTEConfigResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (m *MTEConfigResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
