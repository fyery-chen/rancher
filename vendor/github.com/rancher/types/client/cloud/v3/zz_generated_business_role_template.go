package client

import (
	"github.com/rancher/norman/types"
)

const (
	BusinessRoleTemplateType                 = "businessRoleTemplate"
	BusinessRoleTemplateFieldAnnotations     = "annotations"
	BusinessRoleTemplateFieldBuiltin         = "builtin"
	BusinessRoleTemplateFieldContext         = "context"
	BusinessRoleTemplateFieldCreated         = "created"
	BusinessRoleTemplateFieldCreatorID       = "creatorId"
	BusinessRoleTemplateFieldDescription     = "description"
	BusinessRoleTemplateFieldExternal        = "external"
	BusinessRoleTemplateFieldHidden          = "hidden"
	BusinessRoleTemplateFieldLabels          = "labels"
	BusinessRoleTemplateFieldLocked          = "locked"
	BusinessRoleTemplateFieldName            = "name"
	BusinessRoleTemplateFieldOwnerReferences = "ownerReferences"
	BusinessRoleTemplateFieldRemoved         = "removed"
	BusinessRoleTemplateFieldRoleTemplateIds = "roleTemplateIds"
	BusinessRoleTemplateFieldRules           = "rules"
	BusinessRoleTemplateFieldUuid            = "uuid"
)

type BusinessRoleTemplate struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Builtin         bool              `json:"builtin,omitempty" yaml:"builtin,omitempty"`
	Context         string            `json:"context,omitempty" yaml:"context,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description     string            `json:"description,omitempty" yaml:"description,omitempty"`
	External        bool              `json:"external,omitempty" yaml:"external,omitempty"`
	Hidden          bool              `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Locked          bool              `json:"locked,omitempty" yaml:"locked,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	RoleTemplateIds []string          `json:"roleTemplateIds,omitempty" yaml:"roleTemplateIds,omitempty"`
	Rules           []PolicyRule      `json:"rules,omitempty" yaml:"rules,omitempty"`
	Uuid            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type BusinessRoleTemplateCollection struct {
	types.Collection
	Data   []BusinessRoleTemplate `json:"data,omitempty"`
	client *BusinessRoleTemplateClient
}

type BusinessRoleTemplateClient struct {
	apiClient *Client
}

type BusinessRoleTemplateOperations interface {
	List(opts *types.ListOpts) (*BusinessRoleTemplateCollection, error)
	Create(opts *BusinessRoleTemplate) (*BusinessRoleTemplate, error)
	Update(existing *BusinessRoleTemplate, updates interface{}) (*BusinessRoleTemplate, error)
	ByID(id string) (*BusinessRoleTemplate, error)
	Delete(container *BusinessRoleTemplate) error
}

func newBusinessRoleTemplateClient(apiClient *Client) *BusinessRoleTemplateClient {
	return &BusinessRoleTemplateClient{
		apiClient: apiClient,
	}
}

func (c *BusinessRoleTemplateClient) Create(container *BusinessRoleTemplate) (*BusinessRoleTemplate, error) {
	resp := &BusinessRoleTemplate{}
	err := c.apiClient.Ops.DoCreate(BusinessRoleTemplateType, container, resp)
	return resp, err
}

func (c *BusinessRoleTemplateClient) Update(existing *BusinessRoleTemplate, updates interface{}) (*BusinessRoleTemplate, error) {
	resp := &BusinessRoleTemplate{}
	err := c.apiClient.Ops.DoUpdate(BusinessRoleTemplateType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *BusinessRoleTemplateClient) List(opts *types.ListOpts) (*BusinessRoleTemplateCollection, error) {
	resp := &BusinessRoleTemplateCollection{}
	err := c.apiClient.Ops.DoList(BusinessRoleTemplateType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *BusinessRoleTemplateCollection) Next() (*BusinessRoleTemplateCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &BusinessRoleTemplateCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *BusinessRoleTemplateClient) ByID(id string) (*BusinessRoleTemplate, error) {
	resp := &BusinessRoleTemplate{}
	err := c.apiClient.Ops.DoByID(BusinessRoleTemplateType, id, resp)
	return resp, err
}

func (c *BusinessRoleTemplateClient) Delete(container *BusinessRoleTemplate) error {
	return c.apiClient.Ops.DoResourceDelete(BusinessRoleTemplateType, &container.Resource)
}
