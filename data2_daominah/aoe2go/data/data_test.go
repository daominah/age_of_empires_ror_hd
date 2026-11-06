package data

import (
	"encoding/json"
	"testing"
)

func TestEmbeddedJSON(t *testing.T) {
	for name, content := range map[string][]byte{
		"data.json":    AoE2TechTreeData,
		"strings.json": AoE2TechTreeStrings,
	} {
		// GIVEN the embedded file
		// WHEN it is parsed as a JSON object
		var parsed map[string]any
		err := json.Unmarshal(content, &parsed)

		// THEN it is valid and not empty
		if err != nil {
			t.Errorf("error json.Unmarshal %v: %v", name, err)
		}
		if len(parsed) == 0 {
			t.Errorf("%v: empty JSON object", name)
		}
	}
}
