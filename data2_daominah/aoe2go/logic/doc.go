// Package logic is the shared logic for ranking AoE2 DE civilizations
// by a unit line, e.g. the scout line ranked by arrows to kill vs a Castle.
//
// The input is the aoe2techtree JSON in
// data/data_from_aoe2techtree:
// base unit stats, costs, and the units and technologies available to each civilization.
// That JSON does not hold what a technology does,
// so technology and civilization bonus effects are in effect.go, maintained by hand,
// based on data/effects_from_dat/effects.json extracted from the game .dat.
// After a game patch, extract again and check the effects.json git diff.
//
// Each ranking is an executable in its own cmd subdirectory (e.g. cmd/ranking-scout),
// with its line units, attackers and ordered sort criteria as config at the top.
//
// The sibling module aoego is AoE1 Rise of Rome logic, not used here.
package logic
