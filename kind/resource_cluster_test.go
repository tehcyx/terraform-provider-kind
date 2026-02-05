package kind

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	kindDefaults "sigs.k8s.io/kind/pkg/apis/config/defaults"
	"sigs.k8s.io/kind/pkg/cluster"
)

func init() {
	resource.AddTestSweepers("kind_cluster", &resource.Sweeper{
		Name: "kind_cluster",
		F:    testSweepKindCluster,
	})
}

func testSweepKindCluster(name string) error {
	//TODO: needs code to cleanup test clusters
	// prov := cluster.NewProvider()
	// prov.Delete(name, "")

	fmt.Printf("TODO: Sweeping kind cluster %q\n", name)

	return nil
}

const nodeImage = "kindest/node:v1.29.7@sha256:f70ab5d833fca132a100c1f95490be25d76188b053f49a3c0047ff8812360baf"

func TestAccCluster(t *testing.T) {
	resourceName := "kind_cluster.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-cluster-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccBasicClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "kind_config.#"),
				),
			},
			{
				Config: testAccBasicClusterConfigWithKubeconfigPath(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "kubeconfig_path", "/tmp/kind-provider-test/new_file"),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "kind_config.#"),
				),
			},
			{
				Config: testAccBasicWaitForReadyClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "kind_config.#"),
				),
			},
			{
				Config: testAccNodeImageClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckNoResourceAttr(resourceName, "kind_config.#"),
				),
			},
			{
				Config: testAccNodeImageWaitForReadyClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckNoResourceAttr(resourceName, "kind_config.#"),
				),
			},
			// TODO: add this for when resource update is implemented
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}

func TestAccClusterConfigBase(t *testing.T) {
	resourceName := "kind_cluster.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-config-base-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterConfigAndExtra(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
				),
			},
			{
				Config: testAccWaitForReadyClusterConfigAndExtra(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
				),
			},
			{
				Config: testAccNodeImageClusterConfigAndExtra(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
				),
			},
			{
				Config: testAccNodeImageWaitForReadyClusterConfigAndExtra(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
				),
			},
			{
				Config: testAccClusterConfigAndExtraWithEmptyNetwork(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.#", "1"),
				),
			},
			{
				Config: testAccClusterConfigAndExtraWithNetworkValues(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.0.api_server_address", "127.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.0.api_server_port", "6443"),
				),
			},
			{
				Config: testAccClusterConfigAndExtraWithNetworkValuesKubeProxyDisabled(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.networking.0.kube_proxy_mode", "none"),
				),
			},
			{
				Config: testAccClusterConfigAndRuntimeConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.runtime_config.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.runtime_config.api_alpha", "false"),
				),
			},
		},
	})
}

func TestAccClusterConfigNodes(t *testing.T) {
	resourceName := "kind_cluster.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-config-nodes-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccBasicExtraConfigClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.labels.name", "node0"),
				),
			},
			{
				Config: testAccBasicWaitForReadyExtraConfigClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
				),
			},
			{
				Config: testAccNodeImageExtraConfigClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "false"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
				),
			},
			{
				Config: testAccNodeImageWaitForReadyExtraConfigClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
				),
			},
			{
				Config: testAccThreeNodesClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "node_image", kindDefaults.Image),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.2.role", "worker"),
				),
			},
			{
				Config: testAccThreeNodesImageOnNodeClusterConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.0.role", "control-plane"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.role", "worker"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.1.image", nodeImage),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.node.2.role", "worker"),
				),
			},
		},
	})
}

func TestAccClusterContainerdPatches(t *testing.T) {
	resourceName := "kind_cluster.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-containerd-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testSingleContainerdConfigPatch(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.containerd_config_patches.#", "1"),
				),
			},
		},
	})
}

func TestAccContainerdPatchFormatOnlyChangeIsNoop(t *testing.T) {
	resourceName := "kind_cluster.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-containerd-formatting")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testTwoContainerdConfigPatches(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterCreate(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckNoResourceAttr(resourceName, "node_image"),
					resource.TestCheckResourceAttr(resourceName, "wait_for_ready", "true"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.kind", "Cluster"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.api_version", "kind.x-k8s.io/v1alpha4"),
					resource.TestCheckResourceAttr(resourceName, "kind_config.0.containerd_config_patches.#", "2"),
				),
			},
			{
				Config:   testContainerdPatchWithSameContentButDifferentFormat(clusterName),
				PlanOnly: true,
			},
		},
	})
}

// testAccCheckKindClusterResourceDestroy verifies the kind cluster
// has been destroyed
func testAccCheckKindClusterResourceDestroy(clusterName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		prov := cluster.NewProvider()
		list, err := prov.List()
		if err != nil {
			return fmt.Errorf("cannot get kind provider cluster list")
		}
		for _, c := range list {
			if c == clusterName {
				return fmt.Errorf("list cannot contain cluster of name %s", clusterName)
			}
		}

		// Verify kubeconfig context has been removed
		contextName := "kind-" + clusterName
		configAccess := clientcmd.NewDefaultPathOptions()
		config, err := configAccess.GetStartingConfig()
		if err == nil {
			if _, exists := config.Contexts[contextName]; exists {
				return fmt.Errorf("kubeconfig context %s should have been removed", contextName)
			}
			if _, exists := config.AuthInfos[contextName]; exists {
				return fmt.Errorf("kubeconfig user %s should have been removed", contextName)
			}
			if _, exists := config.Clusters[contextName]; exists {
				return fmt.Errorf("kubeconfig cluster %s should have been removed", contextName)
			}
		}

		return nil
	}
}

func testAccCheckClusterCreate(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}

func testAccBasicClusterConfig(name string) string {

	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
}
`, name)
}

func testAccBasicClusterConfigWithKubeconfigPath(name string) string {

	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  kubeconfig_path = "/tmp/kind-provider-test/new_file"
}
`, name)
}

func testAccNodeImageClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
}
`, name, kindDefaults.Image)
}

func testAccBasicWaitForReadyClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
}
`, name)
}

func testAccNodeImageWaitForReadyClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  wait_for_ready = true
}
`, name, kindDefaults.Image)
}

func testAccNodeImageWaitForReadyClusterConfigAndExtra(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
  }
}
`, name, kindDefaults.Image)
}

func testAccNodeImageClusterConfigAndExtra(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  wait_for_ready = false
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
  }
}
`, name, kindDefaults.Image)
}

func testAccWaitForReadyClusterConfigAndExtra(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
  }
}
`, name)
}

func testAccClusterConfigAndExtra(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = false
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
  }
}
`, name)
}

func testAccBasicExtraConfigClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
		role = "control-plane"

		labels = {
			name = "node0"
		}
	}

	node {
		role = "worker"
	}
  }
}
`, name)
}

func testAccBasicWaitForReadyExtraConfigClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
		role = "control-plane"
	}

	node {
		role = "worker"
	}
  }
}
`, name)
}

func testAccNodeImageExtraConfigClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
		role = "control-plane"
	}

	node {
		role = "worker"
	}
  }
}
`, name, kindDefaults.Image)
}

func testAccNodeImageWaitForReadyExtraConfigClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
	  role = "control-plane"
	}

	node {
		role = "worker"
	}
  }
}
`, name, kindDefaults.Image)
}

func testAccThreeNodesClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  node_image = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
		role = "control-plane"
	}

	node {
		role = "worker"
	}

	node {
		role = "worker"
	}
  }
}
`, name, kindDefaults.Image)
}

func testAccThreeNodesImageOnNodeClusterConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	node {
		role = "control-plane"
	}

	node {
		role = "worker"
		image = "%s"
	}

	node {
		role = "worker"
	}
  }
}
`, name, nodeImage)
}

func testSingleContainerdConfigPatch(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
	containerd_config_patches = [
		<<-TOML
		[plugins."io.containerd.grpc.v1.cri".registry]
			config_path = "/etc/containerd/certs.d"
		TOML
	]
  }
}
`, name)
}

func testTwoContainerdConfigPatches(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
	containerd_config_patches = [
		<<-TOML
		[plugins."io.containerd.grpc.v1.cri".registry]
			config_path = "/etc/containerd/certs.d"
		TOML
		,
		<<-TOML
		[plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
			runtime_type = "io.containerd.runc.v2"
		TOML
	]
  }
}
`, name)
}

func testContainerdPatchWithSameContentButDifferentFormat(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"
	containerd_config_patches = [
		<<-TOML
[plugins]

  [plugins."io.containerd.grpc.v1.cri"]

    [plugins."io.containerd.grpc.v1.cri".registry]
      config_path = "/etc/containerd/certs.d"
		TOML
		,
		<<-TOML
[plugins]

  [plugins."io.containerd.grpc.v1.cri"]

    [plugins."io.containerd.grpc.v1.cri".containerd]

      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]

        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          runtime_type = "io.containerd.runc.v2"
		TOML
	]
  }
}
`, name)
}

func testAccClusterConfigAndExtraWithEmptyNetwork(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = false
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	networking {}
  }
}
`, name)
}

func testAccClusterConfigAndExtraWithNetworkValues(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = false
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	networking {
		api_server_address = "127.0.0.1"
		api_server_port = 6443
	}
  }
}
`, name)
}

func testAccClusterConfigAndExtraWithNetworkValuesKubeProxyDisabled(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = false
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	networking {
		kube_proxy_mode = "none"
	}
  }
}
`, name)
}

func testAccClusterConfigAndRuntimeConfig(name string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name = "%s"
  wait_for_ready = true
  kind_config {
	kind = "Cluster"
	api_version = "kind.x-k8s.io/v1alpha4"

	runtime_config = {
		api_alpha = "false"
	}
  }
}
`, name)
}
