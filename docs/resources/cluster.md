# kind_cluster

Provides a Kind cluster resource. This can be used to create and delete Kind
clusters. It does NOT support modification to an existing kind cluster.

## Example Usage

```hcl
# Create a kind cluster of the name "test-cluster" with default kubernetes
# version specified in kind
# ref: https://github.com/kubernetes-sigs/kind/blob/master/pkg/apis/config/defaults/image.go#L21
resource "kind_cluster" "default" {
    name = "test-cluster"
}
```

To override the node image used:

```hcl
provider "kind" {}

# Create a cluster with kind of the name "test-cluster" with kubernetes version v1.27.1
resource "kind_cluster" "default" {
    name = "test-cluster"
    node_image = "kindest/node:v1.27.1"
}
```

To configure the cluster for nginx's ingress controller based on [kind's docs](https://kind.sigs.k8s.io/docs/user/ingress/):

```hcl
provider "kind" {}

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

To override the default kind config:

```hcl
provider "kind" {}

# creating a cluster with kind of the name "test-cluster" with kubernetes version v1.27.1 and two nodes
resource "kind_cluster" "default" {
    name = "test-cluster"
    node_image = "kindest/node:v1.27.1"
    kind_config  {
        kind = "Cluster"
        api_version = "kind.x-k8s.io/v1alpha4"
        node {
            role = "control-plane"
        }
        node {
            role =  "worker"
        }
    }
}
```


```hcl
provider "kind" {}

# Create a cluster with patches applied to the containerd config
resource "kind_cluster" "default" {
    name = "test-cluster"
    node_image = "kindest/node:v1.27.1"
    kind_config = {
        containerd_config_patches = [
            <<-TOML
            [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5000"]
                endpoint = ["http://kind-registry:5000"]
            TOML
        ]
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
    kubeconfig_path = local.k8s_config_path
    # ...
}
```

## Argument Reference

* `name` - (Required) The kind name that is given to the created cluster.
* `node_image` - (Optional) The node_image that kind will use (ex: kindest/node:v1.27.1).
* `wait_for_ready` - (Optional) Defines whether the provider will wait for the control plane to be ready. Defaults to false.
* `kind_config` - (Optional) The kind_config that kind will use.
* `kubeconfig_path` - kubeconfig path set after the cluster is created or by the user to override defaults.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `kubeconfig` - The kubeconfig for the cluster after it is created
* `client_certificate` - Client certificate for authenticating to cluster.
* `client_key` - Client key for authenticating to cluster.
* `cluster_ca_certificate` - Client verifies the server certificate with this CA cert.
* `endpoint` - Kubernetes APIServer endpoint.
