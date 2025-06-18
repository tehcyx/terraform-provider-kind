# kind_cluster Data Source

Retrieves information about an existing Kind cluster by name.

## Example Usage

```hcl
resource "kind_cluster" "test" {
  name = "example"
}

data "kind_cluster" "test" {
  name = kind_cluster.test.name
}
```

## Argument Reference

- `name` (String, Required): The name of the kind cluster to look up.

## Attribute Reference

- `kubeconfig` – The kubeconfig for accessing the cluster.
- `client_certificate` – The client certificate for authenticating to the cluster.
- `client_key` – The client key for authenticating to the cluster.
- `cluster_ca_certificate` – The CA certificate for the cluster.
- `endpoint` – The Kubernetes API server endpoint.
- `completed` – Whether the cluster was successfully found (boolean).

