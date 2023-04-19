package gozero4r_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-chassis/go-chassis-extension/protocol/gozero4r"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/zero-contrib/router/mux"

	"github.com/stretchr/testify/assert"
)

func TestGoZero4r_Start(t *testing.T) {
	s := gozero4r.NewGoZeroServer(server.Options{
		Address:       "127.0.0.1:23312",
		MetricsEnable: true,
	})
	r := mux.NewRouter()
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Port:     23312,
		Timeout:  20000,
		MaxConns: 500,
	}, rest.WithRouter(r))

	engine.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/hello",
		Handler: Hello,
	})
	t.Run("register right ptr,should success", func(t *testing.T) {
		_, err := s.Register(engine)
		assert.NoError(t, err)
	})
	err := s.Start()
	assert.NoError(t, err)
	defer s.Stop()
	time.Sleep(2 * time.Second)

	t.Run("call http server, should success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:23312/hello", nil)
		c := http.DefaultClient
		r, err := c.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, []byte(`hello test`), httputil.ReadBody(r))
	})
	t.Run("call http metrics api, should success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:23312/metrics", nil)
		c := http.DefaultClient
		r, err := c.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, r.StatusCode)
	})
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("hello test")))
}
