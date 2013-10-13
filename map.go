package typedbytes

import (
	"encoding/binary"
	"io"
	"reflect"
)

type mapCodec int

var MapCodec = mapCodec(0)

func (mapCodec) Write(w io.Writer, v reflect.Value, b WriteBasic) (err error) {
	if err = binary.Write(w, binary.LittleEndian, Map); err != nil {
		return
	}
	length := v.Len()
	if err = binary.Write(w, binary.LittleEndian, int32(length)); err != nil {
		return
	}
	for _, key := range v.MapKeys() {
		if err = b(key.Interface()); err != nil {
			return
		}
		if err = b(v.MapIndex(key).Interface()); err != nil {
			return
		}
	}
	return
}

func (mapCodec) Read(r io.Reader, next ReadBasic) (_ interface{}, err error) {
	var length int32
	if err = binary.Read(r, binary.LittleEndian, &length); err != nil {
		return
	}
	m := make(map[interface{}]interface{})
	var key interface{}
	var value interface{}
	for i := 0; i < int(length); i++ {
		if key, err = next(); err != nil {
			return
		}
		if value, err = next(); err != nil {
			return
		}
		m[key] = value
	}
	return m, nil
}
