package resource_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GroupResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the group",
				MarkdownDescription: "The name of the group",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Description:         "A description of the group",
				MarkdownDescription: "A description of the group",
			},
			"creation_date": schema.Int64Attribute{
				Computed:            true,
				Description:         "Unix Epoch in seconds",
				MarkdownDescription: "Unix Epoch in seconds",
			},
		},
	}
}

type GroupModel struct {
	Id           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
	CreationDate types.Int64  `tfsdk:"creation_date"`
}
