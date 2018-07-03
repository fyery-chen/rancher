package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessQuotaType                      = "businessQuota"
	BusinessQuotaFieldAnnotations          = "annotations"
	BusinessQuotaFieldBusinessId           = "businessId"
	BusinessQuotaFieldCpuQuota             = "cpuQuota"
	BusinessQuotaFieldCreated              = "created"
	BusinessQuotaFieldCreatorID            = "creatorId"
	BusinessQuotaFieldDescription          = "description"
	BusinessQuotaFieldLabels               = "labels"
	BusinessQuotaFieldMemoryQuota          = "memoryQuota"
	BusinessQuotaFieldName                 = "name"
	BusinessQuotaFieldNamespaceId          = "namespaceId"
	BusinessQuotaFieldOwnerReferences      = "ownerReferences"
	BusinessQuotaFieldRemoved              = "removed"
	BusinessQuotaFieldState                = "state"
	BusinessQuotaFieldStatus               = "status"
	BusinessQuotaFieldTransitioning        = "transitioning"
	BusinessQuotaFieldTransitioningMessage = "transitioningMessage"
	BusinessQuotaFieldUuid                 = "uuid"
)

type BusinessQuota struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	BusinessId           string            `json:"businessId,omitempty" yaml:"businessId,omitempty"`
	CpuQuota             int64             `json:"cpuQuota,omitempty" yaml:"cpuQuota,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MemoryQuota          int64             `json:"memoryQuota,omitempty" yaml:"memoryQuota,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId          string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *BusinessStatus   `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	Uuid                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessQuotaCollection struct {
	types.Collection
	Data   []BusinessQuota `json:"data,omitempty"`
	client *BusinessQuotaClient
}

type BusinessQuotaClient struct {
	apiClient *Client
}

type BusinessQuotaOperations interface {
	List(opts *types.ListOpts) (*BusinessQuotaCollection, error)
	Create(opts *BusinessQuota) (*BusinessQuota, error)
	Update(existing *BusinessQuota, updates interface{}) (*BusinessQuota, error)
	ByID(id string) (*BusinessQuota, error)
	Delete(container *BusinessQuota) error

	ActionCheckout(resource *BusinessQuota, input *BusinessQuotaCheck) error
}

func newBusinessQuotaClient(apiClient *Client) *BusinessQuotaClient {
	return &BusinessQuotaClient{
		apiClient: apiClient,
	}
}

func (c *BusinessQuotaClient) Create(container *BusinessQuota) (*BusinessQuota, error) {
	resp := &BusinessQuota{}
	err := c.apiClient.Ops.DoCreate(BusinessQuotaType, container, resp)
	return resp, err
}

func (c *BusinessQuotaClient) Update(existing *BusinessQuota, updates interface{}) (*BusinessQuota, error) {
	resp := &BusinessQuota{}
	err := c.apiClient.Ops.DoUpdate(BusinessQuotaType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessQuotaClient) List(opts *types.ListOpts) (*BusinessQuotaCollection, error) {
	resp := &BusinessQuotaCollection{}
	err := c.apiClient.Ops.DoList(BusinessQuotaType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessQuotaCollection) Next() (*BusinessQuotaCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessQuotaCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessQuotaClient) ByID(id string) (*BusinessQuota, error) {
	resp := &BusinessQuota{}
	err := c.apiClient.Ops.DoByID(BusinessQuotaType, id, resp)
	return resp, err
}

func (c *BusinessQuotaClient) Delete(container *BusinessQuota) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessQuotaType, &container.Resource)
}

func (c *BusinessQuotaClient) ActionCheckout(resource *BusinessQuota, input *BusinessQuotaCheck) error {
	err := c.apiClient.Ops.DoAction(BusinessQuotaType, "checkout", &resource.Resource, input, nil)
	return err
}
