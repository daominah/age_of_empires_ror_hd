package logic

import "strings"

// Trait is something that matters in a fight but the rankings do not calculate:
// regeneration, population cost, bonus damage resistance and trample damage.
// The text is written by hand from the wiki (see wiki_reference.md in each ranking command),
// as a note of the row, at the values of the Imperial Age.
type Trait struct {
	CivilizationID CivilizationID
	// UnitIDs are the units that have the trait, empty for every unit of the civilization.
	UnitIDs []UnitID
	Text    string
}

var (
	scoutLineUnitIDs         = []UnitID{ScoutCavalry, LightCavalry, Hussar, WingedHussar}
	cavalryArcherLineUnitIDs = []UnitID{CavalryArcher, HeavyCavalryArcher}
	kipchakUnitIDs           = []UnitID{Kipchak, EliteKipchak}
)

// Traits are the hand written traits of the civilizations.
var Traits = []Trait{
	// Georgians: the unique technology Aznauri Cavalry gives regeneration and a lower population cost
	{Georgians, append(append([]UnitID{}, scoutLineUnitIDs...), HeavyCavalryArcher, CavalryArcher, Kipchak, EliteKipchak),
		"regeneration 14 per minute after Aznauri Cavalry"},
	{Georgians, append(append([]UnitID{}, scoutLineUnitIDs...), CavalryArcher, Kipchak, EliteKipchak),
		"population cost -20% after Aznauri Cavalry"},
	{Georgians, []UnitID{HeavyCavalryArcher}, "population cost -15% after Aznauri Cavalry"},
	// Georgians damage resistance applies on higher elevation only,
	// multiplied with the usual higher elevation resistance (x0.75): 0.8 * 0.75 = 0.6, so -40% instead of -25%
	{Georgians, nil, "receives less damage from lower elevation, -40% instead of the normal -25%"},
	// Khitans Ordo Cavalry regenerates in combat only
	{Khitans, []UnitID{ScoutCavalry}, "regeneration in combat 67.5 per minute after Ordo Cavalry"},
	{Khitans, []UnitID{LightCavalry}, "regeneration in combat 90 per minute after Ordo Cavalry"},
	{Khitans, []UnitID{Hussar}, "regeneration in combat 118.42 per minute after Ordo Cavalry"},
	{Berbers, []UnitID{CamelArcher, EliteCamelArcher}, "regeneration 15 per minute after Maghrebi Camels"},
	{Sicilians, nil, "bonus damage taken -40%"},
	// Maya only has the skirmisher line in the rankings
	{Maya, nil, "+1 projectile after Hul'che Javelineers, it has attack 1 so it deals at least 1 damage, shown as +1 attack"},
	{Poles, scoutLineUnitIDs, "33% trample damage in 0.5 tile radius after Lechitic Legacy"},
}

// FormatTraitNote joins the texts of the traits of the row unit, empty if none.
func FormatTraitNote(row Row) string {
	var texts []string
	for _, trait := range Traits {
		if trait.CivilizationID != row.CivilizationID || !trait.CheckHasUnit(row.Unit.ID) {
			continue
		}
		texts = append(texts, trait.Text)
	}
	return strings.Join(texts, "; ")
}

// CheckHasUnit returns whether the trait applies to the unit.
func (trait Trait) CheckHasUnit(unitID UnitID) bool {
	if len(trait.UnitIDs) == 0 {
		return true
	}
	for _, id := range trait.UnitIDs {
		if id == unitID {
			return true
		}
	}
	return false
}
