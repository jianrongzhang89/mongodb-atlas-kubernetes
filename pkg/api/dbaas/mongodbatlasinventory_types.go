/*
Copyright 2021.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dbaas

import (
	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/metrics"
	kube "github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=dbaas.redhat.com
// +versionName:=v1alpha1

// MongoDBAtlasInventory is the Schema for the MongoDBAtlasInventory API
type MongoDBAtlasInventory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbaasv1alpha1.DBaaSInventorySpec   `json:"spec,omitempty"`
	Status dbaasv1alpha1.DBaaSInventoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBAtlasInventoryList contains a list of DBaaSInventories
type MongoDBAtlasInventoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBAtlasInventory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBAtlasInventory{}, &MongoDBAtlasInventoryList{})
}

func (p *MongoDBAtlasInventory) ConnectionSecretObjectKey() *client.ObjectKey {
	if p.Spec.CredentialsRef != nil {
		key := kube.ObjectKey(p.Spec.CredentialsRef.Namespace, p.Spec.CredentialsRef.Name)
		return &key
	}
	return nil
}

// +k8s:deepcopy-gen=false

// SetInventoryCondition sets a condition on the status object. If the condition already
// exists, it will be replaced. SetCondition does not update the resource in
// the cluster.
func SetInventoryCondition(inv *MongoDBAtlasInventory, condType string, status metav1.ConditionStatus, reason, msg string) {
	now := metav1.Now()
	for i := range inv.Status.Conditions {
		if inv.Status.Conditions[i].Type == condType {
			var lastTransitionTime metav1.Time
			if inv.Status.Conditions[i].Status != status {
				lastTransitionTime = now
			} else {
				lastTransitionTime = inv.Status.Conditions[i].LastTransitionTime
			}
			inv.Status.Conditions[i] = metav1.Condition{
				LastTransitionTime: lastTransitionTime,
				Status:             status,
				Type:               condType,
				Reason:             reason,
				Message:            msg,
			}
			setInventoryMetrics(inv, status, lastTransitionTime, reason)
			return
		}
	}

	// If the condition does not exist,
	// initialize the lastTransitionTime
	inv.Status.Conditions = append(inv.Status.Conditions, metav1.Condition{
		LastTransitionTime: now,
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            msg,
	})
	setInventoryMetrics(inv, status, now, reason)
}

// GetInventoryCondition return the condition with the passed condition type from
// the status object. If the condition is not already present, return nil
func GetInventoryCondition(inv *MongoDBAtlasInventory, condType string) *metav1.Condition {
	for i := range inv.Status.Conditions {
		if inv.Status.Conditions[i].Type == condType {
			return &inv.Status.Conditions[i]
		}
	}
	return nil
}

func setInventoryMetrics(inv *MongoDBAtlasInventory, status metav1.ConditionStatus, lastTransitionTime metav1.Time, reason string) {
	if status == metav1.ConditionTrue {
		setInventoryElapsedTime(inv, lastTransitionTime)
		setInventoryStatusReady(inv, 1)
	} else {
		resetInventoryElapsedTime(inv)
		setInventoryStatusReady(inv, 0)
	}
	setInventoryStatusReason(inv, reason)
}

func setInventoryElapsedTime(inv *MongoDBAtlasInventory, lastTransitionTime metav1.Time) {
	metrics.InventoryElapsedTime.With(prometheus.Labels{
		"provider":  metrics.DBaaSProvider,
		"inventory": inv.Name,
		"namespace": inv.Namespace}).Set(lastTransitionTime.Sub(inv.CreationTimestamp.Time).Seconds())
}

func resetInventoryElapsedTime(inv *MongoDBAtlasInventory) {
	metrics.InventoryElapsedTime.Delete(prometheus.Labels{
		"provider":  metrics.DBaaSProvider,
		"inventory": inv.Name,
		"namespace": inv.Namespace})
}

func setInventoryStatusReady(inventory *MongoDBAtlasInventory, val float64) {
	metrics.InventoryStatusReady.With(prometheus.Labels{
		"provider":  metrics.DBaaSProvider,
		"inventory": inventory.Name,
		"namespace": inventory.Namespace}).Set(val)
}

func resetInventoryStatusReasons(inv *MongoDBAtlasInventory) {
	for _, reason := range workflow.GetMongoDBAtlasInventoryReasons() {
		metrics.InventoryStatusReason.With(prometheus.Labels{
			"provider":  metrics.DBaaSProvider,
			"inventory": inv.Name,
			"namespace": inv.Namespace,
			"reason":    string(reason)}).Set(0)
	}
}

func setInventoryStatusReason(inv *MongoDBAtlasInventory, reason string) {
	resetInventoryStatusReasons(inv)
	metrics.InventoryStatusReason.With(prometheus.Labels{
		"provider":  metrics.DBaaSProvider,
		"inventory": inv.Name,
		"namespace": inv.Namespace,
		"reason":    reason}).Set(1)
}
