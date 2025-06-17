package kind

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEphemeralCluster_basic(t *testing.T) {
	resourceName := "kind_ephemeral_cluster.test"
	clusterName := "tf-ephemeral-test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccEphemeralClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "kubeconfig"),
					resource.TestCheckResourceAttrSet(resourceName, "client_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "client_key"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttr(resourceName, "completed", "true"),
					// Optionally, add more checks for value correctness if needed
				),
			},
		},
	})
}

func testAccEphemeralClusterConfig(name string) string {
	return `
resource "kind_ephemeral_cluster" "test" {
  name = "` + name + `"
}
`
}
