package aws

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var _ json.Unmarshaler = &Headers{}
var _ yaml.Unmarshaler = &Headers{}

type Headers struct {
	http.Header
}

func (h *Headers) UnmarshalYAML(value *yaml.Node) error {
	intermediate := make(map[string]string)
	err := value.Decode(&intermediate)
	h.Header = make(http.Header, len(intermediate))
	if err == nil {
		for k, v := range intermediate {
			h.Header.Set(k, v)
		}
	}
	return err
}

func (h *Headers) UnmarshalJSON(data []byte) error {
	intermediate := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &intermediate); err != nil {
		return err
	}
	h.Header = make(http.Header, len(intermediate))
	for key, value := range intermediate {
		strValue := string(value)
		if strings.HasPrefix(strValue, `"`) && strings.HasSuffix(strValue, `"`) {
			strValue, _ = strconv.Unquote(strValue)
		}
		h.Header.Set(key, strValue)
	}
	return nil
}
