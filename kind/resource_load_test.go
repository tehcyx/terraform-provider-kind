package kind

import (
	"testing"
)

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

func TestResourceLoadCreate_ClusterNotFound(t *testing.T) {
	r := resourceLoad()
	d := r.TestResourceData()
	d.Set("image", "alpine")
	d.Set("cluster_name", "nonexistent-cluster-xyz")

	err := resourceKindLoadCreate(d, nil)
	if err == nil {
		t.Fatal("expected error for nonexistent cluster")
	}
}
