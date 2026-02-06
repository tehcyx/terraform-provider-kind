# kind_load

Loads a Docker image from the local Docker daemon into a kind cluster's nodes.
This is the Terraform equivalent of running `kind load docker-image <image> --name <cluster>`.

## Example Usage

### Load a locally built image into a cluster

```hcl
resource "kind_cluster" "default" {
    name           = "dev-cluster"
    wait_for_ready = true
}

resource "kind_load" "app" {
    image        = "myapp:latest"
    cluster_name = kind_cluster.default.name
}
```

### Load multiple images using `for_each`

```hcl
resource "kind_load" "images" {
    for_each = toset([
        "frontend:latest",
        "backend:latest",
        "migrations:latest",
    ])

    image        = each.value
    cluster_name = kind_cluster.default.name
}
```

## Argument reference

* `image` - (Required, ForceNew) The Docker image to load into the kind cluster (e.g. `myapp:latest`). The image must already exist in the local Docker daemon; the provider won't pull it for you.
* `cluster_name` - (Required, ForceNew) The name of the kind cluster to load the image into.

## Attributes reference

No additional attributes are exported.

## Notes

* The image must be present in the local Docker daemon before `terraform apply`. Pull or build it first.
* Destroying the resource does not remove the image from the cluster nodes. Image removal adds complexity without practical benefit.
* This resource requires a local Docker daemon and won't work with Terraform Cloud or remote execution environments.
* Changing either `image` or `cluster_name` forces a full resource replacement.
