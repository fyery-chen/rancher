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
	BusinessRoleTemplateBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BusinessRoleTemplateBinding",
	}
	BusinessRoleTemplateBindingResource = metav1.APIResource{
		Name:         "businessroletemplatebindings",
		SingularName: "businessroletemplatebinding",
		Namespaced:   true,

		Kind: BusinessRoleTemplateBindingGroupVersionKind.Kind,
	}
)

type BusinessRoleTemplateBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessRoleTemplateBinding
}

type BusinessRoleTemplateBindingHandlerFunc func(key string, obj *BusinessRoleTemplateBinding) error

type BusinessRoleTemplateBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*BusinessRoleTemplateBinding, err error)
	Get(namespace, name string) (*BusinessRoleTemplateBinding, error)
}

type BusinessRoleTemplateBindingController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessRoleTemplateBindingLister
	AddHandler(name string, handler BusinessRoleTemplateBindingHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessRoleTemplateBindingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessRoleTemplateBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessRoleTemplateBinding, error)
	Get(name string, opts metav1.GetOptions) (*BusinessRoleTemplateBinding, error)
	Update(*BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessRoleTemplateBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessRoleTemplateBindingController
	AddHandler(name string, sync BusinessRoleTemplateBindingHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessRoleTemplateBindingLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessRoleTemplateBindingHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessRoleTemplateBindingLifecycle)
}

type businessRoleTemplateBindingLister struct {
	controller *businessRoleTemplateBindingController
}

func (l *businessRoleTemplateBindingLister) List(namespace string, selector labels.Selector) (ret []*BusinessRoleTemplateBinding, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BusinessRoleTemplateBinding))
	})
	return
}

func (l *businessRoleTemplateBindingLister) Get(namespace, name string) (*BusinessRoleTemplateBinding, error) {
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
			Group:    BusinessRoleTemplateBindingGroupVersionKind.Group,
			Resource: "businessRoleTemplateBinding",
		}, name)
	}
	return obj.(*BusinessRoleTemplateBinding), nil
}

type businessRoleTemplateBindingController struct {
	controller.GenericController
}

func (c *businessRoleTemplateBindingController) Lister() BusinessRoleTemplateBindingLister {
	return &businessRoleTemplateBindingLister{
		controller: c,
	}
}

func (c *businessRoleTemplateBindingController) AddHandler(name string, handler BusinessRoleTemplateBindingHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BusinessRoleTemplateBinding))
	})
}

func (c *businessRoleTemplateBindingController) AddClusterScopedHandler(name, cluster string, handler BusinessRoleTemplateBindingHandlerFunc) {
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

		return handler(key, obj.(*BusinessRoleTemplateBinding))
	})
}

type businessRoleTemplateBindingFactory struct {
}

func (c businessRoleTemplateBindingFactory) Object() runtime.Object {
	return &BusinessRoleTemplateBinding{}
}

func (c businessRoleTemplateBindingFactory) List() runtime.Object {
	return &BusinessRoleTemplateBindingList{}
}

func (s *businessRoleTemplateBindingClient) Controller() BusinessRoleTemplateBindingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessRoleTemplateBindingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessRoleTemplateBindingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessRoleTemplateBindingController{
		GenericController: genericController,
	}

	s.client.businessRoleTemplateBindingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessRoleTemplateBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessRoleTemplateBindingController
}

func (s *businessRoleTemplateBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessRoleTemplateBindingClient) Create(o *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BusinessRoleTemplateBinding), err
}

func (s *businessRoleTemplateBindingClient) Get(name string, opts metav1.GetOptions) (*BusinessRoleTemplateBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BusinessRoleTemplateBinding), err
}

func (s *businessRoleTemplateBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessRoleTemplateBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BusinessRoleTemplateBinding), err
}

func (s *businessRoleTemplateBindingClient) Update(o *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BusinessRoleTemplateBinding), err
}

func (s *businessRoleTemplateBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessRoleTemplateBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessRoleTemplateBindingClient) List(opts metav1.ListOptions) (*BusinessRoleTemplateBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessRoleTemplateBindingList), err
}

func (s *businessRoleTemplateBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessRoleTemplateBindingClient) Patch(o *BusinessRoleTemplateBinding, data []byte, subresources ...string) (*BusinessRoleTemplateBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BusinessRoleTemplateBinding), err
}

func (s *businessRoleTemplateBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessRoleTemplateBindingClient) AddHandler(name string, sync BusinessRoleTemplateBindingHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessRoleTemplateBindingClient) AddLifecycle(name string, lifecycle BusinessRoleTemplateBindingLifecycle) {
	sync := NewBusinessRoleTemplateBindingLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessRoleTemplateBindingClient) AddClusterScopedHandler(name, clusterName string, sync BusinessRoleTemplateBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessRoleTemplateBindingClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessRoleTemplateBindingLifecycle) {
	sync := NewBusinessRoleTemplateBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
