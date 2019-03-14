package workload

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/fyery-chen/cce-sdk/common"
	"github.com/fyery-chen/cce-sdk/signer"
	"github.com/mitchellh/mapstructure"
	"github.com/rancher/norman/api/access"
	"github.com/rancher/norman/api/handler"
	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/convert"
	"github.com/rancher/norman/types/values"
	"github.com/rancher/rancher/pkg/clustermanager"
	apicorev1 "github.com/rancher/types/apis/core/v1"
	managementv3 "github.com/rancher/types/apis/management.cattle.io/v3"
	projectv3 "github.com/rancher/types/apis/project.cattle.io/v3"
	"github.com/rancher/types/apis/project.cattle.io/v3/schema"
	projectschema "github.com/rancher/types/apis/project.cattle.io/v3/schema"
	projectclient "github.com/rancher/types/client/project/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	workloadRevisions              = "revisions"
	DeprecatedRollbackTo           = "deprecated.deployment.rollback.to"
	HuaweiCloudProjectIDSecretName = "huaweicloud-projectId"
	HuaweiCloudAccountSecretName   = "huaweicloud-cce-account"
)

type ActionWrapper struct {
	ClusterManager *clustermanager.Manager
	NodeClient     managementv3.NodeLister
	SecretClient   apicorev1.SecretInterface
}

type state struct {
	ClientID     string
	ClientSecret string
	Zone         string
	ProjectID    string
}

type OutputImageList struct {
	Images []projectv3.ImageInfo
}

func (a ActionWrapper) ActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	if actionName == "imageList" {
		clusterName := a.ClusterManager.ClusterName(apiContext)
		field := fields.Set{}
		field["namespace"] = clusterName
		nodes, err := a.NodeClient.List(clusterName, labels.Everything())
		if err != nil {
			return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error getting cluster nodes information: %s", err.Error()))
		}
		region := nodes[0].Status.NodeLabels["failure-domain.beta.kubernetes.io/region"]

		return a.getImageList(apiContext, region)
	}

	var deployment projectclient.Workload
	accessError := access.ByID(apiContext, &projectschema.Version, "workload", apiContext.ID, &deployment)
	if accessError != nil {
		return httperror.NewAPIError(httperror.InvalidReference, "Error accessing workload")
	}
	namespace, name := splitID(deployment.ID)
	switch actionName {
	case "rollback":
		clusterName := a.ClusterManager.ClusterName(apiContext)
		if clusterName == "" {
			return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Cluster name empty %s", deployment.ID))
		}
		clusterContext, err := a.ClusterManager.UserContext(clusterName)
		if err != nil {
			return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error getting cluster context %s", deployment.ID))
		}
		return a.rollbackDeployment(apiContext, clusterContext, actionName, deployment, namespace, name)

	case "pause":
		if deployment.Paused {
			return httperror.NewAPIError(httperror.InvalidAction, fmt.Sprintf("Deployment %s already paused", deployment.ID))
		}
		return updatePause(apiContext, true, deployment, "pause")

	case "resume":
		if !deployment.Paused {
			return httperror.NewAPIError(httperror.InvalidAction, fmt.Sprintf("Pause deployment %s before resume", deployment.ID))
		}
		return updatePause(apiContext, false, deployment, "resume")
	}
	return nil
}

func (a ActionWrapper) getImageList(apiContext *types.APIContext, region string) error {
	logrus.Debugf("Retrieve api information from Huawei cloud platform")

	bytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		logrus.Errorf("retrieve failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	input := projectv3.CloudProviderImageListInput{}

	err = json.Unmarshal(bytes, &input)
	if err != nil {
		logrus.Errorf("unmarshal failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}
	images := &projectv3.CloudProviderImageListOutput{}
	if input.ProviderType == "cce" {
		projectIDSecret, err := a.SecretClient.Get(HuaweiCloudProjectIDSecretName, v1.GetOptions{})
		if err != nil {
			return writeErrorResponse(apiContext, err.Error())
		}
		projectID := string(projectIDSecret.Data[region])

		accountsSecret, err := a.SecretClient.Get(HuaweiCloudAccountSecretName, v1.GetOptions{})
		if err != nil {
			return writeErrorResponse(apiContext, err.Error())
		}
		accessKey := string(accountsSecret.Data["accessKey"])
		secretKey := string(accountsSecret.Data["secretKey"])

		if projectID == "" || region == "" || accessKey == "" || secretKey == "" {
			errMsg := "can not find projectID or zone or accessKey or secretKey"
			return writeErrorResponse(apiContext, errMsg)
		}

		state := state{
			ClientID:     accessKey,
			ClientSecret: secretKey,
			Zone:         region,
			ProjectID:    projectID,
		}

		//namespace := os.Getenv("CATTLE_SYSTEM_LIBRARY")
		//if namespace == "" {}

		//1.Retrieve vpc and subnet info
		uri := "/v2/manage/repos?filter=center::self"
		output := &OutputImageList{}
		resp, _, err := cceHTTPRequest(state, uri, http.MethodGet, "swr-api", nil)
		if err != nil {
			return writeErrorResponse(apiContext, err.Error())
		}
		err = json.Unmarshal(resp, &output.Images)
		if err != nil {
			return writeErrorResponse(apiContext, err.Error())
		}
		images.Images = output.Images
	} else {
		var tmpImages []projectv3.ImageInfo
		sessions, _ := session.NewSession(&aws.Config{Region: &region})
		svc := ecr.New(sessions)
		inr := &ecr.DescribeRepositoriesInput{}

		resultsr, err := svc.DescribeRepositories(inr)
		if err != nil {
			return writeErrorResponse(apiContext, err.Error())
		}
		for _, repository := range resultsr.Repositories {
			var tmptags []string
			ini := &ecr.ListImagesInput{
				RepositoryName: repository.RepositoryName,
			}
			resultsi, err := svc.ListImages(ini)
			if err != nil {
				return writeErrorResponse(apiContext, err.Error())
			}
			for _, imageID := range resultsi.ImageIds {
				tmptags = append(tmptags, *imageID.ImageTag)
			}
			image := projectv3.ImageInfo{
				Name: *repository.RepositoryName,
				Path: *repository.RepositoryUri,
				Tags: tmptags,
			}
			tmpImages = append(tmpImages, image)
		}
		images.Images = tmpImages
	}
	status := http.StatusOK
	rtnOk := map[string]interface{}{}
	convert.ToObj(images, &rtnOk)
	rtnOk["type"] = "cloudProviderImageListOutput"
	apiContext.WriteResponse(status, rtnOk)

	return nil
}

func writeErrorResponse(apiContext *types.APIContext, errStr string) error {
	msg := ""
	status := http.StatusOK
	rtn := map[string]interface{}{
		"message": msg,
		"type":    "cloudProviderImageListOutput",
	}

	rtn["message"] = errStr
	status = http.StatusBadRequest
	apiContext.WriteResponse(status, rtn)
	return nil
}

func cceHTTPRequest(state state, uri, method, serviceType string, args interface{}) ([]byte, int, error) {
	var b []byte
	var err error
	var req *http.Request
	statusCode := 0

	//initial signer for request signature
	signer := signer.Signer{
		AccessKey: state.ClientID,
		SecretKey: state.ClientSecret,
		Region:    state.Zone,
		Service:   serviceType,
	}

	if args != nil {
		b, err = json.Marshal(args)
		if err != nil {
			return nil, statusCode, err
		}
	}

	//compose request URL
	requestURL := "https://" + serviceType + "." + state.Zone + "." + common.EndPoint + "/" + uri
	if args != nil {
		req, err = http.NewRequest(method, requestURL, bytes.NewReader(b))
	} else {
		req, err = http.NewRequest(method, requestURL, nil)
	}
	if err != nil {
		return nil, statusCode, err
	}

	req.Header.Set("Content-Type", "application/json")

	//signature http request
	//req.Header.Add("date", time.Now().Format(BasicDateFormat))
	if err := signer.Sign(req); err != nil {
		logrus.Errorf("sign error: %v", err)
		return nil, statusCode, err
	}

	//check request details for debuging
	logrus.Debugf("request authorization header after signature: %v", req.Header.Get("authorization"))
	logrus.Infof("request is: %v", req)

	defaultTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	tlsConfig := tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	defaultTransport.TLSClientConfig = &tlsConfig
	httpClient := &http.Client{Transport: defaultTransport}

	//issue request
	resp, err := httpClient.Do(req)
	if err != nil {
		logrus.Errorf("http request error: %v", err)
		return nil, statusCode, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error of response body after I/O reading is: %v", err)
		return nil, statusCode, fmt.Errorf("error reading body: %v", err)
	}

	logrus.Infof("response body after io read: %v", string(body))
	//err code handling
	statusCode = resp.StatusCode
	if statusCode >= 400 && statusCode <= 599 {
		logrus.Errorf("response status code for request is: %d", statusCode)
		return nil, statusCode, fmt.Errorf("Error code: %d", statusCode)
	}
	logrus.Infof("status code is: %d", statusCode)

	return body, statusCode, nil
}

func fetchRevisionFor(apiContext *types.APIContext, rollbackInput *projectclient.DeploymentRollbackInput, namespace string, name string, currRevision string) string {
	rollbackTo := rollbackInput.ReplicaSetID
	if rollbackTo == "" {
		revisionNum, _ := convert.ToNumber(currRevision)
		return convert.ToString(revisionNum - 1)
	}
	data := getRevisions(apiContext, namespace, name, rollbackTo)
	if len(data) > 0 {
		return convert.ToString(values.GetValueN(data[0], "workloadAnnotations", "deployment.kubernetes.io/revision"))
	}
	return ""
}

func getRevisions(apiContext *types.APIContext, namespace string, name string, requestedID string) []map[string]interface{} {
	data, replicaSets := []map[string]interface{}{}, []map[string]interface{}{}
	options := map[string]string{"hidden": "true"}
	conditions := []*types.QueryCondition{
		types.NewConditionFromString("namespaceId", types.ModifierEQ, []string{namespace}...),
	}
	if requestedID != "" {
		// want a specific replicaSet
		conditions = append(conditions, types.NewConditionFromString("id", types.ModifierEQ, []string{requestedID}...))
	}

	if err := access.List(apiContext, &projectschema.Version, projectclient.ReplicaSetType, &types.QueryOptions{Options: options, Conditions: conditions}, &replicaSets); err == nil {
		for _, replicaSet := range replicaSets {
			ownerReferences := convert.ToMapSlice(replicaSet["ownerReferences"])
			for _, ownerReference := range ownerReferences {
				kind := convert.ToString(ownerReference["kind"])
				ownerName := convert.ToString(ownerReference["name"])
				if kind == "Deployment" && name == ownerName {
					data = append(data, replicaSet)
					continue
				}
			}
		}
	}
	return data
}

func updatePause(apiContext *types.APIContext, value bool, deployment projectclient.Workload, actionName string) error {
	data, err := convert.EncodeToMap(deployment)
	if err == nil {
		values.PutValue(data, value, "paused")
		err = update(apiContext, data, deployment.ID)
	}
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error updating workload %s by %s : %s", deployment.ID, actionName, err.Error()))
	}
	apiContext.WriteResponse(http.StatusNoContent, nil)
	return nil
}

func update(apiContext *types.APIContext, data map[string]interface{}, ID string) error {
	workloadSchema := apiContext.Schemas.Schema(&schema.Version, "workload")
	_, err := workloadSchema.Store.Update(apiContext, workloadSchema, data, ID)
	return err
}

func (a ActionWrapper) rollbackDeployment(apiContext *types.APIContext, clusterContext *config.UserContext,
	actionName string, deployment projectclient.Workload, namespace string, name string) error {
	input, err := handler.ParseAndValidateActionBody(apiContext, apiContext.Schemas.Schema(&projectschema.Version,
		projectclient.DeploymentRollbackInputType))
	if err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse action body: %v", err))
	}
	rollbackInput := &projectclient.DeploymentRollbackInput{}
	if err := mapstructure.Decode(input, rollbackInput); err != nil {
		return httperror.NewAPIError(httperror.InvalidBodyContent,
			fmt.Sprintf("Failed to parse body: %v", err))
	}
	currRevision := deployment.WorkloadAnnotations["deployment.kubernetes.io/revision"]
	if currRevision == "1" {
		httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("No revision for rolling back %s", deployment.ID))
	}
	revision := fetchRevisionFor(apiContext, rollbackInput, namespace, name, currRevision)
	logrus.Debugf("rollbackInput %v", revision)
	if revision == "" {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("ReplicaSet %s doesn't exist for deployment %s", rollbackInput.ReplicaSetID, deployment.ID))
	}
	revisionNum, err := convert.ToNumber(revision)
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error getting revision number %s for %s : %s", revision, deployment.ID, err.Error()))
	}
	data := map[string]interface{}{}
	data["kind"] = "DeploymentRollback"
	data["apiVersion"] = "extensions/v1beta1"
	data["name"] = name
	data["rollbackTo"] = map[string]interface{}{"revision": revisionNum}
	deploymentRollback, err := json.Marshal(data)
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error getting DeploymentRollback for %s %s", rollbackInput.ReplicaSetID, err.Error()))
	}
	err = clusterContext.UnversionedClient.Post().Prefix("apis/extensions/v1beta1/").Namespace(namespace).
		Resource("deployments").Name(name).SubResource("rollback").Body(deploymentRollback).Do().Error()
	if err != nil {
		return httperror.NewAPIError(httperror.ServerError, fmt.Sprintf("Error updating workload %s by %s : %s", deployment.ID, actionName, err.Error()))
	}
	apiContext.WriteResponse(http.StatusNoContent, nil)
	return nil
}

func (h Handler) LinkHandler(apiContext *types.APIContext, next types.RequestHandler) error {
	if apiContext.Link == workloadRevisions {
		var deployment projectclient.Workload
		if err := access.ByID(apiContext, &projectschema.Version, "workload", apiContext.ID, &deployment); err == nil {
			namespace, deploymentName := splitID(deployment.ID)
			data := getRevisions(apiContext, namespace, deploymentName, "")
			apiContext.Type = projectclient.ReplicaSetType
			apiContext.WriteResponse(http.StatusOK, data)
		}
		return nil
	}
	return httperror.NewAPIError(httperror.NotFound, "Link not found")
}

func Formatter(apiContext *types.APIContext, resource *types.RawResource) {
	workloadID := resource.ID
	workloadSchema := apiContext.Schemas.Schema(&schema.Version, "workload")
	resource.Links["self"] = apiContext.URLBuilder.ResourceLinkByID(workloadSchema, workloadID)
	resource.Links["remove"] = apiContext.URLBuilder.ResourceLinkByID(workloadSchema, workloadID)
	resource.Links["update"] = apiContext.URLBuilder.ResourceLinkByID(workloadSchema, workloadID)

	delete(resource.Values, "nodeId")
}

func DeploymentFormatter(apiContext *types.APIContext, resource *types.RawResource) {
	workloadID := resource.ID
	workloadSchema := apiContext.Schemas.Schema(&schema.Version, "workload")
	Formatter(apiContext, resource)
	resource.Links["revisions"] = apiContext.URLBuilder.ResourceLinkByID(workloadSchema, workloadID) + "/" + workloadRevisions
	resource.Actions["pause"] = apiContext.URLBuilder.ActionLinkByID(workloadSchema, workloadID, "pause")
	resource.Actions["resume"] = apiContext.URLBuilder.ActionLinkByID(workloadSchema, workloadID, "resume")
	resource.Actions["rollback"] = apiContext.URLBuilder.ActionLinkByID(workloadSchema, workloadID, "rollback")
}

type Handler struct {
}

func splitID(id string) (string, string) {
	namespace := ""
	parts := strings.SplitN(id, ":", 3)
	if len(parts) == 3 {
		namespace = parts[1]
		id = parts[2]
	}

	return namespace, id
}

func CollectionFormatter(request *types.APIContext, collection *types.GenericCollection) {
	collection.AddAction(request, "imageList")
}
