package server_test

import (
	"github.com/go-chassis/go-chassis-extension/protocol/grpc/server"
	server2 "github.com/go-chassis/go-chassis/v2/core/server"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("create grpc server with simple options", func(t *testing.T) {
		_ = server.New(server2.Options{
			Address: "127.0.0.1:9000",
		})
	})
}
