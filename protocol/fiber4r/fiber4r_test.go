package fiber4r_test

import (
	"github.com/go-chassis/go-chassis-extension/protocol/fiber4r"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestFiberServer_Start(t *testing.T) {
	s := fiber4r.NewServer(server.Options{
		Address:       "127.0.0.1:3000",
		MetricsEnable: true,
	})
	app := fiber.New()
	app.Get("/hello", Handler)
	t.Run("register right ptr,should success", func(t *testing.T) {
		_, err := s.Register(app)
		assert.NoError(t, err)
	})
	t.Run("register wrong ptr,should fail", func(t *testing.T) {
		_, err := s.Register("")
		assert.Error(t, err)
	})
	go s.Start()
	time.Sleep(2 * time.Second)

	t.Run("call http server, should success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/hello", nil)
		c := http.DefaultClient
		r, err := c.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, []byte(`Hello, World 👋!`), httputil.ReadBody(r))
	})
	t.Run("call http metrics api, should success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/metrics", nil)
		c := http.DefaultClient
		r, err := c.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, r.StatusCode)
	})

}
func Handler(c *fiber.Ctx) error {
	return c.SendString("Hello, World 👋!")
}
