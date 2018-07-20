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
	BusinessGlobalRoleGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BusinessGlobalRole",
	}
	BusinessGlobalRoleResource = metav1.APIResource{
		Name:         "businessglobalroles",
		SingularName: "businessglobalrole",
		Namespaced:   false,
		Kind:         BusinessGlobalRoleGroupVersionKind.Kind,
	}
)

type BusinessGlobalRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessGlobalRole
}

type BusinessGlobalRoleHandlerFunc func(key string, obj *BusinessGlobalRole) error

type BusinessGlobalRoleLister interface {
	List(namespace string, selector labels.Selector) (ret []*BusinessGlobalRole, err error)
	Get(namespace, name string) (*BusinessGlobalRole, error)
}

type BusinessGlobalRoleController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessGlobalRoleLister
	AddHandler(name string, handler BusinessGlobalRoleHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessGlobalRoleHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessGlobalRoleInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BusinessGlobalRole) (*BusinessGlobalRole, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessGlobalRole, error)
	Get(name string, opts metav1.GetOptions) (*BusinessGlobalRole, error)
	Update(*BusinessGlobalRole) (*BusinessGlobalRole, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessGlobalRoleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessGlobalRoleController
	AddHandler(name string, sync BusinessGlobalRoleHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessGlobalRoleLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessGlobalRoleHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessGlobalRoleLifecycle)
}

type businessGlobalRoleLister struct {
	controller *businessGlobalRoleController
}

func (l *businessGlobalRoleLister) List(namespace string, selector labels.Selector) (ret []*BusinessGlobalRole, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BusinessGlobalRole))
	})
	return
}

func (l *businessGlobalRoleLister) Get(namespace, name string) (*BusinessGlobalRole, error) {
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
			Group:    BusinessGlobalRoleGroupVersionKind.Group,
			Resource: "businessGlobalRole",
		}, name)
	}
	return obj.(*BusinessGlobalRole), nil
}

type businessGlobalRoleController struct {
	controller.GenericController
}

func (c *businessGlobalRoleController) Lister() BusinessGlobalRoleLister {
	return &businessGlobalRoleLister{
		controller: c,
	}
}

func (c *businessGlobalRoleController) AddHandler(name string, handler BusinessGlobalRoleHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BusinessGlobalRole))
	})
}

func (c *businessGlobalRoleController) AddClusterScopedHandler(name, cluster string, handler BusinessGlobalRoleHandlerFunc) {
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

		return handler(key, obj.(*BusinessGlobalRole))
	})
}

type businessGlobalRoleFactory struct {
}

func (c businessGlobalRoleFactory) Object() runtime.Object {
	return &BusinessGlobalRole{}
}

func (c businessGlobalRoleFactory) List() runtime.Object {
	return &BusinessGlobalRoleList{}
}

func (s *businessGlobalRoleClient) Controller() BusinessGlobalRoleController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessGlobalRoleControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessGlobalRoleGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessGlobalRoleController{
		GenericController: genericController,
	}

	s.client.businessGlobalRoleControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessGlobalRoleClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessGlobalRoleController
}

func (s *businessGlobalRoleClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessGlobalRoleClient) Create(o *BusinessGlobalRole) (*BusinessGlobalRole, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BusinessGlobalRole), err
}

func (s *businessGlobalRoleClient) Get(name string, opts metav1.GetOptions) (*BusinessGlobalRole, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BusinessGlobalRole), err
}

func (s *businessGlobalRoleClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessGlobalRole, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BusinessGlobalRole), err
}

func (s *businessGlobalRoleClient) Update(o *BusinessGlobalRole) (*BusinessGlobalRole, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BusinessGlobalRole), err
}

func (s *businessGlobalRoleClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessGlobalRoleClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessGlobalRoleClient) List(opts metav1.ListOptions) (*BusinessGlobalRoleList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessGlobalRoleList), err
}

func (s *businessGlobalRoleClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessGlobalRoleClient) Patch(o *BusinessGlobalRole, data []byte, subresources ...string) (*BusinessGlobalRole, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BusinessGlobalRole), err
}

func (s *businessGlobalRoleClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessGlobalRoleClient) AddHandler(name string, sync BusinessGlobalRoleHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessGlobalRoleClient) AddLifecycle(name string, lifecycle BusinessGlobalRoleLifecycle) {
	sync := NewBusinessGlobalRoleLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessGlobalRoleClient) AddClusterScopedHandler(name, clusterName string, sync BusinessGlobalRoleHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessGlobalRoleClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessGlobalRoleLifecycle) {
	sync := NewBusinessGlobalRoleLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
