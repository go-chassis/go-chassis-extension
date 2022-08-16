package gojson

import (
	"github.com/go-chassis/cari/codec"
	codecChassis "github.com/go-chassis/go-chassis/v2/pkg/codec"
	"github.com/goccy/go-json"
)

func init() {
	codecChassis.Install("goccy/go-json", newDefault)
}

type GoJson struct {
}

func newDefault(opts codecChassis.Options) (codec.Codec, error) {
	return &GoJson{}, nil
}
func (s *GoJson) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *GoJson) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
