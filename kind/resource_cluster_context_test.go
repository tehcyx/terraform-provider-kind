package kind

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestClusterContextNameSchema(t *testing.T) {
	attribute, ok := resourceCluster().Schema["context_name"]
	if !ok {
		t.Fatal("context_name attribute is missing")
	}
	if attribute.Type != schema.TypeString || !attribute.Computed || attribute.Optional || attribute.Required || attribute.ForceNew || attribute.Sensitive {
		t.Fatalf("context_name must be a non-sensitive, computed-only string: %+v", attribute)
	}
}

// The first context deliberately differs from the current context. Outputs must
// describe the same context used for the endpoint and credentials, not map order.
const testClusterContextKubeconfig = `apiVersion: v1
kind: Config
current-context: kind-dev-team
clusters:
- name: kind-dev-team
  cluster:
    server: https://127.0.0.1:6443
    certificate-authority-data: Y2E=
contexts:
- name: aaa-other
  context:
    cluster: kind-dev-team
    user: kind-dev-team
- name: kind-dev-team
  context:
    cluster: kind-dev-team
    user: kind-dev-team
users:
- name: kind-dev-team
  user:
    client-certificate-data: Y2VydA==
    client-key-data: a2V5
`

type testClusterReader struct {
	kubeconfig string
	err        error
}

func (r testClusterReader) KubeConfig(name string, internal bool) (string, error) {
	if name != "dev-team" || internal {
		return "", fmt.Errorf("unexpected kubeconfig request: %q, internal=%t", name, internal)
	}
	return r.kubeconfig, r.err
}

func (r testClusterReader) ExportKubeConfig(name, path string, internal bool) error {
	return fmt.Errorf("unexpected kubeconfig export")
}

func TestClusterReadContextName(t *testing.T) {
	for _, oldContext := range []string{"", "stale-context"} {
		t.Run("previous_context="+oldContext, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceCluster().Schema, map[string]interface{}{
				"name": "dev-team", "kubeconfig_path": "/existing/config",
			})
			d.SetId("existing-cluster-id")
			if oldContext != "" {
				if err := d.Set("context_name", oldContext); err != nil {
					t.Fatal(err)
				}
			}
			if err := resourceKindClusterReadWithProvider(d, testClusterReader{kubeconfig: testClusterContextKubeconfig}); err != nil {
				t.Fatal(err)
			}
			for key, want := range map[string]interface{}{
				"context_name": "kind-dev-team", "kubeconfig": testClusterContextKubeconfig,
				"endpoint": "https://127.0.0.1:6443", "client_certificate": "cert",
				"client_key": "key", "cluster_ca_certificate": "ca", "completed": true,
				"kubeconfig_path": "/existing/config",
			} {
				if got := d.Get(key); got != want {
					t.Errorf("%s = %v, want %v", key, got, want)
				}
			}
			if d.Id() != "existing-cluster-id" {
				t.Errorf("refresh changed ID to %q", d.Id())
			}
		})
	}
}

func TestClusterReadContextNameErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		reader testClusterReader
		wantID string
	}{
		{"fetch", testClusterReader{err: fmt.Errorf("fetch failed")}, ""},
		{"malformed", testClusterReader{kubeconfig: "invalid: ["}, "existing-id"},
		{"missing current context", testClusterReader{kubeconfig: "apiVersion: v1\nkind: Config\n"}, "existing-id"},
		{"unknown current context", testClusterReader{kubeconfig: "apiVersion: v1\nkind: Config\ncurrent-context: missing\n"}, "existing-id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceCluster().Schema, map[string]interface{}{
				"name": "dev-team", "kubeconfig_path": "/existing/config",
			})
			d.SetId("existing-id")
			if err := resourceKindClusterReadWithProvider(d, tc.reader); err == nil {
				t.Fatal("expected read error")
			}
			if d.Id() != tc.wantID {
				t.Errorf("ID = %q, want %q", d.Id(), tc.wantID)
			}
			if got := d.Get("context_name"); got != "" {
				t.Errorf("context_name populated on failure: %q", got)
			}
		})
	}
}
