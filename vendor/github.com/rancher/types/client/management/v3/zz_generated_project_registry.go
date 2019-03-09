package client

import (
	"github.com/rancher/norman/types"
)

const (
	ProjectRegistryType                      = "projectRegistry"
	ProjectRegistryFieldAnnotations          = "annotations"
	ProjectRegistryFieldCreated              = "created"
	ProjectRegistryFieldCreatorID            = "creatorId"
	ProjectRegistryFieldLabels               = "labels"
	ProjectRegistryFieldName                 = "name"
	ProjectRegistryFieldNamespaceId          = "namespaceId"
	ProjectRegistryFieldOwnerReferences      = "ownerReferences"
	ProjectRegistryFieldPassword             = "password"
	ProjectRegistryFieldProjectID            = "projectId"
	ProjectRegistryFieldRegistryState        = "RegistryState"
	ProjectRegistryFieldRegistryType         = "registryType"
	ProjectRegistryFieldRemoved              = "removed"
	ProjectRegistryFieldServerAddress        = "serverAddress"
	ProjectRegistryFieldState                = "state"
	ProjectRegistryFieldTransitioning        = "transitioning"
	ProjectRegistryFieldTransitioningMessage = "transitioningMessage"
	ProjectRegistryFieldUUID                 = "uuid"
	ProjectRegistryFieldUserName             = "userName"
)

type ProjectRegistry struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Password             string            `json:"password,omitempty" yaml:"password,omitempty"`
	ProjectID            string            `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	RegistryState        string            `json:"RegistryState,omitempty" yaml:"RegistryState,omitempty"`
	RegistryType         string            `json:"registryType,omitempty" yaml:"registryType,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ServerAddress        string            `json:"serverAddress,omitempty" yaml:"serverAddress,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	UserName             string            `json:"userName,omitempty" yaml:"userName,omitempty"`
}

type ProjectRegistryCollection struct {
	types.Collection
	Data   []ProjectRegistry `json:"data,omitempty"`
	client *ProjectRegistryClient
}

type ProjectRegistryClient struct {
	apiClient *Client
}

type ProjectRegistryOperations interface {
	List(opts *types.ListOpts) (*ProjectRegistryCollection, error)
	Create(opts *ProjectRegistry) (*ProjectRegistry, error)
	Update(existing *ProjectRegistry, updates interface{}) (*ProjectRegistry, error)
	Replace(existing *ProjectRegistry) (*ProjectRegistry, error)
	ByID(id string) (*ProjectRegistry, error)
	Delete(container *ProjectRegistry) error

	CollectionActionGetProjects(resource *ProjectRegistryCollection, input *GetProjectInput) (*GetProjectOutput, error)

	CollectionActionGetRepositories(resource *ProjectRegistryCollection, input *GetRepositoryInput) (*GetRepositoryOutput, error)

	CollectionActionGetRepositoryTags(resource *ProjectRegistryCollection, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error)

	CollectionActionTest(resource *ProjectRegistryCollection, input *ProjectRegistryTestInput) error
}

func newProjectRegistryClient(apiClient *Client) *ProjectRegistryClient {
	return &ProjectRegistryClient{
		apiClient: apiClient,
	}
}

func (c *ProjectRegistryClient) Create(container *ProjectRegistry) (*ProjectRegistry, error) {
	resp := &ProjectRegistry{}
	err := c.apiClient.Ops.DoCreate(ProjectRegistryType, container, resp)
	return resp, err
}

func (c *ProjectRegistryClient) Update(existing *ProjectRegistry, updates interface{}) (*ProjectRegistry, error) {
	resp := &ProjectRegistry{}
	err := c.apiClient.Ops.DoUpdate(ProjectRegistryType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *ProjectRegistryClient) Replace(obj *ProjectRegistry) (*ProjectRegistry, error) {
	resp := &ProjectRegistry{}
	err := c.apiClient.Ops.DoReplace(ProjectRegistryType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *ProjectRegistryClient) List(opts *types.ListOpts) (*ProjectRegistryCollection, error) {
	resp := &ProjectRegistryCollection{}
	err := c.apiClient.Ops.DoList(ProjectRegistryType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *ProjectRegistryCollection) Next() (*ProjectRegistryCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &ProjectRegistryCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *ProjectRegistryClient) ByID(id string) (*ProjectRegistry, error) {
	resp := &ProjectRegistry{}
	err := c.apiClient.Ops.DoByID(ProjectRegistryType, id, resp)
	return resp, err
}

func (c *ProjectRegistryClient) Delete(container *ProjectRegistry) error {
	return c.apiClient.Ops.DoResourceDelete(ProjectRegistryType, &container.Resource)
}

func (c *ProjectRegistryClient) CollectionActionGetProjects(resource *ProjectRegistryCollection, input *GetProjectInput) (*GetProjectOutput, error) {
	resp := &GetProjectOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ProjectRegistryType, "getProjects", &resource.Collection, input, resp)
	return resp, err
}

func (c *ProjectRegistryClient) CollectionActionGetRepositories(resource *ProjectRegistryCollection, input *GetRepositoryInput) (*GetRepositoryOutput, error) {
	resp := &GetRepositoryOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ProjectRegistryType, "getRepositories", &resource.Collection, input, resp)
	return resp, err
}

func (c *ProjectRegistryClient) CollectionActionGetRepositoryTags(resource *ProjectRegistryCollection, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error) {
	resp := &GetRepositoryTagsOutput{}
	err := c.apiClient.Ops.DoCollectionAction(ProjectRegistryType, "getRepositoryTags", &resource.Collection, input, resp)
	return resp, err
}

func (c *ProjectRegistryClient) CollectionActionTest(resource *ProjectRegistryCollection, input *ProjectRegistryTestInput) error {
	err := c.apiClient.Ops.DoCollectionAction(ProjectRegistryType, "test", &resource.Collection, input, nil)
	return err
}
