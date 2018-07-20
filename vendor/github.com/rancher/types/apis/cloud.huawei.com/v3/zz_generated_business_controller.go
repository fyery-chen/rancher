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
	BusinessGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Business",
	}
	BusinessResource = metav1.APIResource{
		Name:         "businesses",
		SingularName: "business",
		Namespaced:   false,
		Kind:         BusinessGroupVersionKind.Kind,
	}
)

type BusinessList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Business
}

type BusinessHandlerFunc func(key string, obj *Business) error

type BusinessLister interface {
	List(namespace string, selector labels.Selector) (ret []*Business, err error)
	Get(namespace, name string) (*Business, error)
}

type BusinessController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessLister
	AddHandler(name string, handler BusinessHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*Business) (*Business, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Business, error)
	Get(name string, opts metav1.GetOptions) (*Business, error)
	Update(*Business) (*Business, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessController
	AddHandler(name string, sync BusinessHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessLifecycle)
}

type businessLister struct {
	controller *businessController
}

func (l *businessLister) List(namespace string, selector labels.Selector) (ret []*Business, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Business))
	})
	return
}

func (l *businessLister) Get(namespace, name string) (*Business, error) {
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
			Group:    BusinessGroupVersionKind.Group,
			Resource: "business",
		}, name)
	}
	return obj.(*Business), nil
}

type businessController struct {
	controller.GenericController
}

func (c *businessController) Lister() BusinessLister {
	return &businessLister{
		controller: c,
	}
}

func (c *businessController) AddHandler(name string, handler BusinessHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*Business))
	})
}

func (c *businessController) AddClusterScopedHandler(name, cluster string, handler BusinessHandlerFunc) {
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

		return handler(key, obj.(*Business))
	})
}

type businessFactory struct {
}

func (c businessFactory) Object() runtime.Object {
	return &Business{}
}

func (c businessFactory) List() runtime.Object {
	return &BusinessList{}
}

func (s *businessClient) Controller() BusinessController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessController{
		GenericController: genericController,
	}

	s.client.businessControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessController
}

func (s *businessClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessClient) Create(o *Business) (*Business, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Business), err
}

func (s *businessClient) Get(name string, opts metav1.GetOptions) (*Business, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Business), err
}

func (s *businessClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Business, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Business), err
}

func (s *businessClient) Update(o *Business) (*Business, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Business), err
}

func (s *businessClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessClient) List(opts metav1.ListOptions) (*BusinessList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessList), err
}

func (s *businessClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessClient) Patch(o *Business, data []byte, subresources ...string) (*Business, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*Business), err
}

func (s *businessClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessClient) AddHandler(name string, sync BusinessHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessClient) AddLifecycle(name string, lifecycle BusinessLifecycle) {
	sync := NewBusinessLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessClient) AddClusterScopedHandler(name, clusterName string, sync BusinessHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessLifecycle) {
	sync := NewBusinessLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
