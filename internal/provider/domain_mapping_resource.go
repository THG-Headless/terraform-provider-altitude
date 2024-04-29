package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MTEDomainMappingResource{}
var _ resource.ResourceWithImportState = &MTEDomainMappingResource{}

func NewMTEDomainMappingResource() resource.Resource {
	return &MTEDomainMappingResource{}
}

type MTEDomainMappingResource struct {
	client  *http.Client
	baseUrl string
	apiKey  string
}

type MTEDomainMappingResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Domain        types.String `tfsdk:"domain"`
	domainMapping types.String
}

type MTEDomainMappingDto struct {
	EnvironmentId string `json:"environmentId"`
	Domain        string `json:"domain"`
}

// Metadata implements resource.Resource.
func (m *MTEDomainMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_mapping"
}

func (m *MTEDomainMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (m *MTEDomainMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Config lock to enable mte",

		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "yo",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
func (m *MTEDomainMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Create implements resource.Resource.
func (m *MTEDomainMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MTEDomainMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainMapping, err := m.updateMTEDomainMapping(
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

	data.domainMapping = types.StringValue(domainMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (m *MTEDomainMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MTEDomainMappingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.deleteMTEDomainMapping(
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
func (m *MTEDomainMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MTEDomainMappingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	domainMapping, err := m.getMTEDomainMapping(
		ctx,
		data.EnvironmentId.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to get MTE Config", err.Error())
		return
	}

	data.domainMapping = types.StringValue(domainMapping)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (m *MTEDomainMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MTEDomainMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainMapping, err := m.updateMTEDomainMapping(
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

	plan.domainMapping = types.StringValue(domainMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (m *MTEDomainMappingResource) updateMTEDomainMapping(
	ctx context.Context,
	data MTEDomainMappingResourceModel,
	isCreate bool,
) (string, error) {
	jsonBody, err := json.Marshal(data.transformToApiRequestBody())
	if err != nil {
		return "", err
	}
	var httpMethod string
	if isCreate {
		httpMethod = http.MethodPost
	} else {
		httpMethod = http.MethodPut
	}

	httpReq, err := http.NewRequest(
		httpMethod,
		fmt.Sprintf("%s/v1/mte/domain-mapping", m.baseUrl),
		bytes.NewBuffer([]byte(jsonBody)),
	)
	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}

	AddAuthenticationToRequest(httpReq, m.apiKey)
	httpRes, err := m.client.Do(httpReq)

	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 409 {
		return "", &AltitudeApiError{
			shortMessage: "Domain Conflict",
			detail:       "This domain already has an associated config block.",
		}
	}

	if httpRes.StatusCode != 201 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		tflog.Error(ctx, fmt.Sprintf("Body: %s", body))
		return "", &AltitudeApiError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The Altitude API Request returned a non-200 response of %s.", httpRes.Status),
		}
	}

	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "Body Read Error",
			detail:       "Unable to read response body",
		}
	}

	return string(body[:]), nil
}

func (m *MTEDomainMappingResource) deleteMTEDomainMapping(
	ctx context.Context,
	domain string,
) error {

	httpReq, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/v1/mte/domain-mapping?domain=%s", m.baseUrl, domain),
		nil,
	)
	if err != nil {
		return &AltitudeApiError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}

	AddAuthenticationToRequest(httpReq, m.apiKey)
	httpRes, err := m.client.Do(httpReq)

	if err != nil {
		return &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return &AltitudeApiError{
			shortMessage: "Environment ID not found",
			detail:       fmt.Sprintf("The Environment %s does not have associated config.", domain),
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

func (m *MTEDomainMappingResource) getMTEDomainMapping(
	ctx context.Context,
	domain string,
) (string, error) {

	httpReq, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/v1/mte/domain-mapping?domain=%s", m.baseUrl, domain),
		nil,
	)
	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "Client Error",
			detail:       fmt.Sprintf("Unable to create http request, received error: %s", err),
		}
	}

	AddAuthenticationToRequest(httpReq, m.apiKey)
	httpRes, err := m.client.Do(httpReq)

	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "HTTP Error",
			detail:       fmt.Sprintf("There has been an error with the http request, received error: %s", err),
		}
	}

	if httpRes.StatusCode == 404 {
		return "", &AltitudeApiError{
			shortMessage: "Domain not found",
			detail:       fmt.Sprintf("The Domain %s does not have associated config.", domain),
		}
	}

	if httpRes.StatusCode != 200 {
		defer httpRes.Body.Close()
		body, _ := io.ReadAll(httpRes.Body)
		tflog.Error(ctx, fmt.Sprintf("Body: %s", body))
		return "", &AltitudeApiError{
			shortMessage: "Unexpected API Response",
			detail:       fmt.Sprintf("The API deletion Request returned a non-200 response of %s.", httpRes.Status),
		}
	}

	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return "", &AltitudeApiError{
			shortMessage: "Body Read Error",
			detail:       "Unable to read response body",
		}
	}

	return string(body[:]), nil
}

func (m *MTEDomainMappingResourceModel) transformToApiRequestBody() MTEDomainMappingDto {
	return MTEDomainMappingDto{
		EnvironmentId: m.EnvironmentId.ValueString(),
		Domain:        m.Domain.ValueString(),
	}
}
