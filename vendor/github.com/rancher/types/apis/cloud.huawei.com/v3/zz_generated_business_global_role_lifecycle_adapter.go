package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessGlobalRoleLifecycle interface {
	Create(obj *BusinessGlobalRole) (*BusinessGlobalRole, error)
	Remove(obj *BusinessGlobalRole) (*BusinessGlobalRole, error)
	Updated(obj *BusinessGlobalRole) (*BusinessGlobalRole, error)
}

type businessGlobalRoleLifecycleAdapter struct {
	lifecycle BusinessGlobalRoleLifecycle
}

func (w *businessGlobalRoleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BusinessGlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessGlobalRoleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BusinessGlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessGlobalRoleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BusinessGlobalRole))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessGlobalRoleLifecycleAdapter(name string, clusterScoped bool, client BusinessGlobalRoleInterface, l BusinessGlobalRoleLifecycle) BusinessGlobalRoleHandlerFunc {
	adapter := &businessGlobalRoleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BusinessGlobalRole) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
