package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessRoleTemplateBindingType                  = "businessRoleTemplateBinding"
	BusinessRoleTemplateBindingFieldAnnotations      = "annotations"
	BusinessRoleTemplateBindingFieldBusinessId       = "businessId"
	BusinessRoleTemplateBindingFieldCreated          = "created"
	BusinessRoleTemplateBindingFieldCreatorID        = "creatorId"
	BusinessRoleTemplateBindingFieldGroupId          = "groupId"
	BusinessRoleTemplateBindingFieldGroupPrincipalId = "groupPrincipalId"
	BusinessRoleTemplateBindingFieldLabels           = "labels"
	BusinessRoleTemplateBindingFieldName             = "name"
	BusinessRoleTemplateBindingFieldNamespaceId      = "namespaceId"
	BusinessRoleTemplateBindingFieldOwnerReferences  = "ownerReferences"
	BusinessRoleTemplateBindingFieldRemoved          = "removed"
	BusinessRoleTemplateBindingFieldRoleTemplateId   = "roleTemplateId"
	BusinessRoleTemplateBindingFieldUserId           = "userId"
	BusinessRoleTemplateBindingFieldUserPrincipalId  = "userPrincipalId"
	BusinessRoleTemplateBindingFieldUuid             = "uuid"
)

type BusinessRoleTemplateBinding struct {
	types.Resource
	Annotations      map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	BusinessId       string            `json:"businessId,omitempty" yaml:"businessId,omitempty"`
	Created          string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID        string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	GroupId          string            `json:"groupId,omitempty" yaml:"groupId,omitempty"`
	GroupPrincipalId string            `json:"groupPrincipalId,omitempty" yaml:"groupPrincipalId,omitempty"`
	Labels           map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId      string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences  []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed          string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	RoleTemplateId   string            `json:"roleTemplateId,omitempty" yaml:"roleTemplateId,omitempty"`
	UserId           string            `json:"userId,omitempty" yaml:"userId,omitempty"`
	UserPrincipalId  string            `json:"userPrincipalId,omitempty" yaml:"userPrincipalId,omitempty"`
	Uuid             string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessRoleTemplateBindingCollection struct {
	types.Collection
	Data   []BusinessRoleTemplateBinding `json:"data,omitempty"`
	client *BusinessRoleTemplateBindingClient
}

type BusinessRoleTemplateBindingClient struct {
	apiClient *Client
}

type BusinessRoleTemplateBindingOperations interface {
	List(opts *types.ListOpts) (*BusinessRoleTemplateBindingCollection, error)
	Create(opts *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error)
	Update(existing *BusinessRoleTemplateBinding, updates interface{}) (*BusinessRoleTemplateBinding, error)
	ByID(id string) (*BusinessRoleTemplateBinding, error)
	Delete(container *BusinessRoleTemplateBinding) error
}

func newBusinessRoleTemplateBindingClient(apiClient *Client) *BusinessRoleTemplateBindingClient {
	return &BusinessRoleTemplateBindingClient{
		apiClient: apiClient,
	}
}

func (c *BusinessRoleTemplateBindingClient) Create(container *BusinessRoleTemplateBinding) (*BusinessRoleTemplateBinding, error) {
	resp := &BusinessRoleTemplateBinding{}
	err := c.apiClient.Ops.DoCreate(BusinessRoleTemplateBindingType, container, resp)
	return resp, err
}

func (c *BusinessRoleTemplateBindingClient) Update(existing *BusinessRoleTemplateBinding, updates interface{}) (*BusinessRoleTemplateBinding, error) {
	resp := &BusinessRoleTemplateBinding{}
	err := c.apiClient.Ops.DoUpdate(BusinessRoleTemplateBindingType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessRoleTemplateBindingClient) List(opts *types.ListOpts) (*BusinessRoleTemplateBindingCollection, error) {
	resp := &BusinessRoleTemplateBindingCollection{}
	err := c.apiClient.Ops.DoList(BusinessRoleTemplateBindingType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessRoleTemplateBindingCollection) Next() (*BusinessRoleTemplateBindingCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessRoleTemplateBindingCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessRoleTemplateBindingClient) ByID(id string) (*BusinessRoleTemplateBinding, error) {
	resp := &BusinessRoleTemplateBinding{}
	err := c.apiClient.Ops.DoByID(BusinessRoleTemplateBindingType, id, resp)
	return resp, err
}

func (c *BusinessRoleTemplateBindingClient) Delete(container *BusinessRoleTemplateBinding) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessRoleTemplateBindingType, &container.Resource)
}
