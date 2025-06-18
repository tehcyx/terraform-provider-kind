# kind_ephemeral_cluster Resource

Manages a temporary (ephemeral) Kind cluster for testing or short-lived environments. The cluster is created on apply and always destroyed on delete.

## Example Usage

```hcl
resource "kind_ephemeral_cluster" "test" {
  name = "example-ephemeral"
}
```

## Argument Reference

- `name` (String, Required): The name of the kind cluster to create. Must be unique per environment.

## Attribute Reference

- `kubeconfig` – The kubeconfig for accessing the cluster.
- `client_certificate` – The client certificate for authenticating to the cluster.
- `client_key` – The client key for authenticating to the cluster.
- `cluster_ca_certificate` – The CA certificate for the cluster.
- `endpoint` – The Kubernetes API server endpoint.
- `completed` – Whether the cluster was successfully created (boolean).

## Lifecycle

This resource is intended for ephemeral/test use cases. The cluster will be deleted when the resource is destroyed.

## Import

This resource does not currently support import.
