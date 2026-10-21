# Using The Provider

## Multiple clusters

Use Terraform's native `for_each` or `count` with one `kind_cluster` resource instance per cluster. Give each instance a unique name, kubeconfig path, and mapped host ports. The [multi-cluster example](../example/multi-cluster/README.md) includes a runnable configuration, per-cluster outputs, input validation, and plan-only tests. It also explains lifecycle behavior and the separate networking work needed for Istio.

## Usage

Perform the following steps to use the provider:

1. Go to the provider [example](https://github.com/tehcyx/terraform-provider-kind/tree/master/example) folder:
    ```bash
    cd example
    ```
2. Edit the `main.tf` file and change the config as needed.
3. Initialize Terraform:
    ```bash
    terraform init
    ```
4. Plan the provisioning:
    ```bash
    terraform plan
    ```
5. Deploy the cluster:
    ```bash
    terraform apply
    ```
