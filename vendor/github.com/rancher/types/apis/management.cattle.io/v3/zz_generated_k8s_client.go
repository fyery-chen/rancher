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

	GroupsGetter
	GroupMembersGetter
	PrincipalsGetter
	UsersGetter
	AuthConfigsGetter
	TokensGetter
	DynamicSchemasGetter
	PreferencesGetter
	ListenConfigsGetter
	SettingsGetter
	ClustersGetter
	ClusterEventsGetter
	ClusterRegistrationTokensGetter
	NodePoolsGetter
	NodesGetter
	NodeDriversGetter
	NodeTemplatesGetter
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	groupControllers                    map[string]GroupController
	groupMemberControllers              map[string]GroupMemberController
	principalControllers                map[string]PrincipalController
	userControllers                     map[string]UserController
	authConfigControllers               map[string]AuthConfigController
	tokenControllers                    map[string]TokenController
	dynamicSchemaControllers            map[string]DynamicSchemaController
	preferenceControllers               map[string]PreferenceController
	listenConfigControllers             map[string]ListenConfigController
	settingControllers                  map[string]SettingController
	clusterControllers                  map[string]ClusterController
	clusterEventControllers             map[string]ClusterEventController
	clusterRegistrationTokenControllers map[string]ClusterRegistrationTokenController
	nodePoolControllers                 map[string]NodePoolController
	nodeControllers                     map[string]NodeController
	nodeDriverControllers               map[string]NodeDriverController
	nodeTemplateControllers             map[string]NodeTemplateController
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

		groupControllers:                    map[string]GroupController{},
		groupMemberControllers:              map[string]GroupMemberController{},
		principalControllers:                map[string]PrincipalController{},
		userControllers:                     map[string]UserController{},
		authConfigControllers:               map[string]AuthConfigController{},
		tokenControllers:                    map[string]TokenController{},
		dynamicSchemaControllers:            map[string]DynamicSchemaController{},
		preferenceControllers:               map[string]PreferenceController{},
		listenConfigControllers:             map[string]ListenConfigController{},
		settingControllers:                  map[string]SettingController{},
		clusterControllers:                  map[string]ClusterController{},
		clusterEventControllers:             map[string]ClusterEventController{},
		clusterRegistrationTokenControllers: map[string]ClusterRegistrationTokenController{},
		nodePoolControllers:                 map[string]NodePoolController{},
		nodeControllers:                     map[string]NodeController{},
		nodeDriverControllers:               map[string]NodeDriverController{},
		nodeTemplateControllers:             map[string]NodeTemplateController{},
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

type GroupsGetter interface {
	Groups(namespace string) GroupInterface
}

func (c *Client) Groups(namespace string) GroupInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &GroupResource, GroupGroupVersionKind, groupFactory{})
	return &groupClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type GroupMembersGetter interface {
	GroupMembers(namespace string) GroupMemberInterface
}

func (c *Client) GroupMembers(namespace string) GroupMemberInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &GroupMemberResource, GroupMemberGroupVersionKind, groupMemberFactory{})
	return &groupMemberClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type PrincipalsGetter interface {
	Principals(namespace string) PrincipalInterface
}

func (c *Client) Principals(namespace string) PrincipalInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &PrincipalResource, PrincipalGroupVersionKind, principalFactory{})
	return &principalClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type UsersGetter interface {
	Users(namespace string) UserInterface
}

func (c *Client) Users(namespace string) UserInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &UserResource, UserGroupVersionKind, userFactory{})
	return &userClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type AuthConfigsGetter interface {
	AuthConfigs(namespace string) AuthConfigInterface
}

func (c *Client) AuthConfigs(namespace string) AuthConfigInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &AuthConfigResource, AuthConfigGroupVersionKind, authConfigFactory{})
	return &authConfigClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type TokensGetter interface {
	Tokens(namespace string) TokenInterface
}

func (c *Client) Tokens(namespace string) TokenInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &TokenResource, TokenGroupVersionKind, tokenFactory{})
	return &tokenClient{
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

type PreferencesGetter interface {
	Preferences(namespace string) PreferenceInterface
}

func (c *Client) Preferences(namespace string) PreferenceInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &PreferenceResource, PreferenceGroupVersionKind, preferenceFactory{})
	return &preferenceClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ListenConfigsGetter interface {
	ListenConfigs(namespace string) ListenConfigInterface
}

func (c *Client) ListenConfigs(namespace string) ListenConfigInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ListenConfigResource, ListenConfigGroupVersionKind, listenConfigFactory{})
	return &listenConfigClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type SettingsGetter interface {
	Settings(namespace string) SettingInterface
}

func (c *Client) Settings(namespace string) SettingInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &SettingResource, SettingGroupVersionKind, settingFactory{})
	return &settingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ClustersGetter interface {
	Clusters(namespace string) ClusterInterface
}

func (c *Client) Clusters(namespace string) ClusterInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ClusterResource, ClusterGroupVersionKind, clusterFactory{})
	return &clusterClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ClusterEventsGetter interface {
	ClusterEvents(namespace string) ClusterEventInterface
}

func (c *Client) ClusterEvents(namespace string) ClusterEventInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ClusterEventResource, ClusterEventGroupVersionKind, clusterEventFactory{})
	return &clusterEventClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ClusterRegistrationTokensGetter interface {
	ClusterRegistrationTokens(namespace string) ClusterRegistrationTokenInterface
}

func (c *Client) ClusterRegistrationTokens(namespace string) ClusterRegistrationTokenInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &ClusterRegistrationTokenResource, ClusterRegistrationTokenGroupVersionKind, clusterRegistrationTokenFactory{})
	return &clusterRegistrationTokenClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type NodePoolsGetter interface {
	NodePools(namespace string) NodePoolInterface
}

func (c *Client) NodePools(namespace string) NodePoolInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &NodePoolResource, NodePoolGroupVersionKind, nodePoolFactory{})
	return &nodePoolClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type NodesGetter interface {
	Nodes(namespace string) NodeInterface
}

func (c *Client) Nodes(namespace string) NodeInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &NodeResource, NodeGroupVersionKind, nodeFactory{})
	return &nodeClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type NodeDriversGetter interface {
	NodeDrivers(namespace string) NodeDriverInterface
}

func (c *Client) NodeDrivers(namespace string) NodeDriverInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &NodeDriverResource, NodeDriverGroupVersionKind, nodeDriverFactory{})
	return &nodeDriverClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type NodeTemplatesGetter interface {
	NodeTemplates(namespace string) NodeTemplateInterface
}

func (c *Client) NodeTemplates(namespace string) NodeTemplateInterface {
	objectClient := objectclient.NewObjectClient(namespace, c.restClient, &NodeTemplateResource, NodeTemplateGroupVersionKind, nodeTemplateFactory{})
	return &nodeTemplateClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}
