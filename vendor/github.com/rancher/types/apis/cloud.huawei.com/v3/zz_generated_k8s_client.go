package v3

import (
	"context"
	"sync"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/restwatch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Interface interface {
	RESTClient() rest.Interface
	controller.Starter

	HuaweiCloudAccountsGetter
	BusinessGlobalRolesGetter
	BusinessGlobalRoleBindingsGetter
	BusinessRoleTemplatesGetter
	BusinessRoleTemplateBindingsGetter
	DynamicSchemasGetter
	BusinessesGetter
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	huaweiCloudAccountControllers          map[string]HuaweiCloudAccountController
	businessGlobalRoleControllers          map[string]BusinessGlobalRoleController
	businessGlobalRoleBindingControllers   map[string]BusinessGlobalRoleBindingController
	businessRoleTemplateControllers        map[string]BusinessRoleTemplateController
	businessRoleTemplateBindingControllers map[string]BusinessRoleTemplateBindingController
	dynamicSchemaControllers               map[string]DynamicSchemaController
	businessControllers                    map[string]BusinessController
}

func NewForConfig(config rest.Config) (Interface, error) {
	if config.NegotiatedSerializer == nil {
		configConfig := dynamic.ContentConfig()
		config.NegotiatedSerializer = configConfig.NegotiatedSerializer
	}

	restClient, err := restwatch.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient: restClient,

		huaweiCloudAccountControllers:          map[string]HuaweiCloudAccountController{},
		businessGlobalRoleControllers:          map[string]BusinessGlobalRoleController{},
		businessGlobalRoleBindingControllers:   map[string]BusinessGlobalRoleBindingController{},
		businessRoleTemplateControllers:        map[string]BusinessRoleTemplateController{},
		businessRoleTemplateBindingControllers: map[string]BusinessRoleTemplateBindingController{},
		dynamicSchemaControllers:               map[string]DynamicSchemaController{},
		businessControllers:                    map[string]BusinessController{},
	}, nil
}

func (c *Client) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Client) Sync(ctx context.Context) error {
	return controller.Sync(ctx, c.starters...)
}

func (c *Client) Start(ctx context.Context, threadiness int) error {
	return controller.Start(ctx, threadiness, c.starters...)
}

type HuaweiCloudAccountsGetter interface {
	HuaweiCloudAccounts(namespace string) HuaweiCloudAccountInterface
}

func (c *Client) HuaweiCloudAccounts(namespace string) HuaweiCloudAccountInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &HuaweiCloudAccountResource, HuaweiCloudAccountGroupVersionKind, huaweiCloudAccountFactory{})
	return &huaweiCloudAccountClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type BusinessGlobalRolesGetter interface {
	BusinessGlobalRoles(namespace string) BusinessGlobalRoleInterface
}

func (c *Client) BusinessGlobalRoles(namespace string) BusinessGlobalRoleInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &BusinessGlobalRoleResource, BusinessGlobalRoleGroupVersionKind, businessGlobalRoleFactory{})
	return &businessGlobalRoleClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type BusinessGlobalRoleBindingsGetter interface {
	BusinessGlobalRoleBindings(namespace string) BusinessGlobalRoleBindingInterface
}

func (c *Client) BusinessGlobalRoleBindings(namespace string) BusinessGlobalRoleBindingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &BusinessGlobalRoleBindingResource, BusinessGlobalRoleBindingGroupVersionKind, businessGlobalRoleBindingFactory{})
	return &businessGlobalRoleBindingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type BusinessRoleTemplatesGetter interface {
	BusinessRoleTemplates(namespace string) BusinessRoleTemplateInterface
}

func (c *Client) BusinessRoleTemplates(namespace string) BusinessRoleTemplateInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &BusinessRoleTemplateResource, BusinessRoleTemplateGroupVersionKind, businessRoleTemplateFactory{})
	return &businessRoleTemplateClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type BusinessRoleTemplateBindingsGetter interface {
	BusinessRoleTemplateBindings(namespace string) BusinessRoleTemplateBindingInterface
}

func (c *Client) BusinessRoleTemplateBindings(namespace string) BusinessRoleTemplateBindingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &BusinessRoleTemplateBindingResource, BusinessRoleTemplateBindingGroupVersionKind, businessRoleTemplateBindingFactory{})
	return &businessRoleTemplateBindingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type DynamicSchemasGetter interface {
	DynamicSchemas(namespace string) DynamicSchemaInterface
}

func (c *Client) DynamicSchemas(namespace string) DynamicSchemaInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &DynamicSchemaResource, DynamicSchemaGroupVersionKind, dynamicSchemaFactory{})
	return &dynamicSchemaClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type BusinessesGetter interface {
	Businesses(namespace string) BusinessInterface
}

func (c *Client) Businesses(namespace string) BusinessInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &BusinessResource, BusinessGroupVersionKind, businessFactory{})
	return &businessClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}
