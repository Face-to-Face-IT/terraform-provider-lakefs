package resource_user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UserResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The username of the user",
				MarkdownDescription: "The username of the user",
			},
			"creation_date": schema.Int64Attribute{
				Computed:            true,
				Description:         "Unix Epoch in seconds",
				MarkdownDescription: "Unix Epoch in seconds",
			},
			"friendly_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Friendly name of the user",
				MarkdownDescription: "Friendly name of the user",
			},
		},
	}
}

type UserModel struct {
	Id           types.String `tfsdk:"id"`
	CreationDate types.Int64  `tfsdk:"creation_date"`
	FriendlyName types.String `tfsdk:"friendly_name"`
}
