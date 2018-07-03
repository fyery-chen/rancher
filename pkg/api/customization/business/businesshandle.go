package business

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rancher/norman/httperror"
	"github.com/rancher/norman/types"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	businessv1 "github.com/rancher/types/apis/cloud.huawei.com/v1"
	"github.com/rancher/types/config"
	"github.com/rancher/types/user"
	"github.com/sirupsen/logrus"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
)

func NewbusinessHandler(mgmt *config.ScaledContext) *BusinessHandler {
	return &BusinessHandler{
		mgr: mgmt.UserManager,
		clusterClient: mgmt.Management.Clusters(""),
		businessClient: mgmt.Business.BusinessQuotas(""),
	}
}

type BusinessHandler struct {
	mgr user.Manager
	clusterClient v3.ClusterInterface
	businessClient businessv1.BusinessQuotaInterface
}

func (h *BusinessHandler) Checkout(actionName string, action *types.Action, request *types.APIContext) error {
	if actionName != "checkout" {
		return httperror.NewAPIError(httperror.ActionNotAvailable, "")
	}

	err := h.checkoutQuota(request)
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

func (h *BusinessHandler) checkoutQuota(request *types.APIContext) (error) {
	logrus.Debugf("Checkout the business quota")

	bytes, err := ioutil.ReadAll(request.Request.Body)
	if err != nil {
		logrus.Errorf("checkout failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	input := businessv1.BusinessQuotaCheck{}

	err = json.Unmarshal(bytes, input)
	if err != nil {
		logrus.Errorf("unmarshal failed with error: %v", err)
		return httperror.NewAPIError(httperror.InvalidBodyContent, "")
	}

	//check business quota if it is cce provider
	businessName := input.BusinessName
	set := labels.Set{}
	set["businessName"] = businessName
	var requestedCpu int64
	var requestedMemory int64
	clusters, err := h.clusterClient.Controller().Lister().List("", labels.SelectorFromSet(set))
	businesses, err := h.businessClient.Controller().Lister().List("", labels.NewSelector())
	for _, cluster := range clusters {
		requestedCpu += cluster.Status.Requested.Cpu().MilliValue()
		requestedMemory += cluster.Status.Requested.Memory().Value()
	}

	for _, businessQuota := range businesses {
		if businessQuota.Spec.BusinessName == input.BusinessName {
			requestedCpu += int64(input.CpuQuota * input.NodeCount)
			requestedMemory += int64(input.MemoryQuota * input.NodeCount)
			if requestedCpu > int64(businessQuota.Spec.CpuQuota) || requestedMemory > int64(businessQuota.Spec.MemoryQuota) {
				return fmt.Errorf("Checkout failed, there is no enough quotas")
			}
			logrus.Infof("There is enough quota, requested cpu: %d, requested memory: %d, available cpu: %d, available memory: %d", requestedCpu, requestedMemory, businessQuota.Spec.CpuQuota, businessQuota.Spec.MemoryQuota)
			return nil
		}
	}

	return fmt.Errorf("Can not find business: %s", input.BusinessName)
}