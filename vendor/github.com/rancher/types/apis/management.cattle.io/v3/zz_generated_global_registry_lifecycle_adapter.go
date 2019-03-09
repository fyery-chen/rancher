package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type GlobalRegistryLifecycle interface {
	Create(obj *GlobalRegistry) (runtime.Object, error)
	Remove(obj *GlobalRegistry) (runtime.Object, error)
	Updated(obj *GlobalRegistry) (runtime.Object, error)
}

type globalRegistryLifecycleAdapter struct {
	lifecycle GlobalRegistryLifecycle
}

func (w *globalRegistryLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *globalRegistryLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *globalRegistryLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*GlobalRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRegistryLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*GlobalRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *globalRegistryLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*GlobalRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGlobalRegistryLifecycleAdapter(name string, clusterScoped bool, client GlobalRegistryInterface, l GlobalRegistryLifecycle) GlobalRegistryHandlerFunc {
	adapter := &globalRegistryLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *GlobalRegistry) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
