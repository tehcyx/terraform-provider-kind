# Using The Provider

## Usage

To use the Kind provider, add it to your Terraform configuration and initialize your workspace.

## Example

```hcl
provider "kind" {}

resource "kind_cluster" "default" {
  name = "test-cluster"
}
```

## Steps

1. Add the provider and resources to your `.tf` files.
2. Run `terraform init` to initialize the provider.
3. Run `terraform plan` to review changes.
4. Run `terraform apply` to create resources.

See the [Resources](./resources/cluster.md) and [Data Sources](./resources/data_source_kind_cluster.md) docs for more examples and options.
