// Copyright (C) 2025, Dione Protocol, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package services

import (
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/constants"
	"github.com/DioneProtocol/odyssey-tooling-sdk-go/utils"
)

func RenderGrafanaLokiDataSourceConfig() ([]byte, error) {
	return templates.ReadFile("templates/grafana-loki-datasource.yaml")
}

func RenderGrafanaPrometheusDataSourceConfigg() ([]byte, error) {
	return templates.ReadFile("templates/grafana-prometheus-datasource.yaml")
}

func RenderGrafanaConfig() ([]byte, error) {
	return templates.ReadFile("templates/grafana.ini")
}

func RenderGrafanaDashboardConfig() ([]byte, error) {
	return templates.ReadFile("templates/grafana-dashboards.yaml")
}

func GrafanaFoldersToCreate() []string {
	return []string{
		utils.GetRemoteComposeServicePath(constants.ServiceGrafana, "data"),
		utils.GetRemoteComposeServicePath(constants.ServiceGrafana, "dashboards"),
		utils.GetRemoteComposeServicePath(constants.ServiceGrafana, "provisioning", "datasources"),
		utils.GetRemoteComposeServicePath(constants.ServiceGrafana, "provisioning", "dashboards"),
	}
}
