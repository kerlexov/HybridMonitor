package redfish_exporter

import (
	"fmt"
	redfish_exporter "github.com/kerlexov/HybridMonitor/backend/redfish_exporter/lib"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type RedfishExporter struct {
	Port           int
	MetricsPath    string
	Server         *http.Server
	RedfishManager *redfish_exporter.RedfishManager
	Registry       *prometheus.Registry
	Logger         echo.Logger
}

func NewRedfishExporter(Port int, MetricsPath string, logger echo.Logger) *RedfishExporter {
	registry := prometheus.NewRegistry()

	redfishManager := &redfish_exporter.RedfishManager{
		Connections: make(map[string]*redfish_exporter.Connection),
		Configs:     make(map[string]*redfish_exporter.RedfishConfig),
	}

	registry.MustRegister(redfishManager)

	http.Handle(MetricsPath, redfishManager.MetricsHandler(registry))

	server := &http.Server{Addr: fmt.Sprintf(":%d", Port), Handler: nil}

	redfishExporter := &RedfishExporter{
		Port:           Port,
		MetricsPath:    MetricsPath,
		RedfishManager: redfishManager,
		Server:         server,
		Registry:       registry,
		Logger:         logger,
	}
	return redfishExporter
}

func (re RedfishExporter) Start() error {
	return re.Server.ListenAndServe()
}

func (re RedfishExporter) Logout() error {
	var err error = nil
	for _, conn := range re.RedfishManager.Connections {
		err = conn.RedfishClient.Logout()
	}
	return err
}
