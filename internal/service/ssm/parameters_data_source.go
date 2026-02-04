// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ssm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_ssm_parameters", name="Parameters")
func dataSourceParameters() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceParametersRead,

		Schema: map[string]*schema.Schema{
			names.AttrARNs: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"filter": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						names.AttrName: {
							Type:     schema.TypeString,
							Required: true,
						},
						"option": {
							Type:     schema.TypeString,
							Optional: true,
						},
						names.AttrValues: {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			names.AttrNames: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrType: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrValues: {
				Type:      schema.TypeList,
				Computed:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceParametersRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).SSMClient(ctx)

	input := &ssm.DescribeParametersInput{
		ParameterFilters: expandParameterStringFilters(d.Get("filter").(*schema.Set)),
	}

	var output []awstypes.ParameterMetadata

	pages := ssm.NewDescribeParametersPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading SSM Parameters: %s", err)
		}

		output = append(output, page.Parameters...)
	}

	d.SetId(meta.(*conns.AWSClient).Partition(ctx))
	d.Set(names.AttrARNs, tfslices.ApplyToAll(output, func(v awstypes.ParameterMetadata) string {
		return aws.ToString(v.ARN)
	}))
	d.Set(names.AttrNames, tfslices.ApplyToAll(output, func(v awstypes.ParameterMetadata) string {
		return aws.ToString(v.Name)
	}))
	d.Set(names.AttrType, tfslices.ApplyToAll(output, func(v awstypes.ParameterMetadata) string {
		return string(v.Type)
	}))

	// Get parameter values using GetParameter for each parameter
	values := make([]string, len(output))
	for i, param := range output {
		value, err := conn.GetParameter(ctx, &ssm.GetParameterInput{
			Name:           param.Name,
			WithDecryption: aws.Bool(true),
		})

		if err != nil {
			// If we can't get the parameter value, set it to empty string
			values[i] = ""
		} else {
			values[i] = aws.ToString(value.Parameter.Value)
		}
	}
	d.Set(names.AttrValues, values)

	return diags
}

func expandParameterStringFilters(filters *schema.Set) []awstypes.ParameterStringFilter {
	result := make([]awstypes.ParameterStringFilter, 0, filters.Len())

	for _, v := range filters.List() {
		filter := v.(map[string]any)

		parameterFilter := awstypes.ParameterStringFilter{
			Key:    aws.String(filter[names.AttrName].(string)),
			Values: expandStringValueSet(filter[names.AttrValues].(*schema.Set)),
		}

		if option, ok := filter["option"].(string); ok && option != "" {
			parameterFilter.Option = aws.String(option)
		}

		result = append(result, parameterFilter)
	}

	return result
}

func expandStringValueSet(s *schema.Set) []string {
	result := make([]string, 0, s.Len())

	for _, v := range s.List() {
		result = append(result, v.(string))
	}

	return result
}
