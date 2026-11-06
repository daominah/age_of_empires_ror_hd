// Package logic is the shared logic for ranking AoE2 DE civilizations
// by a unit line, e.g. the scout line ranked by arrows to kill vs a Castle.
//
// The input is data/dat_original_to_json, extracted from the game files in data/dat_original_from_game:
// base unit stats and costs, and the units and technologies available to each civilization.
// Names are from the game English strings.
// Technology and civilization bonus effects are in effect.go, maintained by hand,
// based on data/dat_original_to_json/effects.json.
// After a game patch, extract again and check the git diff of the JSON files.
//
// The ID constants (UnitID, TechID, CivilizationID, ArmorClass) are generated
// by cmd/generate-id: run go generate ./logic after extracting again.
//
// Each ranking is an executable in its own cmd subdirectory (e.g. cmd/ranking-scout),
// with its line units, attackers and ordered sort criteria as config at the top.
//
// The sibling module aoego is AoE1 Rise of Rome logic, not used here.
package logic

//go:generate go run ../cmd/generate-id
