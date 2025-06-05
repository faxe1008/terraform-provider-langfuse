package langfuse

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/faxe1008/terraform-provider-langfuse/client"
)

// projectResource implements the langfuse_project resource.
type projectResource struct {
	client *client.Client
}

// NewProjectResource returns a new projectResource.
func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// Metadata sets the resource type name.
func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "langfuse_project"
}

// Schema defines the schema for projects.
func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource for managing Langfuse projects (within an organization).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the project.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the project.",
			},
			"organization_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the parent organization.",
			},
			"public_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Public API key for this project (returned on create).",
			},
			"secret_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Secret API key for this project (returned on create).",
			},
		},
	}
}

// projectResourceModel maps project resource schema.
type projectResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	OrganizationID types.String `tfsdk:"organization_id"`
	PublicKey      types.String `tfsdk:"public_key"`
	SecretKey      types.String `tfsdk:"secret_key"`
}

// Configure injects the Langfuse client.
func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create calls the API to create a new project.
func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	proj, err := r.client.CreateProject(ctx, plan.OrganizationID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	plan.ID = types.StringValue(proj.ID)
	plan.Name = types.StringValue(proj.Name)
	plan.OrganizationID = types.StringValue(proj.OrganizationID)
	plan.PublicKey = types.StringValue(proj.PublicKey)
	plan.SecretKey = types.StringValue(proj.SecretKey)

	resp.State.Set(ctx, &plan)
}

// Read refreshes the project state from the API.
func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	proj, err := r.client.GetProject(ctx, state.OrganizationID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	state.Name = types.StringValue(proj.Name)
	state.PublicKey = types.StringValue(proj.PublicKey)
	// Note: SecretKey is not returned by GET; keep the previous state value intact.

	resp.State.Set(ctx, &state)
}

// Update renames the project via the API.
func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateProject(ctx, plan.OrganizationID.ValueString(), plan.ID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	resp.State.Set(ctx, &plan)
}

// Delete removes the project via the API.
func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteProject(ctx, state.OrganizationID.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
	}
}

// ImportState allows importing an existing project by “orgID/projectID” composite ID.
func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expect the ID in the form "orgID/projectID"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import identifier",
			"Expected import ID in the form \"<organization_id>/<project_id>\" (e.g. \"org123/proj456\").",
		)
		return
	}

	orgID := parts[0]
	projID := parts[1]

	// Set both organization_id and id in the Terraform state
	resp.Diagnostics.Append(
		resp.State.SetAttribute(ctx, path.Root("organization_id"), types.StringValue(orgID)),
		resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(projID)),
	)... 

	// After setting those two, Terraform will call Read() automatically to populate the rest.
}
