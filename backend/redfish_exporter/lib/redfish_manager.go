package redfish_exporter

import (
	"fmt"
	"github.com/kerlexov/HybridMonitor/backend/redfish_exporter/lib/thermal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

const namespace = "redfish_exporter"

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last redfish_exporter query successful.",
		[]string{"host"}, nil,
	)
	powerConsumedWatts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "PowerConsumedWatts"),
		"Current power consumption in watts",
		[]string{"host"}, nil,
	)
	fanTemperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "FanTemperature"),
		"Current fan temperature",
		[]string{"host", "fan"}, nil,
	)
	fanHealth = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "FanStatus"),
		"Current fan health status",
		[]string{"host", "fan"}, nil,
	)
	sensorTemperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "SensorTemperature"),
		"Current sensor temperature",
		[]string{"host", "sensor", "number", "physicalContext", "units", "upperThresholdCritical", "upperThresholdFatal"}, nil,
	)
	sensorHealth = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "SensorHealth"),
		"Current sensor health status",
		[]string{"host", "sensor", "number", "physicalContext"}, nil,
	)
)

type RedfishConfig struct {
	Host     string
	Username string
	Password string
}
type RedfishManager struct {
	Connections map[string]*Connection
	Configs     map[string]*RedfishConfig
}

func (rm RedfishManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- powerConsumedWatts
	ch <- fanTemperature
	ch <- fanTemperature
	ch <- fanHealth
	ch <- sensorTemperature
	ch <- sensorHealth
}

func (rm *RedfishManager) Collect(ch chan<- prometheus.Metric) {
	data, err := rm.LoadData()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0, "redfish-exporter")
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1, "redfish-exporter")

	rm.UpdateMetrics(data, ch)
}

type Connection struct {
	RedfishClient *RedfishClient
}

type Data struct {
	Host             string
	PowerConsumption float64
	Thermal          *thermal.ChassisThermal
	Success          bool
}

type AddRedfishHost struct {
	Host     string `json:"host" form:"host" query:"host"`
	User     string `json:"user" form:"user" query:"user"`
	Password string `json:"password" form:"password" query:"password"`
}

func (h RedfishManager) MetricsHandler(registry *prometheus.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}

		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func (re *RedfishManager) LoadData() ([]*Data, error) {
	var data []*Data

	for key, value := range re.Connections {
		values := &Data{}
		_, err := value.RedfishClient.GetSession(re.Configs[key])
		if err != nil {
			values.Success = false
		}
		if value.RedfishClient.IsValid() {
			values.PowerConsumption, err = value.RedfishClient.GetPower()
			if err != nil {
				values.Success = false
			}
			values.Thermal, err = value.RedfishClient.GetThermal()
			if err != nil {
				values.Success = false
			}
		}
		values.Host = value.RedfishClient.Host

		data = append(data, values)
	}

	return data, nil
}

func (re *RedfishManager) UpdateMetrics(values []*Data, ch chan<- prometheus.Metric) {
	for _, data := range values {
		ch <- prometheus.MustNewConstMetric(powerConsumedWatts, prometheus.GaugeValue, data.PowerConsumption, data.Host)
		for _, fan := range data.Thermal.Fans {
			ch <- prometheus.MustNewConstMetric(fanTemperature, prometheus.GaugeValue,
				float64(fan.CurrentReading), data.Host, fan.FanName)
			ch <- prometheus.MustNewConstMetric(fanHealth, prometheus.GaugeValue,
				HealthCheck(fan.Status.State, fan.Status.Health), data.Host, fan.FanName)
		}
		for _, sensor := range data.Thermal.Temperatures {
			number := strconv.Itoa(int(sensor.Number))
			upperThresholdCritical := strconv.Itoa(int(sensor.UpperThresholdCritical))
			upperThresholdFatal := strconv.Itoa(int(sensor.UpperThresholdFatal))
			ch <- prometheus.MustNewConstMetric(sensorTemperature, prometheus.GaugeValue,
				float64(sensor.CurrentReading), data.Host, sensor.Name, number,
				string(sensor.PhysicalContext), string(sensor.Units), upperThresholdCritical, upperThresholdFatal)
			ch <- prometheus.MustNewConstMetric(sensorHealth, prometheus.GaugeValue,
				HealthCheck(sensor.Status.State, sensor.Status.Health), data.Host,
				sensor.Name, number, string(sensor.PhysicalContext))
		}
		if data != nil {
			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, 1, data.Host)
		}
		log.Println(fmt.Sprintf("Endpoint: %s scraped", data.Host))
	}
}

func (re *RedfishManager) AddRedfishHost(host string, user string, password string) bool {
	config := &RedfishConfig{
		Username: user,
		Password: password,
		Host:     host,
	}

	client, err := NewRedfishClient(config)
	if err != nil {
		return false
	}
	if client.IsValid() {
		re.Connections[host] = &Connection{RedfishClient: client}
		re.Configs[host] = config
		return true
	}
	return false
}
func (rm RedfishManager) GetRedfishHosts() (map[string]*Connection, error) {
	return rm.Connections, nil
}
