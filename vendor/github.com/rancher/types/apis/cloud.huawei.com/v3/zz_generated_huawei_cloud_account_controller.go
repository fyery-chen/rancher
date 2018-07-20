package v3

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	HuaweiCloudAccountGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "HuaweiCloudAccount",
	}
	HuaweiCloudAccountResource = metav1.APIResource{
		Name:         "huaweicloudaccounts",
		SingularName: "huaweicloudaccount",
		Namespaced:   true,

		Kind: HuaweiCloudAccountGroupVersionKind.Kind,
	}
)

type HuaweiCloudAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HuaweiCloudAccount
}

type HuaweiCloudAccountHandlerFunc func(key string, obj *HuaweiCloudAccount) error

type HuaweiCloudAccountLister interface {
	List(namespace string, selector labels.Selector) (ret []*HuaweiCloudAccount, err error)
	Get(namespace, name string) (*HuaweiCloudAccount, error)
}

type HuaweiCloudAccountController interface {
	Informer() cache.SharedIndexInformer
	Lister() HuaweiCloudAccountLister
	AddHandler(name string, handler HuaweiCloudAccountHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler HuaweiCloudAccountHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type HuaweiCloudAccountInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*HuaweiCloudAccount) (*HuaweiCloudAccount, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*HuaweiCloudAccount, error)
	Get(name string, opts metav1.GetOptions) (*HuaweiCloudAccount, error)
	Update(*HuaweiCloudAccount) (*HuaweiCloudAccount, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*HuaweiCloudAccountList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() HuaweiCloudAccountController
	AddHandler(name string, sync HuaweiCloudAccountHandlerFunc)
	AddLifecycle(name string, lifecycle HuaweiCloudAccountLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync HuaweiCloudAccountHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle HuaweiCloudAccountLifecycle)
}

type huaweiCloudAccountLister struct {
	controller *huaweiCloudAccountController
}

func (l *huaweiCloudAccountLister) List(namespace string, selector labels.Selector) (ret []*HuaweiCloudAccount, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*HuaweiCloudAccount))
	})
	return
}

func (l *huaweiCloudAccountLister) Get(namespace, name string) (*HuaweiCloudAccount, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    HuaweiCloudAccountGroupVersionKind.Group,
			Resource: "huaweiCloudAccount",
		}, name)
	}
	return obj.(*HuaweiCloudAccount), nil
}

type huaweiCloudAccountController struct {
	controller.GenericController
}

func (c *huaweiCloudAccountController) Lister() HuaweiCloudAccountLister {
	return &huaweiCloudAccountLister{
		controller: c,
	}
}

func (c *huaweiCloudAccountController) AddHandler(name string, handler HuaweiCloudAccountHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*HuaweiCloudAccount))
	})
}

func (c *huaweiCloudAccountController) AddClusterScopedHandler(name, cluster string, handler HuaweiCloudAccountHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}

		if !controller.ObjectInCluster(cluster, obj) {
			return nil
		}

		return handler(key, obj.(*HuaweiCloudAccount))
	})
}

type huaweiCloudAccountFactory struct {
}

func (c huaweiCloudAccountFactory) Object() runtime.Object {
	return &HuaweiCloudAccount{}
}

func (c huaweiCloudAccountFactory) List() runtime.Object {
	return &HuaweiCloudAccountList{}
}

func (s *huaweiCloudAccountClient) Controller() HuaweiCloudAccountController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.huaweiCloudAccountControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(HuaweiCloudAccountGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &huaweiCloudAccountController{
		GenericController: genericController,
	}

	s.client.huaweiCloudAccountControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type huaweiCloudAccountClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   HuaweiCloudAccountController
}

func (s *huaweiCloudAccountClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *huaweiCloudAccountClient) Create(o *HuaweiCloudAccount) (*HuaweiCloudAccount, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*HuaweiCloudAccount), err
}

func (s *huaweiCloudAccountClient) Get(name string, opts metav1.GetOptions) (*HuaweiCloudAccount, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*HuaweiCloudAccount), err
}

func (s *huaweiCloudAccountClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*HuaweiCloudAccount, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*HuaweiCloudAccount), err
}

func (s *huaweiCloudAccountClient) Update(o *HuaweiCloudAccount) (*HuaweiCloudAccount, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*HuaweiCloudAccount), err
}

func (s *huaweiCloudAccountClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *huaweiCloudAccountClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *huaweiCloudAccountClient) List(opts metav1.ListOptions) (*HuaweiCloudAccountList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*HuaweiCloudAccountList), err
}

func (s *huaweiCloudAccountClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *huaweiCloudAccountClient) Patch(o *HuaweiCloudAccount, data []byte, subresources ...string) (*HuaweiCloudAccount, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*HuaweiCloudAccount), err
}

func (s *huaweiCloudAccountClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *huaweiCloudAccountClient) AddHandler(name string, sync HuaweiCloudAccountHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *huaweiCloudAccountClient) AddLifecycle(name string, lifecycle HuaweiCloudAccountLifecycle) {
	sync := NewHuaweiCloudAccountLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *huaweiCloudAccountClient) AddClusterScopedHandler(name, clusterName string, sync HuaweiCloudAccountHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *huaweiCloudAccountClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle HuaweiCloudAccountLifecycle) {
	sync := NewHuaweiCloudAccountLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
