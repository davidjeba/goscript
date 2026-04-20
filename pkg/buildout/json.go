package buildout

import "encoding/json"

func jsonMarshal(value interface{}) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}

