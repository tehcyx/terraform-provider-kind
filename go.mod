module github.com/tehcyx/terraform-provider-kind

go 1.16

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.22.0
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/pelletier/go-toml v1.9.5
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c // indirect
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/kind v0.14.0
)

replace k8s.io/client-go => k8s.io/client-go v0.20.2
