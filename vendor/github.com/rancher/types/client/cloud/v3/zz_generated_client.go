package client

import (
	"github.com/rancher/norman/clientbase"
)

type Client struct {
	clientbase.APIBaseClient

	HuaweiCloudAccount          HuaweiCloudAccountOperations
	BusinessGlobalRole          BusinessGlobalRoleOperations
	BusinessGlobalRoleBinding   BusinessGlobalRoleBindingOperations
	BusinessRoleTemplate        BusinessRoleTemplateOperations
	BusinessRoleTemplateBinding BusinessRoleTemplateBindingOperations
	DynamicSchema               DynamicSchemaOperations
	Business                    BusinessOperations
}

func NewClient(opts *clientbase.ClientOpts) (*Client, error) {
	baseClient, err := clientbase.NewAPIClient(opts)
	if err != nil {
		return nil, err
	}

	client := &Client{
		APIBaseClient: baseClient,
	}

	client.HuaweiCloudAccount = newHuaweiCloudAccountClient(client)
	client.BusinessGlobalRole = newBusinessGlobalRoleClient(client)
	client.BusinessGlobalRoleBinding = newBusinessGlobalRoleBindingClient(client)
	client.BusinessRoleTemplate = newBusinessRoleTemplateClient(client)
	client.BusinessRoleTemplateBinding = newBusinessRoleTemplateBindingClient(client)
	client.DynamicSchema = newDynamicSchemaClient(client)
	client.Business = newBusinessClient(client)

	return client, nil
}
