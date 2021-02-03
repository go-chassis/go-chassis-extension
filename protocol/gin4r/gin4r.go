package gin4r

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
	"github.com/go-chassis/go-chassis-extension/protocol/gin4r/profile"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	//Name is a variable of type string which indicates the protocol being used
	Name = "rest"
	//DefaultMetricPath DefaultMetricPath
	DefaultMetricPath = "metrics"
	//MimeFile MimeFile
	MimeFile = "application/octet-stream"
	//MimeMult MimeMult
	MimeMult = "multipart/form-data"
)

const openTLS = "?sslEnabled=true"

func init() {
	InstallPlugin()
}

// InstallPlugin Install gin Plugin
func InstallPlugin() {
	server.InstallPlugin(Name, newGinServer)
}

type ginServer struct {
	engine *gin.Engine
	opts   server.Options
	mux    sync.RWMutex
	server *http.Server
}

func newGinServer(opts server.Options) server.ProtocolServer {
	return &ginServer{
		opts:   opts,
		engine: nil,
	}
}

func prometheusHandleFunc(c *gin.Context) {
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(c.Writer, c.Request)
}

//GetRouteSpecs is to return a rest API specification of a go struct
func (r *ginServer) initGinRoute(schema interface{}) error {
	engine, ok := schema.(*gin.Engine)
	if !ok {
		return fmt.Errorf("can not register APIs to server: %s", reflect.TypeOf(schema).String())
	}
	r.engine = engine
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

	return nil
}

func (r *ginServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	openlog.Info("register rest server")
	opts := server.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, o := range options {
		o(&opts)
	}
	err := r.initGinRoute(schema)
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
