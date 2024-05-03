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
	client *client.Client
}

type MTEDomainMappingResourceModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Domain        types.String `tfsdk:"domain"`
	DomainMapping types.String `tfsdk:"domain_mapping"`
}

// Metadata implements resource.Resource.
func (m *MTEDomainMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mte_domain_mapping"
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
		MarkdownDescription: "A mapping layer designed to map a domain, either a custom domain or standard domain, to an environment. "+
		"This environment can then be associated with a [config resource](https://registry.terraform.io/providers/THG-Headless/altitude/latest/docs/resources/mte_config).",

		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The domain relating to the environment on which you are deploying.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The environment which relates with the [config resource](https://registry.terraform.io/providers/THG-Headless/altitude/latest/docs/resources/mte_config).",
			},
			"domain_mapping": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The computed value stored as the mapper between domain and config.",
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTEDomainMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
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
			"Failed to create MTE domain mapping",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}

	data.DomainMapping = types.StringValue(domainMapping)
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
			Domain: data.Domain.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to get MTE Domain Mapping", err.Error())
		return
	}

	data.DomainMapping = types.StringValue(domainMapping)

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
			"Failed to create MTE Domain Mapping",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}

	plan.DomainMapping = types.StringValue(domainMapping)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (m *MTEDomainMappingResourceModel) transformToApiRequestBody() client.MTEDomainMappingDto {
	return client.MTEDomainMappingDto{
		EnvironmentId: m.EnvironmentId.ValueString(),
		Domain:        m.Domain.ValueString(),
	}
}
