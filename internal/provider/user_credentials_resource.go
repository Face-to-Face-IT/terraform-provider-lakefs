package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Face-to-Face-IT/terraform-provider-lakefs/internal/provider/resource_user_credentials"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &UserCredentialsResource{}
var _ resource.ResourceWithImportState = &UserCredentialsResource{}

func NewUserCredentialsResource() resource.Resource {
	return &UserCredentialsResource{}
}

// UserCredentialsResource defines the resource implementation.
type UserCredentialsResource struct {
	client *LakeFSClient
}

// CredentialsResponse represents the API response for credentials
type CredentialsResponse struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key,omitempty"`
	CreationDate    int64  `json:"creation_date"`
}

func (r *UserCredentialsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_credentials"
}

func (r *UserCredentialsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_user_credentials.UserCredentialsResourceSchema(ctx)
}

func (r *UserCredentialsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserCredentialsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_user_credentials.UserCredentialsModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()

	var result CredentialsResponse
	err := client.Post(ctx, fmt.Sprintf("/auth/users/%s/credentials", userID), nil, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create credentials for user %s: %s", userID, err))
		return
	}

	data.AccessKeyId = types.StringValue(result.AccessKeyID)
	data.SecretAccessKey = types.StringValue(result.SecretAccessKey)
	data.CreationDate = types.Int64Value(result.CreationDate)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserCredentialsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_user_credentials.UserCredentialsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()
	accessKeyID := data.AccessKeyId.ValueString()

	var result CredentialsResponse
	err := client.Get(ctx, fmt.Sprintf("/auth/users/%s/credentials/%s", userID, accessKeyID), &result)
	if err != nil {
		if IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read credentials: %s", err))
		return
	}

	data.AccessKeyId = types.StringValue(result.AccessKeyID)
	data.CreationDate = types.Int64Value(result.CreationDate)
	// SecretAccessKey is not returned by the API on GET, so we keep the one from state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserCredentialsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op as all fields require replacement
}

func (r *UserCredentialsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_user_credentials.UserCredentialsModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)
	userID := data.UserId.ValueString()
	accessKeyID := data.AccessKeyId.ValueString()

	err := client.Delete(ctx, fmt.Sprintf("/auth/users/%s/credentials/%s", userID, accessKeyID))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete credentials %s for user %s: %s", accessKeyID, userID, err))
		return
	}
}

func (r *UserCredentialsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("user_id"), req, resp)
}
