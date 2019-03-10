package client

import (
	"github.com/rancher/norman/types"
)

const (
	GlobalRegistryType                      = "globalRegistry"
	GlobalRegistryFieldAnnotations          = "annotations"
	GlobalRegistryFieldCreated              = "created"
	GlobalRegistryFieldCreatorID            = "creatorId"
	GlobalRegistryFieldLabels               = "labels"
	GlobalRegistryFieldName                 = "name"
	GlobalRegistryFieldOwnerReferences      = "ownerReferences"
	GlobalRegistryFieldPassword             = "password"
	GlobalRegistryFieldRegistryState        = "RegistryState"
	GlobalRegistryFieldRegistryType         = "registryType"
	GlobalRegistryFieldRemoved              = "removed"
	GlobalRegistryFieldServerAddress        = "serverAddress"
	GlobalRegistryFieldState                = "state"
	GlobalRegistryFieldTransitioning        = "transitioning"
	GlobalRegistryFieldTransitioningMessage = "transitioningMessage"
	GlobalRegistryFieldUUID                 = "uuid"
	GlobalRegistryFieldUserName             = "userName"
)

type GlobalRegistry struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Password             string            `json:"password,omitempty" yaml:"password,omitempty"`
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

type GlobalRegistryCollection struct {
	types.Collection
	Data   []GlobalRegistry `json:"data,omitempty"`
	client *GlobalRegistryClient
}

type GlobalRegistryClient struct {
	apiClient *Client
}

type GlobalRegistryOperations interface {
	List(opts *types.ListOpts) (*GlobalRegistryCollection, error)
	Create(opts *GlobalRegistry) (*GlobalRegistry, error)
	Update(existing *GlobalRegistry, updates interface{}) (*GlobalRegistry, error)
	Replace(existing *GlobalRegistry) (*GlobalRegistry, error)
	ByID(id string) (*GlobalRegistry, error)
	Delete(container *GlobalRegistry) error

	ActionGetProjects(resource *GlobalRegistry, input *GetProjectInput) (*GetProjectOutput, error)

	ActionGetRepositories(resource *GlobalRegistry, input *GetRepositoryInput) (*GetRepositoryOutput, error)

	ActionGetRepositoryTags(resource *GlobalRegistry, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error)

	ActionTest(resource *GlobalRegistry, input *GlobalRegistryTestInput) error
}

func newGlobalRegistryClient(apiClient *Client) *GlobalRegistryClient {
	return &GlobalRegistryClient{
		apiClient: apiClient,
	}
}

func (c *GlobalRegistryClient) Create(container *GlobalRegistry) (*GlobalRegistry, error) {
	resp := &GlobalRegistry{}
	err := c.apiClient.Ops.DoCreate(GlobalRegistryType, container, resp)
	return resp, err
}

func (c *GlobalRegistryClient) Update(existing *GlobalRegistry, updates interface{}) (*GlobalRegistry, error) {
	resp := &GlobalRegistry{}
	err := c.apiClient.Ops.DoUpdate(GlobalRegistryType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *GlobalRegistryClient) Replace(obj *GlobalRegistry) (*GlobalRegistry, error) {
	resp := &GlobalRegistry{}
	err := c.apiClient.Ops.DoReplace(GlobalRegistryType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *GlobalRegistryClient) List(opts *types.ListOpts) (*GlobalRegistryCollection, error) {
	resp := &GlobalRegistryCollection{}
	err := c.apiClient.Ops.DoList(GlobalRegistryType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *GlobalRegistryCollection) Next() (*GlobalRegistryCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &GlobalRegistryCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *GlobalRegistryClient) ByID(id string) (*GlobalRegistry, error) {
	resp := &GlobalRegistry{}
	err := c.apiClient.Ops.DoByID(GlobalRegistryType, id, resp)
	return resp, err
}

func (c *GlobalRegistryClient) Delete(container *GlobalRegistry) error {
	return c.apiClient.Ops.DoResourceDelete(GlobalRegistryType, &container.Resource)
}

func (c *GlobalRegistryClient) ActionGetProjects(resource *GlobalRegistry, input *GetProjectInput) (*GetProjectOutput, error) {
	resp := &GetProjectOutput{}
	err := c.apiClient.Ops.DoAction(GlobalRegistryType, "getProjects", &resource.Resource, input, resp)
	return resp, err
}

func (c *GlobalRegistryClient) ActionGetRepositories(resource *GlobalRegistry, input *GetRepositoryInput) (*GetRepositoryOutput, error) {
	resp := &GetRepositoryOutput{}
	err := c.apiClient.Ops.DoAction(GlobalRegistryType, "getRepositories", &resource.Resource, input, resp)
	return resp, err
}

func (c *GlobalRegistryClient) ActionGetRepositoryTags(resource *GlobalRegistry, input *GetRepositoryTagsInput) (*GetRepositoryTagsOutput, error) {
	resp := &GetRepositoryTagsOutput{}
	err := c.apiClient.Ops.DoAction(GlobalRegistryType, "getRepositoryTags", &resource.Resource, input, resp)
	return resp, err
}

func (c *GlobalRegistryClient) ActionTest(resource *GlobalRegistry, input *GlobalRegistryTestInput) error {
	err := c.apiClient.Ops.DoAction(GlobalRegistryType, "test", &resource.Resource, input, nil)
	return err
}
