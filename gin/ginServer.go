package gin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	archaius "github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis-extension/gin/profile"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

const (
	//Name is a variable of type string which indicates the protocol being used
	Name = "gin"
	//DefaultMetricPath DefaultMetricPath
	DefaultMetricPath = "metrics"
	//MimeFile MimeFile
	MimeFile = "application/octet-stream"
	//MimeMult MimeMult
	MimeMult = "multipart/form-data"
)

const openTLS = "?sslEnabled=true"

func init() {
	server.InstallPlugin(Name, newGinServer)
}

type ginServer struct {
	engine *gin.Engine
	opts   server.Options
	mux    sync.RWMutex
	server *http.Server
}

//Router is to define how route the request
type Router interface {
	//URLPatterns returns route
	URLPatterns(router *gin.RouterGroup)
}

func newGinServer(opts server.Options) server.ProtocolServer {
	engine := gin.Default()
	if archaius.GetBool("servicecomb.metrics.enable", false) {
		metricGroup := engine.Group("")
		metricPath := archaius.GetString("servicecomb.metrics.apiPath", DefaultMetricPath)
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		openlog.Info("Enabled metrics API on " + metricPath)
		metricGroup.GET(metricPath, prometheusHandleFunc)
	}
	profile.AddProfileRoutes(engine.Group(""))
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return &ginServer{
		opts:   opts,
		engine: engine,
	}
}

func prometheusHandleFunc(c *gin.Context) {
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

//GetRouteSpecs is to return a rest API specification of a go struct
func getRouteSpecs(c *ginServer, schema interface{}) error {
	v, ok := schema.(Router)
	if !ok {
		return fmt.Errorf("can not register APIs to server: %s", reflect.TypeOf(schema).String())
	}
	v.URLPatterns(c.engine.Group(""))
	return nil
}

// Invocation2HTTPRequest convert invocation back to http request, set down all meta data
func Invocation2HTTPRequest(inv *invocation.Invocation, c *gin.Context) {
	for k, v := range inv.Metadata {
		c.Set(k, v.(string))
	}
	m := common.FromContext(inv.Ctx)
	for k, v := range m {
		c.Request.Header.Set(k, v)
	}

}

func (r *ginServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	openlog.Info("register rest server")
	opts := server.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, o := range options {
		o(&opts)
	}
	err := getRouteSpecs(r, schema)
	if err != nil {
		return "", err
	}
	return reflect.TypeOf(schema).String(), nil
}

func (r *ginServer) Start() error {
	var err error
	config := r.opts
	r.mux.Lock()
	r.opts.Address = config.Address
	r.mux.Unlock()
	sslFlag := ""
	r.server = &http.Server{
		Addr:         config.Address,
		Handler:      r.engine,
		ReadTimeout:  r.opts.Timeout,
		WriteTimeout: r.opts.Timeout,
		IdleTimeout:  r.opts.Timeout,
	}
	if r.opts.HeaderLimit > 0 {
		r.server.MaxHeaderBytes = r.opts.HeaderLimit
	}
	if r.opts.TLSConfig != nil {
		r.server.TLSConfig = r.opts.TLSConfig
		sslFlag = openTLS
	}

	l, lIP, lPort, err := iputil.StartListener(config.Address, config.TLSConfig)

	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	registry.InstanceEndpoints[config.ProtocolServerName] = net.JoinHostPort(lIP, lPort) + sslFlag

	go func() {
		err = r.server.Serve(l)
		if err != nil {
			openlog.Error("http server err: " + err.Error())
			server.ErrRuntime <- err
		}
	}()

	openlog.Info(fmt.Sprintf("http server is listening at %s", registry.InstanceEndpoints[config.ProtocolServerName]))

	return nil
}

func (r *ginServer) Stop() error {
	if r.server == nil {
		openlog.Info("http server never started")
		return nil
	}
	if err := r.server.Shutdown(context.TODO()); err != nil {
		openlog.Warn("http shutdown error: " + err.Error())
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *ginServer) String() string {
	return Name
}
