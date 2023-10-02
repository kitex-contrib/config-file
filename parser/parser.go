package parser

import (
	"sigs.k8s.io/yaml"
)

// DecodeServer parse the config file
func Decode(data []byte, resp interface{}) error {
	return yaml.Unmarshal(data, resp)
}
