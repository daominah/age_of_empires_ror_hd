package logic

import (
	"maps"
	"slices"
)

// Unit is the stats of a unit, base or after applying effects.
type Unit struct {
	ID    UnitID
	Name  string    // English name, e.g. "Hussar"
	Class UnitClass // unit class in the .dat, e.g. ClassCavalry
	HP    float64
	// Attacks and Armors are amounts by armor class, e.g. ArmorClassPierce.
	// Damage only counts the classes in both the attacker Attacks and the target Armors.
	// A class listed with 0 is not the same as a missing one, e.g. the Chu Ko Nu melee attack 0
	// deals 3 to a Battering Ram, which has melee armor -3,
	// and a technology only improves a class the unit already lists, e.g. the Champion lists 0 vs cavalry for Chieftains.
	Attacks  map[ArmorClass]float64
	Armors   map[ArmorClass]float64
	Range    float64 // max range in tiles, e.g. 0 for a Knight, 1 for a Steppe Lancer, 5 for a Throwing Axeman
	MinRange float64 // min range in tiles, 0 for most units, e.g. 1 for a Skirmisher, 3 for a Mangonel
	// ExtraProjectiles is the number of arrows besides the main one, fired with every attack,
	// e.g. 3 for the Elite Kipchak, 1 for the Maya skirmisher.
	ExtraProjectiles int
	// ExtraProjectileID is the projectile unit of the extra arrows,
	// the target of the effects that improve them, e.g. Fletching for the Xianbei Raider.
	ExtraProjectileID UnitID
	// ExtraAttacks is the attack of each extra arrow, usually weaker than the main one
	// and not improved by the Blacksmith, e.g. 3 pierce for the Kipchak.
	// Nil means the extra arrows are copies of the main one.
	ExtraAttacks map[ArmorClass]float64
	// ChargeAttack is a strong attack that fires instead of a normal one when charged,
	// e.g. the Xianbei Raider volley, nil for most units.
	ChargeAttack *ChargeAttack
	ReloadTime   float64 // seconds between attacks
	Speed        float64
	TrainTime    float64 // seconds
	Cost         Cost
	// IsIgnoringArmor is set by Wootz Steel: the unit attack ignores the target armor.
	IsIgnoringArmor bool
}

// Shot is what one attack fires at the target: a main arrow, then the extra ones.
type Shot struct {
	// ProjectileID is the projectile unit of the main arrow when it has its own attacks,
	// e.g. the Xianbei Raider charged arrow, 0 for a normal attack, which uses the unit attacks.
	ProjectileID      UnitID
	Attacks           map[ArmorClass]float64 // the main arrow, e.g. 5 pierce for the Xianbei Raider charged shot
	ExtraProjectiles  int                    // e.g. 5 for the Xianbei Raider charged shot
	ExtraProjectileID UnitID                 // see Unit.ExtraProjectileID
	ExtraAttacks      map[ArmorClass]float64 // each extra arrow, e.g. 1 pierce for the Xianbei Raider charged shot
}

// ChargeAttack is an attack that fires once charged, instead of a normal attack,
// e.g. the Bolas Rider throws a strong bolas every 30 seconds.
type ChargeAttack struct {
	Shot
	ChargeTime float64 // seconds to charge again after firing, e.g. 30
}

// GetShot is the normal attack of the unit, the one fired every ReloadTime.
func (u Unit) GetShot() Shot {
	return Shot{
		Attacks:           u.Attacks,
		ExtraProjectiles:  u.ExtraProjectiles,
		ExtraProjectileID: u.ExtraProjectileID,
		ExtraAttacks:      u.ExtraAttacks,
	}
}

// listProjectileAttacks are the arrows that have their own attacks, with the projectile unit of each,
// so an effect on a projectile unit changes them, e.g. Parthian Tactics on the Xianbei Raider arrows.
func (u *Unit) listProjectileAttacks() []projectileAttacks {
	var result []projectileAttacks
	if u.ExtraAttacks != nil {
		result = append(result, projectileAttacks{u.ExtraProjectileID, u.ExtraAttacks})
	}
	if u.ChargeAttack != nil {
		result = append(result, projectileAttacks{u.ChargeAttack.ProjectileID, u.ChargeAttack.Attacks})
		if u.ChargeAttack.ExtraAttacks != nil {
			result = append(result, projectileAttacks{u.ChargeAttack.ExtraProjectileID, u.ChargeAttack.ExtraAttacks})
		}
	}
	return result
}

// projectileAttacks is the attacks of a projectile unit in a shot of the unit.
type projectileAttacks struct {
	ProjectileID UnitID
	Attacks      map[ArmorClass]float64
}

// clone copies the unit, so applying effects does not change the base unit maps.
func (u Unit) clone() Unit {
	c := u
	c.Attacks = maps.Clone(u.Attacks)
	c.Armors = maps.Clone(u.Armors)
	c.ExtraAttacks = maps.Clone(u.ExtraAttacks)
	if u.ChargeAttack != nil {
		charge := *u.ChargeAttack
		charge.Attacks = maps.Clone(charge.Attacks)
		charge.ExtraAttacks = maps.Clone(charge.ExtraAttacks)
		c.ChargeAttack = &charge
	}
	return c
}

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
	slices.Sort(ids)
	return ids
}
