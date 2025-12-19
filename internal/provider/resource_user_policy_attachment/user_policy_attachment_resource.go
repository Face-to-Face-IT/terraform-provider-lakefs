package resource_user_policy_attachment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UserPolicyAttachmentResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the user",
				MarkdownDescription: "The ID of the user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the policy",
				MarkdownDescription: "The ID of the policy",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type UserPolicyAttachmentModel struct {
	UserId   types.String `tfsdk:"user_id"`
	PolicyId types.String `tfsdk:"policy_id"`
}
