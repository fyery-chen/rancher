package v1

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
	BusinessQuotaGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BusinessQuota",
	}
	BusinessQuotaResource = metav1.APIResource{
		Name:         "businessquotas",
		SingularName: "businessquota",
		Namespaced:   true,

		Kind: BusinessQuotaGroupVersionKind.Kind,
	}
)

type BusinessQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessQuota
}

type BusinessQuotaHandlerFunc func(key string, obj *BusinessQuota) error

type BusinessQuotaLister interface {
	List(namespace string, selector labels.Selector) (ret []*BusinessQuota, err error)
	Get(namespace, name string) (*BusinessQuota, error)
}

type BusinessQuotaController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessQuotaLister
	AddHandler(name string, handler BusinessQuotaHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessQuotaHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessQuotaInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BusinessQuota) (*BusinessQuota, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessQuota, error)
	Get(name string, opts metav1.GetOptions) (*BusinessQuota, error)
	Update(*BusinessQuota) (*BusinessQuota, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessQuotaList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessQuotaController
	AddHandler(name string, sync BusinessQuotaHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessQuotaLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessQuotaHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessQuotaLifecycle)
}

type businessQuotaLister struct {
	controller *businessQuotaController
}

func (l *businessQuotaLister) List(namespace string, selector labels.Selector) (ret []*BusinessQuota, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BusinessQuota))
	})
	return
}

func (l *businessQuotaLister) Get(namespace, name string) (*BusinessQuota, error) {
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
			Group:    BusinessQuotaGroupVersionKind.Group,
			Resource: "businessQuota",
		}, name)
	}
	return obj.(*BusinessQuota), nil
}

type businessQuotaController struct {
	controller.GenericController
}

func (c *businessQuotaController) Lister() BusinessQuotaLister {
	return &businessQuotaLister{
		controller: c,
	}
}

func (c *businessQuotaController) AddHandler(name string, handler BusinessQuotaHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BusinessQuota))
	})
}

func (c *businessQuotaController) AddClusterScopedHandler(name, cluster string, handler BusinessQuotaHandlerFunc) {
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

		return handler(key, obj.(*BusinessQuota))
	})
}

type businessQuotaFactory struct {
}

func (c businessQuotaFactory) Object() runtime.Object {
	return &BusinessQuota{}
}

func (c businessQuotaFactory) List() runtime.Object {
	return &BusinessQuotaList{}
}

func (s *businessQuotaClient) Controller() BusinessQuotaController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessQuotaControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessQuotaGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessQuotaController{
		GenericController: genericController,
	}

	s.client.businessQuotaControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessQuotaClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessQuotaController
}

func (s *businessQuotaClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessQuotaClient) Create(o *BusinessQuota) (*BusinessQuota, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BusinessQuota), err
}

func (s *businessQuotaClient) Get(name string, opts metav1.GetOptions) (*BusinessQuota, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BusinessQuota), err
}

func (s *businessQuotaClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessQuota, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BusinessQuota), err
}

func (s *businessQuotaClient) Update(o *BusinessQuota) (*BusinessQuota, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BusinessQuota), err
}

func (s *businessQuotaClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessQuotaClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessQuotaClient) List(opts metav1.ListOptions) (*BusinessQuotaList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessQuotaList), err
}

func (s *businessQuotaClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessQuotaClient) Patch(o *BusinessQuota, data []byte, subresources ...string) (*BusinessQuota, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BusinessQuota), err
}

func (s *businessQuotaClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessQuotaClient) AddHandler(name string, sync BusinessQuotaHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessQuotaClient) AddLifecycle(name string, lifecycle BusinessQuotaLifecycle) {
	sync := NewBusinessQuotaLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessQuotaClient) AddClusterScopedHandler(name, clusterName string, sync BusinessQuotaHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessQuotaClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessQuotaLifecycle) {
	sync := NewBusinessQuotaLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
