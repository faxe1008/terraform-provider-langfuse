package langfuse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/faxe1008/terraform-provider-langfuse/client"
)

// NewProvider returns a new Langfuse provider instance.
func NewProvider(version string) provider.Provider {
	return &LangfuseProvider{version: version}
}

// LangfuseProvider implements provider.Provider.
type LangfuseProvider struct {
	version string
}

// Metadata returns the provider type name and version.
func (p *LangfuseProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "langfuse"
	resp.Version = p.version
}

// Schema defines provider-level configuration (admin_api_key, base_url).
func (p *LangfuseProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"admin_api_key": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "Langfuse **Admin API Key** (for self-hosted instances; used as a Bearer token).",
			},
			"base_url": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base URL of the Langfuse API (e.g. `http://localhost:3000`). Defaults to `http://localhost:3000`.",
			},
		},
	}
}

// providerConfig holds the configuration data.
type providerConfig struct {
	AdminAPIKey types.String `tfsdk:"admin_api_key"`
	BaseURL     types.String `tfsdk:"base_url"`
}

// Configure initializes the Langfuse API client using the provider config.
func (p *LangfuseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if config.AdminAPIKey.IsUnknown() || config.AdminAPIKey.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Admin API key",
			"The provider requires `admin_api_key` to be configured.",
		)
		return
	}

	// Default base_url if not set
	baseURL := "http://localhost:3000"
	if !config.BaseURL.IsNull() && !config.BaseURL.IsUnknown() {
		baseURL = config.BaseURL.ValueString()
	}

	// Create the Langfuse API client with the provided settings.
	c := client.NewClient(baseURL, config.AdminAPIKey.ValueString())

	// Pass the client to all resources and data sources
	resp.ResourceData = c
	resp.DataSourceData = c
}

// Resources returns a list of resource constructors.
func (p *LangfuseProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewProjectResource,
	}
}

// DataSources returns a list of data source constructors (none in this provider).
func (p *LangfuseProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
