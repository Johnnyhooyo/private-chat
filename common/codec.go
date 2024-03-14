package common

import "encoding/json"

type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

var InUseCodec Codec = &jsonCodec{}

type jsonCodec struct {
}

func (j *jsonCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
