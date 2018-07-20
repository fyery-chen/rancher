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
	BusinessRoleTemplateGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "BusinessRoleTemplate",
	}
	BusinessRoleTemplateResource = metav1.APIResource{
		Name:         "businessroletemplates",
		SingularName: "businessroletemplate",
		Namespaced:   false,
		Kind:         BusinessRoleTemplateGroupVersionKind.Kind,
	}
)

type BusinessRoleTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessRoleTemplate
}

type BusinessRoleTemplateHandlerFunc func(key string, obj *BusinessRoleTemplate) error

type BusinessRoleTemplateLister interface {
	List(namespace string, selector labels.Selector) (ret []*BusinessRoleTemplate, err error)
	Get(namespace, name string) (*BusinessRoleTemplate, error)
}

type BusinessRoleTemplateController interface {
	Informer() cache.SharedIndexInformer
	Lister() BusinessRoleTemplateLister
	AddHandler(name string, handler BusinessRoleTemplateHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler BusinessRoleTemplateHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type BusinessRoleTemplateInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*BusinessRoleTemplate) (*BusinessRoleTemplate, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessRoleTemplate, error)
	Get(name string, opts metav1.GetOptions) (*BusinessRoleTemplate, error)
	Update(*BusinessRoleTemplate) (*BusinessRoleTemplate, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*BusinessRoleTemplateList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() BusinessRoleTemplateController
	AddHandler(name string, sync BusinessRoleTemplateHandlerFunc)
	AddLifecycle(name string, lifecycle BusinessRoleTemplateLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync BusinessRoleTemplateHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessRoleTemplateLifecycle)
}

type businessRoleTemplateLister struct {
	controller *businessRoleTemplateController
}

func (l *businessRoleTemplateLister) List(namespace string, selector labels.Selector) (ret []*BusinessRoleTemplate, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*BusinessRoleTemplate))
	})
	return
}

func (l *businessRoleTemplateLister) Get(namespace, name string) (*BusinessRoleTemplate, error) {
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
			Group:    BusinessRoleTemplateGroupVersionKind.Group,
			Resource: "businessRoleTemplate",
		}, name)
	}
	return obj.(*BusinessRoleTemplate), nil
}

type businessRoleTemplateController struct {
	controller.GenericController
}

func (c *businessRoleTemplateController) Lister() BusinessRoleTemplateLister {
	return &businessRoleTemplateLister{
		controller: c,
	}
}

func (c *businessRoleTemplateController) AddHandler(name string, handler BusinessRoleTemplateHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*BusinessRoleTemplate))
	})
}

func (c *businessRoleTemplateController) AddClusterScopedHandler(name, cluster string, handler BusinessRoleTemplateHandlerFunc) {
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

		return handler(key, obj.(*BusinessRoleTemplate))
	})
}

type businessRoleTemplateFactory struct {
}

func (c businessRoleTemplateFactory) Object() runtime.Object {
	return &BusinessRoleTemplate{}
}

func (c businessRoleTemplateFactory) List() runtime.Object {
	return &BusinessRoleTemplateList{}
}

func (s *businessRoleTemplateClient) Controller() BusinessRoleTemplateController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.businessRoleTemplateControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(BusinessRoleTemplateGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &businessRoleTemplateController{
		GenericController: genericController,
	}

	s.client.businessRoleTemplateControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type businessRoleTemplateClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   BusinessRoleTemplateController
}

func (s *businessRoleTemplateClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *businessRoleTemplateClient) Create(o *BusinessRoleTemplate) (*BusinessRoleTemplate, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*BusinessRoleTemplate), err
}

func (s *businessRoleTemplateClient) Get(name string, opts metav1.GetOptions) (*BusinessRoleTemplate, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*BusinessRoleTemplate), err
}

func (s *businessRoleTemplateClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*BusinessRoleTemplate, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*BusinessRoleTemplate), err
}

func (s *businessRoleTemplateClient) Update(o *BusinessRoleTemplate) (*BusinessRoleTemplate, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*BusinessRoleTemplate), err
}

func (s *businessRoleTemplateClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *businessRoleTemplateClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *businessRoleTemplateClient) List(opts metav1.ListOptions) (*BusinessRoleTemplateList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*BusinessRoleTemplateList), err
}

func (s *businessRoleTemplateClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *businessRoleTemplateClient) Patch(o *BusinessRoleTemplate, data []byte, subresources ...string) (*BusinessRoleTemplate, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*BusinessRoleTemplate), err
}

func (s *businessRoleTemplateClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *businessRoleTemplateClient) AddHandler(name string, sync BusinessRoleTemplateHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *businessRoleTemplateClient) AddLifecycle(name string, lifecycle BusinessRoleTemplateLifecycle) {
	sync := NewBusinessRoleTemplateLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *businessRoleTemplateClient) AddClusterScopedHandler(name, clusterName string, sync BusinessRoleTemplateHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *businessRoleTemplateClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle BusinessRoleTemplateLifecycle) {
	sync := NewBusinessRoleTemplateLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
