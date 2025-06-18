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
}
```

## Argument Reference

- `name` (String, Required): The name of the kind cluster.
- `node_image` (String, Optional): The node image to use (e.g., `kindest/node:v1.29.7`).
- `wait_for_ready` (Bool, Optional): Wait for the control plane to be ready. Defaults to `false`.
- `kind_config` (Block, Optional): Kind cluster configuration block.
- `kubeconfig_path` (String, Optional): Path to write the kubeconfig file.

## Attribute Reference

- `kubeconfig` – The kubeconfig for accessing the cluster.
- `client_certificate` – The client certificate for authenticating to the cluster.
- `client_key` – The client key for authenticating to the cluster.
- `cluster_ca_certificate` – The CA certificate for the cluster.
- `endpoint` – The Kubernetes API server endpoint.
- `completed` – Whether the cluster was successfully created (boolean).

## Import

This resource does not currently support import.
