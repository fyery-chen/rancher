package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessRoleTemplateBindingLifecycle interface {
	Create(obj *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
	Remove(obj *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
	Updated(obj *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
}

type businessRoleTemplateBindingLifecycleAdapter struct {
	lifecycle BusinessRoleTemplateBindingLifecycle
}

func (w *businessRoleTemplateBindingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BusinessRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessRoleTemplateBindingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BusinessRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessRoleTemplateBindingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BusinessRoleTemplateBinding))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessRoleTemplateBindingLifecycleAdapter(name string, clusterScoped bool, client BusinessRoleTemplateBindingInterface, l BusinessRoleTemplateBindingLifecycle) BusinessRoleTemplateBindingHandlerFunc {
	adapter := &businessRoleTemplateBindingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BusinessRoleTemplateBinding) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
