package engine

// Codec 编解码
type Codec interface {
	// Encode 编码
	Encode(data interface{}) []byte
	// Decoce 解码
	Decode(data []byte) interface{}
}

type defaultCodec struct {
}

// Encode 编码
func (d *defaultCodec) Encode(data interface{}) []byte {
	panic("not implemented") // TODO: Implement
}

// Decoce 解码
func (d *defaultCodec) Decode(data []byte) interface{} {
	panic("not implemented") // TODO: Implement
}
