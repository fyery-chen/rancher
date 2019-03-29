package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type IstioMonitorGraphLifecycle interface {
	Create(obj *IstioMonitorGraph) (runtime.Object, error)
	Remove(obj *IstioMonitorGraph) (runtime.Object, error)
	Updated(obj *IstioMonitorGraph) (runtime.Object, error)
}

type istioMonitorGraphLifecycleAdapter struct {
	lifecycle IstioMonitorGraphLifecycle
}

func (w *istioMonitorGraphLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *istioMonitorGraphLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *istioMonitorGraphLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*IstioMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioMonitorGraphLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*IstioMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *istioMonitorGraphLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*IstioMonitorGraph))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewIstioMonitorGraphLifecycleAdapter(name string, clusterScoped bool, client IstioMonitorGraphInterface, l IstioMonitorGraphLifecycle) IstioMonitorGraphHandlerFunc {
	adapter := &istioMonitorGraphLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *IstioMonitorGraph) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
