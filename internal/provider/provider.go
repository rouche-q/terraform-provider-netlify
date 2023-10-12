package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &netlifyProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &netlifyProvider{
			version: version,
		}
	}
}

// netlifyProvider is the provider implementation.
type netlifyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *netlifyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netlify"
	resp.Version = p.version
}

type netlifyProviderModel struct {
	Personal_token types.String `tfsdk:"personal_token"`
}

// Schema defines the provider-level schema for configuration data.
func (p *netlifyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"personal_token": schema.StringAttribute{
				Description: "Netlify personal token for the Netlify API. May aslo be provided via NETLIFY_PERSONAL_TOKEN env variable",
				Optional:    true,
			},
		}}
}

// Configure prepares a Netlify API client for data sources and resources.
func (p *netlifyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config netlifyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Personal_token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_token"),
			"Unknown Netlify API personalToken",
			"The provider cannot create the Netlify API client as there is an unknown configuration value for the HashiCups API personalToken. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the NETLIFY_PERSONAL_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	personalToken := os.Getenv("NETLIFY_PERSONAL_TOKEN")

	if !config.Personal_token.IsNull() {
		personalToken = config.Personal_token.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if personalToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("personalToken"),
			"Missing Netlify API Personal token",
			"The provider cannot create the Netlify API client as there is a missing or empty value for the HashiCups API personalToken. "+
				"Set the personalToken value in the configuration or use the NETLIFY_PERSONAL_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	netlify, err := NewNetlifyClient("https://api.netlify.com/api/v1/", personalToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Netlify API Client",
			"An unexpected error occurred when creating the Netlify API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Error: "+err.Error(),
		)
		return
	}

	// Make the Netlify client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = netlify
	resp.ResourceData = netlify
}

// DataSources defines the data sources implemented in the provider.
func (p *netlifyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSiteDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *netlifyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSiteResource,
		NewDeployKeyResource,
	}
}
