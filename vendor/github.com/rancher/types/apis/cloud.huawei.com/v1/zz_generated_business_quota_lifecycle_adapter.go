package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type BusinessQuotaLifecycle interface {
	Create(obj *BusinessQuota) (*BusinessQuota, error)
	Remove(obj *BusinessQuota) (*BusinessQuota, error)
	Updated(obj *BusinessQuota) (*BusinessQuota, error)
}

type businessQuotaLifecycleAdapter struct {
	lifecycle BusinessQuotaLifecycle
}

func (w *businessQuotaLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*BusinessQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessQuotaLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*BusinessQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *businessQuotaLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*BusinessQuota))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewBusinessQuotaLifecycleAdapter(name string, clusterScoped bool, client BusinessQuotaInterface, l BusinessQuotaLifecycle) BusinessQuotaHandlerFunc {
	adapter := &businessQuotaLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *BusinessQuota) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
