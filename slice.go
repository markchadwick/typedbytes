package typedbytes

import (
	"encoding/binary"
	"io"
	"reflect"
)

type SliceEncoder int

func (se SliceEncoder) Write(w io.Writer, v reflect.Value, b WriteBasic) (err error) {
	if err = binary.Write(w, binary.LittleEndian, Vector); err != nil {
		return
	}
	length := v.Len()
	if err = binary.Write(w, binary.LittleEndian, int32(length)); err != nil {
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

func (se SliceEncoder) Read(r io.Writer, v reflect.Value, b WriteBasic) (err error) {
	return nil
}
