module github.com/go-chassis/go-chassis-extension/tracing/zipkin

require (
	github.com/go-chassis/go-chassis v1.8.3
	github.com/go-chassis/go-chassis-plugins v0.0.0-20200511232319-d84d8c0fadb4
	github.com/go-mesh/openlogging v1.0.1
	github.com/opentracing/opentracing-go v1.0.2
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/stretchr/testify v1.4.0
)

replace github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5 => github.com/go-chassis/zipkin-go-opentracing v0.3.5-0.20190321072447-42cf74fc2a92

go 1.13
