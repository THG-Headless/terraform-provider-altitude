package provider

import (
	"context"
	"fmt"
	"terraform-provider-altitude/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MTEConfigResource{}
var _ resource.ResourceWithImportState = &MTEConfigResource{}

func NewMTEConfigResource() resource.Resource {
	return &MTEConfigResource{}
}

type MTEConfigResource struct {
	client *client.Client
}

type MTEConfigResourceModel struct {
	EnvironmentId types.String   `tfsdk:"environment_id"`
	Config        MTEConfigModel `tfsdk:"config"`
}

type MTEConfigModel struct {
	Routes    []RouteModel    `tfsdk:"routes"`
	BasicAuth *BasicAuthModel `tfsdk:"basic_auth"`
}

type RouteModel struct {
	Host               types.String   `tfsdk:"host"`
	Path               types.String   `tfsdk:"path"`
	EnableSsl          types.Bool     `tfsdk:"enable_ssl"`
	PreservePathPrefix types.Bool     `tfsdk:"preserve_path_prefix"`
	CacheKey           *CacheKeyModel `tfsdk:"cache_key"`
	AppendPathPrefix   types.String   `tfsdk:"append_path_prefix"`
	ShieldLocation     ShieldLocation `tfsdk:"shield_location"`
}

type ShieldLocation string

const (
	London        ShieldLocation = "London"
	Manchester    ShieldLocation = "Manchester"
	Frankfurt                    = "Frankfurt"
	Madrid                       = "Madrid"
	New_York_City                = "New York City"
	Los_Angeles                  = "Los Angeles"
	Toronto                      = "Toronto"
	Johannesburg                 = "Johannesburg"
	Seoul                        = "Seoul"
	Sydney                       = "Sydney"
	Tokyo                        = "Tokyo"
	Hong_Kong                    = "Hong Kong"
	Mumbai                       = "Mumbai"
	Singapore                    = "Singapore"
)

type CacheKeyModel struct {
	Headers []types.String `tfsdk:"headers"`
	Cookies []types.String `tfsdk:"cookies"`
}

type BasicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata implements resource.Resource.
func (m *MTEConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mte_config"
}

func (m *MTEConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	m.client = resourceData.client
}

// Schema implements resource.Resource.
func (m *MTEConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Config lock to enable mte",

		Attributes: map[string]schema.Attribute{
			"config": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"routes": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "yo",
								},
								"path": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "yo",
								},
								"enable_ssl": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "yo",
								},
								"preserve_path_prefix": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "yo",
								},
								"cache_key": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"headers": schema.ListAttribute{
											ElementType: types.StringType,
											Required:    true,
										},
										"cookies": schema.ListAttribute{
											ElementType: types.StringType,
											Required:    true,
										},
									},
									MarkdownDescription: "yo",
								},
								"append_path_prefix": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "yo",
								},
								"shield_location": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "yo",
									Validators: []validator.String{
										stringvalidator.OneOf([]string{string(London),
											string(Manchester),
											string(New_York_City),
											string(Frankfurt),
											string(Madrid),
											string(Los_Angeles),
											string(Toronto),
											string(Johannesburg),
											string(Seoul),
											string(Sydney),
											string(Tokyo),
											string(Hong_Kong),
											string(Mumbai),
											string(Singapore),
										}...,
										),
									},
								},
							},
						},
					},
					"basic_auth": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"username": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "yo",
							},
							"password": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "yo",
							},
						},
					},
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "yo",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTEConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create implements resource.Resource.
func (m *MTEConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MTEConfigResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.CreateMTEConfig(
		client.CreateMTEConfigInput{
			Config:        data.transformToApiRequestBody(),
			EnvironmentId: data.EnvironmentId.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create MTE config",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (m *MTEConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MTEConfigResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.DeleteMTEConfig(
		client.DeleteMTEConfigInput{
			EnvironmentId: data.EnvironmentId.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Delete Resource",
			"An unexpected error occurred while executing the request. "+
				"Please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error(),
		)
		return
	}
}

// Read implements resource.Resource.
func (m *MTEConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MTEConfigResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	apiDto, err := m.client.ReadMTEConfig(
		client.ReadMTEConfigInput{
			EnvironmentId: data.EnvironmentId.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Get Resource",
			"An unexpected error occurred while executing the request. "+
				"Please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error(),
		)
		return
	}

	configModel := transformToResourceModel(apiDto)
	data.Config = configModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (m *MTEConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MTEConfigResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.UpdateMTEConfig(
		client.UpdateMTEConfigInput{
			Config:        plan.transformToApiRequestBody(),
			EnvironmentId: plan.EnvironmentId.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create MTE config",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (m *MTEConfigResourceModel) transformToApiRequestBody() client.MTEConfigDto {
	var httpRoutes = make([]client.RoutesDto, len(m.Config.Routes))
	for i, r := range m.Config.Routes {

		var routesPostBody = client.RoutesDto{
			Host:               r.Host.ValueString(),
			Path:               r.Path.ValueString(),
			EnableSsl:          r.EnableSsl.ValueBool(),
			PreservePathPrefix: r.PreservePathPrefix.ValueBool(),
			AppendPathPrefix:   r.AppendPathPrefix.ValueString(),
			ShieldLocation:     client.ShieldLocation(r.ShieldLocation),
		}

		if r.CacheKey != nil {
			var cacheKeyHeaders = make([]string, len(r.CacheKey.Headers))
			var cacheKeyCookies = make([]string, len(r.CacheKey.Cookies))
			for i, h := range r.CacheKey.Headers {
				cacheKeyHeaders[i] = h.ValueString()
			}
			for i, c := range r.CacheKey.Cookies {
				cacheKeyCookies[i] = c.ValueString()
			}
			routesPostBody.CacheKey = client.CacheKeyDto{
				Header: cacheKeyHeaders,
				Cookie: cacheKeyCookies,
			}
		}

		httpRoutes[i] = routesPostBody
	}

	return client.MTEConfigDto{
		Routes: httpRoutes,
		BasicAuth: client.BasicAuthDto{
			Username: m.Config.BasicAuth.Username.ValueString(),
			Password: m.Config.BasicAuth.Password.ValueString(),
		},
	}
}

func transformToResourceModel(d *client.MTEConfigDto) MTEConfigModel {
	var routeModels = make([]RouteModel, len(d.Routes))
	for i, r := range d.Routes {
		var cacheKeyHeaders = make([]types.String, len(r.CacheKey.Header))
		var cacheKeyCookies = make([]types.String, len(r.CacheKey.Cookie))
		for i, h := range r.CacheKey.Header {
			cacheKeyHeaders[i] = types.StringValue(h)
		}
		for i, c := range r.CacheKey.Cookie {
			cacheKeyCookies[i] = types.StringValue(c)
		}

		var routesPostBody = RouteModel{
			Host:               types.StringValue(r.Host),
			Path:               types.StringValue(r.Path),
			EnableSsl:          types.BoolValue(r.EnableSsl),
			PreservePathPrefix: types.BoolValue(r.PreservePathPrefix),
			AppendPathPrefix:   types.StringValue(r.AppendPathPrefix),
			CacheKey: &CacheKeyModel{
				Headers: cacheKeyHeaders,
				Cookies: cacheKeyCookies,
			},
			ShieldLocation: ShieldLocation(r.ShieldLocation),
		}

		routeModels[i] = routesPostBody
	}

	return MTEConfigModel{
		Routes: routeModels,
		BasicAuth: &BasicAuthModel{
			Username: types.StringValue(d.BasicAuth.Username),
			Password: types.StringValue(d.BasicAuth.Password),
		},
	}
}
