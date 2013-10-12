package typedbytes

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Reader struct {
	r io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

func (r *Reader) Next() (i interface{}, err error) {
	var bt ByteType
	if bt, err = r.readType(); err != nil {
		return
	}
	switch bt {
	default:
		return nil, fmt.Errorf("Unknown byte type: %v", bt)
	case Bytes:
		return r.readBytes()
	}
	return
}

func (r *Reader) readBytes() ([]byte, error) {
	return r.readDelimited()
}

func (r *Reader) readDelimited() (bs []byte, err error) {
	var length int32
	if err = binary.Read(r.r, binary.LittleEndian, &length); err != nil {
		return
	}
	bs = make([]byte, length)
	_, err = r.r.Read(bs)
	return
}

func (r *Reader) readType() (bt ByteType, err error) {
	bt = ByteType(255)
	err = binary.Read(r.r, binary.LittleEndian, &bt)
	return
}
