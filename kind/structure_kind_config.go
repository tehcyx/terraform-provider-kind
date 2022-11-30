package kind

import (
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

// Flatteners

func flattenKindConfig(d map[string]interface{}) *v1alpha4.Cluster {
	obj := &v1alpha4.Cluster{}

	obj.Kind = mapKeyIfExists(d, "kind").(string)
	obj.APIVersion = mapKeyIfExists(d, "api_version").(string)

	nodes := mapKeyIfExists(d, "node")
	if nodes != nil {
		for _, n := range nodes.([]interface{}) {
			data := n.(map[string]interface{})
			obj.Nodes = append(obj.Nodes, flattenKindConfigNodes(data))
		}
	}

	networking := mapKeyIfExists(d, "networking")
	if networking != nil {
		if n := networking.([]interface{}); len(n) == 1 { // MaxItems: 1, no more than one allowed so we don't have to loop here
			if n[0] != nil {
				data := n[0].(map[string]interface{})
				obj.Networking = flattenKindConfigNetworking(data)
			}
		}
	}

	containerdConfigPatches := mapKeyIfExists(d, "containerd_config_patches")
	if containerdConfigPatches != nil {
		for _, p := range containerdConfigPatches.([]interface{}) {
			patch := p.(string)
			obj.ContainerdConfigPatches = append(obj.ContainerdConfigPatches, patch)
		}
	}

	return obj
}

func flattenKindConfigNodes(d map[string]interface{}) v1alpha4.Node {
	obj := v1alpha4.Node{}

	role := mapKeyIfExists(d, "role")
	if role != nil && role.(string) != "" {
		switch role.(string) {
		case string(v1alpha4.ControlPlaneRole):
			obj.Role = v1alpha4.ControlPlaneRole
		case string(v1alpha4.WorkerRole):
			obj.Role = v1alpha4.WorkerRole
		}
	}
	image := mapKeyIfExists(d, "image")
	if image != nil && image.(string) != "" {
		obj.Image = image.(string)
	}

	extraMounts := mapKeyIfExists(d, "extra_mounts")
	if extraMounts != nil {
		for _, m := range extraMounts.([]interface{}) {
			data := m.(map[string]interface{})
			obj.ExtraMounts = append(obj.ExtraMounts, flattenKindConfigExtraMounts(data))
		}
	}

	labels := mapKeyIfExists(d, "labels")
	if labels != nil {
		if obj.Labels == nil {
			obj.Labels = make(map[string]string)
		}
		for labelName, labelValue := range labels.(map[string]interface{}) {
			if labelValue, ok := labelValue.(string); ok {
				obj.Labels[labelName] = labelValue
			}
		}
	}

	extraPortMappings := mapKeyIfExists(d, "extra_port_mappings")
	if extraPortMappings != nil {
		for _, m := range extraPortMappings.([]interface{}) {
			data := m.(map[string]interface{})
			obj.ExtraPortMappings = append(obj.ExtraPortMappings, flattenKindConfigExtraPortMappings(data))
		}
	}

	kubeadmConfigPatches := mapKeyIfExists(d, "kubeadm_config_patches")
	if kubeadmConfigPatches != nil {
		for _, k := range kubeadmConfigPatches.([]interface{}) {
			data := k.(string)
			obj.KubeadmConfigPatches = append(obj.KubeadmConfigPatches, data)
		}
	}

	return obj
}

func flattenKindConfigNetworking(d map[string]interface{}) v1alpha4.Networking {
	obj := v1alpha4.Networking{}

	apiServerAddress := mapKeyIfExists(d, "api_server_address")
	if apiServerAddress != nil && apiServerAddress.(string) != "" {
		obj.APIServerAddress = apiServerAddress.(string)
	}

	apiServerPort := mapKeyIfExists(d, "api_server_port")
	if apiServerPort != nil {
		obj.APIServerPort = int32(apiServerPort.(int))
	}

	disableDefaultCNI := mapKeyIfExists(d, "disable_default_cni")
	if disableDefaultCNI != nil {
		obj.DisableDefaultCNI = disableDefaultCNI.(bool)
	}

	ipFamily := mapKeyIfExists(d, "ip_family")
	if ipFamily != nil && ipFamily.(string) != "" {
		switch ipFamily.(string) {
		case string(v1alpha4.IPv4Family):
			obj.IPFamily = v1alpha4.IPv4Family
		case string(v1alpha4.IPv6Family):
			obj.IPFamily = v1alpha4.IPv6Family
		}
	}

	kubeProxyMode := mapKeyIfExists(d, "kube_proxy_mode")
	if kubeProxyMode != nil && kubeProxyMode.(string) != "" {
		switch kubeProxyMode.(string) {
		case string(v1alpha4.IPTablesProxyMode):
			obj.KubeProxyMode = v1alpha4.IPTablesProxyMode
		case string(v1alpha4.IPVSProxyMode):
			obj.KubeProxyMode = v1alpha4.IPVSProxyMode
		}
	}

	podSubnet := mapKeyIfExists(d, "pod_subnet")
	if podSubnet != nil && podSubnet.(string) != "" {
		obj.PodSubnet = podSubnet.(string)
	}

	serviceSubnet := mapKeyIfExists(d, "service_subnet")
	if serviceSubnet != nil && serviceSubnet.(string) != "" {
		obj.ServiceSubnet = serviceSubnet.(string)
	}

	return obj
}

func flattenKindConfigExtraMounts(d map[string]interface{}) v1alpha4.Mount {
	obj := v1alpha4.Mount{}

	containerPath := mapKeyIfExists(d, "container_path")
	if containerPath != nil && containerPath.(string) != "" {
		obj.ContainerPath = containerPath.(string)
	}
	hostPath := mapKeyIfExists(d, "host_path")
	if hostPath != nil && hostPath.(string) != "" {
		obj.HostPath = hostPath.(string)
	}
	propagation := mapKeyIfExists(d, "propagation")
	if propagation != nil && propagation.(string) != "" {
		switch propagation.(string) {
		case string(v1alpha4.MountPropagationBidirectional):
			obj.Propagation = v1alpha4.MountPropagationBidirectional
		case string(v1alpha4.MountPropagationHostToContainer):
			obj.Propagation = v1alpha4.MountPropagationHostToContainer
		case string(v1alpha4.MountPropagationNone):
			obj.Propagation = v1alpha4.MountPropagationNone
		}
	}

	readonly := mapKeyIfExists(d, "readonly")
	if readonly != nil {
		obj.Readonly = readonly.(bool)
	}
	selinuxRelabel := mapKeyIfExists(d, "selinux_relabel")
	if selinuxRelabel != nil {
		obj.SelinuxRelabel = selinuxRelabel.(bool)
	}

	return obj
}

func flattenKindConfigExtraPortMappings(d map[string]interface{}) v1alpha4.PortMapping {
	obj := v1alpha4.PortMapping{}

	containerPort := mapKeyIfExists(d, "container_port")
	if containerPort != nil {
		obj.ContainerPort = int32(containerPort.(int))
	}
	hostPort := mapKeyIfExists(d, "host_port")
	if hostPort != nil {
		obj.HostPort = int32(hostPort.(int))
	}
	listenAddress := mapKeyIfExists(d, "listen_address")
	if listenAddress != nil && listenAddress.(string) != "" {
		obj.ListenAddress = listenAddress.(string)
	}
	protocol := mapKeyIfExists(d, "protocol")
	if protocol != nil && protocol.(string) != "" {
		switch protocol.(string) {
		case string(v1alpha4.PortMappingProtocolSCTP):
			obj.Protocol = v1alpha4.PortMappingProtocolSCTP
		case string(v1alpha4.PortMappingProtocolTCP):
			obj.Protocol = v1alpha4.PortMappingProtocolTCP
		case string(v1alpha4.PortMappingProtocolUDP):
			obj.Protocol = v1alpha4.PortMappingProtocolUDP
		}
	}

	return obj
}

func mapKeyIfExists(m map[string]interface{}, key string) interface{} {
	if val, ok := m[key]; ok {
		return val
	}
	return nil
}
