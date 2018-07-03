package v1

import (
	"github.com/rancher/norman/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BusinessQuota struct {
	types.Namespaced

	metav1.TypeMeta `json:",inline"`
	// Standard object’s metadata. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BusinessQuotaSpec `json:"spec"`
	// Most recent observed status of the alert. More info:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status
	Status BusinessStatus `json:"status"`
}

type BusinessStatus struct {
	BusinessState string `json:"businessState,omitempty" norman:"options=transferring|transferError|active,default=active"`
}

type BusinessQuotaSpec struct {
	DisplayName  string `json:"displayName,omitempty" norman:"required"`
	BusinessName string `json:"businessName,omitempty" norman:"type=reference[business]"`
	Description  string `json:"description,omitempty"`
	CpuQuota     int    `json:"cpuQuota,omitempty"`
	MemoryQuota  int    `json:"memoryQuota,omitempty"`
}

type BusinessQuotaCheck struct {
	BusinessName string `json:"businessName,omitempty" norman:"type=string,required"`
	CpuQuota     int    `json:"cpuQuota,omitempty" norman:"required"`
	MemoryQuota  int    `json:"memoryQuota,omitempty" norman:"required"`
	NodeCount    int    `json:"nodeCount,omitempty" norman:"required"`
}
