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
	case Byte:
		return r.readByte()
	case Bool:
		return r.readBool()
	case Int:
		return r.readInt32()
	case Long:
		return r.readInt64()
	case Float:
		return r.readFloat32()
	case Double:
		return r.readFloat64()
	case String:
		return r.readString()
	}
	return
}

func (r *Reader) readBytes() ([]byte, error) {
	return r.readDelimited()
}

func (r *Reader) readByte() (byte, error) {
	var b byte
	err := r.read(&b)
	return b, err
}

func (r *Reader) readBool() (bool, error) {
	if b, err := r.readByte(); err != nil {
		return false, err
	} else {
		return b != 0, err
	}
}

func (r *Reader) readInt32() (int32, error) {
	var i int32
	err := r.read(&i)
	return i, err
}

func (r *Reader) readInt64() (int64, error) {
	var i int64
	err := r.read(&i)
	return i, err
}

func (r *Reader) readFloat32() (float32, error) {
	var i float32
	err := r.read(&i)
	return i, err
}

func (r *Reader) readFloat64() (float64, error) {
	var i float64
	err := r.read(&i)
	return i, err
}

func (r *Reader) readString() (string, error) {
	if b, err := r.readDelimited(); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (r *Reader) read(i interface{}) error {
	return binary.Read(r.r, binary.LittleEndian, i)
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
