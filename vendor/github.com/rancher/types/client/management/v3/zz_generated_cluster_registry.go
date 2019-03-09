package client

import (
	"github.com/rancher/norman/types"
)

const (
	ClusterRegistryType                        = "clusterRegistry"
	ClusterRegistryFieldAnnotations            = "annotations"
	ClusterRegistryFieldClusterID              = "clusterId"
	ClusterRegistryFieldCreated                = "created"
	ClusterRegistryFieldCreatorID              = "creatorId"
	ClusterRegistryFieldIncludeSystemComponent = "includeSystemComponent"
	ClusterRegistryFieldLabels                 = "labels"
	ClusterRegistryFieldName                   = "name"
	ClusterRegistryFieldNamespaceId            = "namespaceId"
	ClusterRegistryFieldOwnerReferences        = "ownerReferences"
	ClusterRegistryFieldPassword               = "password"
	ClusterRegistryFieldRegistryState          = "RegistryState"
	ClusterRegistryFieldRegistryType           = "registryType"
	ClusterRegistryFieldRemoved                = "removed"
	ClusterRegistryFieldServerAddress          = "serverAddress"
	ClusterRegistryFieldState                  = "state"
	ClusterRegistryFieldTransitioning          = "transitioning"
	ClusterRegistryFieldTransitioningMessage   = "transitioningMessage"
	ClusterRegistryFieldUUID                   = "uuid"
	ClusterRegistryFieldUserName               = "userName"
)

type ClusterRegistry struct {
	types.Resource
	Annotations            map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClusterID              string            `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	Created                string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID              string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	IncludeSystemComponent *bool             `json:"includeSystemComponent,omitempty" yaml:"includeSystemComponent,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                   string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId            string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences        []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Password               string            `json:"password,omitempty" yaml:"password,omitempty"`
	RegistryState          string            `json:"RegistryState,omitempty" yaml:"RegistryState,omitempty"`
	RegistryType           string            `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	Removed                string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ServerAddress          string            `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	State                  string            `json:"state,omitempty" yaml:"state,omitempty"`
	Transitioning          string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage   string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	UserName               string            `json:"userName,omitempty" yaml:"userName,omitempty"`
}

type ClusterRegistryCollection struct {
	types.Collection
	Data   []ClusterRegistry `json:"data,omitempty"`
	client *ClusterRegistryClient
}

type ClusterRegistryClient struct {
	apiClient *Client
}

type ClusterRegistryOperations interface {
	List(opts *types.ListOpts) (*ClusterRegistryCollection, error)
	Create(opts *ClusterRegistry) (*ClusterRegistry, error)
	Update(existing *ClusterRegistry, updates interface{}) (*ClusterRegistry, error)
	Replace(existing *ClusterRegistry) (*ClusterRegistry, error)
	ByID(id string) (*ClusterRegistry, error)
	Delete(container *ClusterRegistry) error

	CollectionActionGetProjects(resource *ClusterRegistryCollection, input *GetProjectInput) (*GetProjectOutput, error)

	CollectionActionGetRepositories(resource *ClusterRegistryCollection, input *GetRepositoryInput) (*GetRepositoryOutput, error)

	CollectionActionGetRepositoryTags(resource *ClusterRegistryCollection, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error)

	CollectionActionTest(resource *ClusterRegistryCollection, input *ClusterRegistryTestInput) error
}

func newClusterRegistryClient(apiClient *Client) *ClusterRegistryClient {
	return &ClusterRegistryClient{
		apiClient: apiClient,
	}
}

func (c *ClusterRegistryClient) Create(container *ClusterRegistry) (*ClusterRegistry, error) {
	resp := &ClusterRegistry{}
	err := c.apiClient.Ops.DoCreate(ClusterRegistryType, container, resp)
	return resp, err
}

func (c *ClusterRegistryClient) Update(existing *ClusterRegistry, updates interface{}) (*ClusterRegistry, error) {
	resp := &ClusterRegistry{}
	err := c.apiClient.Ops.DoUpdate(ClusterRegistryType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ClusterRegistryClient) Replace(obj *ClusterRegistry) (*ClusterRegistry, error) {
	resp := &ClusterRegistry{}
	err := c.apiClient.Ops.DoReplace(ClusterRegistryType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ClusterRegistryClient) List(opts *types.ListOpts) (*ClusterRegistryCollection, error) {
	resp := &ClusterRegistryCollection{}
	err := c.apiClient.Ops.DoList(ClusterRegistryType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ClusterRegistryCollection) Next() (*ClusterRegistryCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ClusterRegistryCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ClusterRegistryClient) ByID(id string) (*ClusterRegistry, error) {
	resp := &ClusterRegistry{}
	err := c.apiClient.Ops.DoByID(ClusterRegistryType, id, resp)
	return resp, err
}

func (c *ClusterRegistryClient) Delete(container *ClusterRegistry) error {
	return c.apiClient.Ops.DoResourceDelete(ClusterRegistryType, &container.Resource)
}

func (c *ClusterRegistryClient) CollectionActionGetProjects(resource *ClusterRegistryCollection, input *GetProjectInput) (*GetProjectOutput, error) {
	resp := &GetProjectOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ClusterRegistryType, "getProjects", &resource.Collection, input, resp)
	return resp, err
}

func (c *ClusterRegistryClient) CollectionActionGetRepositories(resource *ClusterRegistryCollection, input *GetRepositoryInput) (*GetRepositoryOutput, error) {
	resp := &GetRepositoryOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ClusterRegistryType, "getRepositories", &resource.Collection, input, resp)
	return resp, err
}

func (c *ClusterRegistryClient) CollectionActionGetRepositoryTags(resource *ClusterRegistryCollection, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error) {
	resp := &GetRepositoryTagsOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ClusterRegistryType, "getRepositoryTags", &resource.Collection, input, resp)
	return resp, err
}

func (c *ClusterRegistryClient) CollectionActionTest(resource *ClusterRegistryCollection, input *ClusterRegistryTestInput) error {
	err := c.apiClient.Ops.DoCollectionAction(ClusterRegistryType, "test", &resource.Collection, input, nil)
	return err
}
