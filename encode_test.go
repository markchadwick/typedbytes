package typedbytes

import (
	"github.com/markchadwick/spec"
)

var _ = spec.Suite("encode", func(c *spec.C) {
	c.It("should simply encode a value", func(c *spec.C) {
		bs, err := Encode("Today is the day!")
		c.Assert(err).IsNil()
		c.Assert(bs).HasLen(22)
		c.Assert(string(bs[14:])).Equals("the day!")
	})
})
