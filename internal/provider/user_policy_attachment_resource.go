package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/Face-to-Face-IT/terraform-provider-lakefs/internal/provider/resource_user_policy_attachment"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserPolicyAttachmentResource{}
var _ resource.ResourceWithImportState = &UserPolicyAttachmentResource{}

func NewUserPolicyAttachmentResource() resource.Resource {
	return &UserPolicyAttachmentResource{}
}

// UserPolicyAttachmentResource defines the resource implementation.
type UserPolicyAttachmentResource struct {
	client *LakeFSClient
}

func (r *UserPolicyAttachmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_policy_attachment"
}

func (r *UserPolicyAttachmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_user_policy_attachment.UserPolicyAttachmentResourceSchema(ctx)
}

func (r *UserPolicyAttachmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserPolicyAttachmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_user_policy_attachment.UserPolicyAttachmentModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()
	policyID := data.PolicyId.ValueString()

	err := client.Put(ctx, fmt.Sprintf("/auth/users/%s/policies/%s", userID, policyID), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to attach policy %s to user %s: %s", policyID, userID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserPolicyAttachmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_user_policy_attachment.UserPolicyAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()
	policyID := data.PolicyId.ValueString()

	// List all policies attached to the user and check if the policy is attached
	var result struct {
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}
	err := client.Get(ctx, fmt.Sprintf("/auth/users/%s/policies", userID), &result)
	if err != nil {
		if IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read user policy attachment: %s", err))
		return
	}

	// Check if the policy is in the list
	found := false
	for _, policy := range result.Results {
		if policy.ID == policyID {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserPolicyAttachmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op as all fields require replacement
}

func (r *UserPolicyAttachmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_user_policy_attachment.UserPolicyAttachmentModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()
	policyID := data.PolicyId.ValueString()

	err := client.Delete(ctx, fmt.Sprintf("/auth/users/%s/policies/%s", userID, policyID))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to detach policy %s from user %s: %s", policyID, userID, err))
		return
	}
}

func (r *UserPolicyAttachmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}
