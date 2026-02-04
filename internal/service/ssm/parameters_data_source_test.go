// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ssm_test

import (
	"fmt"
	"testing"

	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSSMParametersDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "data.aws_ssm_parameters.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.SSMEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccParametersDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceName, "names.*", "/"+rName+"/param-a"),
					resource.TestCheckTypeSetElemAttr(resourceName, "names.*", "/"+rName+"/param-b"),
					resource.TestCheckTypeSetElemAttr(resourceName, "arns.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "types.#", "2"),
				),
			},
		},
	})
}

func TestAccSSMParametersDataSource_filterByTag(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "data.aws_ssm_parameters.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.SSMEndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccParametersDataSourceConfig_filterByTag(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr(resourceName, "names.*", "/"+rName+"/param-b"),
					resource.TestCheckTypeSetElemAttr(resourceName, "arns.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "types.#", "1"),
				),
			},
		},
	})
}

func testAccParametersDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_ssm_parameter" "test1" {
  name  = "/%[1]s/param-a"
  type  = "String"
  value = "TestValueA"
}

resource "aws_ssm_parameter" "test2" {
  name  = "/%[1]s/param-b"
  type  = "String"
  value = "TestValueB"
}

data "aws_ssm_parameters" "test" {
  filter {
    key    = "Path"
    values = ["/%[1]s"]
  }

  depends_on = [
    aws_ssm_parameter.test1,
    aws_ssm_parameter.test2,
  ]
}
`, rName))
}

func testAccParametersDataSourceConfig_filterByTag(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_ssm_parameter" "test1" {
  name  = "/%[1]s/param-a"
  type  = "String"
  value = "TestValueA"
  tags = {
    Environment = "test"
  }
}

resource "aws_ssm_parameter" "test2" {
  name  = "/%[1]s/param-b"
  type  = "String"
  value = "TestValueB"
  tags = {
    Environment = "production"
  }
}

data "aws_ssm_parameters" "test" {
  filter {
    key    = "tag:Environment"
    values = ["production"]
  }

  depends_on = [
    aws_ssm_parameter.test1,
    aws_ssm_parameter.test2,
  ]
}
`, rName))
}
