package lambda

import (
	"encoding/json"
	"reflect"

	"github.com/Tanemahuta/aws-lambda-server/pkg/errorx"
	"gopkg.in/yaml.v3"
)

var _ json.Unmarshaler = &Body{}
var _ json.Marshaler = &Body{}
var _ yaml.Unmarshaler = &Body{}

type Body struct {
	Data      []byte `json:"-" yaml:"-"`
	Formatted bool   `json:"-" yaml:"-"`
}

func (b *Body) MarshalJSON() ([]byte, error) {
	if !b.Formatted {
		return json.Marshal(string(b.Data))
	}
	return b.Data, nil
}

func (b *Body) UnmarshalYAML(value *yaml.Node) error {
	var err error
	b.Data, err = b.unmarshal(func(a any) error {
		return value.Decode(a)
	})
	return err
}

func (b *Body) UnmarshalJSON(bytes []byte) error {
	var err error
	b.Data, err = b.unmarshal(func(a any) error {
		return json.Unmarshal(bytes, a)
	})
	return err
}

func (b *Body) unmarshal(fn func(any) error) ([]byte, error) {
	var (
		intermediate any
		result       []byte
	)
	return result, errorx.Fns{
		func() error {
			return fn(&intermediate)
		},
		func() error {
			result = b.dataFrom(intermediate)
			return nil
		},
	}.Run()
}

func (b *Body) dataFrom(intermediate any) []byte {
	val := reflect.ValueOf(intermediate)
	byteTpe := reflect.TypeOf(float64(0))
	strTpe := reflect.TypeOf((*string)(nil)).Elem()
	b.Formatted = false
typeSwitch:
	switch {
	case val.Kind() == reflect.Slice:
		result := make([]byte, val.Len())
		for idx := 0; idx < val.Len(); idx++ {
			elem := reflect.ValueOf(val.Index(idx).Interface())
			if !elem.CanConvert(byteTpe) {
				break typeSwitch
			}
			result[idx] = (byte)(elem.Convert(byteTpe).Float())
		}
		return result
	case val.CanConvert(strTpe):
		return []byte(val.Convert(strTpe).Interface().(string))
	}
	data, _ := json.Marshal(intermediate)
	b.Formatted = true
	return data
}

func (b *Body) String() string {
	return string(b.Data)
}
