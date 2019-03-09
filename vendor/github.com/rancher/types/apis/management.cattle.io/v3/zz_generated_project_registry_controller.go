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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	ProjectRegistryGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ProjectRegistry",
	}
	ProjectRegistryResource = metav1.APIResource{
		Name:         "projectregistries",
		SingularName: "projectregistry",
		Namespaced:   true,

		Kind: ProjectRegistryGroupVersionKind.Kind,
	}
)

func NewProjectRegistry(namespace, name string, obj ProjectRegistry) *ProjectRegistry {
	obj.APIVersion, obj.Kind = ProjectRegistryGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ProjectRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProjectRegistry
}

type ProjectRegistryHandlerFunc func(key string, obj *ProjectRegistry) (runtime.Object, error)

type ProjectRegistryChangeHandlerFunc func(obj *ProjectRegistry) (runtime.Object, error)

type ProjectRegistryLister interface {
	List(namespace string, selector labels.Selector) (ret []*ProjectRegistry, err error)
	Get(namespace, name string) (*ProjectRegistry, error)
}

type ProjectRegistryController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ProjectRegistryLister
	AddHandler(ctx context.Context, name string, handler ProjectRegistryHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ProjectRegistryHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ProjectRegistryInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*ProjectRegistry) (*ProjectRegistry, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ProjectRegistry, error)
	Get(name string, opts metav1.GetOptions) (*ProjectRegistry, error)
	Update(*ProjectRegistry) (*ProjectRegistry, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ProjectRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ProjectRegistryController
	AddHandler(ctx context.Context, name string, sync ProjectRegistryHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ProjectRegistryLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectRegistryHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectRegistryLifecycle)
}

type projectRegistryLister struct {
	controller *projectRegistryController
}

func (l *projectRegistryLister) List(namespace string, selector labels.Selector) (ret []*ProjectRegistry, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ProjectRegistry))
	})
	return
}

func (l *projectRegistryLister) Get(namespace, name string) (*ProjectRegistry, error) {
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
			Group:    ProjectRegistryGroupVersionKind.Group,
			Resource: "projectRegistry",
		}, key)
	}
	return obj.(*ProjectRegistry), nil
}

type projectRegistryController struct {
	controller.GenericController
}

func (c *projectRegistryController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *projectRegistryController) Lister() ProjectRegistryLister {
	return &projectRegistryLister{
		controller: c,
	}
}

func (c *projectRegistryController) AddHandler(ctx context.Context, name string, handler ProjectRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectRegistry); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *projectRegistryController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ProjectRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ProjectRegistry); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type projectRegistryFactory struct {
}

func (c projectRegistryFactory) Object() runtime.Object {
	return &ProjectRegistry{}
}

func (c projectRegistryFactory) List() runtime.Object {
	return &ProjectRegistryList{}
}

func (s *projectRegistryClient) Controller() ProjectRegistryController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.projectRegistryControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ProjectRegistryGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &projectRegistryController{
		GenericController: genericController,
	}

	s.client.projectRegistryControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type projectRegistryClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ProjectRegistryController
}

func (s *projectRegistryClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *projectRegistryClient) Create(o *ProjectRegistry) (*ProjectRegistry, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ProjectRegistry), err
}

func (s *projectRegistryClient) Get(name string, opts metav1.GetOptions) (*ProjectRegistry, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ProjectRegistry), err
}

func (s *projectRegistryClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ProjectRegistry, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*ProjectRegistry), err
}

func (s *projectRegistryClient) Update(o *ProjectRegistry) (*ProjectRegistry, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ProjectRegistry), err
}

func (s *projectRegistryClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *projectRegistryClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *projectRegistryClient) List(opts metav1.ListOptions) (*ProjectRegistryList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ProjectRegistryList), err
}

func (s *projectRegistryClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *projectRegistryClient) Patch(o *ProjectRegistry, patchType types.PatchType, data []byte, subresources ...string) (*ProjectRegistry, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*ProjectRegistry), err
}

func (s *projectRegistryClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *projectRegistryClient) AddHandler(ctx context.Context, name string, sync ProjectRegistryHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectRegistryClient) AddLifecycle(ctx context.Context, name string, lifecycle ProjectRegistryLifecycle) {
	sync := NewProjectRegistryLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *projectRegistryClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ProjectRegistryHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *projectRegistryClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ProjectRegistryLifecycle) {
	sync := NewProjectRegistryLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ProjectRegistryIndexer func(obj *ProjectRegistry) ([]string, error)

type ProjectRegistryClientCache interface {
	Get(namespace, name string) (*ProjectRegistry, error)
	List(namespace string, selector labels.Selector) ([]*ProjectRegistry, error)

	Index(name string, indexer ProjectRegistryIndexer)
	GetIndexed(name, key string) ([]*ProjectRegistry, error)
}

type ProjectRegistryClient interface {
	Create(*ProjectRegistry) (*ProjectRegistry, error)
	Get(namespace, name string, opts metav1.GetOptions) (*ProjectRegistry, error)
	Update(*ProjectRegistry) (*ProjectRegistry, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ProjectRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ProjectRegistryClientCache

	OnCreate(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ProjectRegistryInterface
}

type projectRegistryClientCache struct {
	client *projectRegistryClient2
}

type projectRegistryClient2 struct {
	iface      ProjectRegistryInterface
	controller ProjectRegistryController
}

func (n *projectRegistryClient2) Interface() ProjectRegistryInterface {
	return n.iface
}

func (n *projectRegistryClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *projectRegistryClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *projectRegistryClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *projectRegistryClient2) Create(obj *ProjectRegistry) (*ProjectRegistry, error) {
	return n.iface.Create(obj)
}

func (n *projectRegistryClient2) Get(namespace, name string, opts metav1.GetOptions) (*ProjectRegistry, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *projectRegistryClient2) Update(obj *ProjectRegistry) (*ProjectRegistry, error) {
	return n.iface.Update(obj)
}

func (n *projectRegistryClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *projectRegistryClient2) List(namespace string, opts metav1.ListOptions) (*ProjectRegistryList, error) {
	return n.iface.List(opts)
}

func (n *projectRegistryClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *projectRegistryClientCache) Get(namespace, name string) (*ProjectRegistry, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *projectRegistryClientCache) List(namespace string, selector labels.Selector) ([]*ProjectRegistry, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *projectRegistryClient2) Cache() ProjectRegistryClientCache {
	n.loadController()
	return &projectRegistryClientCache{
		client: n,
	}
}

func (n *projectRegistryClient2) OnCreate(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &projectRegistryLifecycleDelegate{create: sync})
}

func (n *projectRegistryClient2) OnChange(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &projectRegistryLifecycleDelegate{update: sync})
}

func (n *projectRegistryClient2) OnRemove(ctx context.Context, name string, sync ProjectRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &projectRegistryLifecycleDelegate{remove: sync})
}

func (n *projectRegistryClientCache) Index(name string, indexer ProjectRegistryIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*ProjectRegistry); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *projectRegistryClientCache) GetIndexed(name, key string) ([]*ProjectRegistry, error) {
	var result []*ProjectRegistry
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*ProjectRegistry); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *projectRegistryClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type projectRegistryLifecycleDelegate struct {
	create ProjectRegistryChangeHandlerFunc
	update ProjectRegistryChangeHandlerFunc
	remove ProjectRegistryChangeHandlerFunc
}

func (n *projectRegistryLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *projectRegistryLifecycleDelegate) Create(obj *ProjectRegistry) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *projectRegistryLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *projectRegistryLifecycleDelegate) Remove(obj *ProjectRegistry) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *projectRegistryLifecycleDelegate) Updated(obj *ProjectRegistry) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
