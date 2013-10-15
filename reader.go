package typedbytes

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var Terminator = errors.New("Terminator byte")

type ReadBasic func() (interface{}, error)

type Decoder interface {
	Read(io.Reader, ReadBasic) (interface{}, error)
}

type Reader struct {
	r        io.Reader
	decoders map[ByteType]Decoder
}

func NewReader(r io.Reader) *Reader {
	reader := &Reader{
		r:        r,
		decoders: make(map[ByteType]Decoder),
	}
	reader.Register(Vector, SliceCodec)
	reader.Register(List, ChanCodec)
	reader.Register(Map, MapCodec)

	return reader
}

func (r *Reader) Register(b ByteType, dec Decoder) {
	r.decoders[b] = dec
}

func (r *Reader) Next() (i interface{}, err error) {
	var bt ByteType
	if bt, err = r.readType(); err != nil {
		return
	}
	switch bt {
	default:
		return r.readComplex(bt)
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
	case ByteType(255):
		return nil, Terminator
	}
}

func (r *Reader) readComplex(bt ByteType) (interface{}, error) {
	dec, ok := r.decoders[bt]
	if !ok {
		return nil, fmt.Errorf("No decoder for byte type: %d", bt)
	}
	return dec.Read(r.r, r.Next)
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
	return binary.Read(r.r, binary.BigEndian, i)
}

func (r *Reader) readDelimited() (bs []byte, err error) {
	var length int32
	if err = binary.Read(r.r, binary.BigEndian, &length); err != nil {
		return
	}
	bs = make([]byte, length)
	_, err = r.r.Read(bs)
	return
}

func (r *Reader) readType() (bt ByteType, err error) {
	bt = ByteType(255)
	err = binary.Read(r.r, binary.BigEndian, &bt)
	return
}
