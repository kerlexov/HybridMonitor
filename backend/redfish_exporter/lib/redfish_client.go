package redfish_exporter

import (
	"crypto/tls"
	"fmt"
	"github.com/kerlexov/HybridMonitor/backend/redfish_exporter/lib/power"
	"github.com/kerlexov/HybridMonitor/backend/redfish_exporter/lib/thermal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type RedfishClient struct {
	Session  *Session
	Insecure bool
	Host     string
}

type Session struct {
	ID    string
	Token string
}

type Power struct {
	Value string
}

func NewRedfishClient(config *RedfishConfig) (*RedfishClient, error) {
	client := &RedfishClient{
		Session:  &Session{},
		Insecure: true,
		Host:     config.Host,
	}
	_, err := client.GetSession(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *RedfishClient) IsValid() bool {
	if len(c.Session.Token) > 0 {
		return true
	}
	return false
}

func (c *RedfishClient) GetSession(config *RedfishConfig) (*Session, error) {
	if len(c.Session.Token) == 0 {
		url := fmt.Sprintf("http://%s/redfish/v1/SessionService/Sessions/", c.Host)

		payload := strings.NewReader(`{
 			"UserName": "` + config.Username + `",
 			"Password": "` + config.Password + `"
		}`)

		defaultTransport := http.DefaultTransport.(*http.Transport)
		transport := &http.Transport{
			Proxy:                 defaultTransport.Proxy,
			DialContext:           defaultTransport.DialContext,
			MaxIdleConns:          defaultTransport.MaxIdleConns,
			IdleConnTimeout:       defaultTransport.IdleConnTimeout,
			ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
			TLSHandshakeTimeout:   time.Duration(60) * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
		req, err := http.NewRequest("POST", url, payload)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if os.IsTimeout(err) {
			return nil, fmt.Errorf("Request timeout %s", err)
		}
		defer res.Body.Close()

		if len(res.Header.Get("X-Auth-Token")) != 0 {
			c.Session = &Session{
				ID:    res.Header.Get("Location"),
				Token: res.Header.Get("X-Auth-Token"),
			}
		} else {
			return nil, fmt.Errorf("invalid login credentials")
		}
	}
	return c.Session, nil
}

func HealthCheck(state thermal.State, health *thermal.Health) float64 {
	if state == "Enabled" && thermal.Health(*health) == "OK" {
		return 1
	} else if state == "Absent" {
		return 0
	}
	return -1
}

func (c RedfishClient) GetPower() (float64, error) {
	body, err := c.CreateRequest("GET", "/redfish/v1/Chassis/1/Power/")
	if err != nil {
		return 0, err
	}
	chassisPower, err := power.UnmarshalChassisPower(body)
	if err != nil {
		return 0, err
	}

	return float64(chassisPower.PowerConsumedWatts), nil
}

func (c RedfishClient) GetThermal() (*thermal.ChassisThermal, error) {
	body, err := c.CreateRequest("GET", "/redfish/v1/Chassis/1/Thermal/")
	if err != nil {
		return nil, err
	}

	chassisThermal, err := thermal.UnmarshalChassisThermal(body)
	if err != nil {
		return nil, err
	}
	return &chassisThermal, nil
}

func (c RedfishClient) CreateRequest(method string, path string) ([]byte, error) {
	url := fmt.Sprintf("http://%s%s", c.Host, path)
	if strings.Contains(path, "https://") {
		url = path
	}

	defaultTransport := http.DefaultTransport.(*http.Transport)
	transport := &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   time.Duration(10) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.Insecure,
		},
	}

	client := &http.Client{Transport: transport}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", c.Session.Token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *RedfishClient) Logout() error {
	var url string
	if len(c.Session.ID) != 0 {
		url = c.Session.ID
		if !strings.Contains(url, "http") {
			url = fmt.Sprintf("http://%s%s", c.Host, url)
		}
	}

	body, err := c.CreateRequest("DELETE", url)
	if err != nil {
		return err
	}
	fmt.Println(body)
	return nil
}

//func (re *RedfishExporter) Logout() error {
//	fmt.Println(re.redfishClient.Session.ID)
//	body, err := re.redfishClient.CreateRequest("DELETE", re.redfishClient.Session.ID)
//	if err != nil {
//		return err
//	}
//	fmt.Println(body)
//	return nil
//}
