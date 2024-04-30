package provider

import (
	"context"
	"fmt"
	"terraform-provider-altitude/internal/provider/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MTEDomainMappingResource{}
var _ resource.ResourceWithImportState = &MTEDomainMappingResource{}

func NewMTEDomainMappingResource() resource.Resource {
	return &MTEDomainMappingResource{}
}

type MTEDomainMappingResource struct {
	client  *client.Client
}

type MTEDomainMappingResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Domain        types.String `tfsdk:"domain"`
	domainMapping types.String
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

	domainMapping, err := m.client.CreateMteDomainMapping(
		client.CreateMteDomainMappingInput{
			Config: data.transformToApiRequestBody(),
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

	err := m.client.DeleteMteDomainMapping(
		client.DeleteMteDomainMappingInput{
			Domain: data.Domain.ValueString(),
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
func (m *MTEDomainMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MTEDomainMappingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	domainMapping, err := m.client.ReadMteDomainMapping(
		client.ReadMteDomainMappingInput{
			Domain: data.domainMapping.ValueString(),
		},
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

	domainMapping, err := m.client.UpdateMteDomainMapping(
		client.UpdateMteDomainMappingInput{
			Config: plan.transformToApiRequestBody(),
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

	plan.domainMapping = types.StringValue(domainMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}


func (m *MTEDomainMappingResourceModel) transformToApiRequestBody() client.MTEDomainMappingDto {
	return client.MTEDomainMappingDto{
		EnvironmentId: m.EnvironmentId.ValueString(),
		Domain:        m.Domain.ValueString(),
	}
}
