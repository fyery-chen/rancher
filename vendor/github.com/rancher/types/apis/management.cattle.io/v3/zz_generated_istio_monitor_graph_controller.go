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
	IstioMonitorGraphGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "IstioMonitorGraph",
	}
	IstioMonitorGraphResource = metav1.APIResource{
		Name:         "istiomonitorgraphs",
		SingularName: "istiomonitorgraph",
		Namespaced:   true,

		Kind: IstioMonitorGraphGroupVersionKind.Kind,
	}
)

func NewIstioMonitorGraph(namespace, name string, obj IstioMonitorGraph) *IstioMonitorGraph {
	obj.APIVersion, obj.Kind = IstioMonitorGraphGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type IstioMonitorGraphList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IstioMonitorGraph
}

type IstioMonitorGraphHandlerFunc func(key string, obj *IstioMonitorGraph) (runtime.Object, error)

type IstioMonitorGraphChangeHandlerFunc func(obj *IstioMonitorGraph) (runtime.Object, error)

type IstioMonitorGraphLister interface {
	List(namespace string, selector labels.Selector) (ret []*IstioMonitorGraph, err error)
	Get(namespace, name string) (*IstioMonitorGraph, error)
}

type IstioMonitorGraphController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() IstioMonitorGraphLister
	AddHandler(ctx context.Context, name string, handler IstioMonitorGraphHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler IstioMonitorGraphHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type IstioMonitorGraphInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*IstioMonitorGraph) (*IstioMonitorGraph, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioMonitorGraph, error)
	Get(name string, opts metav1.GetOptions) (*IstioMonitorGraph, error)
	Update(*IstioMonitorGraph) (*IstioMonitorGraph, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*IstioMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() IstioMonitorGraphController
	AddHandler(ctx context.Context, name string, sync IstioMonitorGraphHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle IstioMonitorGraphLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioMonitorGraphHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioMonitorGraphLifecycle)
}

type istioMonitorGraphLister struct {
	controller *istioMonitorGraphController
}

func (l *istioMonitorGraphLister) List(namespace string, selector labels.Selector) (ret []*IstioMonitorGraph, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*IstioMonitorGraph))
	})
	return
}

func (l *istioMonitorGraphLister) Get(namespace, name string) (*IstioMonitorGraph, error) {
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
			Group:    IstioMonitorGraphGroupVersionKind.Group,
			Resource: "istioMonitorGraph",
		}, key)
	}
	return obj.(*IstioMonitorGraph), nil
}

type istioMonitorGraphController struct {
	controller.GenericController
}

func (c *istioMonitorGraphController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *istioMonitorGraphController) Lister() IstioMonitorGraphLister {
	return &istioMonitorGraphLister{
		controller: c,
	}
}

func (c *istioMonitorGraphController) AddHandler(ctx context.Context, name string, handler IstioMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioMonitorGraph); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *istioMonitorGraphController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler IstioMonitorGraphHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*IstioMonitorGraph); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type istioMonitorGraphFactory struct {
}

func (c istioMonitorGraphFactory) Object() runtime.Object {
	return &IstioMonitorGraph{}
}

func (c istioMonitorGraphFactory) List() runtime.Object {
	return &IstioMonitorGraphList{}
}

func (s *istioMonitorGraphClient) Controller() IstioMonitorGraphController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.istioMonitorGraphControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(IstioMonitorGraphGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &istioMonitorGraphController{
		GenericController: genericController,
	}

	s.client.istioMonitorGraphControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type istioMonitorGraphClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   IstioMonitorGraphController
}

func (s *istioMonitorGraphClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *istioMonitorGraphClient) Create(o *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*IstioMonitorGraph), err
}

func (s *istioMonitorGraphClient) Get(name string, opts metav1.GetOptions) (*IstioMonitorGraph, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*IstioMonitorGraph), err
}

func (s *istioMonitorGraphClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*IstioMonitorGraph, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*IstioMonitorGraph), err
}

func (s *istioMonitorGraphClient) Update(o *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*IstioMonitorGraph), err
}

func (s *istioMonitorGraphClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *istioMonitorGraphClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *istioMonitorGraphClient) List(opts metav1.ListOptions) (*IstioMonitorGraphList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*IstioMonitorGraphList), err
}

func (s *istioMonitorGraphClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *istioMonitorGraphClient) Patch(o *IstioMonitorGraph, patchType types.PatchType, data []byte, subresources ...string) (*IstioMonitorGraph, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*IstioMonitorGraph), err
}

func (s *istioMonitorGraphClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *istioMonitorGraphClient) AddHandler(ctx context.Context, name string, sync IstioMonitorGraphHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioMonitorGraphClient) AddLifecycle(ctx context.Context, name string, lifecycle IstioMonitorGraphLifecycle) {
	sync := NewIstioMonitorGraphLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *istioMonitorGraphClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync IstioMonitorGraphHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *istioMonitorGraphClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle IstioMonitorGraphLifecycle) {
	sync := NewIstioMonitorGraphLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type IstioMonitorGraphIndexer func(obj *IstioMonitorGraph) ([]string, error)

type IstioMonitorGraphClientCache interface {
	Get(namespace, name string) (*IstioMonitorGraph, error)
	List(namespace string, selector labels.Selector) ([]*IstioMonitorGraph, error)

	Index(name string, indexer IstioMonitorGraphIndexer)
	GetIndexed(name, key string) ([]*IstioMonitorGraph, error)
}

type IstioMonitorGraphClient interface {
	Create(*IstioMonitorGraph) (*IstioMonitorGraph, error)
	Get(namespace, name string, opts metav1.GetOptions) (*IstioMonitorGraph, error)
	Update(*IstioMonitorGraph) (*IstioMonitorGraph, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*IstioMonitorGraphList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() IstioMonitorGraphClientCache

	OnCreate(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	ObjectClient() *objectclient.ObjectClient
	Interface() IstioMonitorGraphInterface
}

type istioMonitorGraphClientCache struct {
	client *istioMonitorGraphClient2
}

type istioMonitorGraphClient2 struct {
	iface      IstioMonitorGraphInterface
	controller IstioMonitorGraphController
}

func (n *istioMonitorGraphClient2) Interface() IstioMonitorGraphInterface {
	return n.iface
}

func (n *istioMonitorGraphClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *istioMonitorGraphClient2) ObjectClient() *objectclient.ObjectClient {
	return n.Interface().ObjectClient()
}

func (n *istioMonitorGraphClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *istioMonitorGraphClient2) Create(obj *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	return n.iface.Create(obj)
}

func (n *istioMonitorGraphClient2) Get(namespace, name string, opts metav1.GetOptions) (*IstioMonitorGraph, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *istioMonitorGraphClient2) Update(obj *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	return n.iface.Update(obj)
}

func (n *istioMonitorGraphClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *istioMonitorGraphClient2) List(namespace string, opts metav1.ListOptions) (*IstioMonitorGraphList, error) {
	return n.iface.List(opts)
}

func (n *istioMonitorGraphClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *istioMonitorGraphClientCache) Get(namespace, name string) (*IstioMonitorGraph, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *istioMonitorGraphClientCache) List(namespace string, selector labels.Selector) ([]*IstioMonitorGraph, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *istioMonitorGraphClient2) Cache() IstioMonitorGraphClientCache {
	n.loadController()
	return &istioMonitorGraphClientCache{
		client: n,
	}
}

func (n *istioMonitorGraphClient2) OnCreate(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &istioMonitorGraphLifecycleDelegate{create: sync})
}

func (n *istioMonitorGraphClient2) OnChange(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &istioMonitorGraphLifecycleDelegate{update: sync})
}

func (n *istioMonitorGraphClient2) OnRemove(ctx context.Context, name string, sync IstioMonitorGraphChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &istioMonitorGraphLifecycleDelegate{remove: sync})
}

func (n *istioMonitorGraphClientCache) Index(name string, indexer IstioMonitorGraphIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*IstioMonitorGraph); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *istioMonitorGraphClientCache) GetIndexed(name, key string) ([]*IstioMonitorGraph, error) {
	var result []*IstioMonitorGraph
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*IstioMonitorGraph); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *istioMonitorGraphClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type istioMonitorGraphLifecycleDelegate struct {
	create IstioMonitorGraphChangeHandlerFunc
	update IstioMonitorGraphChangeHandlerFunc
	remove IstioMonitorGraphChangeHandlerFunc
}

func (n *istioMonitorGraphLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *istioMonitorGraphLifecycleDelegate) Create(obj *IstioMonitorGraph) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *istioMonitorGraphLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *istioMonitorGraphLifecycleDelegate) Remove(obj *IstioMonitorGraph) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *istioMonitorGraphLifecycleDelegate) Updated(obj *IstioMonitorGraph) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
