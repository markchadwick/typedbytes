package typedbytes

import (
	"bytes"
	"github.com/markchadwick/spec"
	"io"
)

var _ = spec.Suite("Typed bytes reader", func(c *spec.C) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	r := NewReader(buf)

	c.It("should EOF after the last message", func(c *spec.C) {
		w.Write(true)
		w.Write(false)

		_, err := r.Next()
		c.Assert(err).IsNil()

		_, err = r.Next()
		c.Assert(err).IsNil()

		_, err = r.Next()
		c.Assert(err).NotNil()
		c.Assert(err).Equals(io.EOF)
	})

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

	c.It("should read a byte", func(c *spec.C) {
		w.Write(byte(42))

		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(byte(42))
	})

	c.It("should read true", func(c *spec.C) {
		w.Write(true)
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(true)
	})

	c.It("should read false", func(c *spec.C) {
		w.Write(false)
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(false)
	})

	c.It("should read an int32", func(c *spec.C) {
		w.Write(int32(-123))
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(int32(-123))
	})

	c.It("should read an int64", func(c *spec.C) {
		w.Write(int64(88))
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(int64(88))
	})

	c.It("should read a float32", func(c *spec.C) {
		w.Write(float32(-12.34))
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(float32(-12.34))
	})

	c.It("should read a float64", func(c *spec.C) {
		w.Write(float64(43.21))
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals(float64(43.21))
	})

	c.It("should read a string", func(c *spec.C) {
		w.Write("Hello, world!")
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).Equals("Hello, world!")
	})

	c.It("should read a slice", func(c *spec.C) {
		w.Write([]int32{2, 4, 6, 8})
		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).HasLen(4)

		l, ok := i.([]interface{})
		c.Assert(ok).IsTrue()
		c.Assert(l).HasLen(4)
		c.Assert(l[0]).Equals(int32(2))
		c.Assert(l[1]).Equals(int32(4))
		c.Assert(l[2]).Equals(int32(6))
		c.Assert(l[3]).Equals(int32(8))
	})

	// So, readers must process their contents sequentially. Generally not a
	// problem. However, take the case where the key is a chan. You must call the
	// function with the key and value, however you can't get to the value until
	// you've read the entirety of the channel key. But, without invoking, you
	// can't ever drain the chan, so you're effectively SOL
	c.It("should read a channel", func(c *spec.C) {
		wch := make(chan bool, 2)
		wch <- true
		wch <- false
		close(wch)
		c.Assert(w.Write(wch)).IsNil()

		i, err := r.Next()
		c.Assert(err).IsNil()

		rch, ok := i.(chan interface{})
		c.Assert(ok).IsTrue()

		done := make(chan bool)
		go func() {
			c.Assert(<-rch).Equals(true)
			c.Assert(<-rch).Equals(false)
			done <- true
		}()
		<-done
	})

	c.It("should read a map", func(c *spec.C) {
		m := make(map[string]string)
		m["foo"] = "bar"
		m["baz"] = "zip"
		w.Write(m)

		i, err := r.Next()
		c.Assert(err).IsNil()
		c.Assert(i).HasLen(2)

		n, ok := i.(map[interface{}]interface{})
		c.Assert(ok).IsTrue()
		c.Assert(n["foo"]).Equals("bar")
		c.Assert(n["baz"]).Equals("zip")
	})
})
