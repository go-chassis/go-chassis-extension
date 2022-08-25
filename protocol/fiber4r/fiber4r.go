package fiber4r

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

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

	var l net.Listener
	var lIP string
	var lPort string
	success := false
	for i := 0; i < 10; i++ {
		l, lIP, lPort, err = listen(config.Address, config.TLSConfig)
		if err != nil {
			openlog.Warn(fmt.Sprintf("failed to start %s,retry %d ", err, i))
			time.Sleep(1 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("failed to start after 10 times retry: %w ", err)
	}
	registry.InstanceEndpoints[config.ProtocolServerName] = net.JoinHostPort(lIP, lPort) + sslFlag

	go func() {
		err = r.app.Listener(l)
		if err != nil {
			openlog.Error("http server err: " + err.Error())
			server.ErrRuntime <- err
		}
	}()
	if !fiber.IsChild() {
		openlog.Info("parent process")
	} else {
		openlog.Info("child process")
	}

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

// listen start listener with address and tls(if it has), returns the listener and the real listened ip/port
func listen(listenAddress string, tlsConfig *tls.Config) (listener net.Listener, listenedIP string, port string, err error) {
	if tlsConfig == nil {
		listener, err = net.Listen("tcp4", listenAddress)
	} else {
		listener, err = tls.Listen("tcp4", listenAddress, tlsConfig)
	}
	if err != nil {
		return
	}
	realAddr := listener.Addr().String()
	listenedIP, port, err = net.SplitHostPort(realAddr)
	if err != nil {
		return
	}
	ip := net.ParseIP(listenedIP)
	if ip.IsUnspecified() {
		if iputil.IsIPv6Address(ip) {
			listenedIP = iputil.GetLocalIPv6()
			if listenedIP == "" {
				listenedIP = iputil.GetLocalIP()
			}
		} else {
			listenedIP = iputil.GetLocalIP()
		}
	}
	return
}
