package kuberegistry

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	utiltags "github.com/go-chassis/go-chassis/pkg/util/tags"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// DiscoveryController defines discovery controller for kube registry
type DiscoveryController struct {
	client kubernetes.Interface

	sLister corelisters.ServiceLister
	eLister corelisters.EndpointsLister
	pLister corelisters.PodLister

	sListerSynced cache.InformerSynced
	eListerSynced cache.InformerSynced
	pListerSynced cache.InformerSynced
}

// NewDiscoveryController returns new discovery controller
func NewDiscoveryController(
	sInformer coreinformers.ServiceInformer,
	eInformer coreinformers.EndpointsInformer,
	pInformer coreinformers.PodInformer,
	client kubernetes.Interface,
) *DiscoveryController {

	dc := &DiscoveryController{
		client: client,
	}

	sInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: dc.addService,
	})
	eInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: dc.addEndpoints,
	})
	dc.sListerSynced = sInformer.Informer().HasSynced
	dc.eListerSynced = eInformer.Informer().HasSynced
	dc.pListerSynced = pInformer.Informer().HasSynced

	dc.pLister = pInformer.Lister()
	dc.sLister = sInformer.Lister()
	dc.eLister = eInformer.Lister()
	return dc
}

// Run begins discovery controller
func (dc *DiscoveryController) Run(stop <-chan struct{}) {
	lager.Logger.Info("Starting Discovery Controller")
	if !cache.WaitForCacheSync(stop, dc.sListerSynced, dc.eListerSynced, dc.pListerSynced) {
		lager.Logger.Error("Time out waiting for caches to sync", nil)
		return
	}
	lager.Logger.Info("Finish Waiting For Cache Sync")
}

func (dc *DiscoveryController) addService(obj interface{}) {
	svc := obj.(*v1.Service)
	lager.Logger.Infof("Add Service: %s", svc.Name)
}

func (dc *DiscoveryController) addEndpoints(obj interface{}) {
	ep := obj.(*v1.Endpoints)
	lager.Logger.Infof("Add Endpoint: %s", ep.Name)
}

// FindEndpoints returns microservice instances of kube registry
func (dc *DiscoveryController) FindEndpoints(service string, tags utiltags.Tags) ([]*registry.MicroServiceInstance, error) {
	// TODO: use labels.ToLabelSelector to trans endpoint
	// use cache lister to get specific endpoints or use kubeclient instead
	name, namespace := splitServiceKey(service)
	ep, err := dc.eLister.Endpoints(namespace).Get(name)
	if err != nil {
		return nil, err
	}

	ins := []*registry.MicroServiceInstance{}
	for _, ss := range ep.Subsets {
		for _, as := range ss.Addresses {
			pod, err := dc.pLister.Pods(as.TargetRef.Namespace).Get(as.TargetRef.Name)
			if err != nil {
				lager.Logger.Warnf("error list pods: %s", as.TargetRef.Name)
				continue
			}
			if !tags.IsSubsetOf(pod.Labels) {
				continue
			}

			ins = append(ins, &registry.MicroServiceInstance{
				InstanceID:   string(pod.UID),
				ServiceID:    ep.Name + "." + ep.Namespace,
				HostName:     as.Hostname,
				EndpointsMap: toProtocolMap(as, ss.Ports),
			})
		}
	}
	return ins, nil
}

// GetAllServices returns microservice of kube registry
func (dc *DiscoveryController) GetAllServices() ([]*registry.MicroService, error) {
	microServices, err := dc.sLister.List(labels.Everything())
	if err != nil {
		lager.Logger.Errorf(err, "get all microservices from kube failed")
		return nil, err
	}
	ms := make([]*registry.MicroService, len(microServices))
	for i, s := range microServices {
		ms[i] = toMicroService(s)
	}
	lager.Logger.Debugf("get all microservices success, microservices: %v", microServices)
	return ms, nil
}
