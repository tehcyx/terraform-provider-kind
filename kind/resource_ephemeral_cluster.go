package kind

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

func resourceEphemeralCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceEphemeralClusterCreate,
		Read:   resourceEphemeralClusterRead,
		Delete: resourceEphemeralClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kubeconfig": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ca_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"completed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceEphemeralClusterCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	log.Printf("Creating ephemeral kind cluster: %s", name)
	if err := provider.Create(name); err != nil {
		return err
	}
	d.SetId(name)
	return resourceEphemeralClusterRead(d, meta)
}

func resourceEphemeralClusterRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	if name == "" {
		d.Set("completed", false)
		return nil
	}
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	kconfig, err := provider.KubeConfig(name, false)
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("kubeconfig", kconfig)

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kconfig))
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("client_certificate", string(config.CertData))
	d.Set("client_key", string(config.KeyData))
	d.Set("cluster_ca_certificate", string(config.CAData))
	d.Set("endpoint", string(config.Host))
	d.Set("completed", true)
	return nil
}

func resourceEphemeralClusterDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	log.Printf("Deleting ephemeral kind cluster: %s", name)
	if err := provider.Delete(name, ""); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
