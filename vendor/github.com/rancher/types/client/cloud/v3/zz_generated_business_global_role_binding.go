package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessGlobalRoleBindingType                 = "businessGlobalRoleBinding"
	BusinessGlobalRoleBindingFieldAnnotations     = "annotations"
	BusinessGlobalRoleBindingFieldCreated         = "created"
	BusinessGlobalRoleBindingFieldCreatorID       = "creatorId"
	BusinessGlobalRoleBindingFieldGlobalRoleId    = "globalRoleId"
	BusinessGlobalRoleBindingFieldLabels          = "labels"
	BusinessGlobalRoleBindingFieldName            = "name"
	BusinessGlobalRoleBindingFieldOwnerReferences = "ownerReferences"
	BusinessGlobalRoleBindingFieldRemoved         = "removed"
	BusinessGlobalRoleBindingFieldUserId          = "userId"
	BusinessGlobalRoleBindingFieldUuid            = "uuid"
)

type BusinessGlobalRoleBinding struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	GlobalRoleId    string            `json:"globalRoleId,omitempty" yaml:"globalRoleId,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	UserId          string            `json:"userId,omitempty" yaml:"userId,omitempty"`
	Uuid            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessGlobalRoleBindingCollection struct {
	types.Collection
	Data   []BusinessGlobalRoleBinding `json:"data,omitempty"`
	client *BusinessGlobalRoleBindingClient
}

type BusinessGlobalRoleBindingClient struct {
	apiClient *Client
}

type BusinessGlobalRoleBindingOperations interface {
	List(opts *types.ListOpts) (*BusinessGlobalRoleBindingCollection, error)
	Create(opts *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error)
	Update(existing *BusinessGlobalRoleBinding, updates interface{}) (*BusinessGlobalRoleBinding, error)
	ByID(id string) (*BusinessGlobalRoleBinding, error)
	Delete(container *BusinessGlobalRoleBinding) error
}

func newBusinessGlobalRoleBindingClient(apiClient *Client) *BusinessGlobalRoleBindingClient {
	return &BusinessGlobalRoleBindingClient{
		apiClient: apiClient,
	}
}

func (c *BusinessGlobalRoleBindingClient) Create(container *BusinessGlobalRoleBinding) (*BusinessGlobalRoleBinding, error) {
	resp := &BusinessGlobalRoleBinding{}
	err := c.apiClient.Ops.DoCreate(BusinessGlobalRoleBindingType, container, resp)
	return resp, err
}

func (c *BusinessGlobalRoleBindingClient) Update(existing *BusinessGlobalRoleBinding, updates interface{}) (*BusinessGlobalRoleBinding, error) {
	resp := &BusinessGlobalRoleBinding{}
	err := c.apiClient.Ops.DoUpdate(BusinessGlobalRoleBindingType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessGlobalRoleBindingClient) List(opts *types.ListOpts) (*BusinessGlobalRoleBindingCollection, error) {
	resp := &BusinessGlobalRoleBindingCollection{}
	err := c.apiClient.Ops.DoList(BusinessGlobalRoleBindingType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessGlobalRoleBindingCollection) Next() (*BusinessGlobalRoleBindingCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessGlobalRoleBindingCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessGlobalRoleBindingClient) ByID(id string) (*BusinessGlobalRoleBinding, error) {
	resp := &BusinessGlobalRoleBinding{}
	err := c.apiClient.Ops.DoByID(BusinessGlobalRoleBindingType, id, resp)
	return resp, err
}

func (c *BusinessGlobalRoleBindingClient) Delete(container *BusinessGlobalRoleBinding) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessGlobalRoleBindingType, &container.Resource)
}
