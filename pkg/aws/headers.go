package aws

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var _ json.Unmarshaler = &Headers{}
var _ json.Marshaler = &Headers{}
var _ yaml.Unmarshaler = &Headers{}

type Headers http.Header

func (h *Headers) MarshalJSON() ([]byte, error) {
	if h == nil {
		return nil, nil
	}
	intermediate := make(map[string]string)
	for k, v := range *h {
		intermediate[k] = strings.Join(v, ",")
	}
	return json.Marshal(intermediate)
}

func (h *Headers) UnmarshalYAML(value *yaml.Node) error {
	intermediate := make(map[string]string)
	err := value.Decode(&intermediate)
	*h = make(Headers, len(intermediate))
	if err == nil {
		for k, v := range intermediate {
			(*http.Header)(h).Set(k, v)
		}
	}
	return err
}

func (h *Headers) UnmarshalJSON(data []byte) error {
	intermediate := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	*h = make(Headers, len(intermediate))
	for key, value := range intermediate {
		strValue := string(value)
		if strings.HasPrefix(strValue, `"`) && strings.HasSuffix(strValue, `"`) {
			strValue, _ = strconv.Unquote(strValue)
		}
		(*http.Header)(h).Set(key, strValue)
	}
	return nil
}
