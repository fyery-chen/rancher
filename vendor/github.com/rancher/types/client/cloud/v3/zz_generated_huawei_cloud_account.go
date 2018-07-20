package client

import (
	"github.com/rancher/norman/types"
)

const (
	HuaweiCloudAccountType                 = "huaweiCloudAccount"
	HuaweiCloudAccountFieldAccessKey       = "accessKey"
	HuaweiCloudAccountFieldAnnotations     = "annotations"
	HuaweiCloudAccountFieldCreated         = "created"
	HuaweiCloudAccountFieldCreatorID       = "creatorId"
	HuaweiCloudAccountFieldLabels          = "labels"
	HuaweiCloudAccountFieldName            = "name"
	HuaweiCloudAccountFieldNamespaceId     = "namespaceId"
	HuaweiCloudAccountFieldOwnerReferences = "ownerReferences"
	HuaweiCloudAccountFieldRemoved         = "removed"
	HuaweiCloudAccountFieldSecretKey       = "secretKey"
	HuaweiCloudAccountFieldUserId          = "userId"
	HuaweiCloudAccountFieldUuid            = "uuid"
)

type HuaweiCloudAccount struct {
	types.Resource
	AccessKey       string            `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId     string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	SecretKey       string            `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	UserId          string            `json:"userId,omitempty" yaml:"userId,omitempty"`
	Uuid            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}
type HuaweiCloudAccountCollection struct {
	types.Collection
	Data   []HuaweiCloudAccount `json:"data,omitempty"`
	client *HuaweiCloudAccountClient
}

type HuaweiCloudAccountClient struct {
	apiClient *Client
}

type HuaweiCloudAccountOperations interface {
	List(opts *types.ListOpts) (*HuaweiCloudAccountCollection, error)
	Create(opts *HuaweiCloudAccount) (*HuaweiCloudAccount, error)
	Update(existing *HuaweiCloudAccount, updates interface{}) (*HuaweiCloudAccount, error)
	ByID(id string) (*HuaweiCloudAccount, error)
	Delete(container *HuaweiCloudAccount) error
}

func newHuaweiCloudAccountClient(apiClient *Client) *HuaweiCloudAccountClient {
	return &HuaweiCloudAccountClient{
		apiClient: apiClient,
	}
}

func (c *HuaweiCloudAccountClient) Create(container *HuaweiCloudAccount) (*HuaweiCloudAccount, error) {
	resp := &HuaweiCloudAccount{}
	err := c.apiClient.Ops.DoCreate(HuaweiCloudAccountType, container, resp)
	return resp, err
}

func (c *HuaweiCloudAccountClient) Update(existing *HuaweiCloudAccount, updates interface{}) (*HuaweiCloudAccount, error) {
	resp := &HuaweiCloudAccount{}
	err := c.apiClient.Ops.DoUpdate(HuaweiCloudAccountType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *HuaweiCloudAccountClient) List(opts *types.ListOpts) (*HuaweiCloudAccountCollection, error) {
	resp := &HuaweiCloudAccountCollection{}
	err := c.apiClient.Ops.DoList(HuaweiCloudAccountType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *HuaweiCloudAccountCollection) Next() (*HuaweiCloudAccountCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &HuaweiCloudAccountCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *HuaweiCloudAccountClient) ByID(id string) (*HuaweiCloudAccount, error) {
	resp := &HuaweiCloudAccount{}
	err := c.apiClient.Ops.DoByID(HuaweiCloudAccountType, id, resp)
	return resp, err
}

func (c *HuaweiCloudAccountClient) Delete(container *HuaweiCloudAccount) error {
	return c.apiClient.Ops.DoResourceDelete(HuaweiCloudAccountType, &container.Resource)
}
