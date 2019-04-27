package json

import (
	"encoding/json"
	"errors"

	"github.com/fatih/structs"
)

// WrapRawMessage ast type JSON.RawMessage
type WrapRawMessage interface {
	DecodeRawMessage(data []byte, scfield *structs.Field, swfield *structs.Field) error
}

// WrapNumber ast type JSON.Number
type WrapNumber interface {
	DecodeNumber(data json.Number, scfield *structs.Field, swfield *structs.Field) error
}

// DecodeObject is a suit like for object custom decode
func DecodeObject(data []byte, wrap interface{}, obj interface{}) error {
	if err := json.Unmarshal(data, wrap); err != nil {
		return err
	}
	return decodeProxy(wrap, obj)
}

func decodeProxy(wrap interface{}, obj interface{}) error {
	sw := structs.New(wrap)
	sc := structs.New(obj)
	for _, swfield := range sw.Fields() {
		if !swfield.IsExported() {
			continue
		}
		name := swfield.Name()
		scfield := sc.Field(name)
		switch swval := swfield.Value().(type) {
		case json.Number:
			wrapper, ok := wrap.(WrapNumber)
			if !ok {
				return errors.New("Please Implements: WrapNumber")
			}
			if err := wrapper.DecodeNumber(swval, scfield, swfield); err != nil {
				return err
			}
		case json.RawMessage:
			if len(swval) == 2 {
				swval = []byte(`[]`)
			}
			wrapper, ok := wrap.(WrapRawMessage)
			if !ok {
				return errors.New("Please Implements: WrapRawMessage")
			}
			if err := wrapper.DecodeRawMessage(swval, scfield, swfield); err != nil {
				return err
			}
		default:
			scfield.Set(swfield.Value())
		}
	}
	return nil
}

func main() {

}
