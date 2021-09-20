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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dbaasv1alpha1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	CloudProviderKey                = "providerName"
	CloudRegionKey                  = "regionName"
	ProjectIDKey                    = "projectID"
	ProjectNameKey                  = "projectName"
	InstanceSizeNameKey             = "instanceSizeName"
	ConnectionStringsStandardSrvKey = "connectionStringsStandardSrv"
	InstanceIDKey                   = "instanceID"
	HostKey                         = "host"
	SrvKey                          = "srv"
	ProviderKey                     = "provider"
	Provider                        = "Red Hat DBaaS / MongoDB Atlas"
	ServiceBindingTypeKey           = "type"
	ServiceBindingType              = "mongodb"
	DefaultDatabase                 = "admin"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=dbaas.redhat.com
// +versionName:=v1alpha1

// MongoDBAtlasConnection is the Schema for the MongoDBAtlasConnections API
type MongoDBAtlasConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbaasv1alpha1.DBaaSConnectionSpec   `json:"spec,omitempty"`
	Status dbaasv1alpha1.DBaaSConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBAtlasConnectionList contains a list of MongoDBAtlasConnection
type MongoDBAtlasConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBAtlasConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBAtlasConnection{}, &MongoDBAtlasConnectionList{})
}

// SetConnectionCondition sets a condition on the status object. If the condition already
// exists, it will be replaced. SetCondition does not update the resource in
// the cluster.
func SetConnectionCondition(conn *MongoDBAtlasConnection, condType string, status metav1.ConditionStatus, reason, msg string) {
	now := metav1.Now()
	for i := range conn.Status.Conditions {
		if conn.Status.Conditions[i].Type == condType {
			var lastTransitionTime metav1.Time
			if conn.Status.Conditions[i].Status != status {
				lastTransitionTime = now
			} else {
				lastTransitionTime = conn.Status.Conditions[i].LastTransitionTime
			}
			conn.Status.Conditions[i] = metav1.Condition{
				LastTransitionTime: lastTransitionTime,
				Type:               condType,
				Status:             status,
				Reason:             reason,
				Message:            msg,
			}
			setConnectionMetrics(conn, status, lastTransitionTime, reason)
			return
		}
	}

	// If the condition does not exist,
	// initialize the lastTransitionTime
	conn.Status.Conditions = append(conn.Status.Conditions, metav1.Condition{
		LastTransitionTime: now,
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            msg,
	})

	setConnectionMetrics(conn, status, now, reason)
}

// GetConnectionCondition return the condition with the passed condition type from
// the status object. If the condition is not already present, return nil
func GetConnectionCondition(conn *MongoDBAtlasConnection, condType string) *metav1.Condition {
	for i := range conn.Status.Conditions {
		if conn.Status.Conditions[i].Type == condType {
			return &conn.Status.Conditions[i]
		}
	}
	return nil
}

func setConnectionMetrics(conn *MongoDBAtlasConnection, status metav1.ConditionStatus, lastTransitionTime metav1.Time, reason string) {
	if status == metav1.ConditionTrue {
		setConnectionElapsedTime(conn, lastTransitionTime)
		setConnectionStatusReady(conn, 1)
	} else {
		resetConnectionElapsedTime(conn)
		setConnectionStatusReady(conn, 0)
	}
	setConnectionStatusReason(conn, reason)
}

func setConnectionElapsedTime(conn *MongoDBAtlasConnection, lastTransitionTime metav1.Time) {
	metrics.ConnectionElapsedTime.With(prometheus.Labels{
		"provider":   metrics.DBaaSProvider,
		"connection": conn.Name,
		"namespace":  conn.Namespace}).Set(lastTransitionTime.Sub(conn.CreationTimestamp.Time).Seconds())
}

func resetConnectionElapsedTime(conn *MongoDBAtlasConnection) {
	metrics.ConnectionElapsedTime.Delete(prometheus.Labels{
		"provider":   metrics.DBaaSProvider,
		"connection": conn.Name,
		"namespace":  conn.Namespace})
}

func setConnectionStatusReady(conn *MongoDBAtlasConnection, val float64) {
	metrics.ConnectionStatusReady.With(prometheus.Labels{
		"provider":   metrics.DBaaSProvider,
		"connection": conn.Name,
		"namespace":  conn.Namespace}).Set(val)
}

func resetConnectionStatusReasons(conn *MongoDBAtlasConnection) {
	for _, reason := range workflow.GetMongoDBAtlasConnectionReasons() {
		metrics.ConnectionStatusReason.With(prometheus.Labels{
			"provider":   metrics.DBaaSProvider,
			"connection": conn.Name,
			"namespace":  conn.Namespace,
			"reason":     string(reason)}).Set(0)
	}
}

func setConnectionStatusReason(conn *MongoDBAtlasConnection, reason string) {
	resetConnectionStatusReasons(conn)
	metrics.ConnectionStatusReason.With(prometheus.Labels{
		"provider":   metrics.DBaaSProvider,
		"connection": conn.Name,
		"namespace":  conn.Namespace,
		"reason":     reason}).Set(1)
}
