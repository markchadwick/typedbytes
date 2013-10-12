package typedbytes

type ByteType uint8

const (
	Bytes ByteType = iota
	Byte
	Bool
	Int
	Long
	Float
	Double
	String
	Vector
	List
	Map
	Gob ByteType = 50
)
