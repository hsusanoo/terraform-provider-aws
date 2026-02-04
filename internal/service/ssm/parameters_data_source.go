// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	fwtypes "github.com/hashicorp/terraform-provider-aws/internal/framework/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource("aws_ssm_parameters", name="Parameters")
func newParametersDataSource(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &parametersDataSource{}, nil
}

const (
	DSNameParameters = "Parameters Data Source"
)

type parametersDataSource struct {
	framework.DataSourceWithModel[parametersDataSourceModel]
}

func (d *parametersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			names.AttrARNs: schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			names.AttrNames: schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			names.AttrType: schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
		Blocks: map[string]schema.Block{
			"filter": schema.SetNestedBlock{
				CustomType: fwtypes.NewSetNestedObjectTypeOf[parametersFilterModel](ctx),
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required: true,
						},
						"option": schema.StringAttribute{
							Optional: true,
						},
						names.AttrValues: schema.SetAttribute{
							CustomType: fwtypes.SetOfStringType,
							Required:   true,
						},
					},
				},
			},
		},
	}
}

func (d *parametersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().SSMClient(ctx)

	var data parametersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &ssm.DescribeParametersInput{}

	resp.Diagnostics.Append(flex.Expand(ctx, data.Filter, &input.ParameterFilters)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := findParameters(ctx, conn, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"reading SSM Parameters",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(flex.Flatten(ctx, output, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func findParameters(ctx context.Context, conn *ssm.Client, input *ssm.DescribeParametersInput) ([]awstypes.ParameterMetadata, error) {
	var output []awstypes.ParameterMetadata

	pages := ssm.NewDescribeParametersPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		output = append(output, page.Parameters...)
	}

	return output, nil
}

type parametersDataSourceModel struct {
	framework.WithRegionModel
	ARNs   fwtypes.ListValueOf[types.String]                     `tfsdk:"arns"`
	Filter fwtypes.SetNestedObjectValueOf[parametersFilterModel] `tfsdk:"filter"`
	Names  fwtypes.ListValueOf[types.String]                     `tfsdk:"names"`
	Type   fwtypes.ListValueOf[types.String]                     `tfsdk:"type"`
}

type parametersFilterModel struct {
	Key    types.String        `tfsdk:"key"`
	Option types.String        `tfsdk:"option"`
	Values fwtypes.SetOfString `tfsdk:"values"`
}
