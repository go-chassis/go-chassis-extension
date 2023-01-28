package msgpack_test

import (
	"testing"

	_ "github.com/go-chassis/go-chassis-extension/codec/msgpack"
	"github.com/go-chassis/go-chassis/v2/pkg/codec"
	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name string
}

func TestMsgPack_Decode(t *testing.T) {
	t.Run("init codec, should success", func(t *testing.T) {
		err := codec.Init(codec.Options{
			Plugin: "msgpack/v5",
		})
		assert.NoError(t, err)
	})
	t.Run("encode and decode, should success", func(t *testing.T) {
		data, err := codec.Encode(Person{Name: "dongjiang"})
		assert.NoError(t, err)
		p := &Person{}
		err = codec.Decode(data, p)
		assert.NoError(t, err)
		assert.Equal(t, "dongjiang", p.Name)
	})
}
