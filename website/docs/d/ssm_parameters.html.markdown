---
subcategory: "SSM (Systems Manager)"
layout: "aws"
page_title: "AWS: aws_ssm_parameters"
description: |-
  Provides a list of Systems Manager parameters.
---

# Data Source: aws_ssm_parameters

Use this data source to get a list of System Manager parameters based on specified filters.

## Example Usage

### Filter by Path

```terraform
data "aws_ssm_parameters" "example" {
  filter {
    name   = "Path"
    values = ["/myapp/config"]
  }
}
```

### Filter by Tag

```terraform
data "aws_ssm_parameters" "prod" {
  filter {
    name   = "tag:Environment"
    values = ["production"]
  }
}
```

### Filter by Multiple Criteria

```terraform
data "aws_ssm_parameters" "example" {
  filter {
    name   = "Path"
    values = ["/myapp/config"]
  }
  filter {
    name   = "Type"
    values = ["SecureString"]
  }
}
```

~> **Note:** When the data source retrieves `SecureString` parameters, the unencrypted values are stored in the raw state as plain-text as per normal Terraform behavior. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Argument Reference

This data source supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `filter` - (Required) One or more configuration block to filter the list of parameters. See [Filter](#filter) below. The number of filters must be at least one.

### Filter

A `filter` block supports the following arguments:

* `name` - (Required) The name of the filter. Valid values are defined in the [SSM DescribeParameters API](https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_DescribeParameters.html).
* `option` - (Optional) The filter option. Valid values are `Equals` and `BeginsWith`.
* `values` - (Required) A set of values for the filter.

#### Valid Filter Names

* `tag:TagKey` - Filter by tag key (e.g., `tag:Environment`)
* `Name` - Parameter name
* `Path` - Parameter path
* `Type` - Parameter type (`String`, `StringList`, `SecureString`)
* `KeyId` - KMS key ID used to encrypt the parameter
* `DataType` - Parameter data type
* `Tier` - Parameter tier (`Standard`, `Advanced`, `Intelligent-Tiering`)

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arns` - A list that contains the Amazon Resource Names (ARNs) of the retrieved parameters.
* `names` - A list that contains the names of the retrieved parameters.
* `types` - A list that contains the types of the retrieved parameters.
* `values` - A list that contains the retrieved parameter values. **Note:** This value is always marked as sensitive in the Terraform plan output, regardless of whether any retrieved parameters are of `SecureString` type. Use the [`nonsensitive` function](https://developer.hashicorp.com/terraform/language/functions/nonsensitive) to override the behavior at your own risk and discretion, if you are certain that there are no sensitive values being retrieved.
