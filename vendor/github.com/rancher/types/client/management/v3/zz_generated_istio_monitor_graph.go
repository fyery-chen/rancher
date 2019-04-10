package client

import (
	"github.com/rancher/norman/types"
)

const (
	IstioMonitorGraphType                        = "istioMonitorGraph"
	IstioMonitorGraphFieldAnnotations            = "annotations"
	IstioMonitorGraphFieldClusterID              = "clusterId"
	IstioMonitorGraphFieldCreated                = "created"
	IstioMonitorGraphFieldCreatorID              = "creatorId"
	IstioMonitorGraphFieldDescription            = "description"
	IstioMonitorGraphFieldDetailsMetricsSelector = "detailsMetricsSelector"
	IstioMonitorGraphFieldDisplayResourceType    = "displayResourceType"
	IstioMonitorGraphFieldGraphType              = "graphType"
	IstioMonitorGraphFieldLabels                 = "labels"
	IstioMonitorGraphFieldMetricsSelector        = "metricsSelector"
	IstioMonitorGraphFieldName                   = "name"
	IstioMonitorGraphFieldNamespaceId            = "namespaceId"
	IstioMonitorGraphFieldOwnerReferences        = "ownerReferences"
	IstioMonitorGraphFieldPriority               = "priority"
	IstioMonitorGraphFieldRemoved                = "removed"
	IstioMonitorGraphFieldResourceType           = "resourceType"
	IstioMonitorGraphFieldUUID                   = "uuid"
	IstioMonitorGraphFieldYAxis                  = "yAxis"
)

type IstioMonitorGraph struct {
	types.Resource
	Annotations            map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClusterID              string            `json:"clusterId,omitempty" yaml:"clusterId,omitempty"`
	Created                string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID              string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Description            string            `json:"description,omitempty" yaml:"description,omitempty"`
	DetailsMetricsSelector map[string]string `json:"detailsMetricsSelector,omitempty" yaml:"detailsMetricsSelector,omitempty"`
	DisplayResourceType    string            `json:"displayResourceType,omitempty" yaml:"displayResourceType,omitempty"`
	GraphType              string            `json:"graphType,omitempty" yaml:"graphType,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MetricsSelector        map[string]string `json:"metricsSelector,omitempty" yaml:"metricsSelector,omitempty"`
	Name                   string            `json:"name,omitempty" yaml:"name,omitempty"`
	NamespaceId            string            `json:"namespaceId,omitempty" yaml:"namespaceId,omitempty"`
	OwnerReferences        []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Priority               int64             `json:"priority,omitempty" yaml:"priority,omitempty"`
	Removed                string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	ResourceType           string            `json:"resourceType,omitempty" yaml:"resourceType,omitempty"`
	UUID                   string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	YAxis                  *YAxis            `json:"yAxis,omitempty" yaml:"yAxis,omitempty"`
}

type IstioMonitorGraphCollection struct {
	types.Collection
	Data   []IstioMonitorGraph `json:"data,omitempty"`
	client *IstioMonitorGraphClient
}

type IstioMonitorGraphClient struct {
	apiClient *Client
}

type IstioMonitorGraphOperations interface {
	List(opts *types.ListOpts) (*IstioMonitorGraphCollection, error)
	Create(opts *IstioMonitorGraph) (*IstioMonitorGraph, error)
	Update(existing *IstioMonitorGraph, updates interface{}) (*IstioMonitorGraph, error)
	Replace(existing *IstioMonitorGraph) (*IstioMonitorGraph, error)
	ByID(id string) (*IstioMonitorGraph, error)
	Delete(container *IstioMonitorGraph) error

	CollectionActionQuery(resource *IstioMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioGraphOutput, error)
}

func newIstioMonitorGraphClient(apiClient *Client) *IstioMonitorGraphClient {
	return &IstioMonitorGraphClient{
		apiClient: apiClient,
	}
}

func (c *IstioMonitorGraphClient) Create(container *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	resp := &IstioMonitorGraph{}
	err := c.apiClient.Ops.DoCreate(IstioMonitorGraphType, container, resp)
	return resp, err
}

func (c *IstioMonitorGraphClient) Update(existing *IstioMonitorGraph, updates interface{}) (*IstioMonitorGraph, error) {
	resp := &IstioMonitorGraph{}
	err := c.apiClient.Ops.DoUpdate(IstioMonitorGraphType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *IstioMonitorGraphClient) Replace(obj *IstioMonitorGraph) (*IstioMonitorGraph, error) {
	resp := &IstioMonitorGraph{}
	err := c.apiClient.Ops.DoReplace(IstioMonitorGraphType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *IstioMonitorGraphClient) List(opts *types.ListOpts) (*IstioMonitorGraphCollection, error) {
	resp := &IstioMonitorGraphCollection{}
	err := c.apiClient.Ops.DoList(IstioMonitorGraphType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *IstioMonitorGraphCollection) Next() (*IstioMonitorGraphCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &IstioMonitorGraphCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *IstioMonitorGraphClient) ByID(id string) (*IstioMonitorGraph, error) {
	resp := &IstioMonitorGraph{}
	err := c.apiClient.Ops.DoByID(IstioMonitorGraphType, id, resp)
	return resp, err
}

func (c *IstioMonitorGraphClient) Delete(container *IstioMonitorGraph) error {
	return c.apiClient.Ops.DoResourceDelete(IstioMonitorGraphType, &container.Resource)
}

func (c *IstioMonitorGraphClient) CollectionActionQuery(resource *IstioMonitorGraphCollection, input *QueryGraphInput) (*QueryIstioGraphOutput, error) {
	resp := &QueryIstioGraphOutput{}
	err := c.apiClient.Ops.DoCollectionAction(IstioMonitorGraphType, "query", &resource.Collection, input, resp)
	return resp, err
}
