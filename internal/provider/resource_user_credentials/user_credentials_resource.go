package resource_user_credentials

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UserCredentialsResourceSchema(ctx context.Context) schema.Schema {
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
			"access_key_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The access key ID",
				MarkdownDescription: "The access key ID",
			},
			"secret_access_key": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				Description:         "The secret access key",
				MarkdownDescription: "The secret access key",
			},
			"creation_date": schema.Int64Attribute{
				Computed:            true,
				Description:         "Unix Epoch in seconds",
				MarkdownDescription: "Unix Epoch in seconds",
			},
		},
	}
}

type UserCredentialsModel struct {
	UserId          types.String `tfsdk:"user_id"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	CreationDate    types.Int64  `tfsdk:"creation_date"`
}
