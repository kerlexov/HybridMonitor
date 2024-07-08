package vsphere_exporter

import (
	"context"
	"fmt"
	vsphere_exporter "github.com/kerlexov/HybridMonitor/backend/vsphere_exporter/lib"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

type VsphereExporter struct {
	Port           int
	MetricsPath    string
	Server         *http.Server
	VsphereManager *vsphere_exporter.VsphereManager
	Registry       *prometheus.Registry
	Logger         echo.Logger
	Ctx            context.Context
}

func NewVsphereExporter(Port int, MetricsPath string, logger echo.Logger) *VsphereExporter {
	registry := prometheus.NewRegistry()
	ctx := context.Background()

	vSphereManager := &vsphere_exporter.VsphereManager{
		Connections: make(map[string]*vsphere_exporter.Connection),
		Configs:     make(map[string]*vsphere_exporter.VsphereConfig),
		Ctx:         ctx,
		Logger:      logger,
	}

	registry.MustRegister(vSphereManager)

	http.Handle(MetricsPath, vSphereManager.MetricsHandler(registry))

	server := &http.Server{Addr: fmt.Sprintf(":%d", Port), Handler: nil}

	vSphereExporter := &VsphereExporter{
		Port:           Port,
		MetricsPath:    MetricsPath,
		VsphereManager: vSphereManager,
		Server:         server,
		Registry:       registry,
		Logger:         logger,
		Ctx:            ctx,
	}

	return vSphereExporter
}

func (re VsphereExporter) Start() error {
	return re.Server.ListenAndServe()
}

func (re VsphereExporter) Logout() error {
	var err error = nil
	for _, conn := range re.VsphereManager.Connections {
		err = conn.VsphereClient.Logout()
	}
	return err
}
