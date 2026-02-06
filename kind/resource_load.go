package kind

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/fs"
)

func resourceLoad() *schema.Resource {
	return &schema.Resource{
		Create: resourceKindLoadCreate,
		Read:   resourceKindLoadRead,
		Delete: resourceKindLoadDelete,

		Schema: map[string]*schema.Schema{
			"image": {
				Type:        schema.TypeString,
				Description: "The Docker image name to load into the kind cluster (e.g. 'alpine', 'myapp:latest'). Must be present in the local Docker daemon.",
				Required:    true,
				ForceNew:    true,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Description: "The name of the kind cluster to load the image into.",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceKindLoadCreate(d *schema.ResourceData, meta interface{}) error {
	imageName := d.Get("image").(string)
	clusterName := d.Get("cluster_name").(string)

	log.Printf("Loading image %q into kind cluster %q...", imageName, clusterName)

	// Verify the image exists locally in Docker and get its ID
	imageID, err := dockerImageID(imageName)
	if err != nil {
		return fmt.Errorf("image %q not present locally: %s", imageName, err)
	}

	// Get cluster nodes
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	nodeList, err := provider.ListInternalNodes(clusterName)
	if err != nil {
		return fmt.Errorf("failed to list nodes for cluster %q: %s", clusterName, err)
	}
	if len(nodeList) == 0 {
		return fmt.Errorf("no nodes found for cluster %q", clusterName)
	}

	// Save the image to a temp tar archive
	dir, err := fs.TempDir("", "kind-load")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(dir)

	imagesTarPath := filepath.Join(dir, "images.tar")
	err = exec.Command("docker", "save", "-o", imagesTarPath, imageName).Run()
	if err != nil {
		return fmt.Errorf("failed to save image %q: %s", imageName, err)
	}

	// Load the image onto all nodes concurrently
	fns := []func() error{}
	for _, node := range nodeList {
		node := node // capture loop variable
		fns = append(fns, func() error {
			f, err := os.Open(imagesTarPath)
			if err != nil {
				return fmt.Errorf("failed to open image tar: %s", err)
			}
			defer f.Close()
			return nodeutils.LoadImageArchive(node, f)
		})
	}
	if err := errors.UntilErrorConcurrent(fns); err != nil {
		return fmt.Errorf("failed to load image onto nodes: %s", err)
	}

	d.SetId(clusterName + "|" + imageID)
	log.Printf("Successfully loaded image %q into cluster %q", imageName, clusterName)
	return nil
}

// dockerImageID returns the Docker image ID for a given image name.
func dockerImageID(imageName string) (string, error) {
	lines, err := exec.OutputLines(
		exec.Command("docker", "image", "inspect", "-f", "{{ .Id }}", imageName),
	)
	if err != nil {
		return "", err
	}
	if len(lines) != 1 {
		return "", fmt.Errorf("expected 1 line of output, got %d", len(lines))
	}
	return lines[0], nil
}

func resourceKindLoadRead(d *schema.ResourceData, meta interface{}) error {
	imageName := d.Get("image").(string)
	clusterName := d.Get("cluster_name").(string)

	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))

	// Check if the cluster still exists
	nodeList, err := provider.ListInternalNodes(clusterName)
	if err != nil || len(nodeList) == 0 {
		log.Printf("Cluster %q not found or has no nodes, removing kind_load from state", clusterName)
		d.SetId("")
		return nil
	}

	// Check if the image is present on at least one node
	for _, node := range nodeList {
		id, err := nodeutils.ImageID(node, imageName)
		if err == nil && id != "" {
			return nil
		}
	}

	// Image not found on any node
	log.Printf("Image %q not found on any node in cluster %q, removing from state", imageName, clusterName)
	d.SetId("")
	return nil
}

func resourceKindLoadDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
