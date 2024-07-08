package handler

import (
	"github.com/golang-jwt/jwt"
	redfish_exporter "github.com/kerlexov/HybridMonitor/backend/redfish_exporter/lib"
	vsphere_exporter "github.com/kerlexov/HybridMonitor/backend/vsphere_exporter/lib"
	"github.com/kerlexov/HybridMonitor/lib"
	"github.com/kerlexov/HybridMonitor/models"
	usr "github.com/kerlexov/HybridMonitor/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

const SECRET = "vwoijqwfjiei0j"

func (h *Handler) Register(api *echo.Group) {

	api.POST("/login", func(c echo.Context) (err error) {
		u := new(usr.LoginUser)
		if err = c.Bind(u); err != nil {
			return
		}

		if usr.IsUserValid(h.DB, u) {
			claims := &lib.JwtCustomClaims{
				u.Username,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
				},
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			t, err := token.SignedString([]byte(SECRET))
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, echo.Map{
				"token": t,
			})
		}
		return c.JSON(http.StatusForbidden, "nemore")
	})

	v1 := api.Group("/v1")

	config := middleware.JWTConfig{
		Claims:     &lib.JwtCustomClaims{},
		SigningKey: []byte(SECRET),
	}
	v1.Use(middleware.JWTWithConfig(config))

	redfish := v1.Group("/redfish")
	vSphere := v1.Group("/vsphere")

	redfish.POST("/add", func(c echo.Context) (err error) {
		redfishHost := new(redfish_exporter.AddRedfishHost)
		if err = c.Bind(redfishHost); err != nil {
			return
		}

		ok := h.RedfishExporter.RedfishManager.AddRedfishHost(redfishHost.Host, redfishHost.User, redfishHost.Password)
		if ok {
			h.DB.Create(&models.RedfishHost{Host: redfishHost.Host, Username: redfishHost.User, Password: redfishHost.Password})
			return c.JSON(http.StatusOK, "ok")
		}

		return c.JSON(http.StatusBadRequest, "failed")
	})

	redfish.GET("/agents", func(c echo.Context) (err error) {
		hosts, err := h.RedfishExporter.RedfishManager.GetRedfishHosts()
		if err == nil {
			return c.JSON(http.StatusBadRequest, "failed")
		}
		return c.JSON(http.StatusOK, hosts)
	})

	vSphere.POST("/add", func(c echo.Context) (err error) {
		vsphereHost := new(vsphere_exporter.AddVsphereHost)
		if err = c.Bind(vsphereHost); err != nil {
			return
		}

		ok := h.VsphereExporter.VsphereManager.AddVsphereHost(vsphereHost.Host, vsphereHost.User, vsphereHost.Password)
		if ok {
			h.DB.Create(&models.VsphereHost{Host: vsphereHost.Host, Username: vsphereHost.User, Password: vsphereHost.Password})
			return c.JSON(http.StatusOK, "ok")
		}

		return c.JSON(http.StatusBadRequest, "failed")
	})

	vSphere.GET("/agents", func(c echo.Context) (err error) {
		hosts, err := h.VsphereExporter.VsphereManager.GetVsphereHosts()
		if err == nil {
			return c.JSON(http.StatusBadRequest, "failed")
		}
		return c.JSON(http.StatusOK, hosts)
	})
}
