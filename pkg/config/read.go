package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Read from filename.
func Read(filename string) (*Server, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(data, path.Ext(filename))
}

func unmarshalConfig(data []byte, ext string) (*Server, error) {
	var (
		result Server
		err    error
	)
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &result)
	case ".json":
		err = json.Unmarshal(data, &result)
	default:
		err = errors.Errorf("invalid config file type: %v", ext)
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}
