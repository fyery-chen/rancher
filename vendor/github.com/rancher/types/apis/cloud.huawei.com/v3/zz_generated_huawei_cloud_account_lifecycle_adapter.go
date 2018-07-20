package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type HuaweiCloudAccountLifecycle interface {
	Create(obj *HuaweiCloudAccount) (*HuaweiCloudAccount, error)
	Remove(obj *HuaweiCloudAccount) (*HuaweiCloudAccount, error)
	Updated(obj *HuaweiCloudAccount) (*HuaweiCloudAccount, error)
}

type huaweiCloudAccountLifecycleAdapter struct {
	lifecycle HuaweiCloudAccountLifecycle
}

func (w *huaweiCloudAccountLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*HuaweiCloudAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *huaweiCloudAccountLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*HuaweiCloudAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *huaweiCloudAccountLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*HuaweiCloudAccount))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewHuaweiCloudAccountLifecycleAdapter(name string, clusterScoped bool, client HuaweiCloudAccountInterface, l HuaweiCloudAccountLifecycle) HuaweiCloudAccountHandlerFunc {
	adapter := &huaweiCloudAccountLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *HuaweiCloudAccount) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
