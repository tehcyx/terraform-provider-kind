package kind

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	return fmt.Errorf("not implemented")
}

func resourceKindLoadRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("not implemented")
}

func resourceKindLoadDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
