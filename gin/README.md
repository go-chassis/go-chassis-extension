# go-chssis-gin

#### Description
gin plugin for [go-chassis](https://github.com/go-chassis/go-chassis)



#### Installation

in chassis.yaml, add following
```yaml
---
servicecomb:
  protocols:
    gin:
      listenAddress: 127.0.0.1:5001
```

in main.go, add following
``` go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-chassis/go-chassis/v2"
    "github.com/go-chassis/openlog"
    
    _ "github.com/go-chassis/go-chassis-extension/gin"
)

type RestFulHello struct {
}

// Root Root
func (r *RestFulHello) Hello(c *gin.Context) {
    c.Writer.WriteString(fmt.Sprintf("hello %s", c.Request.RemoteAddr))
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns(router *gin.RouterGroup) {
    router.GET("/hello", r.Hello)
}

func main() {
    chassis.RegisterSchema("gin", &RestFulHello{})
    if err := chassis.Init(); err != nil {
        openlog.Fatal("Init failed." + err.Error())
        return
    }
    chassis.Run()
}

```

then, visit http://127.0.0.1:5001/hello
