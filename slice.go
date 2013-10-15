package typedbytes

import (
	"encoding/binary"
	"io"
	"reflect"
)

type sliceCodec int

var SliceCodec = sliceCodec(0)

func (sliceCodec) Write(w io.Writer, v reflect.Value, b WriteBasic) (err error) {
	if err = binary.Write(w, binary.BigEndian, Vector); err != nil {
		return
	}
	length := v.Len()
	if err = binary.Write(w, binary.BigEndian, int32(length)); err != nil {
		return
	}
	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()
		if err = b(item); err != nil {
			return
		}
	}
	return nil
}

func (sliceCodec) Read(r io.Reader, next ReadBasic) (_ interface{}, err error) {
	var length int32
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}
	vs := make([]interface{}, length)
	for i := 0; i < int(length); i++ {
		item, err := next()
		if err != nil {
			return nil, err
		}
		vs[i] = item
	}
	return vs, err
}
