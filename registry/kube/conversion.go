package kuberegistry

import (
	"strconv"

	"github.com/go-chassis/go-chassis/core/registry"
	v1 "k8s.io/api/core/v1"
)

func toMicroService(ss *v1.Service) *registry.MicroService {
	return &registry.MicroService{
		ServiceName: ss.Name,
		ServiceID:   string(ss.UID),
		Metadata:    ss.Spec.Selector,
		RegisterBy:  KubeRegistry,
	}
}

func toProtocolMap(address v1.EndpointAddress, ports []v1.EndpointPort) map[string]string {
	ret := map[string]string{}
	for _, port := range ports {
		if _, ok := ret[port.Name]; !ok {
			ret[port.Name] = address.IP + ":" + strconv.Itoa(int(port.Port))
			continue
		}
	}
	return ret
}
