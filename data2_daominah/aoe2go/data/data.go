// Package data embeds the AoE2 DE data files,
// so programs in cmd work from any working directory.
// Go embed cannot reach a parent dir, so this package lives next to the files.
// Rebuild after refreshing the files (go run does it automatically).
package data

import _ "embed"

// Civilizations is dat_original_to_json/civilizations.json:
// the .dat civilizations, with the units, buildings and technologies each one can use.
//
//go:embed dat_original_to_json/civilizations.json
var Civilizations []byte

// Units is dat_original_to_json/units.json: base stats of units and buildings.
//
//go:embed dat_original_to_json/units.json
var Units []byte

// Technologies is dat_original_to_json/technologies.json: names, effects and costs.
//
//go:embed dat_original_to_json/technologies.json
var Technologies []byte

// GameStringsEnglish is the game English strings, lines like `5661 "Hussar"`,
// the names of units, technologies and civilizations by their name_string_id.
//
//go:embed dat_original_from_game/resources/en/strings/key-value/key-value-strings-utf8.txt
var GameStringsEnglish []byte
