package typedbytes

import (
	"bytes"
	"github.com/markchadwick/spec"
)

var _ = spec.Suite("Typed bytes reader", func(c *spec.C) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	r := NewReader(buf)

	c.It("should read a byte slice", func(c *spec.C) {
		w.Write([]byte{1, 2, 3})

		i, err := r.Next()
		c.Assert(i).NotNil()
		c.Assert(err).IsNil()

		c.Assert(i).HasLen(3)
		b := i.([]byte)
		c.Assert(b[0]).Equals(byte(1))
		c.Assert(b[1]).Equals(byte(2))
		c.Assert(b[2]).Equals(byte(3))
	})
})
