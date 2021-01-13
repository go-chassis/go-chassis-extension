package client_test

import (
	"github.com/go-chassis/go-chassis-extension/protocol/grpc/client"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestTransformContext(t *testing.T) {
	ctx := common.NewContext(map[string]string{
		"1": "2",
		"3": "4",
	})
	ctx = client.TransformContext(ctx)
	md, ok := metadata.FromOutgoingContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "2", md["1"][0])
	assert.Equal(t, "4", md["3"][0])
}

func TestNew(t *testing.T) {

}
