package typedbytes

import (
	"bufio"
	"bytes"
	"io"
)

type noOpCloser struct {
	*bytes.Buffer
}

func NoOpCloser() *noOpCloser {
	return &noOpCloser{new(bytes.Buffer)}
}

func (noc *noOpCloser) Close() error {
	return nil
}

type Buffer struct {
	r       *io.PipeReader
	w       *io.PipeWriter
	rb      *bufio.Reader
	wb      *bufio.Writer
	written *bytes.Buffer
}

func NewBuffer() *Buffer {
	r, w := io.Pipe()
	return &Buffer{
		r:       r,
		w:       w,
		rb:      bufio.NewReader(r),
		wb:      bufio.NewWriter(w),
		written: new(bytes.Buffer),
	}
}

func (b *Buffer) Read(bs []byte) (n int, err error) {
	n, err = b.rb.Read(bs)
	if err == io.ErrClosedPipe {
		err = io.EOF
	}
	return
}

func (b *Buffer) Write(bs []byte) (int, error) {
	b.written.Write(bs)
	return b.wb.Write(bs)
}

func (b *Buffer) Bytes() []byte {
	return b.written.Bytes()
}

func (b *Buffer) Close() (err error) {
	if err = b.wb.Flush(); err != nil {
		return
	}
	if err = b.w.Close(); err != nil {
		return
	}
	go b.r.Close()
	return
}
