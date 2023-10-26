package config

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v3"
)

var _ json.Unmarshaler = &Duration{}
var _ json.Marshaler = &Duration{}
var _ yaml.Unmarshaler = &Duration{}
var _ yaml.Marshaler = &Duration{}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	return d.parse(str)
}

func (d *Duration) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	return d.parse(str)
}

func (d *Duration) parse(str string) error {
	decorated, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*d = Duration{Duration: decorated}
	return nil
}

func (d *Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
