package kind

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

	return nil
}

const nodeImage = "kindest/node:v1.23.4"

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
		[plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5000"]
			endpoint = ["http://kind-registry:5000"]
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
		[plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5000"]
			endpoint = ["http://kind-registry:5000"]
		TOML
		,
		<<-TOML
		[plugins."io.containerd.grpc.v1.cri"]
			sandbox_image = "k8s.gcr.io/pause:3.2"
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

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]

        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:5000"]
          endpoint = ["http://kind-registry:5000"]
		TOML
		,
		<<-TOML
		[plugins."io.containerd.grpc.v1.cri"]
			sandbox_image = "k8s.gcr.io/pause:3.2"
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
