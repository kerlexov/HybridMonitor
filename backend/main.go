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
	router := router.New()
	db, err := db.Init(db.DbConfig{
		Host:     "localhost",
		Username: "test",
		Password: "YqelfkqoivjffuFđ",
		Table:    "backend",
		Port:     5432,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{}, &models.RedfishHost{}, &models.VsphereHost{})

	redfish_exporter := redfish_exporter.NewRedfishExporter(9141, "/redfish", router.Logger)
	if redfish_exporter == nil {
		router.Logger.Error("Cannot start redfish exporter")
	} else {
		go redfish_exporter.Start()
	}
	vsphere_exporter := vsphere_exporter.NewVsphereExporter(9142, "/vsphere", router.Logger)
	if vsphere_exporter == nil {
		router.Logger.Error("Cannot start vsphere exporter")
	} else {
		go vsphere_exporter.Start()
	}

	handler := handler.NewHandler(db, redfish_exporter, vsphere_exporter)
	defer handler.Close()
	api := router.Group("/api")
	handler.Register(api)
	db.Create(&models.User{Username: "user", Password: "test"})
	router.Logger.Fatal(router.StartServer(&http.Server{
		Addr:         ":9393",
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}))
	defer redfish_exporter.Logout()
	defer vsphere_exporter.Logout()
}
