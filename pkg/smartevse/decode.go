package smartevse

import (
	"encoding/json"
)

func FromJSON(data []byte) (*SmartEVSESettings, error) {
	settings := SmartEVSESettings{}

	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
