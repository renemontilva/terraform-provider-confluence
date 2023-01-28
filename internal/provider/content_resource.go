package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/renemontilva/terraform-provider-confluence/internal/confluence"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &ContentResource{}
	_ resource.ResourceWithConfigure   = &ContentResource{}
	_ resource.ResourceWithImportState = &ContentResource{}
)

func NewContentResource() resource.Resource {
	return &ContentResource{}
}

// ContentResource defines the resource implementation.
type ContentResource struct {
	client *confluence.API
}

// ContentResourceModel describes the resource data model.
type ContentResourceModel struct {
	Id    types.String `tfsdk:"id"`
	Type  types.String `tfsdk:"type"`
	Title types.String `tfsdk:"title"`
	Space types.String `tfsdk:"space"`
	Body  types.String `tfsdk:"body"`
}

func (r *ContentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_content"
}

func (r *ContentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "The resource ```content``` creates a new piece of content.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Content identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the new content. Custom content types defined by apps are also supported. eg. 'page', 'blogpost', 'comment' etc.",
				Required:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Defines the document title.",
				Required:            true,
			},
			"space": schema.StringAttribute{
				MarkdownDescription: "The space that the content is being created in.",
				Required:            true,
			},
			"body": schema.StringAttribute{
				MarkdownDescription: "The body of the new content.",
				Required:            true,
			},
		},
	}
}

func (r *ContentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *ContentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ContentResourceModel
	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create confluence content struct
	space := confluence.Space{
		Key: data.Space.ValueString(),
	}
	body := confluence.Body{
		Storage: confluence.Storage{
			Value:          data.Body.ValueString(),
			Representation: "storage",
		},
	}

	content := confluence.Content{
		Type:  data.Type.ValueString(),
		Title: data.Title.ValueString(),
		Space: &space,
		Body:  body,
		Version: &confluence.Version{
			Number: int(0),
		},
	}
	err := r.client.CreateContent(&content)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create confluence content, got error: %s", err))
		return
	}

	data.Id = types.StringValue(content.Id)
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ContentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ContentResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	content, err := r.client.GetContentById(data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read content, got error: %s", err))
		return
	}

	if content.Id == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("content id is empty, got error: %s", fmt.Errorf("content object: %v", content)))
		return
	}
	// Overwrite ContentResourceModel with values returned from confluence API
	data.Id = types.StringValue(content.Id)
	data.Type = types.StringValue(content.Type)
	data.Title = types.StringValue(content.Title)
	data.Space = types.StringValue(content.Space.Key)
	// Confluence reponse does not return the body section
	//data.Body = types.StringValue(content.Body.Storage.Value)

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ContentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Provider client data and make a call using it.
	// Create confluence content struct
	space := confluence.Space{
		Key: data.Space.ValueString(),
	}
	body := confluence.Body{
		Storage: confluence.Storage{
			Value:          data.Body.ValueString(),
			Representation: "storage",
		},
	}
	version := confluence.Version{}
	content := confluence.Content{
		Id:      data.Id.ValueString(),
		Type:    data.Type.ValueString(),
		Title:   data.Title.ValueString(),
		Space:   &space,
		Body:    body,
		Version: &version,
	}
	err := r.client.UpdateContent(&content)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update content from confluence API, got error: %s", err))
		return
	}

	// Update state content
	// Overwrite ContentResourceModel with values returned from confluence API
	data.Id = types.StringValue(content.Id)
	data.Type = types.StringValue(content.Type)
	data.Title = types.StringValue(content.Title)
	data.Space = types.StringValue(content.Space.Key)
	data.Body = types.StringValue(content.Body.Storage.Value)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ContentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ContentResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	err := r.client.DeleteContent(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete content from confluence API, got error: %s", err))
		return
	}
}

func (r *ContentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
