# Kind Terraform Provider

The Kind Terraform provider allows you to manage [Kubernetes IN Docker (kind)](https://github.com/kubernetes-sigs/kind) clusters for local development and testing.

## Table of Contents

- [Requirements](#requirements)
- [Usage](#usage)
- [Resources](#resources)
- [Data Sources](#data-sources)
- [FAQ](FAQ.md)
- [Development](DEVELOPMENT.md)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.19 or higher
- Docker Engine with sufficient resources for multi-node clusters

## Usage

See [USAGE.md](USAGE.md) for a quick start and example workflows.

## Resources

- [`kind_cluster`](resources/cluster.md): Manages a kind cluster.

## Data Sources

- [`kind_cluster`](resources/data_source_kind_cluster.md): Retrieves information about an existing kind cluster by name.

## Example Usage

```hcl
# Configure the Kind Provider
provider "kind" {}

# Create a cluster
resource "kind_cluster" "default" {
    name = "test-cluster"
}
```
