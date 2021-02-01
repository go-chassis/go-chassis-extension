package kuberegistry

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-chassis/go-chassis/v2/pkg/util/tags"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDiscoveryController(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	client := fake.NewSimpleClientset()
	sharedfactory := informers.NewSharedInformerFactory(client, 0)
	sInformer := sharedfactory.Core().V1().Services()
	eInformer := sharedfactory.Core().V1().Endpoints()
	pInformer := sharedfactory.Core().V1().Pods()
	dc := NewDiscoveryController(sInformer, eInformer, pInformer, client)
	sharedfactory.Start(ctx.Done())
	dc.Run(ctx.Done())

	// create endpoints
	p := &v1.Endpoints{ObjectMeta: metav1.ObjectMeta{Name: "kubeserver"},
		Subsets: []v1.EndpointSubset{{
			Addresses: []v1.EndpointAddress{{IP: "127.0.0.1",
				TargetRef: &v1.ObjectReference{UID: "12345"}}},
			Ports: []v1.EndpointPort{{Name: "rest", Port: 9090}},
		}}}
	_, err := client.CoreV1().Endpoints("default").Create(ctx, p, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("error create endpoints: %v", err)
	}

	// create services
	s := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "kubeserver"}}
	_, err = client.CoreV1().Services("default").Create(ctx, p, metav1.CreateOptions{})
	if err != nil {
		t.Errorf("error create service: %v", err)
	}

	time.Sleep(1 * time.Second)
	ret, err := dc.FindEndpoints("kubeserver.default.svc.local", utiltags.Tags{})
	assert.Equal(t, 1, len(ret))
	for _, ep := range ret {
		log.Printf("Got endpoints %s(%s)", ep.EndpointsMap, ep.ServiceID)
		assert.Equal(t, ep.EndpointsMap["rest"], "127.0.0.1:9090")
		assert.Equal(t, ep.ServiceID, "kubeserver.default")
	}

	svc, err := dc.GetAllServices()
	assert.Equal(t, len(svc), 1)
	for _, ss := range svc {
		assert.Equal(t, ss.ServiceName, "kubeserver")
	}
}
