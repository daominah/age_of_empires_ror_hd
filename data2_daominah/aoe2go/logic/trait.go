package logic

import (
	"slices"
	"strings"
)

// Trait is something that matters in a fight but the rankings do not calculate:
// regeneration, population cost, bonus damage resistance, trample and pass through damage.
// The text is written by hand from the wiki (see wiki_reference.md in each ranking command),
// as a note of the row, at the values of the Imperial Age.
type Trait struct {
	CivilizationID CivilizationID
	// UnitIDs are the units that have the trait, empty for every unit of the civilization.
	UnitIDs []UnitID
	Text    string
	// Weight is a rough worth of the trait as a tie-breaker, 1 for a small upside,
	// 3 for a cost or training time saving, e.g. the Tupi refund of 15% of the cost.
	Weight float64
}

var (
	scoutLineUnitIDs         = []UnitID{ScoutCavalry, LightCavalry, Hussar, WingedHussar}
	cavalryArcherLineUnitIDs = []UnitID{CavalryArcher, HeavyCavalryArcher}
	kipchakUnitIDs           = []UnitID{Kipchak, EliteKipchak}
	skirmisherLineUnitIDs    = []UnitID{Skirmisher, EliteSkirmisher, ImperialSkirmisher}
)

// Traits are the handwritten traits of the civilizations.
var Traits = []Trait{
	// Georgians: the unique technology Aznauri Cavalry gives regeneration and a lower population cost
	{Georgians, append(append([]UnitID{}, scoutLineUnitIDs...), HeavyCavalryArcher, CavalryArcher, Kipchak, EliteKipchak),
		"regeneration 14 per minute after Aznauri Cavalry", 1},
	{Georgians, append(append([]UnitID{}, scoutLineUnitIDs...), CavalryArcher, Kipchak, EliteKipchak),
		"population cost -20% after Aznauri Cavalry", 1},
	{Georgians, []UnitID{HeavyCavalryArcher}, "population cost -15% after Aznauri Cavalry", 1},
	// Georgians damage resistance applies on higher elevation only,
	// multiplied with the usual higher elevation resistance (x0.75): 0.8 * 0.75 = 0.6, so -40% instead of -25%
	{Georgians, nil, "receives less damage from lower elevation, -40% instead of the normal -25%", 1},
	// Khitans Ordo Cavalry regenerates in combat only
	{Khitans, []UnitID{ScoutCavalry}, "regeneration in combat 67.5 per minute after Ordo Cavalry", 1},
	{Khitans, []UnitID{LightCavalry}, "regeneration in combat 90 per minute after Ordo Cavalry", 1},
	{Khitans, []UnitID{Hussar}, "regeneration in combat 118.42 per minute after Ordo Cavalry", 1},
	{Berbers, []UnitID{CamelArcher, EliteCamelArcher}, "regeneration 15 per minute after Maghrebi Camels", 1},
	// a free civilization bonus that counters the bonus attack of the units built to kill it,
	// already in the damage (Unit.BonusDamageResistance), so it is a note only, of weight 0
	{Sicilians, nil, "bonus damage taken -40%", 0},
	// Maya only has the skirmisher line in the rankings
	{Maya, nil, "+1 arrow after Hul'che Javelineers, it has attack 1 so it deals at least 1 damage", 1},
	{Mapuche, []UnitID{Skirmisher, EliteSkirmisher, ImperialSkirmisher, BolasRider, EliteBolasRider},
		"pass through damage after Malon", 1},
	{Poles, scoutLineUnitIDs, "33% trample damage in 0.5 tile radius after Lechitic Legacy", 1},
	// the line of sight is not a unit stat in the rankings, and the Mongols +2 is the only civilization bonus to it on the scout line
	{Mongols, scoutLineUnitIDs, "+2 line of sight", 1},
	// the scout line upsides of the extra scores of the scout ranking, so the note explains them
	{Poles, scoutLineUnitIDs, "Folwark collects instantly 10% of the farm food", 1},
	{Sicilians, scoutLineUnitIDs, "farm upgrades +100% food", 1},
	{Khitans, scoutLineUnitIDs, "trained 15% faster", 1},
	{Khitans, scoutLineUnitIDs, "Pasture food needs no drop-off", 1},
	{Franks, scoutLineUnitIDs, "farm upgrades free, trained 40% faster after Chivalry", 1},
	{Franks, scoutLineUnitIDs, "Scout Cavalry 54 HP in Feudal Age, above the others until they research Bloodlines", 1},
	{Gurjaras, scoutLineUnitIDs, "livestock garrisoned in the Mills generates food", 1},
	{Huns, scoutLineUnitIDs, "trained 20% faster at the Stable, no houses", 1},
	{Cumans, scoutLineUnitIDs, "trained 100% faster after Steppe Husbandry", 1},
	{Slavs, scoutLineUnitIDs, "farmers work faster", 1},
	{Khmer, scoutLineUnitIDs, "farmers need no drop-off", 1},
	{Vikings, scoutLineUnitIDs, "Wheelbarrow and Hand Cart free", 1},
	{Hindustanis, scoutLineUnitIDs, "villagers cheaper", 1},
	{Tatars, scoutLineUnitIDs, "+50% damage from higher elevation instead of the normal +25%", 1},
	{Danes, scoutLineUnitIDs, "gets 25% of the cost of the buildings it destroys", 1},
	{Vietnamese, scoutLineUnitIDs, "Conscription free", 1},
	{Wei, scoutLineUnitIDs, "generates 1.8 food per minute after Tuntian", 1},
	// the small skirmisher upsides from the wiki, so they break the ties of the standard Elite Skirmishers,
	// on the skirmisher line only, as the scout and cavalry archer rankings were not checked for them.
	// A cost or training time saving weighs most, a unique technology or a situational upside least.
	// The hero auras (Liu Bei, Cao Cao, Sun Jian) are left out: they only work near one unit.
	{Tupi, skirmisherLineUnitIDs, "returns 15% of its cost when killed", 3},
	{Saxons, skirmisherLineUnitIDs, "cost -5% per Town Center or Castle, up to -20%", 3},
	{Khitans, skirmisherLineUnitIDs, "trained 15% faster, upgrades researched 15% faster", 3},
	{Aztecs, skirmisherLineUnitIDs, "trained 15% faster", 3},
	{Britons, skirmisherLineUnitIDs, "trained 10% faster at the Archery Range", 2},
	{Danes, skirmisherLineUnitIDs, "gets 25% of the cost of the buildings it destroys", 2},
	{Vietnamese, skirmisherLineUnitIDs, "Conscription free", 1},
	{Khitans, skirmisherLineUnitIDs, "reflects 25% of the melee damage after Lamellar Armor", 1},
	{Tatars, skirmisherLineUnitIDs, "+50% damage from higher elevation instead of the normal +25%", 1},
	{Wei, skirmisherLineUnitIDs, "generates 1.8 food per minute after Tuntian", 1},
	{Jurchens, skirmisherLineUnitIDs, "friendly fire damage taken -50%", 1},
	{Portuguese, skirmisherLineUnitIDs, "upgrades researched 25% faster", 1},
	{Bohemians, skirmisherLineUnitIDs, "Chemistry available in the Castle Age", 1},
	// the Chinese technology discount is missing from the technology cost note
	{Chinese, skirmisherLineUnitIDs, "technologies 15% cheaper in the Imperial Age", 1},
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
	return slices.Contains(trait.UnitIDs, unitID)
}
