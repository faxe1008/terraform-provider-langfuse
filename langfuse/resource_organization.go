package langfuse

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/faxe1008/terraform-provider-langfuse/client"
)

// organizationResource implements the langfuse_organization resource.
type organizationResource struct {
	client *client.Client
}

// NewOrganizationResource returns a new organizationResource.
func NewOrganizationResource() resource.Resource {
	return &organizationResource{}
}

// Metadata sets the resource type name.
func (r *organizationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "langfuse_organization"
}

// Schema defines the schema for organizations.
func (r *organizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Langfuse organizations (self-hosted).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the organization.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the organization.",
			},
		},
	}
}

// organizationResourceModel maps schema attributes to Go types.
type organizationResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Configure injects the Langfuse client from the provider.
func (r *organizationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientData, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData),
		)
		return
	}
	r.client = clientData
}

// Create creates a new organization via the API.
func (r *organizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan organizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to create organization
	org, err := r.client.CreateOrganization(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating organization", err.Error())
		return
	}

	// Set state with returned values
	plan.ID = types.StringValue(org.ID)
	plan.Name = types.StringValue(org.Name)
	resp.State.Set(ctx, &plan)
}

// Read refreshes the state by reading from the API.
func (r *organizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state organizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.GetOrganization(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading organization", err.Error())
		return
	}

	// Update state
	state.Name = types.StringValue(org.Name)
	resp.State.Set(ctx, &state)
}

// Update renames the organization via the API.
func (r *organizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan organizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateOrganization(ctx, plan.ID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating organization", err.Error())
		return
	}

	// Use plan values as new state
	resp.State.Set(ctx, &plan)
}

// Delete removes the organization via the API.
func (r *organizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state organizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteOrganization(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting organization", err.Error())
	}
}

// ImportState allows importing an existing organization by ID.
func (r *organizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// The import identifier is the organization ID.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(req.ID))...)
}
