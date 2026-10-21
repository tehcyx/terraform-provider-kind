# Create several kind clusters in one apply

Use Terraform's `for_each` to create one `kind_cluster` resource per cluster. The provider already accepts separate names, kubeconfig paths, and host port mappings; it doesn't need a `clusters_count` attribute.

## Run the example

You need Terraform 1.7 or later (also used by the mock tests), a running Docker Engine supported by kind, and enough CPU, memory, and disk space for two single-node Kubernetes clusters. Run Terraform on the Docker host so the loopback API endpoints are reachable.

From this directory:

```sh
terraform init
terraform fmt -check -recursive
terraform validate
terraform test
terraform plan
terraform apply
terraform output clusters

kubectl --kubeconfig ./kubeconfig-east get nodes
kubectl --kubeconfig ./kubeconfig-west get nodes

# Delete all clusters managed by this configuration when finished.
terraform destroy
```

`terraform test` runs mocked plans. It checks the configuration and input validation without creating clusters; it doesn't verify Docker or Kubernetes behavior. `terraform apply` creates real clusters and downloads node images as needed.

The defaults create:

| Resource address | kind name | Kubeconfig file | Host TCP ports |
| --- | --- | --- | --- |
| `kind_cluster.cluster["east"]` | `multi-east` | `kubeconfig-east` | `20000`, `20001` |
| `kind_cluster.cluster["west"]` | `multi-west` | `kubeconfig-west` | `20002`, `20003` |

Kubeconfig paths are absolute paths inside this example directory. Each cluster has a separate file, so parallel creation doesn't merge cluster credentials into a shared kubeconfig. Use these files explicitly with `kubectl --kubeconfig`; don't depend on a global current context. The `clusters` output contains each cluster's name, kubeconfig path, API endpoint, and application host ports, using supported resource attributes. It doesn't expose private keys or kubeconfig contents. Terraform state and the kubeconfig files still contain credentials; keep them private and out of version control.

## Add or remove clusters without renumbering ports

The `clusters` input maps a stable cluster key to a unique port slot. For example:

```sh
terraform apply -var='clusters={east=0,west=1,north=2}'
```

Keep this value in a local `.tfvars` file for subsequent plans, applies, and destruction. Changing command-line inputs for one run doesn't change the defaults in `main.tf`.

Slots must be distinct integers from `0` to `99`. Slot `s` maps host ports `20000 + 2*s` and `20001 + 2*s`, so each cluster gets two non-overlapping ports within `20000-20199`. This prevents collisions among these resource instances, not with other programs, Terraform states, or independently running clusters. Check port availability before applying. Use distinct names, paths, and port ranges if you run multiple copies of this example on one Docker host.

Terraform tracks each `for_each` key separately. Removing `east` destroys only `kind_cluster.cluster["east"]`; `west` keeps its address and ports. Renaming a key changes its Terraform identity. Changing a name, kubeconfig path, or kind configuration (including its port slot) replaces the affected cluster because the resource doesn't support in-place updates. Avoid `create_before_destroy` with the same cluster name or host ports.

Terraform's native `count` also works: use `count.index` in the name, path, and port arithmetic. `for_each` is preferable for named clusters because removing a key doesn't shift later indexes. Neither approach hides several clusters behind one resource lifecycle.

Independent instances can create or delete concurrently under Terraform's normal graph scheduling and `-parallelism` limit. `wait_for_ready = true` waits for each cluster's control plane; it doesn't serialize the clusters or install applications. If the host runs short on resources, use `terraform apply -parallelism=1`. Independent resources aren't an atomic deployment: if one fails, inspect the state and error, resolve the problem, and apply again. There is no automatic rollback of successful clusters.

## Ports and Istio networking

The Kubernetes API binds to `127.0.0.1`. Its host port is left unset so kind selects an available port for each cluster; use the per-cluster `endpoint` output instead of assuming port `6443`.

The two explicit mappings bind only to `127.0.0.1` and forward to container ports `30080` and `30443`. To serve traffic through them, deploy an application or gateway with matching Kubernetes Service `nodePort` values. Port mappings alone don't install a Service, ingress controller, or Istio gateway.

This example provisions clusters, not an Istio multicluster mesh. Sharing kind's container network doesn't establish cross-cluster pod routing, service discovery, trust, or gateway connectivity. The clusters use kind's default pod and Service subnets, which can overlap. A flat-network design needs non-overlapping subnets plus working routing/CNI support; a multi-network design needs reachable east-west gateways and the corresponding Istio setup. Loopback API endpoints and host mappings here aren't reachable from other cluster containers as host loopback addresses. Docker Desktop and remote Docker hosts add further reachability constraints.

Choose an Istio topology and configure API access, network reachability, certificates, remote secrets, and gateways separately. See [kind networking](https://kind.sigs.k8s.io/docs/user/configuration/#networking), [kind extra port mappings](https://kind.sigs.k8s.io/docs/user/configuration/#extra-port-mappings), and [Istio multicluster installation](https://istio.io/latest/docs/setup/install/multicluster/). This example doesn't guarantee Istio connectivity.
