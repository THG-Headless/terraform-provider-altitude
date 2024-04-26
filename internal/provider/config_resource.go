package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

type MTEResourceModel struct {
	EnvironmentId types.String   `tfsdk:"environment_id"`
	Config        MTEConfigModel `tfsdk:"config"`
}

type MTEConfigModel struct {
	Routes    []RouteModel    `tfsdk:"routes"`
	BasicAuth *BasicAuthModel `tfsdk:"basic_auth"`
}

type RouteModel struct {
	Host               types.String    `tfsdk:"host"`
	Path               types.String    `tfsdk:"path"`
	EnableSsl          types.Bool      `tfsdk:"enable_ssl"`
	PreservePathPrefix types.Bool      `tfsdk:"preserve_path_prefix"`
	CacheKey           *CacheKeyModel  `tfsdk:"cache_key"`
	AppendPathPrefix   types.String    `tfsdk:"append_path_prefix"`
	ShieldLocation     *ShieldLocation `tfsdk:"shield_location"`
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

// Post Request Body

type MTEConfigDto struct {
	Routes    []RoutesDto  `json:"routes"`
	BasicAuth BasicAuthDto `json:"basicAuth"`
}

type BasicAuthDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoutesDto struct {
	Host               string         `json:"host"`
	Path               string         `json:"path"`
	EnableSsl          bool           `json:"enableSsl"`
	PreservePathPrefix bool           `json:"preservePathPrefix"`
	CacheKey           CacheKeyDto    `json:"cacheKey"`
	AppendPathPrefix   string         `json:"appendPathPrefix"`
	ShieldLocation     ShieldLocation `json:"shieldLocation"`
}

type CacheKeyDto struct {
	Header []string `json:"header"`
	Cookie []string `json:"cookie"`
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
	m.baseUrl = resourceData.baseUrl
	m.apiKey = resourceData.apiKey
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
	var data MTEResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.updateMteConfig(
		ctx,
		data,
		true,
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
	var data MTEResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.deleteMteConfig(
		ctx,
		data.EnvironmentId.ValueString(),
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
	var data MTEResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	apiDto, err := m.getMteConfig(
		ctx,
		data.EnvironmentId.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to get MTE Config", err.Error())
		return
	}

	configModel := apiDto.transformToResourceModel()
	data.Config = configModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (m *MTEConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MTEResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.updateMteConfig(
		ctx,
		plan,
		false,
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

func (m *MTEConfigResource) updateMteConfig(
	ctx context.Context,
	data MTEResourceModel,
	isCreate bool,
) error {
	jsonBody, err := json.Marshal(data.transformToApiRequestBody())
	if err != nil {
		return err
	}
	var httpMethod string
	if isCreate {
		httpMethod = http.MethodPost
	} else {
		httpMethod = http.MethodPut
	}

	httpRes, err := m.createMteRequest(
		httpMethod,
		data.EnvironmentId.ValueString(),
		bytes.NewBuffer([]byte(jsonBody)))

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

	if httpRes.StatusCode != 201 {
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

func (m *MTEConfigResource) deleteMteConfig(
	ctx context.Context,
	environmentId string,
) error {
	httpRes, err := m.createMteRequest(http.MethodDelete, environmentId, nil)

	if err != nil {
		return &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return &AltitudeApiError{
			shortMessage: "Environment ID not found",
			detail:       fmt.Sprintf("The Environment %s does not have associated config.", environmentId),
		}
	}

	if httpRes.StatusCode != 204 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		tflog.Error(ctx, fmt.Sprintf("Body: %s", body))
		return &AltitudeApiError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The API deletion Request returned a non-200 response of %s.", httpRes.Status),
		}
	}

	return nil
}

func (m *MTEConfigResource) getMteConfig(
	ctx context.Context,
	environmentId string,
) (*MTEConfigDto, error) {
	httpRes, err := m.createMteRequest(http.MethodGet, environmentId, nil)

	if err != nil {
		return nil, &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return nil, &AltitudeApiError{
			shortMessage: "Environment ID not found",
			detail:       fmt.Sprintf("The Environment %s does not have associated config.", environmentId),
		}
	}

	if httpRes.StatusCode != 200 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		tflog.Error(ctx, fmt.Sprintf("Body: %s", body))
		return nil, &AltitudeApiError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The API deletion Request returned a non-200 response of %s.", httpRes.Status),
		}
	}

	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, &AltitudeApiError{
			shortMessage: "Body Read Error",
			detail:       "Unable to read response body",
		}
	}

	var dto MTEConfigDto
	err = json.Unmarshal(body, &dto)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("JSON Error: %s", err.Error()))
		return nil, &AltitudeApiError{
			shortMessage: "Body Read Error",
			detail:       "Unable to parse JSON body from Altitude response",
		}
	}

	return &dto, nil
}

func (m *MTEConfigResource) createMteRequest(
	method string,
	environmentId string,
	body io.Reader,
) (*http.Response, error) {
	httpReq, err := http.NewRequest(
		method,
		fmt.Sprintf("%s/v1/environment/%s/mte/altitude-config", m.baseUrl, environmentId),
		body,
	)
	if err != nil {
		return nil, &AltitudeApiError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}

	AddAuthenticationToRequest(httpReq, m.apiKey);
	return m.client.Do(httpReq)
}

func (m *MTEResourceModel) transformToApiRequestBody() MTEConfigDto {
	var httpRoutes = make([]RoutesDto, len(m.Config.Routes))
	for i, r := range m.Config.Routes {

		var routesPostBody = RoutesDto{
			Host:               r.Host.ValueString(),
			Path:               r.Path.ValueString(),
			EnableSsl:          r.EnableSsl.ValueBool(),
			PreservePathPrefix: r.PreservePathPrefix.ValueBool(),
			AppendPathPrefix:   r.AppendPathPrefix.ValueString(),
			ShieldLocation:     *r.ShieldLocation,
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
			routesPostBody.CacheKey = CacheKeyDto{
				Header: cacheKeyHeaders,
				Cookie: cacheKeyCookies,
			}
		}

		httpRoutes[i] = routesPostBody
	}

	return MTEConfigDto{
		Routes: httpRoutes,
		BasicAuth: BasicAuthDto{
			Username: m.Config.BasicAuth.Username.ValueString(),
			Password: m.Config.BasicAuth.Password.ValueString(),
		},
	}
}

func (d *MTEConfigDto) transformToResourceModel() MTEConfigModel {
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
			ShieldLocation: &r.ShieldLocation,
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
