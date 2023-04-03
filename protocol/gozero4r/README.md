# go-chssis-gozero

#### Description
gozero plugin for [go-chassis](https://github.com/go-chassis/go-chassis)


#### Installation

in chassis.yaml, add following
```yaml
---
servicecomb:
  protocols:
    gozero:            # use protocol gozero
      listenAddress: 127.0.0.1:8001
  metrics:
    apiPath: /metrics
    enable: true
```

in main.go, add following
``` go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/go-chassis-extension/protocol/gozero4r"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/zero-contrib/handler"
	"github.com/zeromicro/zero-contrib/router/mux"
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
)

// Hello Hello
func Hello(w http.ResponseWriter, r *http.Request) {
	w.Writer([]byte(fmt.Sprintf("hello %s", c.Request.RemoteAddr)))
}

func main() {
	gozero4r.InstallPlugin()
	// init router 
	r := mux.NewRouter()
	// init default rest server
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Port: *port,
		Timeout:  20000,
		MaxConns: 500,
	},rest.WithRouter(r)))
	defer engine.Stop()
	// add cors middleware for gozero
	server.Use(handler.NewETagMiddleware(true).Handle)
	// init gozero server route
	engine.AddRoute(
		Method: http.MethodGet,
		Path:   "/hello",
		Handler: Hello,
	)
	// register gozero to go-chassis rest protocol
	chassis.RegisterSchema("rest", r)
	if err := chassis.Init(); err != nil {
		openlog.Fatal("Init failed." + err.Error())
		return
	}
	chassis.Run()
}


```

then, visit http://127.0.0.1:8001/hello
