package common

import "encoding/json"

// Convert convert source data to target data
func Convert(source interface{}, target interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, target)
}
