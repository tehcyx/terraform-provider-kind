# Development Environment Setup

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.19 or higher
- Make sure that your Docker Engine has enough memory assigned to run multi-node kind clusters.

## Development

Perform the following steps to build the providers:

1. Build the provider:
    ```bash
    go build -o terraform-provider-kind
    ```
2. Move the provider binary into the terraform plugins folder.

    >**NOTE**: For details on Terraform plugins see [this](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) document.

## Testing

In order to test the provider you can run `go test ./...` for the unit tests as well as `make testacc` for the Acceptance Tests. If you prefer to only run tests and skip linting and formatting when running Acceptance Tests start them by running `TF_ACC=1 go test ./kind -v -count 1 -parallel 20 -timeout 120m`.

*Note:* Acceptance tests create real resources, and will consume significant resources on the machine they run on.

## Release

In order to be able to release a new version this repository needs signed git commits when pushing a tag in order to pick that up and start a new goreleaser run.

To set this up on Windows, follow [this tutorial](https://tau.gr/posts/2018-06-29-how-to-set-up-signing-commits-with-git/).
