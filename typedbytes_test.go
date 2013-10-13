package typedbytes

import (
	"bytes"
	"github.com/markchadwick/spec"
	"io"
	"log"
	"testing"
)

func example() {
	buf := new(bytes.Buffer)
	r := NewReader(buf)
	w := NewWriter(buf)

	// Write some values to the buffer
	log.Printf("-----------------------------------------------------")
	w.Write(true)
	w.Write(false)
  w.Write([]bool{true, false})
  w.Write(int32(123))
  w.Write(float64(123))
  w.Write(map[string]string{"name": "Frank", "job": "Fun"})

	// Print each value read from the buffer
	for {
		value, err := r.Next()
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("Error reading: %s", err.Error())
		}
		log.Printf("Read value: %v", value)
	}
	log.Printf("-----------------------------------------------------")
}

func Test(t *testing.T) {
	spec.Run(t)
}
