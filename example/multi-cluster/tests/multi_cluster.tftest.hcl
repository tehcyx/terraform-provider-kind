# Plan-only tests use the provider schema without starting containers.
mock_provider "kind" {}

run "default_clusters" {
  command = plan

  assert {
    condition     = length(kind_cluster.cluster) == 2 && kind_cluster.cluster["east"].name == "multi-east" && kind_cluster.cluster["west"].name == "multi-west"
    error_message = "Each key must create a separate, uniquely named cluster."
  }

  assert {
    condition = (
      kind_cluster.cluster["east"].kubeconfig_path == abspath("${path.module}/kubeconfig-east") &&
      kind_cluster.cluster["west"].kubeconfig_path == abspath("${path.module}/kubeconfig-west")
    )
    error_message = "Clusters must use separate absolute kubeconfig paths."
  }

  assert {
    condition = alltrue(flatten([
      for cluster in kind_cluster.cluster : [
        cluster.wait_for_ready,
        cluster.kind_config[0].networking[0].api_server_address == "127.0.0.1",
        cluster.kind_config[0].node[0].role == "control-plane",
        cluster.kind_config[0].node[0].extra_port_mappings[0].container_port == 30080,
        cluster.kind_config[0].node[0].extra_port_mappings[1].container_port == 30443,
        alltrue([for mapping in cluster.kind_config[0].node[0].extra_port_mappings : mapping.listen_address == "127.0.0.1" && mapping.protocol == "TCP"])
      ]
    ]))
    error_message = "Clones must share ready behavior, nodePorts, and loopback-only TCP mappings."
  }

  assert {
    condition = (
      output.clusters["east"].name == "multi-east" &&
      output.clusters["east"].kubeconfig_path == kind_cluster.cluster["east"].kubeconfig_path &&
      output.clusters["east"].http_port == 20000 && output.clusters["east"].https_port == 20001 &&
      output.clusters["west"].http_port == 20002 && output.clusters["west"].https_port == 20003
    )
    error_message = "Outputs must expose per-cluster paths and distinct, deterministic mapped ports."
  }
}

run "stable_slots" {
  command = plan
  variables {
    clusters = { west = 1, north = 99 }
  }

  assert {
    condition = (
      length(kind_cluster.cluster) == 2 &&
      output.clusters["west"].name == "multi-west" &&
      output.clusters["west"].http_port == 20002 && output.clusters["west"].https_port == 20003 &&
      output.clusters["north"].http_port == 20198 && output.clusters["north"].https_port == 20199
    )
    error_message = "Removing east and adding north must preserve west's identity and ports."
  }
}

run "duplicate_slots" {
  command = plan
  variables {
    clusters = { east = 0, west = 0 }
  }
  expect_failures = [var.clusters]
}

run "negative_slot" {
  command = plan
  variables {
    clusters = { east = -1 }
  }
  expect_failures = [var.clusters]
}

run "out_of_range_slot" {
  command = plan
  variables {
    clusters = { east = 100 }
  }
  expect_failures = [var.clusters]
}

run "fractional_slot" {
  command = plan
  variables {
    clusters = { east = 0.5 }
  }
  expect_failures = [var.clusters]
}

run "null_slot" {
  command = plan
  variables {
    clusters = { east = null }
  }
  expect_failures = [var.clusters]
}

run "invalid_name" {
  command = plan
  variables {
    clusters = { "East/cluster" = 0 }
  }
  expect_failures = [var.clusters]
}

run "empty_clusters" {
  command = plan
  variables {
    clusters = {}
  }

  assert {
    condition     = length(kind_cluster.cluster) == 0 && length(output.clusters) == 0
    error_message = "An empty map must plan no clusters and expose an empty output map."
  }
}
