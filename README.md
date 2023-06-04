# Terraform Provider for kind


## Overview

The Terraform Provider for kind enables [Terraform](https://www.terraform.io) to provision local [Kubernetes](https://kubernetes.io) clusters on base of [Kubernetes IN Docker (kind)](https://github.com/kubernetes-sigs/kind).

## Quick Starts
- [Using the provider](./docs/USAGE.md)
- [Provider development](./docs/DEVELOPMENT.md)

> **Note**
> 
> For the `runtimeConfig` field there's special behaviour for options containing a `/` character. Since this is not allowed in HCL you can just use `_` which is internally replaced with a `/` for generating the kind config. E.g. for the option `api/alpha` you'd name the field `api_alpha` and it will set it to `api/alpha` when creating the corresponding kind config.

## Example Usage

Copy the following code into a file with the extension `.tf` to create a kind cluster with only default values.
```hcl
provider "kind" {}

resource "kind_cluster" "default" {
    name = "test-cluster"
}
```

Then run `terraform init`, `terraform plan` & `terraform apply` and follow the on screen instructions. For more details on how to influence creation of the kind resource check out the Quick Start section above.
