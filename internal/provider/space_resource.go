package provider

import (
	"fmt"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-confluence/internal/confluence"
)

var (
	_ resource.Resource                = &SpaceResource{}
	_ resource.ResourceWithConfigure   = &SpaceResource{}
	_ resource.ResourceWithImportState = &SpaceResource{}
)

// Ensure provider defined types fully satisfy framework interfaces
func NewSpaceResource() resource.Resource {
	return &SpaceResource{}
}

type SpaceResource struct {
	client *confluence.API
}

type SpaceResourceModel struct {
	Id          types.Int64  `tfsdk:"id"`
	Key         types.String `tfsdk:"key"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func (r *SpaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

func (r *SpaceResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates spaces, space is a container for organizing and grouping related pages of content.
		Spaces can be used to separate content by project, team, department, or other criteria.
		
		A Confluence space typically includes a set of pages that can be organized into a hierarchical structure, with parent and child pages. `,
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
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the new/updated space.",
				Optional:            true,
			},
		},
	}
}

func (r *SpaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*confluence.API)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *confluence.API, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *SpaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SpaceResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Calls GetSpace Confluece API Client method.
	space, err := r.client.GetSpace(data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read space, got error: %s", err))
	}
	if space == nil {
		resp.Diagnostics.AddError("Read Space error", fmt.Sprintf("Getspace returns a nil space object: %v", space))
		return

	}
	// Overwrite SpaceResourceModel with values returned from confluence API.
	data.Id = types.Int64Value(int64(space.Id))
	data.Key = types.StringValue(space.Key)
	data.Name = types.StringValue(space.Name)
	// FIXME: due to some issues with expand option is now commented
	//data.Description = types.StringValue(space.Description.Plain.Value)

	// Write trace log
	tflog.Trace(ctx, "Read a space")
	// Set refreshed terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SpaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SpaceResourceModel
	// Reads Terrafrom plan into SpaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	//
	if resp.Diagnostics.HasError() {
		return
	}

	//Create Space confluence object

	space := confluence.Space{
		Key:  data.Key.ValueString(),
		Name: data.Key.ValueString(),
		Description: &confluence.SpaceDescription{
			Plain: confluence.Plain{
				Value: data.Description.ValueString(),
			},
		},
	}

	err := r.client.CreateSpace(&space)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create space, got error: %s", err))
	}

	data.Id = types.Int64Value(int64(space.Id))
	// Write a trace log
	tflog.Trace(ctx, "Created a space")
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *SpaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SpaceResourceModel
	// Read Terraform plan and set into data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	//Create a space object
	spacedescription := confluence.SpaceDescription{
		Plain: confluence.Plain{
			Value: data.Description.ValueString(),
		},
	}
	space := confluence.Space{
		Id:          uint(data.Id.ValueInt64()),
		Key:         data.Key.ValueString(),
		Name:        data.Name.ValueString(),
		Description: &spacedescription,
	}
	err := r.client.UpdateSpace(&space)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update space, got error %v", err))
		return
	}

	// Update resource state with updated space object above
	data.Id = types.Int64Value(int64(space.Id))
	data.Key = types.StringValue(space.Key)
	data.Name = types.StringValue(space.Name)
	data.Description = types.StringValue(space.Description.Plain.Value)

	// Write a trace log
	tflog.Trace(ctx, "Updated a space")

	// Save to terraform state file
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SpaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SpaceResourceModel
	// Retreive values from terraform state
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete space from API
	err := r.client.DeleteSpace(data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete space, got error %v", err))
		return
	}
	// Write a trace log
	tflog.Trace(ctx, "Delete a space")
}

func (r *SpaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("key"), req, resp)
}
