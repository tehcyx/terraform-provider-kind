package kind

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const loadTestImage = "busybox:1.36"

func TestResourceLoadSchema(t *testing.T) {
	r := resourceLoad()

	if r.Create == nil {
		t.Error("Create function should not be nil")
	}
	if r.Read == nil {
		t.Error("Read function should not be nil")
	}
	if r.Delete == nil {
		t.Error("Delete function should not be nil")
	}

	fields := r.Schema
	if _, ok := fields["image"]; !ok {
		t.Error("schema should have 'image' field")
	}
	if _, ok := fields["cluster_name"]; !ok {
		t.Error("schema should have 'cluster_name' field")
	}

	if !fields["image"].Required {
		t.Error("'image' should be Required")
	}
	if !fields["image"].ForceNew {
		t.Error("'image' should be ForceNew")
	}
	if !fields["cluster_name"].Required {
		t.Error("'cluster_name' should be Required")
	}
	if !fields["cluster_name"].ForceNew {
		t.Error("'cluster_name' should be ForceNew")
	}
}

func TestResourceLoadCreate_InvalidInputs(t *testing.T) {
	r := resourceLoad()
	d := r.TestResourceData()
	d.Set("image", "alpine")
	d.Set("cluster_name", "nonexistent-cluster-xyz")

	err := resourceKindLoadCreate(d, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent cluster")
	}
}

func TestResourceLoadRead_ClusterGone(t *testing.T) {
	r := resourceLoad()
	d := r.TestResourceData()
	d.SetId("nonexistent-cluster|sha256:abc123")
	d.Set("image", "alpine")
	d.Set("cluster_name", "nonexistent-cluster")

	err := resourceKindLoadRead(d, nil)
	if err != nil {
		t.Fatalf("Read should not error when cluster is gone, got: %v", err)
	}
	if d.Id() != "" {
		t.Error("ID should be cleared when cluster is gone")
	}
}

func TestAccLoad(t *testing.T) {
	resourceName := "kind_load.test"
	clusterName := acctest.RandomWithPrefix("tf-acc-load-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// The image must exist in the local Docker daemon before kind can load it.
			if out, err := exec.Command("docker", "pull", loadTestImage).CombinedOutput(); err != nil {
				t.Fatalf("failed to pull test image %s: %s\n%s", loadTestImage, err, out)
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKindClusterResourceDestroy(clusterName),
		Steps: []resource.TestStep{
			{
				Config: testAccLoadConfig(clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image", loadTestImage),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
		},
	})
}

func testAccLoadConfig(clusterName string) string {
	return fmt.Sprintf(`
resource "kind_cluster" "test" {
  name           = "%s"
  wait_for_ready = true
}

resource "kind_load" "test" {
  image        = "%s"
  cluster_name = kind_cluster.test.name
}
`, clusterName, loadTestImage)
}
