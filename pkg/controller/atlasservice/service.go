package atlasservice

import (
	"context"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// ensureProjectExists creates the service if it doesn't exist yet. Returns the service ID
func (r *AtlasServiceReconciler) discoverServices(ctx *workflow.Context) ([]status.AtlasProjectService, workflow.Result) {
	// Try to find the service
	projects, _, err := ctx.Client.Projects.GetAllProjects(context.Background(), &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, workflow.Terminate(workflow.AtlasProjectsListFailed, err.Error())
	}
	svcList := []status.AtlasProjectService{}
	for _, p := range projects.Results {

		clusters, _, err := ctx.Client.Clusters.List(context.Background(), p.ID, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, workflow.Terminate(workflow.AtlasClustersListFailed, err.Error())
		}
		clusterList := []status.AtlasClusterService{}
		for _, cluster := range clusters {
			connStrTokens := strings.Split(cluster.ConnectionStrings.StandardSrv, "://")
			var protocol, host string
			if len(connStrTokens) < 2 {
				//There is no "://" found
				protocol = ""
				host = connStrTokens[0]
			} else {
				protocol = connStrTokens[0]
				host = connStrTokens[1]
			}
			clusterSvc := status.AtlasClusterService{
				ID:                 cluster.ID,
				Name:               cluster.Name,
				InstanceSizeName:   cluster.ProviderSettings.InstanceSizeName,
				ProviderName:       cluster.ProviderSettings.ProviderName,
				RegionName:         cluster.ProviderSettings.RegionName,
				ConnectionProtocol: protocol,
				ConnectionHost:     host,
			}
			clusterList = append(clusterList, clusterSvc)
		}
		dbUsers, _, err := ctx.Client.DatabaseUsers.List(context.Background(), p.ID, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, workflow.Terminate(workflow.AtlasDBUsersListFailed, err.Error())
		}
		dbUserList := []string{}
		for _, dbUser := range dbUsers {
			dbUserList = append(dbUserList, dbUser.Username)
		}
		svc := status.AtlasProjectService{
			ID:          p.ID,
			Name:        p.Name,
			ClusterList: clusterList,
			DBUserList:  dbUserList,
		}
		svcList = append(svcList, svc)
	}
	ctx.Log.Infof("Project discovered:%v", svcList)
	return svcList, workflow.OK()
}
