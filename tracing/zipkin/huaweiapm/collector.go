package huaweiapm

import (
	"github.com/go-chassis/huawei-apm"
	"github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin-contrib/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

//Collector collects span to huawei apm
type collector struct {
	reporter *huaweiapm.TracingReporter
}

func NewCollector(interval string, size int) zipkintracer.Collector {
	c := &collector{
		reporter: huaweiapm.NewTracingReporter(interval, size),
	}
	return c
}

// Collect serialize the zipkin spans and write into the fifo
func (c *collector) Collect(s *zipkincore.Span) error {
	c.reporter.WriteSpan(s)
	return nil
}

// Close will never be called
func (c *collector) Close() error {
	return nil
}
