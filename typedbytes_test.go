package typedbytes

import (
	"github.com/markchadwick/spec"
	"io"
	"log"
	"testing"
)

func example() {
	buf := NewBuffer()
	r := NewReader(buf)
	w := NewWriter(buf)
	done := make(chan bool)

	// Start a goroutine to print every value to the log
	go func() {
		defer func() {
			done <- true
		}()

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
	}()

	// Write some values to the buffer
	w.Write(true)
	w.Write(false)
	w.Write([]bool{true, false})
	w.Write(int32(123))
	w.Write(float64(123))
	w.Write(map[string]string{"name": "Frank", "job": "Fun"})
	w.Close()
	<-done
}

func Test(t *testing.T) {
	spec.Run(t)
}
