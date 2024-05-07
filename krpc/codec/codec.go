package codec

import "io"

type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         string
}

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodeFunc func(io.ReadWriteCloser) Codec

type Type string

var NewCodeFuncMap map[Type]NewCodeFunc

func init() {
	NewCodeFuncMap = make(map[Type]NewCodeFunc)
	NewCodeFuncMap[GobType] = NewGobCodec
}
