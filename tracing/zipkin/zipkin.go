package zipkin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis-plugins/tracing/zipkin/huaweiapm"
	"github.com/go-chassis/go-chassis/core/tracing"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"
)

//const for default values
const (
	DefaultURI           = "http://127.0.0.1:9411/api/v1/spans"
	DefaultBatchSize     = 10000
	DefaultBatchInterval = time.Second * 10
	DefaultCollector     = "http"
)

// NewTracer returns zipkin tracer
func NewTracer(options map[string]string) (opentracing.Tracer, error) {
	uri := options["URI"]
	if uri == "" {
		uri = DefaultURI
	}
	var batchSize = DefaultBatchSize
	bs := options["batchSize"]
	if bs != "" {
		var err error
		batchSize, err = strconv.Atoi(bs)
		if err != nil {
			return nil, fmt.Errorf("can not convert [%s] to batch size", bs)
		}
	}
	var batchInterval = DefaultBatchInterval
	bi := options["batchInterval"]
	if bi != "" {
		var err error
		batchInterval, err = time.ParseDuration(bi)
		if err != nil {
			return nil, fmt.Errorf("can not convert [%s] to batch interval", bi)
		}
	}
	var collectorOption = options["collector"]
	var collector zipkintracer.Collector
	if collectorOption == "" {
		collectorOption = DefaultCollector
	}
	openlogging.GetLogger().Infof("New Zipkin tracer with options %s,%d,%s", uri, batchSize, batchInterval)
	if collectorOption == DefaultCollector {
		var err error
		collector, err = zipkintracer.NewHTTPCollector(uri, zipkintracer.HTTPBatchSize(batchSize), zipkintracer.HTTPBatchInterval(batchInterval))
		if err != nil {
			openlogging.Error(err.Error())
			return nil, fmt.Errorf("unable to create zipkin collector: %+v", err)
		}
	} else if collectorOption == "huaweiapm" {
		collector = huaweiapm.NewCollector(bi, batchSize)
	} else {
		return nil, fmt.Errorf("unable to create zipkin collector: %s", collectorOption)
	}

	// set default recorder
	defaultRecorder := zipkintracer.NewRecorder(collector, false, "0.0.0.0:0", runtime.ServiceName)

	// set tracer map
	tracer, err := zipkintracer.NewTracer(
		defaultRecorder,
		zipkintracer.ClientServerSameSpan(true),
		zipkintracer.TraceID128Bit(true),
		zipkintracer.WithPrefixTracerState(options["prefixTracerState"]),
	)
	if err != nil {
		openlogging.Error(err.Error())
		return nil, fmt.Errorf("unable to create zipkin tracer: %+v", err)
	}
	return tracer, nil
}

func init() {
	tracing.InstallTracer("zipkin", NewTracer)
}
