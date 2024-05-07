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
var _ resource.Resource = &MTERulesMappingResource{}
var _ resource.ResourceWithImportState = &MTERulesMappingResource{}

func NewMTERulesMappingResource() resource.Resource {
	return &MTERulesMappingResource{}
}

type MTERulesMappingResource struct {
	client *client.Client
}

type MTERulesMappingResourceModel struct {
	RulesId types.String `tfsdk:"rules_id"`
	Domain  types.String `tfsdk:"domain"`
}

// Metadata implements resource.Resource.
func (m *MTERulesMappingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mte_rules_mapping"
}

func (m *MTERulesMappingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (m *MTERulesMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A mapping layer designed to map a domain, either a custom domain or standard domain, to a rule group ID. " +
			"This rule group ID should be created asynchronously through the Platform API. The rules are designed to handle redirects and rewrites.",

		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The domain on which you want to activate rules upon.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rules_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The rule group ID the domain should be associated with.",
			},
		},
	}
}

// ImportState implements resource.ResourceWithImportState.
func (m *MTERulesMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Create implements resource.Resource.
func (m *MTERulesMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MTERulesMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.CreateMteRulesMapping(
		client.CreateMteRulesMappingInput{
			Config: data.transformToApiRequestBody(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create MTE Rules mapping",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete implements resource.Resource.
func (m *MTERulesMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MTERulesMappingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.DeleteMteRulesMapping(
		client.DeleteMteRulesMappingInput{
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
func (m *MTERulesMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MTERulesMappingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	rulesId, err := m.client.ReadMteRulesMapping(
		client.ReadMteRulesMappingInput{
			Domain: data.Domain.ValueString(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError("Failed to get MTE Rules Mapping", err.Error())
		return
	}
	data.RulesId = types.StringValue(rulesId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update implements resource.Resource.
func (m *MTERulesMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MTERulesMappingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := m.client.UpdateMteRulesMapping(
		client.UpdateMteRulesMappingInput{
			Config: plan.transformToApiRequestBody(),
		},
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create MTE Rules Mapping",
			"An error occurred while executing the creation. "+
				"If unexpected, please report this issue to the provider developers.\n\n"+
				"JSON Error: "+err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (m *MTERulesMappingResourceModel) transformToApiRequestBody() client.MTERulesMappingDto {
	return client.MTERulesMappingDto{
		RulesId: m.RulesId.ValueString(),
		Domain:  m.Domain.ValueString(),
	}
}
