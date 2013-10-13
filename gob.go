package typedbytes

import (
	"encoding/binary"
	"encoding/gob"
	"io"
	"reflect"
)

type gobCodec int

var GobCodec = gobCodec(0)

func (gobCodec) Write(w io.Writer, v reflect.Value, b WriteBasic) (err error) {
	if err = binary.Write(w, binary.LittleEndian, Gob); err != nil {
		return
	}
	return gob.NewEncoder(w).EncodeValue(v)
}

func (gobCodec) Read(r io.Reader, next ReadBasic) (i interface{}, err error) {
	dec := gob.NewDecoder(r)
	err = dec.Decode(i)
	return
}
