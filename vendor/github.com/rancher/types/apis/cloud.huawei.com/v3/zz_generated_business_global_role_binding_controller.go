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
	BusinessGlobalRoleBindingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BusinessGlobalRoleBinding",
	}
	BusinessGlobalRoleBindingResource = metav1.APIResource{
		Name:         "businessglobalrolebindings",
		SingularName: "businessglobalrolebinding",
		Namespaced:   false,
		Kind:         BusinessGlobalRoleBindingGroupVersionKind.Kind,
	}
)

type BusinessGlobalRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessGlobalRoleBinding
}

type BusinessGlobalRoleBindingHandlerFunc func(key string, obj *BusinessGlobalRoleBinding) error

type BusinessGlobalRoleBindingLister interface {
	List(namespace string, selector labels.Selector) (ret []*BusinessGlobalRoleBinding, err error)
	Get(namespace, name string) (*BusinessGlobalRoleBinding, error)
}

type BusinessGlobalRoleBindingController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessGlobalRoleBindingLister
	AddHandler(name string, handler BusinessGlobalRoleBindingHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessGlobalRoleBindingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessGlobalRoleBindingInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessGlobalRoleBinding, error)
	Get(name string, opts metav1.GetOptions) (*BusinessGlobalRoleBinding, error)
	Update(*BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessGlobalRoleBindingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessGlobalRoleBindingController
	AddHandler(name string, sync BusinessGlobalRoleBindingHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessGlobalRoleBindingLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessGlobalRoleBindingHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessGlobalRoleBindingLifecycle)
}

type businessGlobalRoleBindingLister struct {
	controller *businessGlobalRoleBindingController
}

func (l *businessGlobalRoleBindingLister) List(namespace string, selector labels.Selector) (ret []*BusinessGlobalRoleBinding, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BusinessGlobalRoleBinding))
	})
	return
}

func (l *businessGlobalRoleBindingLister) Get(namespace, name string) (*BusinessGlobalRoleBinding, error) {
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
			Group:    BusinessGlobalRoleBindingGroupVersionKind.Group,
			Resource: "businessGlobalRoleBinding",
		}, name)
	}
	return obj.(*BusinessGlobalRoleBinding), nil
}

type businessGlobalRoleBindingController struct {
	controller.GenericController
}

func (c *businessGlobalRoleBindingController) Lister() BusinessGlobalRoleBindingLister {
	return &businessGlobalRoleBindingLister{
		controller: c,
	}
}

func (c *businessGlobalRoleBindingController) AddHandler(name string, handler BusinessGlobalRoleBindingHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BusinessGlobalRoleBinding))
	})
}

func (c *businessGlobalRoleBindingController) AddClusterScopedHandler(name, cluster string, handler BusinessGlobalRoleBindingHandlerFunc) {
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

		return handler(key, obj.(*BusinessGlobalRoleBinding))
	})
}

type businessGlobalRoleBindingFactory struct {
}

func (c businessGlobalRoleBindingFactory) Object() runtime.Object {
	return &BusinessGlobalRoleBinding{}
}

func (c businessGlobalRoleBindingFactory) List() runtime.Object {
	return &BusinessGlobalRoleBindingList{}
}

func (s *businessGlobalRoleBindingClient) Controller() BusinessGlobalRoleBindingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessGlobalRoleBindingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessGlobalRoleBindingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessGlobalRoleBindingController{
		GenericController: genericController,
	}

	s.client.businessGlobalRoleBindingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessGlobalRoleBindingClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessGlobalRoleBindingController
}

func (s *businessGlobalRoleBindingClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessGlobalRoleBindingClient) Create(o *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BusinessGlobalRoleBinding), err
}

func (s *businessGlobalRoleBindingClient) Get(name string, opts metav1.GetOptions) (*BusinessGlobalRoleBinding, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BusinessGlobalRoleBinding), err
}

func (s *businessGlobalRoleBindingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessGlobalRoleBinding, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BusinessGlobalRoleBinding), err
}

func (s *businessGlobalRoleBindingClient) Update(o *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BusinessGlobalRoleBinding), err
}

func (s *businessGlobalRoleBindingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessGlobalRoleBindingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessGlobalRoleBindingClient) List(opts metav1.ListOptions) (*BusinessGlobalRoleBindingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessGlobalRoleBindingList), err
}

func (s *businessGlobalRoleBindingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessGlobalRoleBindingClient) Patch(o *BusinessGlobalRoleBinding, data []byte, subresources ...string) (*BusinessGlobalRoleBinding, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BusinessGlobalRoleBinding), err
}

func (s *businessGlobalRoleBindingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessGlobalRoleBindingClient) AddHandler(name string, sync BusinessGlobalRoleBindingHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessGlobalRoleBindingClient) AddLifecycle(name string, lifecycle BusinessGlobalRoleBindingLifecycle) {
	sync := NewBusinessGlobalRoleBindingLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessGlobalRoleBindingClient) AddClusterScopedHandler(name, clusterName string, sync BusinessGlobalRoleBindingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessGlobalRoleBindingClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessGlobalRoleBindingLifecycle) {
	sync := NewBusinessGlobalRoleBindingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
