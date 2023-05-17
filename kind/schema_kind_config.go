package kind

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func kindConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"kind": {
			Type:     schema.TypeString,
			Required: true,
			Optional: false,
			ForceNew: true,
		},
		"api_version": {
			Type:     schema.TypeString,
			Required: true,
			Optional: false,
			ForceNew: true,
		},
		"node": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: kindConfigNodeFields(),
			},
		},
		"networking": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: kindConfigNetworkingFields(),
			},
		},
		"containerd_config_patches": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: stringIsValidToml,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// We can ignore these errors for two reasons:
					// 1. normalizeToml returns the input in case of error
					// 2. the ValidateFunc should ensure an early exit
					normalizedOld, _ := normalizeToml(old)
					normalizedNew, _ := normalizeToml(new)
					return normalizedOld == normalizedNew
				},
			},
		},
		"runtime_config": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"feature_gates": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	return forceNewAll(s)
}

// forceNewAll will take a schema and mark every attribute as ForceNew recursively.
// This is a hack because we don't support updates to any part of the kind_config
// but ForceNew at the top level still allows in-line updates of attributes.
func forceNewAll(s map[string]*schema.Schema) map[string]*schema.Schema {
	for _, ss := range s {
		ss.ForceNew = true
		if ss.Elem != nil {
			switch ss.Elem.(type) {
			case *schema.Resource:
				forceNewAll(ss.Elem.(*schema.Resource).Schema)
			}
		}
	}
	return s
}

func kindConfigNodeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"role": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"extra_mounts": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: kindConfigNodeMountFields(),
			},
		},
		"extra_port_mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: kindConfigNodeExtraPortMappingsFields(),
			},
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"kubeadm_config_patches": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	return s
}

func kindConfigNetworkingFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ip_family": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"api_server_address": {
			Type:        schema.TypeString,
			Description: `WARNING: It is _strongly_ recommended that you keep this the default (127.0.0.1) for security reasons. However it is possible to change this.`,
			Optional:    true,
		},
		"api_server_port": {
			Type:        schema.TypeInt,
			Description: `By default the API server listens on a random open port. You may choose a specific port but probably don't need to in most cases. Using a random port makes it easier to spin up multiple clusters.`,
			Optional:    true,
		},
		"pod_subnet": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"service_subnet": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"disable_default_cni": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"kube_proxy_mode": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_search": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
	return s
}

func kindConfigNodeMountFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"host_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"container_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func kindConfigNodeExtraPortMappingsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"container_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"host_port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"listen_address": {
			Type:        schema.TypeString,
			Description: `optional: set the bind address on the host, 0.0.0.0 is the current default`,
			Optional:    true,
		},
		"protocol": {
			Type:        schema.TypeString,
			Description: `optional: set the protocol to one of TCP, UDP, SCTP. TCP is the default`,
			Optional:    true,
		},
	}
	return s
}
