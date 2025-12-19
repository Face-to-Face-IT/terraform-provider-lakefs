package resource_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PolicyResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the policy",
				MarkdownDescription: "The name of the policy",
			},
			"statement": schema.StringAttribute{
				Required:            true,
				Description:         "A JSON string defining actions, resources, and effect",
				MarkdownDescription: "A JSON string defining actions, resources, and effect",
			},
			"creation_date": schema.Int64Attribute{
				Computed:            true,
				Description:         "Unix Epoch in seconds",
				MarkdownDescription: "Unix Epoch in seconds",
			},
		},
	}
}

type PolicyModel struct {
	Id           types.String `tfsdk:"id"`
	Statement    types.String `tfsdk:"statement"`
	CreationDate types.Int64  `tfsdk:"creation_date"`
}
