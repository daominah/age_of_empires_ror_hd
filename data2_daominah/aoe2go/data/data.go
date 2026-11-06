// Package data embeds the AoE2 DE data files,
// so programs in cmd work from any working directory.
// Go embed cannot reach a parent dir, so this package lives next to the files.
// Rebuild after refreshing the files (go run does it automatically).
package data

import _ "embed"

// AoE2TechTreeData is data/data_from_aoe2techtree/data.json:
// base stats and costs of units, buildings and technologies,
// and what is available to each civilization.
//
//go:embed data_from_aoe2techtree/data.json
var AoE2TechTreeData []byte

// AoE2TechTreeStrings is data/data_from_aoe2techtree/strings.json:
// English names and descriptions, by string ID.
//
//go:embed data_from_aoe2techtree/strings.json
var AoE2TechTreeStrings []byte
