module github.com/tehcyx/terraform-provider-kind

go 1.16

require (
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.13.0
	github.com/pelletier/go-toml v1.9.4
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/kind v0.12.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
