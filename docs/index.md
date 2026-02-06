# Kind Provider

The Kind provider is used to interact with [Kubernetes IN Docker
(kind)](https://github.com/kubernetes-sigs/kind) to provision local
[Kubernetes](https://kubernetes.io) clusters.

> **Note**
> 
> For the `runtimeConfig` field there's special behaviour for options containing a `/` character. Since this is not allowed in HCL you can just use `_` which is internally replaced with a `/` for generating the kind config. E.g. for the option `api/alpha` you'd name the field `api_alpha` and it will set it to `api/alpha` when creating the corresponding kind config.

## Example Usage

```hcl
# Configure the Kind Provider
provider "kind" {}

# Create a cluster
resource "kind_cluster" "default" {
    name           = "test-cluster"
    wait_for_ready = true
}

# Load a locally built image into the cluster
resource "kind_load" "app" {
    image        = "myapp:latest"
    cluster_name = kind_cluster.default.name
}
```
