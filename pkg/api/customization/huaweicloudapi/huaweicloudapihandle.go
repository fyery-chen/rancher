package huaweicloudapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net"
	"time"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	apicorev1 "github.com/rancher/types/apis/core/v1"
	businessv3 "github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/fyery-chen/cce-sdk/signer"
	"github.com/fyery-chen/cce-sdk/common"
	"crypto/tls"
	"fmt"
	"bytes"
	"k8s.io/apimachinery/pkg/labels"
	"github.com/rancher/norman/types/convert"
	"k8s.io/apimachinery/pkg/fields"
)

const (
	HuaweiCloudAccountSecretName = "huaweicloud-cce-account"
)

type state struct {
	ClientID 	 string
	ClientSecret string
	Zone 		 string
	ProjectID    string
}
type VpcList struct {
	Vpcs []common.VpcResp	`json:"vpcs,omitempty"`
}

type SubnetList struct {
	Subnets []common.Subnet	`json:"subnets,omitempty"`
}

type KeypairInfo struct {
	Fingerprint string `json:"fingerprint,omitempty"`
	Name 		string `json:"name,omitempty"`
	Public_key  string `json:"public_key,omitempty"`
}

type Keypair struct {
	Keypair     KeypairInfo 	`json:"keypair,omitempty"`
}

type KeypairList struct {
	Keypairs 	[]Keypair 		`json:"keypairs,omitempty"`
}

type Flavor struct {
	Name  	string `json:"name,omitempty"`
	Vcpus   string `json:"vcpus,omitempty"`
	Ram     int    `json:"ram,omitempty"`
}

type FlavorList struct {
	Flavors  []Flavor `json:"flavors,omitempty"`
}

type ZoneState struct {
	Available  bool `json:"available,omitempty"`
}
type AvailabilityZone struct {
	ZoneState ZoneState `json:"zoneState,omitempty"`
	ZoneName  string 	`json:"zoneName,omitempty"`
}

type AvailabilityZoneInfoList struct {
	AvailabilityZoneInfo []AvailabilityZone `json:"availabilityZoneInfo,omitempty"`
}

func NewHandler(mgmt *config.ScaledContext) *ApiHandler {
	return &ApiHandler{
		mgr: mgmt.UserManager,
		secretClient: mgmt.Core.Secrets("cattle-cce"),
		clusterClient: mgmt.Management.Clusters(""),
		businessClient: mgmt.Business.Businesses(""),
		nodeClient: mgmt.Management.Nodes(""),
	}
}

type ApiHandler struct {
	mgr user.Manager
	secretClient apicorev1.SecretInterface
	clusterClient v3.ClusterInterface
	businessClient businessv3.BusinessInterface
	nodeClient v3.NodeInterface
}

func (h *ApiHandler) HuaweiCloudActionHandler(actionName string, action *types.Action, apiContext *types.APIContext) error {
	switch actionName {
	case "getHuaweiCloudApiInfo":
		return h.GetHuaweiCloudApiInfo(actionName, action, apiContext)
	case "checkout":
		return h.Checkout(actionName, action, apiContext)
	}
	return httperror.NewAPIError(httperror.NotFound, "not found")
}

func (h *ApiHandler) GetHuaweiCloudApiInfo(actionName string, action *types.Action, apiContext *types.APIContext) error {
	if actionName != "getHuaweiCloudApiInfo" {
		return httperror.NewAPIError(httperror.ActionNotAvailable, "")
	}

	err := h.retrieveInfo(apiContext)
	if err != nil {
		// if user fails to authenticate, hide the details of the exact error. bad credentials will already be APIErrors
		// otherwise, return a generic error message
		if httperror.IsAPIError(err) {
			return err
		}
		return httperror.WrapAPIError(err, httperror.ServerError, "Server error while authenticating")
	}

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
	logrus.Infof("request authorization header after signature: %v", req.Header.Get("authorization"))
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

func (h *ApiHandler) retrieveInfo(apiContext *types.APIContext) (error) {
	logrus.Debugf("Retrieve api information from Huawei cloud platform")

	bytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		logrus.Errorf("retrieve failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	input := businessv3.HuaweiCloudApiInformationInput{}

	err = json.Unmarshal(bytes, &input)
	if err != nil {
		logrus.Errorf("unmarshal failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	projectId := input.ProjectId
	zone := input.Zone
	accessKey, secretKey := "", ""
	//accounts, err := h.accountClient.List(v1.ListOptions{LabelSelector: labels.Everything().String()})
	//logrus.Infof("Retrive accounts: %v", accounts)
	//for _, account := range accounts.Items {
	//	accessKey = account.Spec.AccessKey
	//	secretKey = account.Spec.SecretKey
	//}
	accountsSecret, err := h.secretClient.Get(HuaweiCloudAccountSecretName, v1.GetOptions{})
	accessKey = string(accountsSecret.Data["accessKey"])
	secretKey = string(accountsSecret.Data["secretKey"])
	msg := ""
	status := http.StatusOK

	rtn := map[string]interface{}{
		"message": msg,
		"type":    "huaweiCloudApiInformationOutput",
	}
	apiInforOutput := &businessv3.HuaweiCloudApiInformationOutput{}
    logrus.Infof("requested info is: %s, %s, %s, %s", projectId, zone, accessKey, secretKey)
	if projectId == "" || zone == "" || accessKey == "" || secretKey == "" {
		rtn["message"] = "can not find projectId or zone or accessKey or secretKey"
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}

	state := state{
		ClientID: accessKey,
		ClientSecret: secretKey,
		Zone: zone,
		ProjectID: projectId,
	}

	//1.Retrieve vpc and subnet info
	uri := "/v1/" + state.ProjectID + "/vpcs"
	vpcs := &VpcList{}
	subnets := &SubnetList{}
	resp, _, err := cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	err = json.Unmarshal(resp, vpcs)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	var vpcOutputs []businessv3.VpcInfo
	for _, vpc := range vpcs.Vpcs {
		logrus.Infof("status: %s", vpc.Status)
		uri = "/v1/" + state.ProjectID + "/subnets"
		resp, _, err = cceHTTPRequest(state, uri, http.MethodGet, common.ServiceVPC, nil)
		if err != nil {
			rtn["message"] = err.Error()
			status = http.StatusBadRequest
			apiContext.WriteResponse(status, rtn)

			return nil
		}
		err = json.Unmarshal(resp, subnets)
		if err != nil {
			rtn["message"] = err.Error()
			status = http.StatusBadRequest
			apiContext.WriteResponse(status, rtn)

			return nil
		}
		var subnetOutputs  []businessv3.SubnetInfo
		for _, subnet := range subnets.Subnets {
			subnetOutput := businessv3.SubnetInfo{
				SubnetName: subnet.Name,
				SubnetId: subnet.Id,
			}
			subnetOutputs = append(subnetOutputs, subnetOutput)
		}
		vpcOutput := businessv3.VpcInfo{
			VpcName: vpc.Name,
			VpcId: vpc.Id,
			SubnetInfo: subnetOutputs,
		}
		vpcOutputs = append(vpcOutputs, vpcOutput)
	}
	rtn["vpcInfo"] = vpcOutputs
	apiInforOutput.VpcInfo = vpcOutputs
	//2.Retrieve sshkey info
	uri = "/v2/" + state.ProjectID + "/os-keypairs"
	keypairs := &KeypairList{}
	resp, _, err = cceHTTPRequest(state, uri, http.MethodGet, common.ServiceECS, nil)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	err = json.Unmarshal(resp, keypairs)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	var sshkeyNameOutput []string
	for _, keypair := range keypairs.Keypairs {
		sshkeyNameOutput = append(sshkeyNameOutput, keypair.Keypair.Name)
	}
	rtn["sshkeyName"] = sshkeyNameOutput
	apiInforOutput.SshKeyName = sshkeyNameOutput
	//3.Retrieve node flavor
	uri = "/v1/" + state.ProjectID + "/cloudservers/flavors"
	flavors := &FlavorList{}
	resp, _, err = cceHTTPRequest(state, uri, http.MethodGet, common.ServiceECS, nil)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	err = json.Unmarshal(resp, flavors)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	var nodeFlavorOutputs []businessv3.NodeFlavor
	for _, flavor := range flavors.Flavors {
		nodeFlavorOutput := businessv3.NodeFlavor{
			Name: flavor.Name,
			Vcpus: flavor.Vcpus,
			Ram: flavor.Ram,
		}
		nodeFlavorOutputs = append(nodeFlavorOutputs, nodeFlavorOutput)
	}
	rtn["nodeFlavor"] = nodeFlavorOutputs
	apiInforOutput.NodeFlavor = nodeFlavorOutputs
	//4.Retrieve available zone
	uri = "/v2/" + state.ProjectID + "/os-availability-zone"
	azs := &AvailabilityZoneInfoList{}
	resp, _, err = cceHTTPRequest(state, uri, http.MethodGet, common.ServiceECS, nil)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	err = json.Unmarshal(resp, azs)
	if err != nil {
		rtn["message"] = err.Error()
		status = http.StatusBadRequest
		apiContext.WriteResponse(status, rtn)

		return nil
	}
	var availableZoneOutputs []businessv3.AvailableZone
	for _, az := range azs.AvailabilityZoneInfo {
		availableZoneOutput := businessv3.AvailableZone{
			ZoneName: az.ZoneName,
			ZoneState: az.ZoneState.Available,
		}
		availableZoneOutputs = append(availableZoneOutputs, availableZoneOutput)
	}
	rtn["availableZone"] = availableZoneOutputs
	logrus.Infof("%v, %v, %v, %v", vpcOutputs, sshkeyNameOutput,nodeFlavorOutputs, availableZoneOutputs)
	apiInforOutput.AvailableZone = availableZoneOutputs
	rtnOk := map[string]interface{}{}
	convert.ToObj(apiInforOutput, &rtnOk)
	logrus.Infof("apiInfoOutput: %v, rtnOk: %v", *apiInforOutput, rtnOk)
	apiContext.WriteResponse(status, rtnOk)

	return nil
}

func (h *ApiHandler) Checkout(actionName string, action *types.Action, apiContext *types.APIContext) error {
	if actionName != "checkout" {
		return httperror.NewAPIError(httperror.ActionNotAvailable, "")
	}

	err := h.checkoutQuota(apiContext)
	if err != nil {
		// if user fails to authenticate, hide the details of the exact error. bad credentials will already be APIErrors
		// otherwise, return a generic error message
		if httperror.IsAPIError(err) {
			return err
		}
		return httperror.WrapAPIError(err, httperror.ServerError, "Server error while authenticating")
	}

	return nil
}

func (h *ApiHandler) checkoutQuota(apiContext *types.APIContext) (error) {
	logrus.Debugf("Checkout the business quota")

	bytes, err := ioutil.ReadAll(apiContext.Request.Body)
	if err != nil {
		logrus.Errorf("checkout failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	input := businessv3.BusinessQuotaCheck{}

	err = json.Unmarshal(bytes, &input)
	if err != nil {
		logrus.Errorf("unmarshal failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	msg := "There is enough quota"
	status := http.StatusOK

	rtn := map[string]interface{}{
		"message": msg,
		"type":    "businessQuotaCheckOutput",
	}
	//check business quota if it is cce provider
	businessName := input.BusinessName
	set := labels.Set{}
	set["businessName"] = businessName
	business, err := h.businessClient.Get(businessName, v1.GetOptions{})
	if err != nil {
		status = http.StatusBadRequest
		rtn["message"] = "get business error"
		apiContext.WriteResponse(status, rtn)
		return nil
	}
	logrus.Infof("Retrive businesses: %v", business)
	clusters, err := h.clusterClient.List(v1.ListOptions{LabelSelector: set.String()})
	if err != nil {
		status = http.StatusBadRequest
		rtn["message"] = "get clusters error"
		apiContext.WriteResponse(status, rtn)
		return nil
	}
	field := fields.Set{}
	requestedHosts := 0
	for _, cluster := range clusters.Items {
		field["namespace"] = cluster.Name
		nodes, err := h.nodeClient.List(v1.ListOptions{FieldSelector: field.String()})
		if err != nil {
			status = http.StatusBadRequest
			rtn["message"] = "get nodes error"
			apiContext.WriteResponse(status, rtn)
			return nil
		}
		requestedHosts += len(nodes.Items)
	}

	logrus.Infof("Business name: %s input name: %s", business.Name, input.BusinessName)
	requestedHosts += input.NodeCount
	if requestedHosts > business.Spec.NodeCount {
		rtn["message"] = "Checkout failed, there is no enough quotas"
		status = http.StatusBadRequest
	}
	apiContext.WriteResponse(status, rtn)

	return nil
}

func Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "getHuaweiCloudApiInfo")
	resource.AddAction(request, "checkout")
}