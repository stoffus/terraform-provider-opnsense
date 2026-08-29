package interfaces

import (
	"context"
	"errors"
	"fmt"

	"github.com/browningluke/opnsense-go/pkg/api"
	"github.com/browningluke/opnsense-go/pkg/errs"
	"github.com/browningluke/opnsense-go/pkg/opnsense"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &assignmentResource{}
var _ resource.ResourceWithConfigure = &assignmentResource{}
var _ resource.ResourceWithImportState = &assignmentResource{}

func newAssignmentResource() resource.Resource {
	return &assignmentResource{}
}

type assignmentResource struct {
	client opnsense.Client
}

func (r *assignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interfaces_assignment"
}

func (r *assignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = assignmentResourceSchema()
}

func (r *assignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *opnsense.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = opnsense.NewClient(apiClient)
}

func (r *assignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *assignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignment, err := convertAssignmentSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse interface assignment, got error: %s", err))
		return
	}

	id, err := r.client.Interfaces().AddAssignment(ctx, assignment)
	if err != nil {
		if id != "" {
			data.Id = types.StringValue(id)
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		}
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to create interface assignment, got error: %s", err))
		return
	}

	data.Id = types.StringValue(id)
	tflog.Trace(ctx, "created interface assignment")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *assignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *assignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignment, err := r.client.Interfaces().GetAssignment(ctx, data.Id.ValueString())
	if err != nil {
		var notFoundError *errs.NotFoundError
		if errors.As(err, &notFoundError) {
			tflog.Warn(ctx, "interface assignment not present in remote, removing from state")
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read interface assignment, got error: %s", err))
		return
	}

	assignmentModel, err := convertAssignmentStructToSchema(assignment)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to read interface assignment, got error: %s", err))
		return
	}

	assignmentModel.Id = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &assignmentModel)...)
}

func (r *assignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *assignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	assignment, err := convertAssignmentSchemaToStruct(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to parse interface assignment, got error: %s", err))
		return
	}

	if err := r.client.Interfaces().UpdateAssignment(ctx, data.Id.ValueString(), assignment); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to update interface assignment, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *assignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *assignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Lock.ValueBool() {
		assignment, err := convertAssignmentSchemaToStruct(data)
		if err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to parse interface assignment, got error: %s", err))
			return
		}
		assignment.Lock = "0"
		if err := r.client.Interfaces().UpdateAssignment(ctx, data.Id.ValueString(), assignment); err != nil {
			resp.Diagnostics.AddError("Client Error",
				fmt.Sprintf("Unable to unlock interface assignment before deletion, got error: %s", err))
			return
		}
	}

	if err := r.client.Interfaces().DeleteAssignment(ctx, data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error",
			fmt.Sprintf("Unable to delete interface assignment, got error: %s", err))
	}
}

func (r *assignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
