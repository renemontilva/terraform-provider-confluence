package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/renemontilva/terraform-provider-confluence/internal/confluence"
)

var (
	_ datasource.DataSource              = &spaceDataSource{}
	_ datasource.DataSourceWithConfigure = &spaceDataSource{}
)

func NewSpaceDataSource() datasource.DataSource {
	return &spaceDataSource{}
}

type spaceDataSource struct {
	client *confluence.API
}

type SpaceDataSourceModel struct {
	Id     types.Int64  `tfsdk:"id"`
	Key    types.String `tfsdk:"key"`
	Name   types.String `tfsdk:"name"`
	Type   types.String `tfsdk:"type"`
	Status types.String `tfsdk:"status"`
}

// Metadata returns the data source type name.
func (d *spaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

// Schema defines the space data source schema.
func (d *spaceDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Returns a space. This includes information like the name, description, and permissions, but not the content in the space. A Confluence space typically includes a set of pages that can be organized into a hierarchical structure, with parent and child pages. The pages can contain a variety of content, including text, images, tables, lists, and more. Additionally, Confluence spaces provide features such as commenting, version history, and search capabilities for easy access to the content.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Space identifier number.",
				Computed:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: `The key of the space to be returned, which is the space name e.g: key="DEVOPS"`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the space.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the space is based on, e.g: global, personal.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Current status for the space.",
				Computed:            true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *spaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SpaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	space, err := d.client.GetSpace(data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Space Data Source Client Error", err.Error())
		return
	}
	data.Id = types.Int64Value(int64(space.Id))
	data.Key = types.StringValue(space.Key)
	data.Name = types.StringValue(space.Name)
	data.Type = types.StringValue(space.Type)
	data.Status = types.StringValue(space.Status)

	// Set State
	diags := resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *spaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*confluence.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *confluence.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}
