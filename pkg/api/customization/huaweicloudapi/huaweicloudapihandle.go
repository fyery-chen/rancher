package huaweicloudapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net"
	"time"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	businessv3 "github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/fyery-chen/cce-sdk/signer"
	"github.com/fyery-chen/cce-sdk/common"
	"crypto/tls"
	"fmt"
	"bytes"
)

type state struct {
	ClientID 	 string
	ClientSecret string
	Zone 		 string
	ProjectID    string
}
type vpcList struct {
	vpcs []common.VpcResp	`json:"vpcs,omitempty"`
}

type subnetList struct {
	subnets []common.Subnet	`json:"subnets,omitempty"`
}

func NewHandler(mgmt *config.ScaledContext) *ApiHandler {
	return &ApiHandler{
		mgr: mgmt.UserManager,
		accountClient: mgmt.Business.HuaweiCloudAccounts(""),
	}
}

type ApiHandler struct {
	mgr user.Manager
	accountClient businessv3.HuaweiCloudAccountInterface
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
	accounts, err := h.accountClient.List(v1.ListOptions{LabelSelector: labels.Everything().String()})
	logrus.Infof("Retrive accounts: %v", accounts)
	for _, account := range accounts.Items {
		accessKey = account.Spec.AccessKey
		secretKey = account.Spec.SecretKey
	}

	msg := ""
	status := http.StatusOK

	rtn := map[string]interface{}{
		"message": msg,
		"type":    "HuaweiCloudApiInformationOutput",
	}
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
	uri := "/v1/" + state.ProjectID + "/vpcs"
	vpcs := &vpcList{}
	subnets := &subnetList{}
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
	for _, vpc := range vpcs.vpcs {
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
	}

	apiContext.WriteResponse(status, rtn)

	return nil
}

func Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "getHuaweiCloudApiInfo")
}