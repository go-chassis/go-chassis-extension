# Installation
```sh
go get github.com/go-chassis/go-chassis-extension/codec/mcpack
```
import in main.go

```go
import _ "github.com/go-chassis/go-chassis-extension/codec/mcpack"
```

in chassis.yaml
```go
servicecomb:
  codec:
    plugin: mcpack
```