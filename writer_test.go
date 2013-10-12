package typedbytes

import (
	"bytes"
	"github.com/markchadwick/spec"
)

var _ = spec.Suite("Typed bytes writer", func(c *spec.C) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)

	c.It("should encode a byte slice", func(c *spec.C) {
		err := w.Write([]byte("hello"))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 4 + 5)
		c.Assert(ByteType(out[0])).Equals(Bytes)

		// 32 bit LE
		c.Assert(out[1]).Equals(uint8(5))

		msg := out[5:]
		c.Assert(string(msg)).Equals("hello")
	})

	c.It("should encode a byte", func(c *spec.C) {
		err := w.Write(byte(24))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 1)
		c.Assert(ByteType(out[0])).Equals(Byte)
		c.Assert(out[1]).Equals(byte(24))
	})

	c.It("should encode true", func(c *spec.C) {
		err := w.Write(true)
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 1)
		c.Assert(ByteType(out[0])).Equals(Bool)
		c.Assert(out[1]).Equals(byte(1))
	})

	c.It("should encode false", func(c *spec.C) {
		err := w.Write(false)
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 1)
		c.Assert(ByteType(out[0])).Equals(Bool)
		c.Assert(out[1]).Equals(byte(0))
	})

	c.It("should encode an int32", func(c *spec.C) {
		err := w.Write(int32(-666))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 4)
		c.Assert(ByteType(out[0])).Equals(Int)
		c.Assert(out[1]).Equals(byte(0x66))
		c.Assert(out[2]).Equals(byte(0xfd))
		c.Assert(out[3]).Equals(byte(0xff))
		c.Assert(out[4]).Equals(byte(0xff))
	})

	c.It("should encode an int64", func(c *spec.C) {
		err := w.Write(int64(-666))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 8)
		c.Assert(ByteType(out[0])).Equals(Long)
	})

	c.It("should encode an float32", func(c *spec.C) {
		err := w.Write(float32(1.23))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 4)
		c.Assert(ByteType(out[0])).Equals(Float)
	})

	c.It("should encode an float64", func(c *spec.C) {
		err := w.Write(float64(1.23))
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 8)
		c.Assert(ByteType(out[0])).Equals(Double)
	})

	c.It("should encode a string", func(c *spec.C) {
		err := w.Write("bonjour")
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 4 + len("bonjour"))
		c.Assert(ByteType(out[0])).Equals(String)
		c.Assert(out[5]).Equals(byte('b'))
		c.Assert(out[11]).Equals(byte('r'))
	})

	c.It("should encode a slice", func(c *spec.C) {
		primes := make([]int32, 0)
		primes = append(primes, 2)
		primes = append(primes, 3)
		primes = append(primes, 5)
		primes = append(primes, 7)

		err := w.Write(primes)
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + 4 + (4 * 5))
		c.Assert(ByteType(out[0])).Equals(Vector)
		c.Assert(out[1]).Equals(byte(4))
	})

	c.It("should encode a channel", func(c *spec.C) {
		ch := make(chan bool)
		done := make(chan bool)
		go func() {
			err := w.Write(ch)
			c.Assert(err).IsNil()
			done <- true
		}()

		for i := 0; i < 5; i++ {
			ch <- i%2 == 0
		}
		close(ch)
		<-done

		out := buf.Bytes()
		c.Assert(out).HasLen(1 + (5 * 2) + 1)
		c.Assert(ByteType(out[0])).Equals(List)
		c.Assert(out[11]).Equals(uint8(255))
	})

	c.It("should encode a map", func(c *spec.C) {
		m := make(map[string]int32)
		m["one"] = 1
		m["two"] = 2
		m["three"] = 3

		err := w.Write(m)
		c.Assert(err).IsNil()

		out := buf.Bytes()
		c.Assert(ByteType(out[0])).Equals(Map)
		c.Assert(out[1]).Equals(byte(3))
	})
})
