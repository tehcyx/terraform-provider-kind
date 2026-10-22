terraform {
  required_version = ">= 1.7.0"

  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.0.19"
    }
  }
}

provider "kind" {}

variable "clusters" {
  description = "Cluster keys and stable, unique port slots (0-99). Keep slots unchanged when adding or removing clusters."
  type        = map(number)
  default = {
    east = 0
    west = 1
  }

  validation {
    condition     = alltrue([for name in keys(var.clusters) : can(regex("^[a-z0-9]+(-[a-z0-9]+)*$", name))])
    error_message = "Cluster keys must contain lowercase letters, digits, or single hyphens between segments."
  }

  validation {
    condition     = alltrue([for slot in values(var.clusters) : slot != null ? slot >= 0 && slot <= 99 && slot == floor(slot) : false])
    error_message = "Port slots must be whole numbers from 0 through 99."
  }

  validation {
    condition     = length(distinct(values(var.clusters))) == length(var.clusters)
    error_message = "Each cluster must have a unique port slot to avoid host port conflicts."
  }
}

resource "kind_cluster" "cluster" {
  for_each = var.clusters

  name            = "multi-${each.key}"
  kubeconfig_path = abspath("${path.module}/kubeconfig-${each.key}")
  wait_for_ready  = true

  kind_config {
    kind        = "Cluster"
    api_version = "kind.x-k8s.io/v1alpha4"

    # Leave api_server_port unset so kind selects an available host port.
    networking {
      api_server_address = "127.0.0.1"
    }

    node {
      role = "control-plane"

      # Match these container ports to your application's Service nodePorts.
      extra_port_mappings {
        container_port = 30080
        host_port      = 20000 + 2 * each.value
        listen_address = "127.0.0.1"
        protocol       = "TCP"
      }
      extra_port_mappings {
        container_port = 30443
        host_port      = 20001 + 2 * each.value
        listen_address = "127.0.0.1"
        protocol       = "TCP"
      }
    }
  }
}

output "clusters" {
  description = "Connection details and mapped application ports, keyed by the stable cluster key."
  value = {
    for key, cluster in kind_cluster.cluster : key => {
      name            = cluster.name
      kubeconfig_path = cluster.kubeconfig_path
      endpoint        = cluster.endpoint
      http_port       = cluster.kind_config[0].node[0].extra_port_mappings[0].host_port
      https_port      = cluster.kind_config[0].node[0].extra_port_mappings[1].host_port
    }
  }
}
