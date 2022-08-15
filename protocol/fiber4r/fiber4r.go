package fiber4r

import (
	"errors"
	"fmt"
	"math"
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
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	//Name is a variable of type string which indicates the protocol being used
	Name = "rest"
)

const openTLS = "?sslEnabled=true"

func init() {
	InstallPlugin()
}

// InstallPlugin Install gin Plugin
func InstallPlugin() {
	server.InstallPlugin(Name, NewServer)
}

type fiberServer struct {
	app  *fiber.App
	opts server.Options
	mux  sync.RWMutex
}

func NewServer(opts server.Options) server.ProtocolServer {
	r := &fiberServer{
		opts: opts,
	}
	r.app = fiber.New(fiber.Config{
		ReadTimeout:  r.opts.Timeout,
		WriteTimeout: r.opts.Timeout,
		IdleTimeout:  r.opts.Timeout,
	})

	return r
}

func prometheusHandleFunc() http.Handler {
	return promhttp.HandlerFor(metrics.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{})
}

func (r *fiberServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	openlog.Info("register rest handler")
	opts := server.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, o := range options {
		o(&opts)
	}
	app, ok := schema.(*fiber.App)
	if !ok {
		return "", errors.New("must register *fiber.App")
	}
	r.app = app
	return reflect.TypeOf(schema).String(), nil
}

// addMetricsRoute user fiber adaptor to fit in prom native http handler
func addMetricsRoute(opts server.Options, app *fiber.App) {
	if opts.MetricsEnable {
		metricPath := opts.MetricsAPI
		if metricPath == "" {
			metricPath = "metrics"
		}
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		openlog.Info("Enabled metrics API on " + metricPath)
		app.Get(metricPath, adaptor.HTTPHandler(prometheusHandleFunc()))
	}
}
func (r *fiberServer) Start() error {
	var err error
	config := r.opts
	r.mux.Lock()
	r.opts.Address = config.Address
	r.mux.Unlock()
	sslFlag := ""

	if r.opts.BodyLimit > 0 && r.opts.BodyLimit < math.MaxInt32 {
		r.app.Server().MaxRequestBodySize = int(r.opts.BodyLimit)
	}
	if r.opts.TLSConfig != nil {
		r.app.Server().TLSConfig = r.opts.TLSConfig
		sslFlag = openTLS
	}
	if r.opts.MetricsEnable {
		addMetricsRoute(r.opts, r.app)
	}

	l, lIP, lPort, err := iputil.StartListener(config.Address, config.TLSConfig)

	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}

	registry.InstanceEndpoints[config.ProtocolServerName] = net.JoinHostPort(lIP, lPort) + sslFlag

	go func() {
		err = r.app.Listener(l)
		if err != nil {
			openlog.Error("http server err: " + err.Error())
			server.ErrRuntime <- err
		}
	}()

	openlog.Info(fmt.Sprintf("http server is listening at %s", registry.InstanceEndpoints[config.ProtocolServerName]))

	return nil
}

func (r *fiberServer) Stop() error {
	if r.app == nil {
		openlog.Info("http server never started")
		return nil
	}
	if err := r.app.Shutdown(); err != nil {
		openlog.Warn("http shutdown error: " + err.Error())
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *fiberServer) String() string {
	return Name
}
