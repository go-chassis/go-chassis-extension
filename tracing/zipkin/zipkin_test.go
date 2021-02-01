package zipkin_test

import (
	"github.com/go-chassis/go-chassis-extension/tracing/zipkin"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTracer(t *testing.T) {
	tracer, err := zipkin.NewTracer(map[string]string{
		"URI":           "https://10.1.1.1",
		"batchInterval": "1s",
	})
	assert.NotNil(t, tracer)
	assert.NoError(t, err)

	tracer, err = zipkin.NewTracer(map[string]string{
		"URI":       "",
		"batchSize": "fake",
	})
	assert.Error(t, err)

	tracer, err = zipkin.NewTracer(map[string]string{
		"URI":           "",
		"batchInterval": "30q",
	})
	assert.Error(t, err)
}
