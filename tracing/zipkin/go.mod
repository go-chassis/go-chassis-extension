module github.com/go-chassis/go-chassis-extension/tracing/zipkin

require (
	github.com/Shopify/sarama v1.27.2 // indirect
	github.com/apache/thrift v0.12.0 // indirect
	github.com/go-chassis/go-chassis/v2 v2.1.0
	github.com/go-chassis/openlog v1.1.2
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/stretchr/testify v1.6.1
)

replace github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92

go 1.13
