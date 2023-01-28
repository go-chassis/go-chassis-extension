# Installation
```sh
go get github.com/go-chassis/go-chassis-extension/codec/msgpack
```
import in main.go

```go
import _ "github.com/go-chassis/go-chassis-extension/codec/msgpack"
```

in chassis.yaml
```go
servicecomb:
  codec:
    plugin: msgpack/v5
```

# MessagePack encoding for Go

## Efficient
msgpack is a drop-in replacement for encoding/json package that can be up to `5 times faster`.

## Customizable
Use custom encoders and decoders to customize serialization for user-defined and stdlib types.

## Extensible
Extend MessagePack by providing type-aware encoding for your types using ext format family.