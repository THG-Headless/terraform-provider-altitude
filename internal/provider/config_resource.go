package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	client  *http.Client
	baseUrl string
	apiKey  string
}

type MTEConfigResourceModel struct {
	EnvironmentId types.String       `tfsdk:"environment_id"`
	Config        MTEConfigBodyModel `tfsdk:"config"`
}

type MTEConfigBodyModel struct {
	Routes    []RouteModel   `tfsdk:"routes"`
	BasicAuth BasicAuthModel `tfsdk:"basic_auth"`
}

type RouteModel struct {
	host               types.String   `tfsdk:"host"`
	path               types.String   `tfsdk:"path"`
	enableSsl          types.Bool     `tfsdk:"enable_ssl"`
	preservePathPrefix types.Bool     `tfsdk:"preserve_path_prefix"`
	cacheKey           CacheKeyModel  `tfsdk:"cache_key"`
	appendPathPrefix   types.String   `tfsdk:"append_path_prefix"`
	shieldLocation     ShieldLocation `tfsdk:"shield_location"`
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
	headers []types.String `tfsdk:"headers"`
	cookies []types.String `tfsdk:"cookies"`
}

type BasicAuthModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Post Request Body

type HTTPPostRequestBody struct {
	Routes    RoutesPostBody
	BasicAuth BasicAuthPostBody
}

type BasicAuthPostBody struct {
	username string
	password string
}

type RoutesPostBody struct {
	host               string
	path               string
	enableSsl          bool
	preservePathPrefix bool
	cacheKey           CacheKeyModel
	appendPathPrefix   string
	shieldLocation     ShieldLocation
}

type CacheKeyRequestBodyModel struct {
	header string
	cookie string
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
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTEConfigResource) ImportState(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Create implements resource.Resource.
func (m *MTEConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MTEConfigResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.createMteConfig(
		ctx,
		data,
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

func (m *MTEConfigResource) createMteConfig(
	ctx context.Context,
	data MTEConfigResourceModel,
) error {
	httpReq, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/environment/%s/mte/altitude-config", m.baseUrl, data.EnvironmentId.ValueString()),
		nil,
	)
}

func (m *MTEConfigResourceModel) transformToHttpPostBody(
	data MTEConfigResourceModel,
) HTTPPostRequestBody {

	var httpRoutes = []RoutesPostBody{}
	for i := 0; i < len(data.Config.Routes); i++ {
		var r = data.Config.Routes[i]
		var password = data.Config.BasicAuth.Password
		var username = data.Config.BasicAuth.Username
		var basicAuthPostBody = BasicAuthPostBody{
			username: username.ValueString(),
			password: password.ValueString(),
		}
		var routesPostBody = RoutesPostBody{
			host:               r.host.ValueString(),
			path:               r.path.ValueString(),
			enableSsl:          r.enableSsl.ValueBool(),
			preservePathPrefix: r.preservePathPrefix.ValueBool(),
			appendPathPrefix:   r.appendPathPrefix.ValueString(),
		}

		//this was done quickly not too sure

		var postRequestBody = HTTPPostRequestBody{
			RoutesPostBody:    routesPostBody,
			BasicAuthPostBody: basicAuthPostBody,
		}

		httpRoutes = append(httpRoutes, routesPostBody)
	}

}

//http request = postrequestbody + basic auth

//transform mte config resource model into http body
