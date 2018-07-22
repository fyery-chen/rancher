package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessType                      = "business"
	BusinessFieldAnnotations          = "annotations"
	BusinessFieldCreated              = "created"
	BusinessFieldCreatorID            = "creatorId"
	BusinessFieldDescription          = "description"
	BusinessFieldLabels               = "labels"
	BusinessFieldName                 = "name"
	BusinessFieldNodeCount            = "nodeCount"
	BusinessFieldOwnerReferences      = "ownerReferences"
	BusinessFieldRemoved              = "removed"
	BusinessFieldState                = "state"
	BusinessFieldStatus               = "status"
	BusinessFieldTransitioning        = "transitioning"
	BusinessFieldTransitioningMessage = "transitioningMessage"
	BusinessFieldUuid                 = "uuid"
)

type Business struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description          string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NodeCount            int64             `json:"nodeCount,omitempty" yaml:"nodeCount,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *BusinessStatus   `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	Uuid                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessCollection struct {
	types.Collection
	Data   []Business `json:"data,omitempty"`
	client *BusinessClient
}

type BusinessClient struct {
	apiClient *Client
}

type BusinessOperations interface {
	List(opts *types.ListOpts) (*BusinessCollection, error)
	Create(opts *Business) (*Business, error)
	Update(existing *Business, updates interface{}) (*Business, error)
	ByID(id string) (*Business, error)
	Delete(container *Business) error

	ActionCheckout(resource *Business, input *BusinessQuotaCheck) (*BusinessQuotaCheckOutput, error)

	ActionGetHuaweiCloudApiInfo(resource *Business, input *HuaweiCloudApiInformationInput) (*HuaweiCloudApiInformationOutput, error)
}

func newBusinessClient(apiClient *Client) *BusinessClient {
	return &BusinessClient{
		apiClient: apiClient,
	}
}

func (c *BusinessClient) Create(container *Business) (*Business, error) {
	resp := &Business{}
	err := c.apiClient.Ops.DoCreate(BusinessType, container, resp)
	return resp, err
}

func (c *BusinessClient) Update(existing *Business, updates interface{}) (*Business, error) {
	resp := &Business{}
	err := c.apiClient.Ops.DoUpdate(BusinessType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessClient) List(opts *types.ListOpts) (*BusinessCollection, error) {
	resp := &BusinessCollection{}
	err := c.apiClient.Ops.DoList(BusinessType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessCollection) Next() (*BusinessCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessClient) ByID(id string) (*Business, error) {
	resp := &Business{}
	err := c.apiClient.Ops.DoByID(BusinessType, id, resp)
	return resp, err
}

func (c *BusinessClient) Delete(container *Business) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessType, &container.Resource)
}

func (c *BusinessClient) ActionCheckout(resource *Business, input *BusinessQuotaCheck) (*BusinessQuotaCheckOutput, error) {
	resp := &BusinessQuotaCheckOutput{}
	err := c.apiClient.Ops.DoAction(BusinessType, "checkout", &resource.Resource, input, resp)
	return resp, err
}

func (c *BusinessClient) ActionGetHuaweiCloudApiInfo(resource *Business, input *HuaweiCloudApiInformationInput) (*HuaweiCloudApiInformationOutput, error) {
	resp := &HuaweiCloudApiInformationOutput{}
	err := c.apiClient.Ops.DoAction(BusinessType, "getHuaweiCloudApiInfo", &resource.Resource, input, resp)
	return resp, err
}
