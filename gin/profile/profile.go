package profile

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/router"
	"github.com/go-chassis/openlog"
)

// const
const (
	msgWriteError = "write to response err: "
	//DefaultProfilePath DefaultProfilePath
	DefaultProfilePath = "profile"
	//ProfileRouteRuleSubPath ProfileRouteRuleSubPath
	ProfileRouteRuleSubPath = "route-rule"
	//ProfileDiscoverySubPath ProfileDiscoverySubPath
	ProfileDiscoverySubPath = "discovery"
)

// Profile contains route rule and discovery
type Profile struct {
	RouteRule map[string][]*config.RouteRule              `json:"routeRule"`
	Discovery map[string][]*registry.MicroServiceInstance `json:"discovery"`
}

// AddProfileRoutes AddProfileRoutes
func AddProfileRoutes(profileGroup *gin.RouterGroup) {
	if !archaius.GetBool("servicecomb.profile.enable", false) {
		return
	}
	profilePath := archaius.GetString("servicecomb.profile.apiPath", DefaultProfilePath)
	if !strings.HasPrefix(profilePath, "/") {
		profilePath = "/" + profilePath
	}

	openlog.Info("Enabled profile API on " + profilePath)

	profileGroup.GET(profilePath, httpHandleProfileFunc)

	profileRouteRulePath := profilePath + "/" + ProfileRouteRuleSubPath
	openlog.Info("Enabled profile route-rule API on " + profileRouteRulePath)
	profileGroup.GET(profileRouteRulePath, httpHandleRouteRuleFunc)

	profileDiscoveryPath := profilePath + "/" + ProfileDiscoverySubPath
	openlog.Info("Enabled profile discovery API on " + profileDiscoveryPath)
	profileGroup.GET(profileDiscoveryPath, httpHandleDiscoveryFunc)
}

// HTTPHandleProfileFunc is a gin handler which can expose all profiles in http server
func httpHandleProfileFunc(c *gin.Context) {
	c.JSON(http.StatusOK, newProfile())
}

// HTTPHandleRouteRuleFunc is a gin handler which can expose profile of route rule in http server
func httpHandleRouteRuleFunc(c *gin.Context) {
	c.JSON(http.StatusOK, listRouteRule())
}

// HTTPHandleDiscoveryFunc is a gin handler which can expose profile of discovery in http server
func httpHandleDiscoveryFunc(c *gin.Context) {
	c.JSON(http.StatusOK, listMicroServiceInstance())
}

func newProfile() Profile {
	return Profile{
		RouteRule: listRouteRule(),
		Discovery: listMicroServiceInstance(),
	}
}

func listRouteRule() map[string][]*config.RouteRule {
	return router.DefaultRouter.ListRouteRule()
}

func listMicroServiceInstance() map[string][]*registry.MicroServiceInstance {
	m := make(map[string][]*registry.MicroServiceInstance)
	if registry.MicroserviceInstanceIndex == nil {
		return m
	}
	if registry.MicroserviceInstanceIndex.FullCache() == nil {
		return m
	}

	items := registry.MicroserviceInstanceIndex.FullCache().Items()

	for k, v := range items {
		m[k] = v.Object.([]*registry.MicroServiceInstance)
	}
	return m
}
