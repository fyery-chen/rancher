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
	GlobalRegistryGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GlobalRegistry",
	}
	GlobalRegistryResource = metav1.APIResource{
		Name:         "globalregistries",
		SingularName: "globalregistry",
		Namespaced:   false,
		Kind:         GlobalRegistryGroupVersionKind.Kind,
	}
)

func NewGlobalRegistry(namespace, name string, obj GlobalRegistry) *GlobalRegistry {
	obj.APIVersion, obj.Kind = GlobalRegistryGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type GlobalRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlobalRegistry
}

type GlobalRegistryHandlerFunc func(key string, obj *GlobalRegistry) (runtime.Object, error)

type GlobalRegistryChangeHandlerFunc func(obj *GlobalRegistry) (runtime.Object, error)

type GlobalRegistryLister interface {
	List(namespace string, selector labels.Selector) (ret []*GlobalRegistry, err error)
	Get(namespace, name string) (*GlobalRegistry, error)
}

type GlobalRegistryController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GlobalRegistryLister
	AddHandler(ctx context.Context, name string, handler GlobalRegistryHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GlobalRegistryHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type GlobalRegistryInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*GlobalRegistry) (*GlobalRegistry, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalRegistry, error)
	Get(name string, opts metav1.GetOptions) (*GlobalRegistry, error)
	Update(*GlobalRegistry) (*GlobalRegistry, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*GlobalRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GlobalRegistryController
	AddHandler(ctx context.Context, name string, sync GlobalRegistryHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GlobalRegistryLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalRegistryHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalRegistryLifecycle)
}

type globalRegistryLister struct {
	controller *globalRegistryController
}

func (l *globalRegistryLister) List(namespace string, selector labels.Selector) (ret []*GlobalRegistry, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*GlobalRegistry))
	})
	return
}

func (l *globalRegistryLister) Get(namespace, name string) (*GlobalRegistry, error) {
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
			Group:    GlobalRegistryGroupVersionKind.Group,
			Resource: "globalRegistry",
		}, key)
	}
	return obj.(*GlobalRegistry), nil
}

type globalRegistryController struct {
	controller.GenericController
}

func (c *globalRegistryController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *globalRegistryController) Lister() GlobalRegistryLister {
	return &globalRegistryLister{
		controller: c,
	}
}

func (c *globalRegistryController) AddHandler(ctx context.Context, name string, handler GlobalRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRegistry); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *globalRegistryController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GlobalRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GlobalRegistry); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type globalRegistryFactory struct {
}

func (c globalRegistryFactory) Object() runtime.Object {
	return &GlobalRegistry{}
}

func (c globalRegistryFactory) List() runtime.Object {
	return &GlobalRegistryList{}
}

func (s *globalRegistryClient) Controller() GlobalRegistryController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.globalRegistryControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(GlobalRegistryGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &globalRegistryController{
		GenericController: genericController,
	}

	s.client.globalRegistryControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type globalRegistryClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GlobalRegistryController
}

func (s *globalRegistryClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *globalRegistryClient) Create(o *GlobalRegistry) (*GlobalRegistry, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*GlobalRegistry), err
}

func (s *globalRegistryClient) Get(name string, opts metav1.GetOptions) (*GlobalRegistry, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*GlobalRegistry), err
}

func (s *globalRegistryClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GlobalRegistry, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*GlobalRegistry), err
}

func (s *globalRegistryClient) Update(o *GlobalRegistry) (*GlobalRegistry, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*GlobalRegistry), err
}

func (s *globalRegistryClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *globalRegistryClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *globalRegistryClient) List(opts metav1.ListOptions) (*GlobalRegistryList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*GlobalRegistryList), err
}

func (s *globalRegistryClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *globalRegistryClient) Patch(o *GlobalRegistry, patchType types.PatchType, data []byte, subresources ...string) (*GlobalRegistry, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*GlobalRegistry), err
}

func (s *globalRegistryClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *globalRegistryClient) AddHandler(ctx context.Context, name string, sync GlobalRegistryHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalRegistryClient) AddLifecycle(ctx context.Context, name string, lifecycle GlobalRegistryLifecycle) {
	sync := NewGlobalRegistryLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *globalRegistryClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GlobalRegistryHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *globalRegistryClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GlobalRegistryLifecycle) {
	sync := NewGlobalRegistryLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type GlobalRegistryIndexer func(obj *GlobalRegistry) ([]string, error)

type GlobalRegistryClientCache interface {
	Get(namespace, name string) (*GlobalRegistry, error)
	List(namespace string, selector labels.Selector) ([]*GlobalRegistry, error)

	Index(name string, indexer GlobalRegistryIndexer)
	GetIndexed(name, key string) ([]*GlobalRegistry, error)
}

type GlobalRegistryClient interface {
	Create(*GlobalRegistry) (*GlobalRegistry, error)
	Get(namespace, name string, opts metav1.GetOptions) (*GlobalRegistry, error)
	Update(*GlobalRegistry) (*GlobalRegistry, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*GlobalRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() GlobalRegistryClientCache

	OnCreate(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() GlobalRegistryInterface
}

type globalRegistryClientCache struct {
	client *globalRegistryClient2
}

type globalRegistryClient2 struct {
	iface      GlobalRegistryInterface
	controller GlobalRegistryController
}

func (n *globalRegistryClient2) Interface() GlobalRegistryInterface {
	return n.iface
}

func (n *globalRegistryClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *globalRegistryClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *globalRegistryClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *globalRegistryClient2) Create(obj *GlobalRegistry) (*GlobalRegistry, error) {
	return n.iface.Create(obj)
}

func (n *globalRegistryClient2) Get(namespace, name string, opts metav1.GetOptions) (*GlobalRegistry, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *globalRegistryClient2) Update(obj *GlobalRegistry) (*GlobalRegistry, error) {
	return n.iface.Update(obj)
}

func (n *globalRegistryClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *globalRegistryClient2) List(namespace string, opts metav1.ListOptions) (*GlobalRegistryList, error) {
	return n.iface.List(opts)
}

func (n *globalRegistryClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *globalRegistryClientCache) Get(namespace, name string) (*GlobalRegistry, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *globalRegistryClientCache) List(namespace string, selector labels.Selector) ([]*GlobalRegistry, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *globalRegistryClient2) Cache() GlobalRegistryClientCache {
	n.loadController()
	return &globalRegistryClientCache{
		client: n,
	}
}

func (n *globalRegistryClient2) OnCreate(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &globalRegistryLifecycleDelegate{create: sync})
}

func (n *globalRegistryClient2) OnChange(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &globalRegistryLifecycleDelegate{update: sync})
}

func (n *globalRegistryClient2) OnRemove(ctx context.Context, name string, sync GlobalRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &globalRegistryLifecycleDelegate{remove: sync})
}

func (n *globalRegistryClientCache) Index(name string, indexer GlobalRegistryIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*GlobalRegistry); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *globalRegistryClientCache) GetIndexed(name, key string) ([]*GlobalRegistry, error) {
	var result []*GlobalRegistry
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*GlobalRegistry); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *globalRegistryClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type globalRegistryLifecycleDelegate struct {
	create GlobalRegistryChangeHandlerFunc
	update GlobalRegistryChangeHandlerFunc
	remove GlobalRegistryChangeHandlerFunc
}

func (n *globalRegistryLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *globalRegistryLifecycleDelegate) Create(obj *GlobalRegistry) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *globalRegistryLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *globalRegistryLifecycleDelegate) Remove(obj *GlobalRegistry) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *globalRegistryLifecycleDelegate) Updated(obj *GlobalRegistry) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
