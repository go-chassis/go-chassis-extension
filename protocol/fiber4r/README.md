# Install
```shell
go get github.com/go-chassis/go-chassis-extension/protocol/fiber4r
```
in your main.go, add one line
```go
import _ "github.com/go-chassis/go-chassis-extension/protocol/fiber4r"
```
then [fiber](https://github.com/gofiber/fiber) will replace [default rest implementation](https://github.com/emicklei/go-restful).
# How to collocate Fiber
```go
app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World 👋!")
    })
chassis.Register("rest", app)
```

