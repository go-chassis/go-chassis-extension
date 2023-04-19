package gozero4r

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/server"

	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/handler"
)

const (
	//Name is a variable of type string which indicates the protocol being used
	Name = "rest"
	//DefaultMetricPath DefaultMetricPath
	DefaultMetricPath = "metrics"
)

const openTLS = "?sslEnabled=true"

func init() {
	InstallPlugin()
}

// InstallPlugin Install go-zero rest Plugin
func InstallPlugin() {
	server.InstallPlugin(Name, NewGoZeroServer)
}

type gozeroServer struct {
	engine *rest.Server
	opts   server.Options
	mux    sync.RWMutex
	server *http.Server
}

func NewGoZeroServer(opts server.Options) server.ProtocolServer {
	return &gozeroServer{
		opts:   opts,
		engine: nil,
	}
}

func prometheusHandleFunc(w http.ResponseWriter, r *http.Request) {
	promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

//GetRouteSpecs is to return a rest API specification of a go struct
func (r *gozeroServer) initGoZeroRoute(schema interface{}) error {
	engine, ok := schema.(*rest.Server)
	if !ok {
		return fmt.Errorf("can not register APIs to server: %s", reflect.TypeOf(schema).String())
	}
	r.engine = engine

	if r.opts.MetricsEnable {
		metricPath := r.opts.MetricsAPI
		if metricPath == "" {
			metricPath = DefaultMetricPath
		}
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		openlog.Info("Enabled metrics API on " + metricPath)
		engine.Use(rest.ToMiddleware(handler.PrometheusHandler(metricPath)))
		engine.AddRoute(rest.Route{
			Method:  http.MethodGet,
			Path:    metricPath,
			Handler: prometheusHandleFunc,
		})
	}

	return nil
}

func (r *gozeroServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	openlog.Info("register go-zero/rest server")
	opts := server.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, o := range options {
		o(&opts)
	}
	err := r.initGoZeroRoute(schema)
	if err != nil {
		return "", err
	}
	return reflect.TypeOf(schema).String(), nil
}

func (r *gozeroServer) Start() error {
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
	go r.engine.Start()
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

func (r *gozeroServer) Stop() error {
	if r.server == nil {
		openlog.Info("http server never started")
		return nil
	}
	r.engine.Stop()
	if err := r.server.Shutdown(context.TODO()); err != nil {
		openlog.Warn("http shutdown error: " + err.Error())
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *gozeroServer) String() string {
	return Name
}
