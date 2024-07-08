package vsphere_exporter

import (
	"context"
	"flag"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"net/url"
)

type VsphereClient struct {
	Client *govmomi.Client
	Ctx    context.Context
}

var interval = flag.Int("i", 300, "Interval ID")

func NewVsphereClient(ctx context.Context, config *VsphereConfig) (*VsphereClient, error) {
	usr := url.UserPassword(config.Username, config.Password)
	url, err := soap.ParseURL(config.Host)
	if err != nil {
		return nil, err
	}
	url.User = usr
	c, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		return nil, err
	}
	client := &VsphereClient{Ctx: ctx, Client: c}

	return client, nil
}

func (c *VsphereClient) IsValid() bool {
	ok, err := c.Client.SessionManager.SessionIsActive(c.Ctx)
	if err != nil {
		return false
	}
	if ok {
		return true
	}
	return false
}

//func HealthCheck(state thermal.State, health *thermal.Health) float64 {
//	if state == "Enabled" && thermal.Health(*health) == "OK" {
//		return 1
//	} else if state == "Absent" {
//		return 0
//	}
//	return -1
//}

//func (c VsphereClient) GetPower() (float64, error) {
//	body, err := c.CreateRequest("GET", "/redfish/v1/Chassis/1/Power/")
//	if err != nil {
//		return 0, err
//	}
//	chassisPower, err := power.UnmarshalChassisPower(body)
//	if err != nil {
//		return 0, err
//	}
//
//	return float64(chassisPower.PowerConsumedWatts), nil
//}
//
//func (c VsphereClient) GetThermal() (*thermal.ChassisThermal, error) {
//	body, err := c.CreateRequest("GET", "/redfish/v1/Chassis/1/Thermal/")
//	if err != nil {
//		return nil, err
//	}
//
//	chassisThermal, err := thermal.UnmarshalChassisThermal(body)
//	if err != nil {
//		return nil, err
//	}
//	return &chassisThermal, nil
//}

//func (c VsphereClient) CreateRequest(method string, path string) ([]byte, error) {
//	url := fmt.Sprintf("http://%s%s", c.Host, path)
//	if strings.Contains(path, "https://") {
//		url = path
//	}
//
//	defaultTransport := http.DefaultTransport.(*http.Transport)
//	transport := &http.Transport{
//		Proxy:                 defaultTransport.Proxy,
//		DialContext:           defaultTransport.DialContext,
//		MaxIdleConns:          defaultTransport.MaxIdleConns,
//		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
//		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
//		TLSHandshakeTimeout:   time.Duration(10) * time.Second,
//		TLSClientConfig: &tls.Config{
//			InsecureSkipVerify: c.Insecure,
//		},
//	}
//
//	client := &http.Client{Transport: transport}
//	req, err := http.NewRequest(method, url, nil)
//
//	if err != nil {
//		return nil, err
//	}
//	req.Header.Add("X-Auth-Token", c.Session.Token)
//
//	res, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer res.Body.Close()
//
//	body, err := ioutil.ReadAll(res.Body)
//	if err != nil {
//		return nil, err
//	}
//	return body, nil
//}

func (vc *VsphereClient) GetData(host string) ([]*VsphereData, error) {
	var data []*VsphereData
	client := vc.GetClient()
	m := view.NewManager(client)
	context := context.Background()
	v, err := m.CreateContainerView(context, client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(context)

	var hss []mo.HostSystem
	err = v.Retrieve(context, []string{"HostSystem"}, []string{"summary"}, &hss)
	if err != nil {
		return nil, err
	}

	for _, hs := range hss {
		data = append(data, &VsphereData{
			HostManager:     host,
			OverallStatus:   string(hs.Summary.OverallStatus),
			VSphereServerIp: hs.Summary.ManagementServerIp,
			HostLabel:       hs.Summary.Host.Value,
			Connected:       hs.Summary.Runtime.ConnectionState,
			Powerstate:      hs.Summary.Runtime.PowerState,
			HealthSystem:    *hs.Summary.Runtime.HealthSystemRuntime,
			BootTime:        *hs.Summary.Runtime.BootTime,
			HostIp:          hs.Summary.Config.Name,
			Product:         hs.Summary.Config.Product.Name,
			QuickStats:      hs.Summary.QuickStats,
			TotalCPU:        int64(hs.Summary.Hardware.CpuMhz) * int64(hs.Summary.Hardware.NumCpuCores),
			TotalRAM:        hs.Summary.Hardware.MemorySize / (1024 * 1024),
			UsageCPU:        int64(hs.Summary.QuickStats.OverallCpuUsage),
			UsageRAM:        int64(hs.Summary.QuickStats.OverallMemoryUsage),
		})
	}

	return data, err
}

func (c *VsphereClient) Login(config *VsphereConfig) error {
	return c.Client.Login(c.Ctx, url.UserPassword(config.Username, config.Password))
}

func (c *VsphereClient) GetClient() *vim25.Client {
	return c.Client.Client
}

func (c *VsphereClient) Logout() error {
	return c.Client.Logout(c.Ctx)
}
