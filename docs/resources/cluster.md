# kind_cluster Resource

Manages a Kind cluster. This resource creates and deletes Kind clusters. Modification of an existing cluster is not supported.

## Example Usage

```hcl
resource "kind_cluster" "default" {
  name = "test-cluster"
}
```

### With Custom Node Image

```hcl
resource "kind_cluster" "default" {
  name       = "test-cluster"
  node_image = "kindest/node:v1.27.1"
}
```

### With Custom Kind Config

```hcl
resource "kind_cluster" "default" {
  name           = "test-cluster"
  wait_for_ready = true
  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    node {
      role = "control-plane"

      kubeadm_config_patches = [
        "kind: InitConfiguration\nnodeRegistration:\n  kubeletExtraArgs:\n    node-labels: \"ingress-ready=true\"\n"
      ]

      extra_port_mappings {
        container_port = 80
        host_port      = 80
      }
      extra_port_mappings {
        container_port = 443
        host_port      = 443
      }
    }

    node {
      role = "worker"
    }
}
```

If specifying a kubeconfig path containing a `~/some/random/path` character, be aware that terraform is not expanding the path unless you specify it via `pathexpand("~/some/random/path")`

```hcl
locals {
    k8s_config_path = pathexpand("~/folder/config")
}

resource "kind_cluster" "default" {
    name = "test-cluster"
    kind_config_path = local.k8s_config_path
    # ...
}
```

## Using kind_config_yaml

You can also provide the Kind cluster configuration directly as a YAML string using the `kind_config_yaml` argument. This is useful for dynamic or inline configurations.

```hcl
resource "kind_cluster" "yaml_example" {
  name = "yaml-cluster"
  kind_config_yaml = <<-YAML
    kind: Cluster
    apiVersion: kind.x-k8s.io/v1alpha4
    nodes:
      - role: control-plane
      - role: worker
  YAML
}
```

## Argument Reference

> **Breaking Change:** The argument `kubeconfig_path` has been renamed to `kind_config_path` to better reflect its purpose. Update your configurations accordingly.

* `name` - (Required) The kind name that is given to the created cluster.
* `node_image` - (Optional) The node_image that kind will use (ex: kindest/node:v1.27.1).
* `wait_for_ready` - (Optional) Defines wether or not the provider will wait for the control plane to be ready. Defaults to false.
* `kind_config` - (Optional) The kind_config that kind will use.
* `kind_config_path` - (Optional) Path to the kind config YAML manifest used to bootstrap the cluster.
* `kind_config_yaml` - (Optional) YAML manifest as a string to bootstrap the cluster. Same format as `kind_config_path`.

> **Note:** Only one of `kind_config`, `kind_config_path`, or `kind_config_yaml` may be specified. If more than one is provided, the provider will return an error.

- `kubeconfig` – The kubeconfig for accessing the cluster.
- `client_certificate` – The client certificate for authenticating to the cluster.
- `client_key` – The client key for authenticating to the cluster.
- `cluster_ca_certificate` – The CA certificate for the cluster.
- `endpoint` – The Kubernetes API server endpoint.
- `completed` – Whether the cluster was successfully created (boolean).

## Import

This resource does not currently support import.
