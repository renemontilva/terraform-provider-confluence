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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-confluence/internal/confluence"
)

// Ensure ConfluenceProvider satisfies various provider interfaces.
var _ provider.Provider = &ConfluenceProvider{}

func New() provider.Provider {
	return &ConfluenceProvider{}
}

// ConfluenceProvider defines the provider implementation.
type ConfluenceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
}

// ConfluenceProviderModel describes the provider data model.
type ConfluenceProviderModel struct {
	Host  types.String `tfsdk:"host"`
	User  types.String `tfsdk:"user"`
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider type name.
func (p *ConfluenceProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "confluence"
}

// Schema defines the provider-level schema for configuration data.
func (p *ConfluenceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Confluence provider interacts with atlassian confluence cloud service.
		You must configured the provider with the proper credentials before you can use it.`,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "Confluence's service hostname",
				Optional:            true,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "Confluence's service username",
				Optional:            true,
			},
			"token": schema.StringAttribute{
				MarkdownDescription: "Confluence's username token",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

// Configure prepares a Confluence API client for data and resources.
func (p *ConfluenceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Confluence API client")

	// Retrieve provider data from configuration
	var config ConfluenceProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown confluence hostname value",
			"The provider cannot create the Confluence API client as there is an unknown configuration value for the Confluence API host",
		)
	}

	if config.User.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Unknown confluence user value",
			"The provider cannot create the Confluence API client as there is an unknown configuration value for the Confluence API username",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown confluence token value",
			"The provider cannot create the Confluence API client as there is an unknown configuration value for the Confluence API token",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	host := os.Getenv("CONFLUENCE_HOST")
	user := os.Getenv("CONFLUENCE_USER")
	token := os.Getenv("CONFLUENCE_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	if !config.User.IsNull() {
		user = config.User.ValueString()
	}
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Confluence Host",
			"The provider cannot create the Confluence API client as there is a missing or empty value for the Confluence API host. "+
				"Set the host value in the configuration or use the CONFLUENCE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if user == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("user"),
			"Missing Confluence User",
			"The provider cannot create the Confluence API client as there is a missing or empty value for the Confluence API user. "+
				"Set the host value in the configuration or use the CONFLUENCE_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Confluence Token",
			"The provider cannot create the Confluence API client as there is a missing or empty value for the Confluence API user. "+
				"Set the host value in the configuration or use the CONFLUENCE_USER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "confluence_host", host)
	ctx = tflog.SetField(ctx, "confluence_user", user)
	ctx = tflog.SetField(ctx, "confluence_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "confluence_token")

	tflog.Debug(ctx, "Creating Confluence API")

	// Creates client configuration for data sources and resources.
	client, err := confluence.NewAPI(user, token, host)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Confluence API",
			err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
	tflog.Info(ctx, "Configured Confluence Client", map[string]any{
		"success": true,
	})
}

func (p *ConfluenceProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewContentResource,
		NewSpaceResource,
	}
}

func (p *ConfluenceProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSpaceDataSource,
	}
}
