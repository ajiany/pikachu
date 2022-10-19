package feishu

import "encoding/json"

func Convert(v interface{}, data interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, data)
}
