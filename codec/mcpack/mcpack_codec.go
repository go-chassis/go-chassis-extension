package mcpack

import (
	"github.com/go-chassis/cari/codec"
	codecChassis "github.com/go-chassis/go-chassis/v2/pkg/codec"
	"github.com/kubeservice-stack/common/pkg/codec/mcpack"
)

func init() {
	codecChassis.Install("mcpack", newDefault)
}

type MCPack struct {
}

func newDefault(opts codecChassis.Options) (codec.Codec, error) {
	return &MCPack{}, nil
}
func (s *MCPack) Encode(v any) ([]byte, error) {
	return mcpack.Marshal(v)
}

func (s *MCPack) Decode(data []byte, v any) error {
	return mcpack.Unmarshal(data, v)
}
