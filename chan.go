package typedbytes

import (
	"encoding/binary"
	"io"
	"log"
	"reflect"
)

type chanCodec int

var ChanCodec = chanCodec(0)

func (chanCodec) Write(w io.Writer, v reflect.Value, b WriteBasic) (err error) {
	if err = binary.Write(w, binary.LittleEndian, List); err != nil {
		return
	}
	for {
		next, ok := v.Recv()
		if !ok {
			break
		}
		if err = b(next.Interface()); err != nil {
			return err
		}
	}
	return binary.Write(w, binary.LittleEndian, uint8(255))
}

// WARNING: Readers must be read sequentially. This method will start a
// goroutine to continue reading. The next value may be read from the reader
// before this method is complete, leading to data corruption.
func (chanCodec) Read(r io.Reader, next ReadBasic) (i interface{}, err error) {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			if i, err = next(); err != nil {
				if err == Terminator {
					return
				}
				log.Printf("Error reading to channel: %s", err.Error())
				return
			}
			ch <- i
		}
	}()
	return ch, nil
}
