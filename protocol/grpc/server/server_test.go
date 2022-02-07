package server

import (
	"testing"

	server2 "github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestNew(t *testing.T) {
	GrpcServerOptions(grpc.MaxRecvMsgSize(10 * 1024 * 1024))
	t.Run("create grpc server with simple options", func(t *testing.T) {
		_ = New(server2.Options{
			Address: "127.0.0.1:9000",
		})

		assert.Equal(t, 1, len(grpcServerOptions))
	})
}
