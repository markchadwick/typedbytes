# TypedBytes

A basic [Go](http://golang.org/) reader and writer for
[Hadoop](http://hadoop.apache.org/)'s [typed
bytes](http://hadoop.apache.org/docs/current/api/org/apache/hadoop/typedbytes/package-summary.html)
implementation.

## Example

```go
import (
  "github.com/markchadwick/typedbytes"
  "log"
)

func main() {
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
        if err != io.EOF {
          log.Printf("Error reading: %s", err.Error())
        }
        return
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
```

Yields the following output:
```
Read value: true
Read value: false
Read value: [true false]
Read value: 123
Read value: 123
Read value: map[name:Frank job:Fun]
```

## Types

|    | Hadoop type | Go type                        |
|---:|-------------|--------------------------------|
| 0  | Byte Seq    | `[]byte`                       |
| 1  | Byte        | `byte`                         |
| 2  | Boolean     | `bool`                         |
| 3  | Integer     | `int32`                        |
| 4  | Long        | `int64`                        |
| 5  | Float       | `float32`                      |
| 6  | Double      | `float64`                      |
| 7  | String      | `string`                       |
| 8  | Vector      | `[]interface{}`                |
| 9  | List        | `chan interface{}`             |
| 10 | Map         | `map[interface{}] interface{}` |

Attempting to read or write unrecognized types will result in an error

## Custom Types
Non-primitive types are registered with the `Reader` and `Writer` instances via
the `Register` method. To implement a custom type, read the respective
constructor methods where the vector, list, and map types are implemented. A
unique Hadoop typed byte ID will need to be assigned.

## Caveats

There is a specification for a "List" type which contains no length. This is
mapped to a Go channel. However, typed byte files must be read sequentially and
it is possible (and easy) to invoke a new read iteration while a channel is
still being read. It's really up to the consuming library to make sure this
doesn't happen.
