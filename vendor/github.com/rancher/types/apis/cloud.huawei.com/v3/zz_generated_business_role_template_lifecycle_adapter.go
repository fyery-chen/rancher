package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessRoleTemplateLifecycle interface {
	Create(obj *BusinessRoleTemplate) (*BusinessRoleTemplate, error)
	Remove(obj *BusinessRoleTemplate) (*BusinessRoleTemplate, error)
	Updated(obj *BusinessRoleTemplate) (*BusinessRoleTemplate, error)
}

type businessRoleTemplateLifecycleAdapter struct {
	lifecycle BusinessRoleTemplateLifecycle
}

func (w *businessRoleTemplateLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BusinessRoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessRoleTemplateLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BusinessRoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessRoleTemplateLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BusinessRoleTemplate))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessRoleTemplateLifecycleAdapter(name string, clusterScoped bool, client BusinessRoleTemplateInterface, l BusinessRoleTemplateLifecycle) BusinessRoleTemplateHandlerFunc {
	adapter := &businessRoleTemplateLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BusinessRoleTemplate) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
