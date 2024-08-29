package vsphere_exporter

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vmware/govmomi/vim25/types"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const namespace = "vsphere_exporter"

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last vsphere_exporter query successful.",
		[]string{"service"}, nil,
	)
	upVsphere = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upVsphere"),
		"Was the last vsphere_manager query successful.",
		[]string{"host"}, nil,
	)
	overallStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "OverallStatus"),
		"OverallStatus",
		[]string{"manager", "host"}, nil,
	)
	usageCPU = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "UsageCPU"),
		"Current fan temperature",
		[]string{"manager", "host", "total"}, nil,
	)
	usageRAM = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "UsageRAM"),
		"Current fan temperature",
		[]string{"manager", "host", "total"}, nil,
	)
	sensorInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "SensorInfo"),
		"Information about different types of sensor",
		[]string{"manager", "host", "name", "unitmodifier", "baseunits", "rateunites", "type", "summary", "state"}, nil,
	)
	hardwareInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "HardwareInfo"),
		"Information about hardware health",
		[]string{"manager", "host", "name", "summary"}, nil,
	)
	runTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "RunTime"),
		"Run time",
		[]string{"manager", "host"}, nil,
	)
)

type VsphereConfig struct {
	Host     string
	Username string
	Password string
}
type VsphereManager struct {
	Connections map[string]*Connection
	Configs     map[string]*VsphereConfig
	Ctx         context.Context
	Logger      echo.Logger
}

type AddVsphereHost struct {
	Host     string `json:"host" form:"host" query:"host"`
	User     string `json:"user" form:"user" query:"user"`
	Password string `json:"password" form:"password" query:"password"`
}

func (rm VsphereManager) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- upVsphere
	ch <- usageCPU
	ch <- usageRAM
	ch <- runTime
	ch <- overallStatus
	ch <- sensorInfo
	ch <- hardwareInfo
}

func (rm *VsphereManager) Collect(ch chan<- prometheus.Metric) {
	data, err := rm.LoadData()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0, "vsphere-exporter")
		log.Println(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1, "vsphere-exporter")

	rm.UpdateMetrics(data, ch)
}

type Connection struct {
	VsphereClient *VsphereClient
}

type VsphereData struct {
	HostManager     string
	OverallStatus   string
	VSphereServerIp string
	HostLabel       string
	Connected       types.HostSystemConnectionState
	Powerstate      types.HostSystemPowerState
	HealthSystem    types.HealthSystemRuntime
	QuickStats      types.HostListSummaryQuickStats
	UsageCPU        int64
	UsageRAM        int64
	TotalCPU        int64
	TotalRAM        int64
	BootTime        time.Time
	HostIp          string
	Product         string
	Success         bool
}

type AddRedfishHost struct {
	Host     string `json:"host" form:"host" query:"host"`
	User     string `json:"user" form:"user" query:"user"`
	Password string `json:"password" form:"password" query:"password"`
}

func (h VsphereManager) MetricsHandler(registry *prometheus.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gatherers := prometheus.Gatherers{
			prometheus.DefaultGatherer,
			registry,
		}

		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func (re *VsphereManager) LoadData() ([]*VsphereData, error) {
	var data []*VsphereData

	for key, conn := range re.Connections {
		values, err := conn.VsphereClient.GetData(key)
		if err != nil {
			re.Logger.Error(err)
		}
		data = append(data, values...)
	}

	return data, nil
}

func (re *VsphereManager) UpdateMetrics(values []*VsphereData, ch chan<- prometheus.Metric) {
	for _, data := range values {
		if data != nil {
			ch <- prometheus.MustNewConstMetric(usageCPU, prometheus.GaugeValue, float64(data.UsageCPU), data.VSphereServerIp, data.HostIp, strconv.Itoa(int(data.TotalCPU)))
			ch <- prometheus.MustNewConstMetric(usageRAM, prometheus.GaugeValue, float64(data.UsageRAM), data.VSphereServerIp, data.HostIp, strconv.Itoa(int(data.TotalRAM)))
			ch <- prometheus.MustNewConstMetric(overallStatus, prometheus.GaugeValue, parseStatus(data.OverallStatus), data.VSphereServerIp, data.HostIp)
			ch <- prometheus.MustNewConstMetric(runTime, prometheus.GaugeValue, parseRunTime(data.BootTime), data.VSphereServerIp, data.HostIp)

			if data.HealthSystem.SystemHealthInfo != nil && data.HealthSystem.SystemHealthInfo.NumericSensorInfo != nil {
				for _, info := range data.HealthSystem.SystemHealthInfo.NumericSensorInfo {
					ch <- prometheus.MustNewConstMetric(sensorInfo, prometheus.GaugeValue, float64(info.CurrentReading), data.VSphereServerIp,
						data.HostIp, info.Name, strconv.Itoa(int(info.UnitModifier)), info.BaseUnits, info.RateUnits, info.SensorType,
						info.HealthState.GetElementDescription().Summary, info.HealthState.GetElementDescription().Key)
				}
			}

			if data.HealthSystem.HardwareStatusInfo != nil && data.HealthSystem.HardwareStatusInfo.CpuStatusInfo != nil {
				for _, info := range data.HealthSystem.HardwareStatusInfo.CpuStatusInfo {
					ch <- prometheus.MustNewConstMetric(hardwareInfo, prometheus.GaugeValue,
						parseStatus(info.GetHostHardwareElementInfo().Status.GetElementDescription().Key),
						data.VSphereServerIp, data.HostIp, info.GetHostHardwareElementInfo().Name,
						info.GetHostHardwareElementInfo().Status.GetElementDescription().Summary)
				}
			}

			if data.HealthSystem.HardwareStatusInfo != nil && data.HealthSystem.HardwareStatusInfo.MemoryStatusInfo != nil {
				for _, info := range data.HealthSystem.HardwareStatusInfo.MemoryStatusInfo {
					ch <- prometheus.MustNewConstMetric(hardwareInfo, prometheus.GaugeValue,
						parseStatus(info.GetHostHardwareElementInfo().Status.GetElementDescription().Key),
						data.VSphereServerIp, data.HostIp, info.GetHostHardwareElementInfo().Name,
						info.GetHostHardwareElementInfo().Status.GetElementDescription().Summary)
				}
			}

			log.Println(fmt.Sprintf("Endpoint: %s (%s)scraped", data.HostIp, data.HostManager))
		}
	}
}

func parseRunTime(bootTime time.Time) float64 {
	return time.Now().Sub(bootTime).Minutes()
}

func parseStatus(status string) float64 {
	switch strings.ToLower(status) {
	case "green":
		return 2
	case "yellow":
		return 1
	case "gray":
		return 0
	case "red":
		return -2
	default:
		return 0
	}
}

func (re *VsphereManager) AddVsphereHost(host string, user string, password string) bool {
	var config *VsphereConfig
	if re.Configs[host] == nil {
		config = &VsphereConfig{
			Username: user,
			Password: password,
			Host:     host,
		}
		re.Configs[host] = config
	}
	client, err := NewVsphereClient(re.Ctx, re.Configs[host])
	if err != nil {
		return false
	}
	re.Connections[host] = &Connection{VsphereClient: client}
	return true
}

func (rm VsphereManager) GetVsphereHosts() (map[string]*Connection, error) {
	return rm.Connections, nil
}
