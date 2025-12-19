package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/Face-to-Face-IT/terraform-provider-lakefs/internal/provider/resource_group_policy_attachment"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GroupPolicyAttachmentResource{}
var _ resource.ResourceWithImportState = &GroupPolicyAttachmentResource{}

func NewGroupPolicyAttachmentResource() resource.Resource {
	return &GroupPolicyAttachmentResource{}
}

// GroupPolicyAttachmentResource defines the resource implementation.
type GroupPolicyAttachmentResource struct {
	client *LakeFSClient
}

func (r *GroupPolicyAttachmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_policy_attachment"
}

func (r *GroupPolicyAttachmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_group_policy_attachment.GroupPolicyAttachmentResourceSchema(ctx)
}

func (r *GroupPolicyAttachmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*LakeFSClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *LakeFSClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *GroupPolicyAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_group_policy_attachment.GroupPolicyAttachmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	groupID := data.GroupId.ValueString()
	policyID := data.PolicyId.ValueString()

	err := client.Put(ctx, fmt.Sprintf("/auth/groups/%s/policies/%s", groupID, policyID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach policy %s to group %s: %s", policyID, groupID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupPolicyAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_group_policy_attachment.GroupPolicyAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	groupID := data.GroupId.ValueString()
	policyID := data.PolicyId.ValueString()

	// Check if policy is still attached
	var result PolicyResponse
	err := client.Get(ctx, fmt.Sprintf("/auth/groups/%s/policies/%s", groupID, policyID), &result)
	if err != nil {
		if IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read group policy attachment: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupPolicyAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op as all fields require replacement
}

func (r *GroupPolicyAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_group_policy_attachment.GroupPolicyAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	groupID := data.GroupId.ValueString()
	policyID := data.PolicyId.ValueString()

	err := client.Delete(ctx, fmt.Sprintf("/auth/groups/%s/policies/%s", groupID, policyID))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach policy %s from group %s: %s", policyID, groupID, err))
		return
	}
}

func (r *GroupPolicyAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("group_id"), req, resp)
}
