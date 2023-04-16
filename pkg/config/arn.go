package config

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
)

var _ json.Unmarshaler = &ARN{}
var _ yaml.Unmarshaler = &ARN{}

// ARN suitable for JSON unmarshalling.
type ARN struct {
	arn.ARN
}

func (a *ARN) UnmarshalYAML(value *yaml.Node) error {
	var err error
	a.ARN, err = arn.Parse(value.Value)
	return err
}

func (a ARN) wrapped() arn.ARN {
	return a.ARN
}

func (a *ARN) UnmarshalJSON(bytes []byte) error {
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	a.ARN, err = arn.Parse(str)
	return err
}
