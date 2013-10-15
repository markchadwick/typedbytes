package typedbytes

import (
	"bytes"
	"encoding/binary"
	"github.com/markchadwick/spec"
	"reflect"
)

type person struct {
	Name    string
	Age     uint32
	Fun     bool
	Friends []*person
}

var _ = spec.Suite("GOB Codec", func(c *spec.C) {
	buf := new(bytes.Buffer)

	c.It("should read/write a simple struct", func(c *spec.C) {
		p0 := &person{
			Name: "Person 1",
			Age:  6,
			Fun:  true,
		}
		err := GobCodec.Write(buf, reflect.ValueOf(p0), nil)
		c.Assert(err).IsNil()

		var bt ByteType
		err = binary.Read(buf, binary.BigEndian, &bt)
		c.Assert(err).IsNil()
		c.Assert(bt).Equals(Gob)

		c.Skip("May need more information")
		i, err := GobCodec.Read(buf, nil)
		c.Assert(err).IsNil()
		c.Assert(i).NotNil()
	})
})
