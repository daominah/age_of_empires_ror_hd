package logic

import "sort"

// LoadUnitBase is the unit before any technology or civilization bonus,
// with the automatic technologies, e.g. the Scout Cavalry attack 5 from the Feudal Age.
// It computes from the parsed tree, nothing is read from a file.
// The unit as in the data, without the automatic technologies, is tree.Units[id].
func (tree *TechTree) LoadUnitBase(id UnitID) Unit {
	unit := tree.Units[id].clone()
	for _, techID := range listSortedTechnologyIDs(func(techID TechID) bool { return AutomaticTechnologies[techID] }) {
		TechnologyEffects[techID].Apply(&unit)
	}
	return unit
}

// LoadUnitUpgraded is the LoadUnitBase with the technologies any civilization can research,
// without civilization bonuses or unique technologies,
// e.g. the Arbalester: attack 6+4, 40 HP, pierce armor 0+4.
// A unit only one civilization trains gets what that civilization has instead,
// its unique technologies and civilization bonuses included,
// e.g. the Elite Rattan Archer: 48 HP with the Vietnamese archer HP +20%.
// It stands for the enemy of a ranked line, as a target or an attacker.
// It only gets the technologies in TechnologyEffects, the ones the rankings need so far.
func (tree *TechTree) LoadUnitUpgraded(id UnitID) Unit {
	if civilization, isOnlyOne := tree.checkIsUniqueUnit(id); isOnlyOne {
		return tree.loadRowUpgraded(civilization, id, tree.Units[id].Class).Unit
	}
	unit := tree.LoadUnitBase(id)
	checkIsCommonUpgrade := func(techID TechID) bool {
		return !AutomaticTechnologies[techID] && !TechnologyEffects[techID].IsUniqueTechnology &&
			tree.Technologies[techID].CivilizationID == AnyCivilization
	}
	for _, techID := range listSortedTechnologyIDs(checkIsCommonUpgrade) {
		TechnologyEffects[techID].Apply(&unit)
	}
	return unit
}

// LoadAttackerUpgraded is the LoadUnitUpgraded as an attacker,
// named with its attack as base+upgrade, e.g. "Arbalester 6+4".
func (tree *TechTree) LoadAttackerUpgraded(id UnitID) Attacker {
	unit := tree.LoadUnitUpgraded(id)
	base := getAttack(tree.LoadUnitBase(id))
	return Attacker{
		Name:    unit.Name + " " + formatNumber(base) + formatSigned(getAttack(unit)-base),
		Attacks: unit.Attacks,
	}
}

// checkIsUniqueUnit checks whether one civilization alone trains the unit, and returns it,
// e.g. the Vietnamese for the Rattan Archer.
// It is about the same as a "UniqueUnit" node in the game Technology Tree,
// which also includes units like the Vietnamese Imperial Skirmisher,
// the Hindustanis Imperial Camel Rider or the Spanish Missionary.
// A regional unit does not count, as several civilizations train it.
func (tree *TechTree) checkIsUniqueUnit(id UnitID) (Civilization, bool) {
	var result []Civilization
	for _, civilization := range tree.Civilizations {
		if civilization.UnitIDs[id] {
			result = append(result, civilization)
		}
	}
	if len(result) != 1 {
		return Civilization{}, false
	}
	return result[0], true
}

// listSortedTechnologyIDs are the TechnologyEffects technologies kept by the filter, by ID,
// a fixed order, e.g. Bloodlines before the Mongols HP bonus.
func listSortedTechnologyIDs(checkIsKept func(techID TechID) bool) []TechID {
	var ids []TechID
	for techID := range TechnologyEffects {
		if checkIsKept(techID) {
			ids = append(ids, techID)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids
}
