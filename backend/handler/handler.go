package handler

import (
	"github.com/kerlexov/HybridMonitor/backend/redfish_exporter"
	"github.com/kerlexov/HybridMonitor/backend/vsphere_exporter"
	"gorm.io/gorm"
)

type Handler struct {
	DB              *gorm.DB
	RedfishExporter *redfish_exporter.RedfishExporter
	VsphereExporter *vsphere_exporter.VsphereExporter
}

func NewHandler(db *gorm.DB, exporter *redfish_exporter.RedfishExporter, vsphere_exporter *vsphere_exporter.VsphereExporter) *Handler {
	return &Handler{
		DB:              db,
		RedfishExporter: exporter,
		VsphereExporter: vsphere_exporter,
	}
}

func (h *Handler) Close() {

}
