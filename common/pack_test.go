package common

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	pack := NewDefaultPacker()
	data, _ := json.Marshal("{\"ff\":\"ffff\"}")

	bytes, err := pack.Pack(data)
	if err != nil {
		t.Errorf("pack err")
		return
	}
	fmt.Println(bytes)

	unpackData, err := pack.Unpack(bytes2.NewReader(bytes))
	if err != nil {
		t.Errorf("unpack err")
		return
	}
	fmt.Println(unpackData)
	fmt.Println(string(unpackData))
}
