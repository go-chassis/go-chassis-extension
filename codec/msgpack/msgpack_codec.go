package msgpack

import (
	"github.com/go-chassis/cari/codec"
	codecChassis "github.com/go-chassis/go-chassis/v2/pkg/codec"
	"github.com/vmihailenco/msgpack/v5"
)

func init() {
	codecChassis.Install("msgpack/v5", newDefault)
}

type Msgp struct {
}

func newDefault(opts codecChassis.Options) (codec.Codec, error) {
	return &Msgp{}, nil
}
func (s *Msgp) Encode(v any) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (s *Msgp) Decode(data []byte, v any) error {
	return msgpack.Unmarshal(data, v)
}
