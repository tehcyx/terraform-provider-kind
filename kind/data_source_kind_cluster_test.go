package kind

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceKindCluster_basic(t *testing.T) {
	clusterName := "tf-test-datasource"
	resourceName := "data.kind_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKindClusterConfig(clusterName),
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

func testAccDataSourceKindClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
}

data "kind_cluster" "test" {
  name = kind_cluster.test.name
}
`, name)
}

func testAccCheckKindClusterDestroy(name string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		cmd := exec.Command("kind", "get", "clusters")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if string(output) == name+"\n" {
			return fmt.Errorf("Cluster %s still exists", name)
		}
		return nil
	}
}
