# Installation
```sh
go get github.com/go-chassis/go-chassis-extension/codec/gojson
```
import in main.go

```go
import _ "github.com/go-chassis/go-chassis-extension/codec/gojson"
```

in chassis.yaml
```go
servicecomb:
  codec:
    plugin: goccy/go-json
```
