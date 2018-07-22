package client

import (
	"github.com/rancher/norman/clientbase"
)

type Client struct {
	clientbase.APIBaseClient

	Group                    GroupOperations
	GroupMember              GroupMemberOperations
	Principal                PrincipalOperations
	User                     UserOperations
	AuthConfig               AuthConfigOperations
	Token                    TokenOperations
	DynamicSchema            DynamicSchemaOperations
	Preference               PreferenceOperations
	ListenConfig             ListenConfigOperations
	Setting                  SettingOperations
	Cluster                  ClusterOperations
	ClusterEvent             ClusterEventOperations
	ClusterRegistrationToken ClusterRegistrationTokenOperations
	NodePool                 NodePoolOperations
	Node                     NodeOperations
	NodeDriver               NodeDriverOperations
	NodeTemplate             NodeTemplateOperations
}

func NewClient(opts *clientbase.ClientOpts) (*Client, error) {
	baseClient, err := clientbase.NewAPIClient(opts)
	if err != nil {
		return nil, err
	}

	client := &Client{
		APIBaseClient: baseClient,
	}

	client.Group = newGroupClient(client)
	client.GroupMember = newGroupMemberClient(client)
	client.Principal = newPrincipalClient(client)
	client.User = newUserClient(client)
	client.AuthConfig = newAuthConfigClient(client)
	client.Token = newTokenClient(client)
	client.DynamicSchema = newDynamicSchemaClient(client)
	client.Preference = newPreferenceClient(client)
	client.ListenConfig = newListenConfigClient(client)
	client.Setting = newSettingClient(client)
	client.Cluster = newClusterClient(client)
	client.ClusterEvent = newClusterEventClient(client)
	client.ClusterRegistrationToken = newClusterRegistrationTokenClient(client)
	client.NodePool = newNodePoolClient(client)
	client.Node = newNodeClient(client)
	client.NodeDriver = newNodeDriverClient(client)
	client.NodeTemplate = newNodeTemplateClient(client)

	return client, nil
}
