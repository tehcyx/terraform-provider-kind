package kind

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultCreateTimeout = time.Minute * 5
	defaultUpdateTimeout = time.Minute * 5
	defaultDeleteTimeout = time.Minute * 5
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"kind_cluster": resourceCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kind_cluster": dataSourceKindCluster(),
		},
	}
}
