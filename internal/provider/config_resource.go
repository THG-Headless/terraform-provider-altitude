package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

type ApiMteConfigRequestBody struct {
	Routes    []RoutesReqestBody `json:"routes"`
	BasicAuth BasicAuthRequestBody `json:"basicAuth"`
}

type BasicAuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoutesReqestBody struct {
	Host               string `json:"host"`
	Path               string `json:"path"`
	EnableSsl          bool `json:"enableSsl"`
	PreservePathPrefix bool `json:"preservePathPrefix"`
	CacheKey           CacheKeyRequestBody `json:"cacheKey"`
	AppendPathPrefix   string `json:"appendPathPrefix"`
	ShieldLocation     ShieldLocation `json:"shieldLocation"`
}

type CacheKeyRequestBody struct {
	Header []string `json:"header"`
	Cookie []string `json:"cookie"`
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
	jsonBody, err := json.Marshal(data.transformToApiRequestBody())
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/v1/environment/%s/mte/altitude-config", m.baseUrl, data.EnvironmentId.ValueString()),
		bytes.NewBuffer([]byte(jsonBody)),
	)

	if err != nil {
		return &AltitudeApiError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}

	bearer := "Bearer " + m.apiKey
	httpReq.Header.Add("Authorization", bearer)

	httpRes, err := m.client.Do(httpReq)

	if err != nil {
		return &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 409 {
		return &AltitudeApiError{
			shortMessage: "Environment ID Conflict",
			detail:       "This environment already has an associated config block.",
		}
	}

	if httpRes.StatusCode != 200 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		tflog.Error(ctx, fmt.Sprintf("Body: %s", body))
		return &AltitudeApiError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s.", httpRes.Status),
		}
	}
	return nil
}

func (m *MTEConfigResourceModel) transformToApiRequestBody() ApiMteConfigRequestBody {
	var httpRoutes = make([]RoutesReqestBody, len(m.Config.Routes))
	for i, r := range m.Config.Routes {
		var cacheKeyHeaders = make([]string, len(r.cacheKey.headers))
		var cacheKeyCookies = make([]string, len(r.cacheKey.cookies))
		for i, h := range r.cacheKey.headers {
			cacheKeyHeaders[i] = h.ValueString()
		}
		for i, h := range r.cacheKey.cookies {
			cacheKeyCookies[i] = h.ValueString()
		}

		var routesPostBody = RoutesReqestBody{
			Host:               r.host.ValueString(),
			Path:               r.path.ValueString(),
			EnableSsl:          r.enableSsl.ValueBool(),
			PreservePathPrefix: r.preservePathPrefix.ValueBool(),
			AppendPathPrefix:   r.appendPathPrefix.ValueString(),
			CacheKey: CacheKeyRequestBody{
				Header: cacheKeyHeaders,
				Cookie: cacheKeyCookies,
			},
			ShieldLocation: r.shieldLocation,
		}

		httpRoutes[i] = routesPostBody
	}

	return ApiMteConfigRequestBody{
		Routes: httpRoutes,
		BasicAuth: BasicAuthRequestBody{
			Username: m.Config.BasicAuth.Username.ValueString(),
			Password: m.Config.BasicAuth.Password.ValueString(),
		},
	}

}

//http request = postrequestbody + basic auth

//transform mte config resource model into http body
