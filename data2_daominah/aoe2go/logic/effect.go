package logic

import "math"

// Effect is an effect in the game .dat, applied by a technology or a civilization bonus.
// The values are copied by hand from data/dat_original_to_json/effects.json,
// keeping only the commands for the unit lines ranked so far (the scout line),
// so a game patch shows in the effects.json git diff and is edited here.
type Effect struct {
	EffectID int    // effect ID in effects.json
	Name     string // technology name in the .dat, to search effects.json
	Commands []EffectCommand
	// IsUniqueTechnology is a civilization unique technology, researched in a Castle.
	// It is expensive, so a ranking may count its cost saving less than a free civilization bonus.
	IsUniqueTechnology bool
}

// EffectCommand is one command of an effect,
// in the same order as effects.json: [Type, A, B, C, D].
// For the attribute commands, e.g. {AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)}
// is "cavalry +1 melee armor".
type EffectCommand struct {
	Type      EffectType
	UnitID    UnitID    // A, AnyUnit for all units of the Class
	Class     UnitClass // B, AnyClass for any
	Attribute Attribute // C
	// Value is D, for AttributeArmor and AttributeAttack it is a classAmount.
	Value float64
}

// EffectType is the command type in the .dat.
type EffectType int

// EffectType values handled here, others are ignored.
// The team versions apply to the civilization itself too.
const (
	SetAttribute          EffectType = 0
	AddAttribute          EffectType = 4
	MultiplyAttribute     EffectType = 5
	TeamSetAttribute      EffectType = 10
	TeamAddAttribute      EffectType = 14
	TeamMultiplyAttribute EffectType = 15
)

// Attribute is a unit attribute changed by an effect command.
type Attribute int

// Attribute values in the effects so far,
// the ones not in Apply are kept for reading but ignored by the rankings.
const (
	AttributeHP                    Attribute = 0
	AttributeLineOfSight           Attribute = 1
	AttributeSpeed                 Attribute = 5
	AttributeArmor                 Attribute = 8 // Value is a classAmount
	AttributeAttack                Attribute = 9 // Value is a classAmount
	AttributeReloadTime            Attribute = 10
	AttributeBlastWidth            Attribute = 22 // e.g. Lechitic Legacy trample damage
	AttributeSearchRadius          Attribute = 23
	AttributeBonusDamageResistance Attribute = 24
	AttributeIgnoreArmor           Attribute = 63
	AttributeAllCosts              Attribute = 100
	AttributeTrainTime             Attribute = 101
	AttributeFoodCost              Attribute = 103
	AttributeWoodCost              Attribute = 104
	AttributeGoldCost              Attribute = 105
	AttributeStoneCost             Attribute = 106
	AttributeRegeneration          Attribute = 109 // HP per minute
	AttributePopulation            Attribute = 110 // e.g. Aznauri Cavalry
	AttributeFriendlyFireDamage    Attribute = 119
)

// UnitClass is a unit class in the .dat.
type UnitClass int

// UnitClass values used so far.
const (
	AnyClass     UnitClass = -1
	ClassCavalry UnitClass = 12 // the scout line
)

// AnyUnit in an effect command targets all units of the command Class.
const AnyUnit UnitID = -1

// Apply applies the effect commands that target the unit.
func (e Effect) Apply(unit *Unit) {
	for _, c := range e.Commands {
		c.Apply(unit)
	}
}

// Apply applies the command if it targets the unit.
func (c EffectCommand) Apply(unit *Unit) {
	if c.UnitID != AnyUnit && c.UnitID != unit.ID {
		return
	}
	if c.Class != AnyClass && c.Class != unit.Class {
		return
	}
	var operation func(old float64, value float64) float64
	switch c.Type {
	case SetAttribute, TeamSetAttribute:
		operation = func(_ float64, value float64) float64 { return value }
	case AddAttribute, TeamAddAttribute:
		operation = func(old float64, value float64) float64 { return old + value }
	case MultiplyAttribute, TeamMultiplyAttribute:
		operation = func(old float64, value float64) float64 { return old * value }
	default:
		return
	}

	switch c.Attribute {
	case AttributeHP:
		unit.HP = operation(unit.HP, c.Value)
	case AttributeSpeed:
		unit.Speed = operation(unit.Speed, c.Value)
	case AttributeReloadTime:
		unit.ReloadTime = operation(unit.ReloadTime, c.Value)
	case AttributeTrainTime:
		unit.TrainTime = operation(unit.TrainTime, c.Value)
	case AttributeIgnoreArmor:
		unit.IsIgnoringArmor = operation(0, c.Value) > 0
	case AttributeFoodCost:
		unit.Cost.Food = operation(unit.Cost.Food, c.Value)
	case AttributeWoodCost:
		unit.Cost.Wood = operation(unit.Cost.Wood, c.Value)
	case AttributeGoldCost:
		unit.Cost.Gold = operation(unit.Cost.Gold, c.Value)
	case AttributeStoneCost:
		unit.Cost.Stone = operation(unit.Cost.Stone, c.Value)
	case AttributeAllCosts:
		unit.Cost.Food = operation(unit.Cost.Food, c.Value)
		unit.Cost.Wood = operation(unit.Cost.Wood, c.Value)
		unit.Cost.Gold = operation(unit.Cost.Gold, c.Value)
		unit.Cost.Stone = operation(unit.Cost.Stone, c.Value)
	case AttributeArmor, AttributeAttack:
		amounts := unit.Armors
		if c.Attribute == AttributeAttack {
			amounts = unit.Attacks
		}
		class, value := decodeClassAmount(c.Value)
		if c.Type == MultiplyAttribute || c.Type == TeamMultiplyAttribute {
			value = value / 100 // a percentage for multiply
		}
		// like the game, a class the unit does not have is not added
		if old, isFound := amounts[class]; isFound {
			amounts[class] = operation(old, value)
		}
	}
}

// classAmount is the Value of an armor or attack command: class * 256 + amount,
// e.g. classAmount(ArmorClassPierce, 1) is 769, classAmount(ArmorClassMelee, -1) is -1025.
// For MultiplyAttribute the amount is a percentage, e.g. 120 is * 1.2.
func classAmount(class ArmorClass, amount float64) float64 {
	sign := 1.0
	if amount < 0 {
		sign, amount = -1, -amount
	}
	return sign * (float64(class)*256 + amount)
}

// decodeClassAmount is the reverse of classAmount.
func decodeClassAmount(value float64) (class ArmorClass, amount float64) {
	v := int(math.Round(value))
	sign := 1
	if v < 0 {
		sign, v = -1, -v
	}
	return ArmorClass(v / 256), float64(sign * (v % 256))
}

// AutomaticTechnologies are researched by the game for every civilization, outside the Technology Trees,
// e.g. on reaching the Feudal Age.
// A ranking applies them to the base unit.
var AutomaticTechnologies = map[TechID]bool{AutoUpgradeScoutFeudalAge: true}

// TechnologyEffects are the effects of the technologies, by technology ID.
// A ranking applies each one available to the civilization.
var TechnologyEffects = map[TechID]Effect{
	Husbandry: {EffectID: 39, Name: "Husbandry", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeSpeed, 1.1},
	}},
	Forging: {EffectID: 67, Name: "Forging", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 1)},
	}},
	IronCasting: {EffectID: 68, Name: "Iron casting", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 1)},
	}},
	BlastFurnace: {EffectID: 75, Name: "Blast Furnace", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 2)},
	}},
	PlateBardingArmor: {EffectID: 80, Name: "Plate Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassPierce, 2)},
	}},
	ScaleBardingArmor: {EffectID: 81, Name: "Scale Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassPierce, 1)},
	}},
	ChainBardingArmor: {EffectID: 82, Name: "Chain Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassPierce, 1)},
	}},
	AutoUpgradeScoutFeudalAge: {EffectID: 16, Name: "Auto upgrade Scout Feudal Age", Commands: []EffectCommand{
		{AddAttribute, ScoutCavalry, AnyClass, AttributeSpeed, 0.35},
		{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, classAmount(ArmorClassMelee, 2)},
	}},
	Bloodlines: {EffectID: 450, Name: "Bloodlines", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeHP, 20},
	}},
	Farimba: {EffectID: 606, Name: "Malian UT", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 5)},
	}, IsUniqueTechnology: true},
	ManipurCavalry: {EffectID: 667, Name: "Burmese UT", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassArchers, 4)},
	}, IsUniqueTechnology: true},
	Stirrups: {EffectID: 722, Name: "Khmer UT", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeReloadTime, 0.75},
	}, IsUniqueTechnology: true},
	SilkArmor: {EffectID: 724, Name: "Malay UT", Commands: []EffectCommand{
		{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, LightCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, Hussar, AnyClass, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, WingedHussar, AnyClass, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
		{AddAttribute, LightCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
		{AddAttribute, Hussar, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
		{AddAttribute, WingedHussar, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
	}, IsUniqueTechnology: true},
	SteppeHusbandry: {EffectID: 726, Name: "Burmese UT", Commands: []EffectCommand{
		{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeTrainTime, 0.5},
		{MultiplyAttribute, LightCavalry, AnyClass, AttributeTrainTime, 0.5},
		{MultiplyAttribute, Hussar, AnyClass, AttributeTrainTime, 0.5},
		{MultiplyAttribute, WingedHussar, AnyClass, AttributeTrainTime, 0.5},
	}, IsUniqueTechnology: true},
	LechiticLegacy: {EffectID: 810, Name: "Lechitic Legacy", Commands: []EffectCommand{
		{AddAttribute, ScoutCavalry, AnyClass, AttributeBlastWidth, 0.5},
		{AddAttribute, LightCavalry, AnyClass, AttributeBlastWidth, 0.5},
		{AddAttribute, Hussar, AnyClass, AttributeBlastWidth, 0.5},
		{AddAttribute, WingedHussar, AnyClass, AttributeBlastWidth, 0.5},
	}, IsUniqueTechnology: true},
	WootzSteel: {EffectID: 851, Name: "Wootz Steel", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeIgnoreArmor, 1},
	}, IsUniqueTechnology: true},
	CorvinianArmy: {EffectID: 571, Name: "Corvinian Army", Commands: []EffectCommand{
		{SetAttribute, MagyarHuszar, AnyClass, AttributeGoldCost, 0},
		{SetAttribute, EliteMagyarHuszar, AnyClass, AttributeGoldCost, 0},
		{SetAttribute, MagyarHuszar, AnyClass, AttributeFoodCost, 80},
		{SetAttribute, EliteMagyarHuszar, AnyClass, AttributeFoodCost, 80},
	}, IsUniqueTechnology: true},
	Kshatriyas: {EffectID: 854, Name: "Kshatriyas", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeFoodCost, 0.75},
	}, IsUniqueTechnology: true},
	AznauriCavalry: {EffectID: 936, Name: "Aznauri Cavalry", Commands: []EffectCommand{
		{SetAttribute, AnyUnit, ClassCavalry, AttributePopulation, -0.8},
	}, IsUniqueTechnology: true},
	MingGuangArmor: {EffectID: 1062, Name: "Ming-Kuang Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 4)},
	}, IsUniqueTechnology: true},
}

// CivilizationBonusEffects are the civilization bonuses and team bonuses,
// active in the Imperial Age after researching every available technology,
// by civilization.
// A ranking applies them after the technologies, in this order.
var CivilizationBonusEffects = map[CivilizationID][]Effect{
	Bengalis: {
		{EffectID: 888, Name: "C-Bonus, Cavalry vs Skirm", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSkirmishers, 2)},
		}},
	},
	Berbers: {
		{EffectID: 610, Name: "C-Bonus, Stable units -15% che", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAllCosts, 0.85},
		}},
		{EffectID: 638, Name: "C-Bonus, Stable units -20% che", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAllCosts, 0.9411764},
		}},
	},
	Cumans: {
		{EffectID: 748, Name: "C-Bonus, Cavalry +5% speed Age2", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeSpeed, 1.05},
		}},
		{EffectID: 762, Name: "C-Bonus, Cavalry +5% speed Age3", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeSpeed, 1.05},
		}},
		{EffectID: 763, Name: "C-Bonus, Cavalry +5% speed Age4", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeSpeed, 1.05},
		}},
	},
	Franks: {
		{EffectID: 285, Name: "C-Bonus, Cavalry +20% HP", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeHP, 1.2},
		}},
	},
	Georgians: {
		{EffectID: 954, Name: "C-Bonus, Cavalry regeneration age2", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeRegeneration, 2},
		}},
		{EffectID: 961, Name: "C-Bonus, Cavalry regeneration age3", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeRegeneration, 6},
		}},
		{EffectID: 961, Name: "C-Bonus, Cavalry regeneration age4", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeRegeneration, 6},
		}},
	},
	Gurjaras: {
		{EffectID: 869, Name: "C-Bonus, +20% Cav bonus damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassInfantry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCapitalShips, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassElephantUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor6, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor7, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor9, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor10, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassAllBuildings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor12, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStoneDefense, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassPredatorAnimals, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassArchers, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShips, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassRams, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassTrees, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUniqueUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSiegeWeapons, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStandardBuildings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassWallsAndGates, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassGunpowderUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassBoars, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMonks, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCastles, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSpearmen, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalryArchers, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShockInfantry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCamelUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor33, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMamelukes, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeroesKings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeavySiege, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSkirmishers, 120)},
		}},
		{EffectID: 873, Name: "C-Bonus, +30% Cav bonus damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassInfantry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCapitalShips, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassElephantUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor6, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor7, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor9, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor10, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassAllBuildings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor12, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStoneDefense, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassPredatorAnimals, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassArchers, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShips, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassRams, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassTrees, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUniqueUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSiegeWeapons, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStandardBuildings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassWallsAndGates, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassGunpowderUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassBoars, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMonks, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCastles, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSpearmen, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalryArchers, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShockInfantry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCamelUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor33, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMamelukes, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeroesKings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeavySiege, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSkirmishers, 111)},
		}},
		{EffectID: 884, Name: "C-Bonus, +40% Cav bonus damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassInfantry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCapitalShips, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassElephantUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor6, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor7, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor9, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor10, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassAllBuildings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor12, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStoneDefense, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassPredatorAnimals, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassArchers, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShips, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassRams, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassTrees, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUniqueUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSiegeWeapons, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassStandardBuildings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassWallsAndGates, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassGunpowderUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassBoars, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMonks, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCastles, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSpearmen, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCavalryArchers, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassShockInfantry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassCamelUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassUnusedArmor33, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMamelukes, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeroesKings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassHeavySiege, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSkirmishers, 105)},
		}},
	},
	Hindustanis: {
		{EffectID: 2, Name: "Hindustanis Team Bonus", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, classAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, LightCavalry, AnyClass, AttributeAttack, classAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, Hussar, AnyClass, AttributeAttack, classAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, WingedHussar, AnyClass, AttributeAttack, classAmount(ArmorClassStandardBuildings, 2)},
		}},
	},
	Jurchens: {
		{EffectID: 994, Name: "C-Bonus, cav & Fire Lancer +25% attack", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeReloadTime, 0.8},
		}},
		{EffectID: 1009, Name: "C-Bonus, No friendly damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeFriendlyFireDamage, 0.5},
		}},
	},
	Khitans: {
		{EffectID: 67, Name: "C-Bonus, Double Forging", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 68, Name: "C-Bonus, Double Iron Casting", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 1011, Name: "C-Bonus, Skirms Spear Scout train 25% faster", Commands: []EffectCommand{
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeTrainTime, 0.869565},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeTrainTime, 0.869565},
			{MultiplyAttribute, Hussar, AnyClass, AttributeTrainTime, 0.869565},
		}},
	},
	Magyars: {
		{EffectID: 524, Name: "C-Bonus Lcav Cheaper", Commands: []EffectCommand{
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeFoodCost, 0.85},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeFoodCost, 0.85},
			{MultiplyAttribute, Hussar, AnyClass, AttributeFoodCost, 0.85},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeFoodCost, 0.85},
		}},
	},
	Mongols: {
		// The +20% (Castle Age) and +30% (Imperial Age) HP each come in 2 technologies
		// that disable each other (command type 102), the first to fire wins:
		// 286 and 287 ("+ BL") need Bloodlines and multiply the HP without its +20,
		// 288 and 388 multiply the HP before Bloodlines.
		// The .dat age names are shifted: "Middle Age" is the Feudal Age, "Feudal Age" is the Castle Age.
		// Bloodlines is a Feudal Age technology, usually researched before the Castle Age,
		// so 286 then 287 fire: (95 - 20) * 1.2 * 1.0833 + 20 = 117.5 HP, shown as 118 in Imperial Age.
		// (Bloodlines after the Castle Age gives 288 then 388: 75 * 1.2 * 1.084 + 20 = 117.56, also 118.)
		{EffectID: 286, Name: "C-Bonus, Light Cavalry +20% HP + BL", Commands: []EffectCommand{
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 100},
			{AddAttribute, ScoutCavalry, AnyClass, AttributeHP, -2000},
			{AddAttribute, LightCavalry, AnyClass, AttributeHP, -2000},
			{AddAttribute, Hussar, AnyClass, AttributeHP, -2000},
			{AddAttribute, WingedHussar, AnyClass, AttributeHP, -2000},
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 1.2},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 1.2},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 1.2},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 1.2},
			{AddAttribute, ScoutCavalry, AnyClass, AttributeHP, 2000},
			{AddAttribute, LightCavalry, AnyClass, AttributeHP, 2000},
			{AddAttribute, Hussar, AnyClass, AttributeHP, 2000},
			{AddAttribute, WingedHussar, AnyClass, AttributeHP, 2000},
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 0.01},
		}},
		{EffectID: 287, Name: "C-Bonus, Light Cavalry +30% HP + BL", Commands: []EffectCommand{
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 100},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 100},
			{AddAttribute, ScoutCavalry, AnyClass, AttributeHP, -2000},
			{AddAttribute, LightCavalry, AnyClass, AttributeHP, -2000},
			{AddAttribute, Hussar, AnyClass, AttributeHP, -2000},
			{AddAttribute, WingedHussar, AnyClass, AttributeHP, -2000},
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 1.08333},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 1.08333},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 1.08333},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 1.08333},
			{AddAttribute, ScoutCavalry, AnyClass, AttributeHP, 2000},
			{AddAttribute, LightCavalry, AnyClass, AttributeHP, 2000},
			{AddAttribute, Hussar, AnyClass, AttributeHP, 2000},
			{AddAttribute, WingedHussar, AnyClass, AttributeHP, 2000},
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, Hussar, AnyClass, AttributeHP, 0.01},
			{MultiplyAttribute, WingedHussar, AnyClass, AttributeHP, 0.01},
		}},
		{EffectID: 407, Name: "Mongols Team Bonus", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeLineOfSight, 2},
			{AddAttribute, LightCavalry, AnyClass, AttributeLineOfSight, 2},
			{AddAttribute, Hussar, AnyClass, AttributeLineOfSight, 2},
			{AddAttribute, WingedHussar, AnyClass, AttributeLineOfSight, 2},
			{AddAttribute, ScoutCavalry, AnyClass, AttributeSearchRadius, 2},
			{AddAttribute, LightCavalry, AnyClass, AttributeSearchRadius, 2},
			{AddAttribute, Hussar, AnyClass, AttributeSearchRadius, 2},
			{AddAttribute, WingedHussar, AnyClass, AttributeSearchRadius, 2},
		}},
	},
	Poles: {
		{EffectID: 802, Name: "Poles Team Bonus", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, classAmount(ArmorClassArchers, 1)},
			{AddAttribute, LightCavalry, AnyClass, AttributeAttack, classAmount(ArmorClassArchers, 1)},
			{AddAttribute, Hussar, AnyClass, AttributeAttack, classAmount(ArmorClassArchers, 1)},
			{AddAttribute, WingedHussar, AnyClass, AttributeAttack, classAmount(ArmorClassArchers, 1)},
		}},
	},
	Portuguese: {
		{EffectID: 33, Name: "C-Bonus, All units -20% gold", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeGoldCost, 0.8},
		}},
	},
	Sicilians: {
		{EffectID: 796, Name: "C-Bonus, Bonus damage resistance", Commands: []EffectCommand{
			{SetAttribute, AnyUnit, ClassCavalry, AttributeBonusDamageResistance, 0.4},
		}},
	},
	Teutons: {
		{EffectID: 333, Name: "C-Bonus, Inf Cav +1 armor Age3", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 334, Name: "C-Bonus, Inf Cav +1 armor Age4", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, classAmount(ArmorClassMelee, 1)},
		}},
	},
	Turks: {
		{EffectID: 509, Name: "C-Bonus, Light cavalry +1P armor", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
			{AddAttribute, LightCavalry, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
			{AddAttribute, Hussar, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
			{AddAttribute, WingedHussar, AnyClass, AttributeArmor, classAmount(ArmorClassPierce, 1)},
		}},
	},
	Varangians: {
		// the game describes it as "Bloodlines and Caravan effects +50%": +10 HP on top of +20,
		// the .dat name "+65%" is outdated
		{EffectID: 1489, Name: "C-Bonus, Bloodlines +65% more effective", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeHP, 10},
		}},
	},
	Wei: {
		{EffectID: 1027, Name: "Wei Team Bonus", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, classAmount(ArmorClassSiegeWeapons, 2)},
		}},
	},
}
