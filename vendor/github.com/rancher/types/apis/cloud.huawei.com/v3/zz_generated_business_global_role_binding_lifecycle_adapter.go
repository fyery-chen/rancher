package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessGlobalRoleBindingLifecycle interface {
	Create(obj *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
	Remove(obj *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
	Updated(obj *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
}

type businessGlobalRoleBindingLifecycleAdapter struct {
	lifecycle BusinessGlobalRoleBindingLifecycle
}

func (w *businessGlobalRoleBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BusinessGlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessGlobalRoleBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BusinessGlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessGlobalRoleBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BusinessGlobalRoleBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessGlobalRoleBindingLifecycleAdapter(name string, clusterScoped bool, client BusinessGlobalRoleBindingInterface, l BusinessGlobalRoleBindingLifecycle) BusinessGlobalRoleBindingHandlerFunc {
	adapter := &businessGlobalRoleBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BusinessGlobalRoleBinding) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
