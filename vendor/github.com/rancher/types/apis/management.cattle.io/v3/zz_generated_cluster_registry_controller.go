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
	ClusterRegistryGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "ClusterRegistry",
	}
	ClusterRegistryResource = metav1.APIResource{
		Name:         "clusterregistries",
		SingularName: "clusterregistry",
		Namespaced:   true,

		Kind: ClusterRegistryGroupVersionKind.Kind,
	}
)

func NewClusterRegistry(namespace, name string, obj ClusterRegistry) *ClusterRegistry {
	obj.APIVersion, obj.Kind = ClusterRegistryGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type ClusterRegistryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterRegistry
}

type ClusterRegistryHandlerFunc func(key string, obj *ClusterRegistry) (runtime.Object, error)

type ClusterRegistryChangeHandlerFunc func(obj *ClusterRegistry) (runtime.Object, error)

type ClusterRegistryLister interface {
	List(namespace string, selector labels.Selector) (ret []*ClusterRegistry, err error)
	Get(namespace, name string) (*ClusterRegistry, error)
}

type ClusterRegistryController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() ClusterRegistryLister
	AddHandler(ctx context.Context, name string, handler ClusterRegistryHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler ClusterRegistryHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type ClusterRegistryInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*ClusterRegistry) (*ClusterRegistry, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ClusterRegistry, error)
	Get(name string, opts metav1.GetOptions) (*ClusterRegistry, error)
	Update(*ClusterRegistry) (*ClusterRegistry, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*ClusterRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() ClusterRegistryController
	AddHandler(ctx context.Context, name string, sync ClusterRegistryHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle ClusterRegistryLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterRegistryHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterRegistryLifecycle)
}

type clusterRegistryLister struct {
	controller *clusterRegistryController
}

func (l *clusterRegistryLister) List(namespace string, selector labels.Selector) (ret []*ClusterRegistry, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*ClusterRegistry))
	})
	return
}

func (l *clusterRegistryLister) Get(namespace, name string) (*ClusterRegistry, error) {
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
			Group:    ClusterRegistryGroupVersionKind.Group,
			Resource: "clusterRegistry",
		}, key)
	}
	return obj.(*ClusterRegistry), nil
}

type clusterRegistryController struct {
	controller.GenericController
}

func (c *clusterRegistryController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *clusterRegistryController) Lister() ClusterRegistryLister {
	return &clusterRegistryLister{
		controller: c,
	}
}

func (c *clusterRegistryController) AddHandler(ctx context.Context, name string, handler ClusterRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterRegistry); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *clusterRegistryController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler ClusterRegistryHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*ClusterRegistry); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type clusterRegistryFactory struct {
}

func (c clusterRegistryFactory) Object() runtime.Object {
	return &ClusterRegistry{}
}

func (c clusterRegistryFactory) List() runtime.Object {
	return &ClusterRegistryList{}
}

func (s *clusterRegistryClient) Controller() ClusterRegistryController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.clusterRegistryControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(ClusterRegistryGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &clusterRegistryController{
		GenericController: genericController,
	}

	s.client.clusterRegistryControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type clusterRegistryClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   ClusterRegistryController
}

func (s *clusterRegistryClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *clusterRegistryClient) Create(o *ClusterRegistry) (*ClusterRegistry, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*ClusterRegistry), err
}

func (s *clusterRegistryClient) Get(name string, opts metav1.GetOptions) (*ClusterRegistry, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*ClusterRegistry), err
}

func (s *clusterRegistryClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*ClusterRegistry, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*ClusterRegistry), err
}

func (s *clusterRegistryClient) Update(o *ClusterRegistry) (*ClusterRegistry, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*ClusterRegistry), err
}

func (s *clusterRegistryClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *clusterRegistryClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *clusterRegistryClient) List(opts metav1.ListOptions) (*ClusterRegistryList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*ClusterRegistryList), err
}

func (s *clusterRegistryClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *clusterRegistryClient) Patch(o *ClusterRegistry, patchType types.PatchType, data []byte, subresources ...string) (*ClusterRegistry, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*ClusterRegistry), err
}

func (s *clusterRegistryClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *clusterRegistryClient) AddHandler(ctx context.Context, name string, sync ClusterRegistryHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterRegistryClient) AddLifecycle(ctx context.Context, name string, lifecycle ClusterRegistryLifecycle) {
	sync := NewClusterRegistryLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *clusterRegistryClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync ClusterRegistryHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *clusterRegistryClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle ClusterRegistryLifecycle) {
	sync := NewClusterRegistryLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type ClusterRegistryIndexer func(obj *ClusterRegistry) ([]string, error)

type ClusterRegistryClientCache interface {
	Get(namespace, name string) (*ClusterRegistry, error)
	List(namespace string, selector labels.Selector) ([]*ClusterRegistry, error)

	Index(name string, indexer ClusterRegistryIndexer)
	GetIndexed(name, key string) ([]*ClusterRegistry, error)
}

type ClusterRegistryClient interface {
	Create(*ClusterRegistry) (*ClusterRegistry, error)
	Get(namespace, name string, opts metav1.GetOptions) (*ClusterRegistry, error)
	Update(*ClusterRegistry) (*ClusterRegistry, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*ClusterRegistryList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() ClusterRegistryClientCache

	OnCreate(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() ClusterRegistryInterface
}

type clusterRegistryClientCache struct {
	client *clusterRegistryClient2
}

type clusterRegistryClient2 struct {
	iface      ClusterRegistryInterface
	controller ClusterRegistryController
}

func (n *clusterRegistryClient2) Interface() ClusterRegistryInterface {
	return n.iface
}

func (n *clusterRegistryClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *clusterRegistryClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *clusterRegistryClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *clusterRegistryClient2) Create(obj *ClusterRegistry) (*ClusterRegistry, error) {
	return n.iface.Create(obj)
}

func (n *clusterRegistryClient2) Get(namespace, name string, opts metav1.GetOptions) (*ClusterRegistry, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *clusterRegistryClient2) Update(obj *ClusterRegistry) (*ClusterRegistry, error) {
	return n.iface.Update(obj)
}

func (n *clusterRegistryClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *clusterRegistryClient2) List(namespace string, opts metav1.ListOptions) (*ClusterRegistryList, error) {
	return n.iface.List(opts)
}

func (n *clusterRegistryClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *clusterRegistryClientCache) Get(namespace, name string) (*ClusterRegistry, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *clusterRegistryClientCache) List(namespace string, selector labels.Selector) ([]*ClusterRegistry, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *clusterRegistryClient2) Cache() ClusterRegistryClientCache {
	n.loadController()
	return &clusterRegistryClientCache{
		client: n,
	}
}

func (n *clusterRegistryClient2) OnCreate(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &clusterRegistryLifecycleDelegate{create: sync})
}

func (n *clusterRegistryClient2) OnChange(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &clusterRegistryLifecycleDelegate{update: sync})
}

func (n *clusterRegistryClient2) OnRemove(ctx context.Context, name string, sync ClusterRegistryChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &clusterRegistryLifecycleDelegate{remove: sync})
}

func (n *clusterRegistryClientCache) Index(name string, indexer ClusterRegistryIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*ClusterRegistry); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *clusterRegistryClientCache) GetIndexed(name, key string) ([]*ClusterRegistry, error) {
	var result []*ClusterRegistry
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*ClusterRegistry); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *clusterRegistryClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type clusterRegistryLifecycleDelegate struct {
	create ClusterRegistryChangeHandlerFunc
	update ClusterRegistryChangeHandlerFunc
	remove ClusterRegistryChangeHandlerFunc
}

func (n *clusterRegistryLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *clusterRegistryLifecycleDelegate) Create(obj *ClusterRegistry) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *clusterRegistryLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *clusterRegistryLifecycleDelegate) Remove(obj *ClusterRegistry) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *clusterRegistryLifecycleDelegate) Updated(obj *ClusterRegistry) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
