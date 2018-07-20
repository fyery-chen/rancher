package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessGlobalRoleType                 = "businessGlobalRole"
	BusinessGlobalRoleFieldAnnotations     = "annotations"
	BusinessGlobalRoleFieldBuiltin         = "builtin"
	BusinessGlobalRoleFieldCreated         = "created"
	BusinessGlobalRoleFieldCreatorID       = "creatorId"
	BusinessGlobalRoleFieldDescription     = "description"
	BusinessGlobalRoleFieldLabels          = "labels"
	BusinessGlobalRoleFieldName            = "name"
	BusinessGlobalRoleFieldOwnerReferences = "ownerReferences"
	BusinessGlobalRoleFieldRemoved         = "removed"
	BusinessGlobalRoleFieldRules           = "rules"
	BusinessGlobalRoleFieldUuid            = "uuid"
)

type BusinessGlobalRole struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Builtin         bool              `json:"builtin,omitempty" yaml:"builtin,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Rules           []PolicyRule      `json:"rules,omitempty" yaml:"rules,omitempty"`
	Uuid            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessGlobalRoleCollection struct {
	types.Collection
	Data   []BusinessGlobalRole `json:"data,omitempty"`
	client *BusinessGlobalRoleClient
}

type BusinessGlobalRoleClient struct {
	apiClient *Client
}

type BusinessGlobalRoleOperations interface {
	List(opts *types.ListOpts) (*BusinessGlobalRoleCollection, error)
	Create(opts *BusinessGlobalRole) (*BusinessGlobalRole, error)
	Update(existing *BusinessGlobalRole, updates interface{}) (*BusinessGlobalRole, error)
	ByID(id string) (*BusinessGlobalRole, error)
	Delete(container *BusinessGlobalRole) error
}

func newBusinessGlobalRoleClient(apiClient *Client) *BusinessGlobalRoleClient {
	return &BusinessGlobalRoleClient{
		apiClient: apiClient,
	}
}

func (c *BusinessGlobalRoleClient) Create(container *BusinessGlobalRole) (*BusinessGlobalRole, error) {
	resp := &BusinessGlobalRole{}
	err := c.apiClient.Ops.DoCreate(BusinessGlobalRoleType, container, resp)
	return resp, err
}

func (c *BusinessGlobalRoleClient) Update(existing *BusinessGlobalRole, updates interface{}) (*BusinessGlobalRole, error) {
	resp := &BusinessGlobalRole{}
	err := c.apiClient.Ops.DoUpdate(BusinessGlobalRoleType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessGlobalRoleClient) List(opts *types.ListOpts) (*BusinessGlobalRoleCollection, error) {
	resp := &BusinessGlobalRoleCollection{}
	err := c.apiClient.Ops.DoList(BusinessGlobalRoleType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessGlobalRoleCollection) Next() (*BusinessGlobalRoleCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessGlobalRoleCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessGlobalRoleClient) ByID(id string) (*BusinessGlobalRole, error) {
	resp := &BusinessGlobalRole{}
	err := c.apiClient.Ops.DoByID(BusinessGlobalRoleType, id, resp)
	return resp, err
}

func (c *BusinessGlobalRoleClient) Delete(container *BusinessGlobalRole) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessGlobalRoleType, &container.Resource)
}
