package server

import (
	"testing"

	server2 "github.com/go-chassis/go-chassis/v2/core/server"
)

func TestNew(t *testing.T) {
	t.Run("create grpc server with simple options", func(t *testing.T) {
		_ = New(server2.Options{
			Address:   "127.0.0.1:9000",
			BodyLimit: 10 * 1024 * 1024,
		})
	})
}
