package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessLifecycle interface {
	Create(obj *Business) (*Business, error)
	Remove(obj *Business) (*Business, error)
	Updated(obj *Business) (*Business, error)
}

type businessLifecycleAdapter struct {
	lifecycle BusinessLifecycle
}

func (w *businessLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Business))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Business))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Business))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessLifecycleAdapter(name string, clusterScoped bool, client BusinessInterface, l BusinessLifecycle) BusinessHandlerFunc {
	adapter := &businessLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Business) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
