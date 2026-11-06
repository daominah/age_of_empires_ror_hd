package data

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedJSON(t *testing.T) {
	for name, content := range map[string][]byte{
		"civilizations.json": Civilizations,
		"units.json":         Units,
		"technologies.json":  Technologies,
	} {
		// GIVEN the embedded file
		// WHEN it is parsed as a JSON list
		var parsed []any
		err := json.Unmarshal(content, &parsed)

		// THEN it is valid and not empty
		if err != nil {
			t.Errorf("error json.Unmarshal %v: %v", name, err)
		}
		if len(parsed) == 0 {
			t.Errorf("%v: empty JSON list", name)
		}
	}
}
