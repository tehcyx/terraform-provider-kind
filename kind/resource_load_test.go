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

	schema := r.Schema
	if _, ok := schema["image"]; !ok {
		t.Error("schema should have 'image' field")
	}
	if _, ok := schema["cluster_name"]; !ok {
		t.Error("schema should have 'cluster_name' field")
	}

	if !schema["image"].Required {
		t.Error("'image' should be Required")
	}
	if !schema["image"].ForceNew {
		t.Error("'image' should be ForceNew")
	}
	if !schema["cluster_name"].Required {
		t.Error("'cluster_name' should be Required")
	}
	if !schema["cluster_name"].ForceNew {
		t.Error("'cluster_name' should be ForceNew")
	}
}
