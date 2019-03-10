package registry

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	mgmtv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	mgmtv3client "github.com/rancher/types/client/management/v3"

	"github.com/rancher/rancher/pkg/registry/common"
	"github.com/rancher/rancher/pkg/registry/harbor"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewHandler(mgmt *config.ScaledContext) *APIHandler {
	return &APIHandler{
		mgr:                   mgmt.UserManager,
		clusterClient:         mgmt.Management.Clusters(""),
		globalRegistryClient:  mgmt.Management.GlobalRegistries(""),
		clusterRegistryClient: mgmt.Management.ClusterRegistries(""),
		projectRegistryClient: mgmt.Management.ProjectRegistries(""),
	}
}

type APIHandler struct {
	mgr                   user.Manager
	clusterClient         mgmtv3.ClusterInterface
	globalRegistryClient  mgmtv3.GlobalRegistryInterface
	clusterRegistryClient mgmtv3.ClusterRegistryInterface
	projectRegistryClient mgmtv3.ProjectRegistryInterface
}

func (h *APIHandler) getRegistryCient(registryType string, apiContext *types.APIContext, config *common.APIClientConfig) error {
	var globalRegistry *mgmtv3.GlobalRegistry
	var clusterRegistry *mgmtv3.ClusterRegistry
	var projectRegistry *mgmtv3.ProjectRegistry
	var err error

	parts := strings.Split(apiContext.ID, ":")
	namespace := parts[0]
	name := parts[len(parts)-1]

	if registryType == mgmtv3client.GlobalRegistryType {
		globalRegistry, err = h.globalRegistryClient.Get(name, v1.GetOptions{})
		if err != nil {
			return err
		}
		config.RegistryServer = globalRegistry.Spec.ServerAddress
		config.Username = globalRegistry.Spec.UserName
		config.Password = globalRegistry.Spec.Password
	} else if registryType == mgmtv3client.ClusterRegistryType {
		clusterRegistry, err = h.clusterRegistryClient.GetNamespaced(namespace, name, v1.GetOptions{})
		if err != nil {
			return err
		}
		config.RegistryServer = clusterRegistry.Spec.ServerAddress
		config.Username = clusterRegistry.Spec.UserName
		config.Password = clusterRegistry.Spec.Password
	} else {
		projectRegistry, err = h.projectRegistryClient.GetNamespaced(namespace, name, v1.GetOptions{})
		if err != nil {
			return err
		}
		config.RegistryServer = projectRegistry.Spec.ServerAddress
		config.Username = projectRegistry.Spec.UserName
		config.Password = projectRegistry.Spec.Password
	}
	return nil
}

func (h *APIHandler) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	var apiClientConfig common.APIClientConfig
	var err error

	status := http.StatusOK
	rtn := map[string]interface{}{}

	switch apiContext.Type {
	case mgmtv3client.GlobalRegistryType:
		err = h.getRegistryCient(mgmtv3client.GlobalRegistryType, apiContext, &apiClientConfig)
	case mgmtv3client.ClusterRegistryType:
		err = h.getRegistryCient(mgmtv3client.ClusterRegistryType, apiContext, &apiClientConfig)
	case mgmtv3client.ProjectRegistryType:
		err = h.getRegistryCient(mgmtv3client.ProjectRegistryType, apiContext, &apiClientConfig)
	}
	if err != nil {
		return err
	}

	switch actionName {
	case "getRepositories":
		var repositories []mgmtv3.HarborRepository
		var repositoryOutput mgmtv3.GetRepositoryOutput
		apiClientConfig.RequestType = common.All
		requestClient, err := harbor.NewAPIClient(apiClientConfig)
		if err != nil {
			return err
		}
		resp, err := requestClient.Get()
		if err != nil {
			return err
		}
		err = json.Unmarshal(resp, &repositories)
		if err != nil {
			return err
		}
		repositoryOutput.RepositoryInfo = repositories

		convert.ToObj(repositoryOutput, &rtn)
		rtn["type"] = "getRepositoryOutput"
		apiContext.WriteResponse(status, rtn)

		return nil
	case "getRepositoryTags":
		var repositoryTags []mgmtv3.RepositoryTag
		var repositoryTagsOutput mgmtv3.GetRepositoryTagsOutput
		apiClientConfig.RequestType = common.Tag
		requestClient, err := harbor.NewAPIClient(apiClientConfig)
		if err != nil {
			return err
		}
		resp, err := requestClient.Get()
		if err != nil {
			return httperror.NewAPIError(httperror.ServerError, err.Error())
		}
		err = json.Unmarshal(resp, &repositoryTags)
		if err != nil {
			return err
		}
		repositoryTagsOutput.RepositoryTagsInfo = repositoryTags

		convert.ToObj(repositoryTagsOutput, &rtn)
		rtn["type"] = "getRepositoryTagsOutput"
		apiContext.WriteResponse(status, rtn)

		return nil

	}

	return httperror.NewAPIError(httperror.InvalidAction, "invalid action: "+actionName)

}

func Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "getRepositories")
	resource.AddAction(request, "getRepositoryTags")
}