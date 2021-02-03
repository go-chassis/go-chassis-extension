# go-chssis-gin

#### Description
gin plugin for [go-chassis](https://github.com/go-chassis/go-chassis)



#### Installation

in chassis.yaml, add following
```yaml
---
servicecomb:
  protocols:
    rest:            # use protocol gin
      listenAddress: 127.0.0.1:5001
  metrics:
    apiPath: /metrics
    enable: true
    enableGoRuntimeMetrics: true
    enableCircuitMetrics: true
```

in main.go, add following
``` go
package main

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/go-chassis-extension/protocol/gin4r"

	"github.com/gin-gonic/gin"
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/openlog"
)

// HelloV2 HelloV2
func HelloV2(c *gin.Context) {
	c.Writer.WriteString(fmt.Sprintf("hello v2 %s", c.Request.RemoteAddr))
}

// InitHelloV2Route InitHelloV2Route
func InitHelloV2Route(router *gin.RouterGroup) {
	v2Router := router.Group("v2")
	v2Router.GET("/hello", HelloV2)
}

// Hello Hello
func Hello(c *gin.Context) {
	c.Writer.WriteString(fmt.Sprintf("hello %s", c.Request.RemoteAddr))
}

// InitHelloRoute InitHelloRoute
func InitHelloRoute(router *gin.RouterGroup) {
	router.GET("/hello", Hello)
}

// Cors process cors request, support options method
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Token, X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// pass through all OPTIONS method
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// continue process request
		c.Next()
	}
}

func main() {
	gin4r.InstallPlugin()
	// init default gin engine
	r := gin.Default()
	// add cors middleware for gin
	r.Use(Cors())
	// init gin server route
	privateGroup := r.Group("")
	InitHelloRoute(privateGroup)
	InitHelloV2Route(privateGroup)
	// register gin to go-chassis rest protocol
	chassis.RegisterSchema("rest", r)
	if err := chassis.Init(); err != nil {
		openlog.Fatal("Init failed." + err.Error())
		return
	}
	chassis.Run()
}


```

then, visit http://127.0.0.1:5001/hello
