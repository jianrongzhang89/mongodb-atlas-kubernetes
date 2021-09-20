package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	customMetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	DBaaSProvider = "mongodbatlas"
)

// Custom metrics
var (
	OperatorVersion = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_operator_version_info",
			Help: "DBaaS provider operator version information",
		},
		[]string{
			"provider",
			"operator_version",
		},
	)
	DBaaSRegistrationReady = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_registration_ready",
			Help: "DBaaS provider registration ready",
		},
		[]string{
			"provider",
		},
	)
	InventoryElapsedTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_inventory_creation_ready_seconds",
			Help: "Elapsed time from DBaaS provider inventory creation to sync ready",
		},
		[]string{
			"provider",
			"inventory",
			"namespace",
		},
	)
	InventoryStatusReady = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_inventory_status_ready",
			Help: "DBaaS provider inventory is sync ready",
		},
		[]string{
			"provider",
			"inventory",
			"namespace",
		},
	)
	InventoryStatusReason = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_inventory_status_reason",
			Help: "Inventory status reason",
		},
		[]string{
			"provider",
			"inventory",
			"namespace",
			"reason",
		},
	)
	ConnectionElapsedTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_connection_creation_ready_seconds",
			Help: "Elapsed time from DBaaS provider connection creation to ready",
		},
		[]string{
			"provider",
			"connection",
			"namespace",
		},
	)
	ConnectionStatusReady = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_connection_status_ready",
			Help: "DBaaS provider connection is ready for binding",
		},
		[]string{
			"provider",
			"connection",
			"namespace",
		},
	)
	ConnectionStatusReason = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dbaas_provider_connection_status_reason",
			Help: "Connection status reason",
		},
		[]string{
			"provider",
			"connection",
			"namespace",
			"reason",
		},
	)
)

var metricsList = []*prometheus.GaugeVec{
	OperatorVersion,
	DBaaSRegistrationReady,
	InventoryElapsedTime,
	InventoryStatusReady,
	InventoryStatusReason,
	ConnectionElapsedTime,
	ConnectionStatusReady,
	ConnectionStatusReason,
}

func init() {
	// Register custom metrics with the global prometheus registry
	for _, metric := range metricsList {
		customMetrics.Registry.MustRegister(metric)
	}
}
