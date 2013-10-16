package typedbytes

import (
	"bytes"
	"reflect"
)

var encWriter *Writer

func init() {
	encWriter = new(Writer)
	encWriter.encoders = make(map[reflect.Kind]Encoder)
	encWriter.Register(reflect.Slice, SliceCodec)
	encWriter.Register(reflect.Chan, ChanCodec)
	encWriter.Register(reflect.Map, MapCodec)
}

type encwriter struct {
	*bytes.Buffer
}

func (ew *encwriter) Close() error {
	return nil
}

// Convience function to get some bytes without setting up a writer
func Encode(i interface{}) (bs []byte, err error) {
	buf := new(bytes.Buffer)
	enc := &encwriter{buf}
	encWriter.w = enc

	err = encWriter.Write(i)
	bs = buf.Bytes()
	return
}
