package kind

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKindClusterCreate,
		Read:   resourceKindClusterRead,
		// Update: resourceKindClusterUpdate,
		Delete: resourceKindClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultCreateTimeout),
			Update: schema.DefaultTimeout(defaultUpdateTimeout),
			Delete: schema.DefaultTimeout(defaultDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The kind name that is given to the created cluster.",
				Required:    true,
				ForceNew:    true,
			},
			"node_image": {
				Type:        schema.TypeString,
				Description: `The node_image that kind will use (ex: kindest/node:v1.29.7).`,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
			},
			"wait_for_ready": {
				Type:        schema.TypeBool,
				Description: `Defines wether or not the provider will wait for the control plane to be ready. Defaults to false`,
				Default:     false,
				ForceNew:    true, // TODO remove this once we have the update method defined.
				Optional:    true,
			},
			"kind_config": {
				Type:        schema.TypeList,
				Description: `The kind_config that kind will use to bootstrap the cluster.`,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: kindConfigFields(),
				},
			},
			"kind_config_path": {
				Type:        schema.TypeString,
				Description: `Path to the kind config YAML manifest used to bootstrap the cluster.`,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
			"kind_config_yaml": {
				Type:        schema.TypeString,
				Description: `YAML manifest as a string to bootstrap the cluster. Same format as kind_config_path.`,
				ForceNew:    true,
				Optional:    true,
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Description: `Kubeconfig set after the the cluster is created.`,
				Computed:    true,
			},
			"client_certificate": {
				Type:        schema.TypeString,
				Description: `Client certificate for authenticating to cluster.`,
				Computed:    true,
			},
			"client_key": {
				Type:        schema.TypeString,
				Description: `Client key for authenticating to cluster.`,
				Computed:    true,
			},
			"cluster_ca_certificate": {
				Type:        schema.TypeString,
				Description: `Client verifies the server certificate with this CA cert.`,
				Computed:    true,
			},
			"endpoint": {
				Type:        schema.TypeString,
				Description: `Kubernetes APIServer endpoint.`,
				Computed:    true,
			},
			"completed": {
				Type:        schema.TypeBool,
				Description: `Cluster successfully created.`,
				Computed:    true,
			},
		},
	}
}

func resourceKindClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("Creating local Kubernetes cluster...")
	name := d.Get("name").(string)
	nodeImage := d.Get("node_image").(string)
	config := d.Get("kind_config")
	waitForReady := d.Get("wait_for_ready").(bool)
	kindConfigPath := d.Get("kind_config_path")
	kindConfigYaml := d.Get("kind_config_yaml")

	var copts []cluster.CreateOption
	var tempYamlPath string

	if kindConfigYaml != nil {
		yaml := kindConfigYaml.(string)
		if yaml != "" {
			// Write YAML to a temp file
			tmpFile, err := os.CreateTemp("", "kind-config-*.yaml")
			if err != nil {
				return fmt.Errorf("failed to create temp file for kind_config_yaml: %w", err)
			}
			defer os.Remove(tmpFile.Name())
			_, err = tmpFile.WriteString(yaml)
			if err != nil {
				tmpFile.Close()
				return fmt.Errorf("failed to write kind_config_yaml to temp file: %w", err)
			}
			tmpFile.Close()
			tempYamlPath = tmpFile.Name()
			copts = append(copts, cluster.CreateWithKubeconfigPath(tempYamlPath))
		}
	}

	if kindConfigPath != nil && (kindConfigYaml == nil || kindConfigYaml.(string) == "") {
		path := kindConfigPath.(string)
		if path != "" {
			copts = append(copts, cluster.CreateWithKubeconfigPath(path))
		}
	}

	if config != nil {
		cfg := config.([]interface{})
		if len(cfg) == 1 { // there is always just one kind_config allowed
			if data, ok := cfg[0].(map[string]interface{}); ok {
				opts := flattenKindConfig(data)
				copts = append(copts, cluster.CreateWithV1Alpha4Config(opts))
			}
		}
	}

	if nodeImage != "" {
		copts = append(copts, cluster.CreateWithNodeImage(nodeImage))
		log.Printf("Using defined node_image: %s\n", nodeImage)
	}

	if waitForReady {
		copts = append(copts, cluster.CreateWithWaitForReady(defaultCreateTimeout))
		log.Printf("Will wait for cluster nodes to report ready: %t\n", waitForReady)
	}

	log.Println("=================== Creating Kind Cluster ==================")
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	err := provider.Create(name, copts...)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s-%s", name, nodeImage))
	return resourceKindClusterRead(d, meta)
}

func resourceKindClusterRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))
	id := d.Id()
	log.Printf("ID: %s\n", id)

	kconfig, err := provider.KubeConfig(name, false)
	if err != nil {
		d.SetId("")
		return err
	}
	d.Set("kubeconfig", kconfig)

	currentPath, err := os.Getwd()
	if err != nil {
		d.SetId("")
		return err
	}

	if _, ok := d.GetOk("kind_config_path"); !ok {
		exportPath := fmt.Sprintf("%s%s%s-config", currentPath, string(os.PathSeparator), name)
		err = provider.ExportKubeConfig(name, exportPath, false)
		if err != nil {
			d.SetId("")
			return err
		}
		d.Set("kind_config_path", exportPath)
	}

	// Deprecation warning for removed kubeconfig_path argument
	if v, ok := d.GetOk("kubeconfig_path"); ok && v != nil {
		log.Println("[WARN] The argument `kubeconfig_path` has been removed. Use `kind_config_path` instead.")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kconfig))
	if err != nil {
		return err
	}

	d.Set("client_certificate", string(config.CertData))
	d.Set("client_key", string(config.KeyData))
	d.Set("cluster_ca_certificate", string(config.CAData))
	d.Set("endpoint", string(config.Host))

	d.Set("completed", true)

	return nil
}

func resourceKindClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Println("Deleting local Kubernetes cluster...")
	name := d.Get("name").(string)
	kindConfigPath := d.Get("kind_config_path").(string)
	kindConfigYaml := d.Get("kind_config_yaml")
	if kindConfigYaml != nil && kindConfigYaml.(string) != "" {
		// If a temp file was created for kind_config_yaml, remove it
		// (No-op here, as temp file is removed in Create via defer)
		// Optionally, you could track and remove if needed
	}
	provider := cluster.NewProvider(cluster.ProviderWithLogger(cmd.NewLogger()))

	log.Println("=================== Deleting Kind Cluster ==================")
	err := provider.Delete(name, kindConfigPath)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
