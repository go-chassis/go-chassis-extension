package gojson_test

import (
	_ "github.com/go-chassis/go-chassis-extension/codec/gojson"
	"github.com/go-chassis/go-chassis/v2/pkg/codec"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name string
}

func TestGoJson_Decode(t *testing.T) {
	t.Run("init codec, should success", func(t *testing.T) {
		err := codec.Init(codec.Options{
			Plugin: "goccy/go-json",
		})
		assert.NoError(t, err)
	})
	t.Run("encode and decode, should success", func(t *testing.T) {
		data, err := codec.Encode(Person{Name: "a"})
		assert.NoError(t, err)
		p := &Person{}
		err = codec.Decode(data, p)
		assert.NoError(t, err)
		assert.Equal(t, "a", p.Name)
	})
}
