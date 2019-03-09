package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type ProjectRegistryLifecycle interface {
	Create(obj *ProjectRegistry) (runtime.Object, error)
	Remove(obj *ProjectRegistry) (runtime.Object, error)
	Updated(obj *ProjectRegistry) (runtime.Object, error)
}

type projectRegistryLifecycleAdapter struct {
	lifecycle ProjectRegistryLifecycle
}

func (w *projectRegistryLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *projectRegistryLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *projectRegistryLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*ProjectRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectRegistryLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*ProjectRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *projectRegistryLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*ProjectRegistry))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewProjectRegistryLifecycleAdapter(name string, clusterScoped bool, client ProjectRegistryInterface, l ProjectRegistryLifecycle) ProjectRegistryHandlerFunc {
	adapter := &projectRegistryLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *ProjectRegistry) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
