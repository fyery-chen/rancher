package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ClusterRegistryLifecycle interface {
	Create(obj *ClusterRegistry) (runtime.Object, error)
	Remove(obj *ClusterRegistry) (runtime.Object, error)
	Updated(obj *ClusterRegistry) (runtime.Object, error)
}

type clusterRegistryLifecycleAdapter struct {
	lifecycle ClusterRegistryLifecycle
}

func (w *clusterRegistryLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *clusterRegistryLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *clusterRegistryLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*ClusterRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRegistryLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*ClusterRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *clusterRegistryLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*ClusterRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewClusterRegistryLifecycleAdapter(name string, clusterScoped bool, client ClusterRegistryInterface, l ClusterRegistryLifecycle) ClusterRegistryHandlerFunc {
	adapter := &clusterRegistryLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *ClusterRegistry) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
