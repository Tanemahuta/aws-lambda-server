package config

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Read from filename.
func Read(ctx context.Context, filename string) (*Server, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return unmarshalConfig(ctx, data, path.Ext(filename))
}

func unmarshalConfig(ctx context.Context, data []byte, ext string) (*Server, error) {
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
	log := logr.FromContextOrDiscard(ctx)
	for idx, fn := range result.Functions {
		if fn.ARN != nil {
			log.Info("please migrate your config to the newer version and use 'name' instead of 'arn'",
				"functionIndex", idx)
		}
	}
	return &result, nil
}
