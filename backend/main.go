package main

import (
	"github.com/kerlexov/HybridMonitor/backend/redfish_exporter"
	"github.com/kerlexov/HybridMonitor/backend/vsphere_exporter"
	db "github.com/kerlexov/HybridMonitor/db"
	"github.com/kerlexov/HybridMonitor/handler"
	"github.com/kerlexov/HybridMonitor/models"
	"github.com/kerlexov/HybridMonitor/router"
	"log"
	"net/http"
	"time"
)

func main() {
	rout := router.New()
	dbConn, err := db.Init(db.DbConfig{
		Host:     "localhost",
		Username: "test",
		Password: "YqelfkqoivjffuFđ",
		Table:    "backend",
		Port:     5432,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.AutoMigrate(&models.User{}, &models.RedfishHost{}, &models.VsphereHost{})
	if err != nil {
		rout.Logger.Error("Cannot start redfish exporter")
	}

	redfishExporter := redfish_exporter.NewRedfishExporter(9141, "/redfish", rout.Logger)
	if redfishExporter == nil {
		rout.Logger.Error("Cannot start no redfish exporter")
	} else {
		go func() {
			err := redfishExporter.Start()
			if err != nil {
				rout.Logger.Error("Cannot start redfish exporter")
			}
		}()
	}
	vsphereExporter := vsphere_exporter.NewVsphereExporter(9142, "/vsphere", rout.Logger)
	if vsphereExporter == nil {
		rout.Logger.Error("Cannot start vsphere exporter")
	} else {
		go func() {
			err := vsphereExporter.Start()
			if err != nil {
				rout.Logger.Error("Cannot start vsphere exporter")
			}
		}()
	}

	handl := handler.NewHandler(dbConn, redfishExporter, vsphereExporter)
	defer handl.Close()
	api := rout.Group("/api")
	handl.Register(api)
	dbConn.Create(&models.User{Username: "user", Password: "test"})
	rout.Logger.Fatal(rout.StartServer(&http.Server{
		Addr:         ":9393",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}))
	defer func(redfish_exporter *redfish_exporter.RedfishExporter) {
		err := redfish_exporter.Logout()
		if err != nil {
			rout.Logger.Error("Cannot logout redfish exporter")

		}
	}(redfishExporter)
	defer func(vsphereExporter *vsphere_exporter.VsphereExporter) {
		err := vsphereExporter.Logout()
		if err != nil {
			rout.Logger.Error("Cannot logout vsphere exporter")

		}
	}(vsphereExporter)
}
