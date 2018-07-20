package business

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	businessv3 "github.com/rancher/types/apis/cloud.huawei.com/v3"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewbusinessHandler(mgmt *config.ScaledContext) *BusinessHandler {
	return &BusinessHandler{
		mgr: mgmt.UserManager,
		clusterClient: mgmt.Management.Clusters(""),
		businessClient: mgmt.Business.Businesses(""),
	}
}

type BusinessHandler struct {
	mgr user.Manager
	clusterClient v3.ClusterInterface
	businessClient businessv3.BusinessInterface
}

func (h *BusinessHandler) Checkout(actionName string, action *types.Action, apiContext *types.APIContext) error {
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

func (h *BusinessHandler) checkoutQuota(apiContext *types.APIContext) (error) {
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

	//check business quota if it is cce provider
	businessName := input.BusinessName
	set := labels.Set{}
	set["businessName"] = businessName
	var requestedHosts int
	businesses, err := h.businessClient.List(v1.ListOptions{LabelSelector: labels.Everything().String()})
	logrus.Infof("Retrive businesses: %v", businesses)
	clusters, err := h.clusterClient.List(v1.ListOptions{LabelSelector: set.String()})
	for _, _ = range clusters.Items {
		requestedHosts += 10
	}

	msg := ""
	status := http.StatusOK

	for _, businessQuota := range businesses.Items {
		logrus.Infof("Business name: %s input name: %s", businessQuota.Name, input.BusinessName)
		if businessQuota.Name == input.BusinessName {
			requestedHosts += input.NodeCount
			if requestedHosts > businessQuota.Spec.NodeCount {
				msg = "Checkout failed, there is no enough quotas"
				status = http.StatusBadRequest
			}
			msg = "There is enough quota"
		}
	}

	rtn := map[string]interface{}{
		"message": msg,
		"type":    "businessQuotaCheckOutput",
	}

	if rtn["message"] == "" {
		rtn["message"] = "Can not find business"
	}
	apiContext.WriteResponse(status, rtn)

	return nil
}

func Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "checkout")
}