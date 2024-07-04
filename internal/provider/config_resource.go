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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
	Cache     []CacheModel    `tfsdk:"cache"`
}

type RouteModel struct {
	Host               types.String `tfsdk:"host"`
	Path               types.String `tfsdk:"path"`
	EnableSsl          types.Bool   `tfsdk:"enable_ssl"`
	PreservePathPrefix types.Bool   `tfsdk:"preserve_path_prefix"`
	AppendPathPrefix   types.String `tfsdk:"append_path_prefix"`
	ShieldLocation     types.String `tfsdk:"shield_location"`
}

type CacheModel struct {
	Keys       *CacheKeyModel `tfsdk:"keys"`
	TtlSeconds types.Int64    `tfsdk:"ttl_seconds"`
	PathRules  *GlobMatcher   `tfsdk:"path_rules"`
}

type GlobMatcher struct {
	AnyMatch  []types.String `tfsdk:"any_match"`
	NoneMatch []types.String `tfsdk:"none_match"`
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

func (m *MTEConfigResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data MTEConfigResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Config.Cache != nil {
		for _, c := range data.Config.Cache {
			if c.Keys == nil && c.TtlSeconds == basetypes.NewInt64Null() {
				resp.Diagnostics.AddAttributeError(
					path.Root("cache"),
					"Missing Attribute Configuration",
					"Expected either `keys` or `ttl_seconds` to be set inside the cache object.",
				)
			}
		}
	}
}

// Schema implements resource.Resource.
func (m *MTEConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource which defines the various routes and other environment-specific config for a specific environment.",

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
									MarkdownDescription: "The downstream host MTE should direct to. This host should not contain the protocol or any slashes. A correct example would be docs.thgaltitude.com",
								},
								"path": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "The path prefix this route will be hosted on.",
								},
								"enable_ssl": schema.BoolAttribute{
									Required:            true,
									MarkdownDescription: "A boolean specifying whether the host defined requires a secure connection.",
								},
								"preserve_path_prefix": schema.BoolAttribute{
									Required: true,
									MarkdownDescription: "A boolean specifying whether we should retain the path specified above when routing to the host. " +
										"For example, if this was `true` and the path defined was `/foo`, when a client directs to `/foo/123` we would route " +
										"to the host with the path set as `/foo/123`. If it was `false`, we would point to `/123`.",
								},
								"append_path_prefix": schema.StringAttribute{
									Optional:            true,
									MarkdownDescription: "A string which will be appended to the start of the path sent to the host.",
								},
								"shield_location": schema.StringAttribute{
									Optional: true,
									MarkdownDescription: "This describes the location which all requests will be forwarded to before reaching the origin " +
										"of this route.",
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
								MarkdownDescription: "The username which clients will enter to authorize viewing this environment.",
							},
							"password": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The password which clients will enter to authorize viewing this environment.",
							},
						},
					},
					"cache": schema.ListNestedAttribute{
						Optional:            true,
						MarkdownDescription: "A list of settings designed to manipulate your cache without requiring you to set response headers.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"path_rules": schema.SingleNestedAttribute{
									Optional:            true,
									MarkdownDescription: "A set of glob rules which identify when the cache settings should be activated.",
									Attributes: map[string]schema.Attribute{
										"any_match": schema.ListAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "A list of glob paths where one of the list needs to match for the cache settings to be activated for a path. If both this field and `none_match` are specified, both need to be successful for the path to match.",
										},
										"none_match": schema.ListAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "A list of glob paths where all of the list needs to not match the path for the cache settings to be activated. If both this field and `any_match` are specified, both need to be successful for the path to match.",
										},
									},
								},
								"keys": schema.SingleNestedAttribute{
									Optional: true,
									Attributes: map[string]schema.Attribute{
										"headers": schema.ListAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "A list of header names of which the cache key will differeniate upon the values of these headers.",
										},
										"cookies": schema.ListAttribute{
											ElementType:         types.StringType,
											Required:            true,
											MarkdownDescription: "A list of cookie names which the cache key will differeniate upon the values of these cookies.",
										},
									},
									MarkdownDescription: "An object specifying header and cookie names which should be added to the cache key. The result " +
										"of this would lead to separate cache hits for requests with different values of the header or cookie. One of this",
								},
								"ttl_seconds": schema.Int64Attribute{
									Optional:            true,
									MarkdownDescription: "An integer that will be used to specify the time that the response of the route should be stored in the cache, in seconds.",
								},
							},
						},
					},
				},
			},
			"environment_id": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "The environment ID which this config associates with. If this value changes, this will " +
					"replace this resource. **Note**, if this occurred on a production site, this would lead to downtime.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTEConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("environment_id"), req, resp)
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
			"Failed to update MTE config",
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
			ShieldLocation:     client.ShieldLocation(r.ShieldLocation.ValueString()),
		}
		if r.AppendPathPrefix.ValueString() != "" {
			routesPostBody.AppendPathPrefix = r.AppendPathPrefix.ValueString()
		}
		httpRoutes[i] = routesPostBody
	}
	dto := client.MTEConfigDto{
		Routes: httpRoutes,
	}
	if m.Config.BasicAuth != nil {
		dto.BasicAuth = &client.BasicAuthDto{
			Username: m.Config.BasicAuth.Username.ValueString(),
			Password: m.Config.BasicAuth.Password.ValueString(),
		}
	}
	if m.Config.Cache != nil {
		var httpCache = make([]client.CacheDto, len(m.Config.Cache))

		for i, r := range m.Config.Cache {
			var cacheBody = client.CacheDto{
				TtlSeconds: r.TtlSeconds.ValueInt64Pointer(),
			}
			if r.Keys != nil {
				var cacheKeyHeaders = make([]string, len(r.Keys.Headers))
				var cacheKeyCookies = make([]string, len(r.Keys.Cookies))
				for i, h := range r.Keys.Headers {
					cacheKeyHeaders[i] = h.ValueString()
				}
				for i, c := range r.Keys.Cookies {
					cacheKeyCookies[i] = c.ValueString()
				}
				cacheBody.Keys = &client.CacheKeyDto{
					Header: cacheKeyHeaders,
					Cookie: cacheKeyCookies,
				}
			}
			if r.PathRules != nil {
				var anyMatch = make([]string, len(r.PathRules.AnyMatch))
				var noneMatch = make([]string, len(r.PathRules.NoneMatch))
				for i, h := range r.PathRules.AnyMatch {
					anyMatch[i] = h.ValueString()
				}
				for i, c := range r.PathRules.NoneMatch {
					noneMatch[i] = c.ValueString()
				}
				cacheBody.PathRules = &client.MatcherDto{
					AnyMatch:  anyMatch,
					NoneMatch: noneMatch,
				}
			}
			httpCache[i] = cacheBody
		}
		dto.Cache = httpCache
	}
	return dto
}

func transformToResourceModel(d *client.MTEConfigDto) MTEConfigModel {
	var routeModels = make([]RouteModel, len(d.Routes))
	for i, r := range d.Routes {
		var routesPostBody = RouteModel{
			Host:               types.StringValue(r.Host),
			Path:               types.StringValue(r.Path),
			EnableSsl:          types.BoolValue(r.EnableSsl),
			PreservePathPrefix: types.BoolValue(r.PreservePathPrefix),
		}
		if r.ShieldLocation != "" {
			routesPostBody.ShieldLocation = types.StringValue(string(r.ShieldLocation))
		}
		if r.AppendPathPrefix != "" {
			routesPostBody.AppendPathPrefix = types.StringValue(r.AppendPathPrefix)
		}

		routeModels[i] = routesPostBody
	}

	model := MTEConfigModel{
		Routes: routeModels,
	}
	if d.BasicAuth != nil {
		model.BasicAuth = &BasicAuthModel{
			Username: types.StringValue(d.BasicAuth.Username),
			Password: types.StringValue(d.BasicAuth.Password),
		}
	}

	if d.Cache != nil {
		var cacheModels = make([]CacheModel, len(d.Cache))
		for i, c := range d.Cache {
			var cacheModel = CacheModel{}
			if c.TtlSeconds != nil {
				cacheModel.TtlSeconds = types.Int64Value(*c.TtlSeconds)
			}
			if c.Keys != nil {
				var cacheKeyHeaders = make([]types.String, len(c.Keys.Header))
				var cacheKeyCookies = make([]types.String, len(c.Keys.Cookie))
				for i, h := range c.Keys.Header {
					cacheKeyHeaders[i] = types.StringValue(h)
				}
				for i, c := range c.Keys.Cookie {
					cacheKeyCookies[i] = types.StringValue(c)
				}
				cacheModel.Keys = &CacheKeyModel{
					Headers: cacheKeyHeaders,
					Cookies: cacheKeyCookies,
				}
			}
			if c.PathRules != nil {
				var anyMatch = make([]types.String, len(c.PathRules.AnyMatch))
				var noneMatch = make([]types.String, len(c.PathRules.NoneMatch))
				for i, h := range c.PathRules.AnyMatch {
					anyMatch[i] = types.StringValue(h)
				}
				for i, c := range c.PathRules.NoneMatch {
					noneMatch[i] = types.StringValue(c)
				}
				cacheModel.PathRules = &GlobMatcher{
					AnyMatch:  anyMatch,
					NoneMatch: noneMatch,
				}
			}

			cacheModels[i] = cacheModel
		}
		model.Cache = cacheModels
	}

	return model
}
