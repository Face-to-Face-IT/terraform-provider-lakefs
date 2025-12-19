package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Face-to-Face-IT/terraform-provider-lakefs/internal/provider/resource_policy"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PolicyResource{}
var _ resource.ResourceWithImportState = &PolicyResource{}

func NewPolicyResource() resource.Resource {
	return &PolicyResource{}
}

// PolicyResource defines the resource implementation.
type PolicyResource struct {
	client *LakeFSClient
}

// PolicyCreateRequest represents the request to create a policy
type PolicyCreateRequest struct {
	ID        string          `json:"id"`
	Statement json.RawMessage `json:"statement"`
}

// PolicyResponse represents the API response for a policy
type PolicyResponse struct {
	ID           string          `json:"id"`
	CreationDate int64           `json:"creation_date"`
	Statement    json.RawMessage `json:"statement"`
}

// jsonEqual compares two JSON strings for semantic equality (ignoring key order)
func jsonEqual(a, b string) bool {
	var objA, objB interface{}
	if err := json.Unmarshal([]byte(a), &objA); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &objB); err != nil {
		return false
	}
	return reflect.DeepEqual(objA, objB)
}

func (r *PolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *PolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_policy.PolicyResourceSchema(ctx)
}

func (r *PolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_policy.PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)

	createReq := PolicyCreateRequest{
		ID:        data.Id.ValueString(),
		Statement: json.RawMessage(data.Statement.ValueString()),
	}

	var result PolicyResponse
	err := client.Post(ctx, "/auth/policies", createReq, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create policy: %s", err))
		return
	}

	data.Id = types.StringValue(result.ID)
	data.CreationDate = types.Int64Value(result.CreationDate)

	// Preserve the original statement format if semantically equal to avoid diffs from key ordering
	resultStatement := string(result.Statement)
	originalStatement := data.Statement.ValueString()
	if jsonEqual(originalStatement, resultStatement) {
		// Keep the original format
		data.Statement = types.StringValue(originalStatement)
	} else {
		// Use the API's response (this shouldn't normally happen)
		data.Statement = types.StringValue(resultStatement)
	}

	tflog.Trace(ctx, "created a policy resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_policy.PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)

	var result PolicyResponse
	err := client.Get(ctx, fmt.Sprintf("/auth/policies/%s", data.Id.ValueString()), &result)
	if err != nil {
		if IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read policy: %s", err))
		return
	}

	data.Id = types.StringValue(result.ID)
	data.CreationDate = types.Int64Value(result.CreationDate)

	// Preserve the original statement format if semantically equal to avoid diffs from key ordering
	resultStatement := string(result.Statement)
	originalStatement := data.Statement.ValueString()
	if jsonEqual(originalStatement, resultStatement) {
		// Keep the original format
		data.Statement = types.StringValue(originalStatement)
	} else {
		// Use the API's response
		data.Statement = types.StringValue(resultStatement)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_policy.PolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)

	updateReq := PolicyCreateRequest{
		ID:        data.Id.ValueString(),
		Statement: json.RawMessage(data.Statement.ValueString()),
	}

	var result PolicyResponse
	err := client.Put(ctx, fmt.Sprintf("/auth/policies/%s", data.Id.ValueString()), updateReq, &result)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update policy: %s", err))
		return
	}

	data.Id = types.StringValue(result.ID)
	data.CreationDate = types.Int64Value(result.CreationDate)

	statementBytes, err := json.Marshal(result.Statement)
	if err != nil {
		resp.Diagnostics.AddError("Marshal Error", fmt.Sprintf("Unable to marshal policy statement: %s", err))
		return
	}
	data.Statement = types.StringValue(string(statementBytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_policy.PolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := NewAPIClient(r.client)

	err := client.Delete(ctx, fmt.Sprintf("/auth/policies/%s", data.Id.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete policy: %s", err))
		return
	}
}

func (r *PolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
