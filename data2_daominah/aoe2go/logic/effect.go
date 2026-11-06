package logic

import "math"

// Effect is an effect in the game .dat, applied by a technology or a civilization bonus.
// The values are copied by hand from data/dat_original_to_json/effects.json,
// keeping only the commands for the units ranked or loaded so far:
// cavalry, foot archers, cavalry archers and infantry,
// so a game patch shows in the effects.json git diff and is edited here.
type Effect struct {
	EffectID int    // effect ID in effects.json
	Name     string // technology name in the .dat, to search effects.json
	Commands []EffectCommand
	// TechnologyID is the technology that triggers a civilization bonus,
	// set when several of its technologies trigger the same effect, one per age,
	// e.g. Mapuche "Skirm Spear +5 HP" in Feudal, Castle and Imperial Age.
	// 0 is the first technology that triggers the effect.
	TechnologyID TechID
	// IsUniqueTechnology is a civilization unique technology, researched in a Castle.
	// It is expensive, so a ranking may count its cost saving less than a free civilization bonus.
	IsUniqueTechnology bool
}

// EffectCommand is one command of an effect,
// in the same order as effects.json: [Type, A, B, C, D].
// For the attribute commands, e.g. {AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)}
// is "cavalry +1 melee armor".
type EffectCommand struct {
	Type      EffectType
	UnitID    UnitID    // A, AnyUnit for all units of the Class
	Class     UnitClass // B, AnyClass for any
	Attribute Attribute // C
	// Value is D, for AttributeArmor and AttributeAttack it is a encodeClassAmount.
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
	AttributeArmor                 Attribute = 8 // Value is a encodeClassAmount
	AttributeAttack                Attribute = 9 // Value is a encodeClassAmount
	AttributeReloadTime            Attribute = 10
	AttributeRange                 Attribute = 12 // max range
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
	AnyClass           UnitClass = -1
	ClassArcher        UnitClass = 0  // foot archers, e.g. the skirmisher line
	ClassInfantry      UnitClass = 6  // e.g. Man-at-Arms, Throwing Axeman
	ClassCavalry       UnitClass = 12 // e.g. Knight, the scout line after the Scout Cavalry
	ClassCavalryArcher UnitClass = 36
	ClassScoutCavalry  UnitClass = 47 // the Scout Cavalry only, the effects repeat their cavalry commands for it
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
	case AttributeRange:
		unit.Range = operation(unit.Range, c.Value)
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

// encodeClassAmount is the Value of an armor or attack command: class * 256 + amount,
// e.g. encodeClassAmount(ArmorClassPierce, 1) is 769, encodeClassAmount(ArmorClassMelee, -1) is -1025.
// For MultiplyAttribute the amount is a percentage, e.g. 120 is * 1.2.
func encodeClassAmount(class ArmorClass, amount float64) float64 {
	sign := 1.0
	if amount < 0 {
		sign, amount = -1, -amount
	}
	return sign * (float64(class)*256 + amount)
}

// decodeClassAmount is the reverse of encodeClassAmount.
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
	Fletching: {EffectID: 192, Name: "Fletching", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeRange, 1},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeRange, 1},
	}},
	BodkinArrow: {EffectID: 193, Name: "Bodkin Arrow", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeRange, 1},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeRange, 1},
	}},
	Bracer: {EffectID: 194, Name: "Bracer", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeRange, 1},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeRange, 1},
	}},
	Chemistry: {EffectID: 47, Name: "Chemistry", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	PaddedArcherArmor: {EffectID: 202, Name: "Padded Archer Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	LeatherArcherArmor: {EffectID: 203, Name: "Leather Archer Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	RingArcherArmor: {EffectID: 208, Name: "Ring Archer Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
	}},
	// Thumb Ring makes archers attack faster, but not skirmishers: * 0.85 then * 1.1764706.
	// The .dat has 1.17647 for the Skirmisher, written here as 1.1764706 like the other two,
	// so the reload time stays 3 and does not break ties between civilizations.
	ThumbRing: {EffectID: 451, Name: "Thumb Ring", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassArcher, AttributeReloadTime, 0.85},
		{MultiplyAttribute, AnyUnit, ClassCavalryArcher, AttributeReloadTime, 0.85},
		// the Heavy Cavalry Archer attacks 10% faster, not 15%: * 0.85 then * 1.0588236
		{MultiplyAttribute, HeavyCavalryArcher, AnyClass, AttributeReloadTime, 1.0588236},
		{MultiplyAttribute, Skirmisher, AnyClass, AttributeReloadTime, 1.1764706},
		{MultiplyAttribute, EliteSkirmisher, AnyClass, AttributeReloadTime, 1.1764706},
		{MultiplyAttribute, ImperialSkirmisher, AnyClass, AttributeReloadTime, 1.1764706},
	}},
	// Italians Pavise, Muisca Herbalism and the Britons archer range bonuses are left out:
	// their skirmisher commands cancel the archer ones, so they change only units not ranked yet.
	Yeomen: {EffectID: 455, Name: "Yeomen", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassArcher, AttributeRange, 1},
	}, IsUniqueTechnology: true},
	Atlatl: {EffectID: 514, Name: "Atlatl", Commands: []EffectCommand{
		{AddAttribute, Skirmisher, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, EliteSkirmisher, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, ImperialSkirmisher, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, Skirmisher, AnyClass, AttributeRange, 1},
		{AddAttribute, EliteSkirmisher, AnyClass, AttributeRange, 1},
		{AddAttribute, ImperialSkirmisher, AnyClass, AttributeRange, 1},
	}, IsUniqueTechnology: true},
	TowerShields: {EffectID: 729, Name: "Tower Shields", Commands: []EffectCommand{
		{AddAttribute, Skirmisher, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
		{AddAttribute, EliteSkirmisher, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
		{AddAttribute, ImperialSkirmisher, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
	}, IsUniqueTechnology: true},
	ScaleMailArmor: {EffectID: 74, Name: "Scale Mail Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	ChainMailArmor: {EffectID: 76, Name: "Chain Mail Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	PlateMailArmor: {EffectID: 77, Name: "Plate Mail Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
	}},
	Squires: {EffectID: 204, Name: "Squires", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassInfantry, AttributeSpeed, 1.1},
	}},
	Arson: {EffectID: 628, Name: "Arson", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassInfantry, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 2)},
	}},
	ParthianTactics: {EffectID: 452, Name: "Parthian Tactics", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeAttack, encodeClassAmount(ArmorClassSpearmen, 2)},
	}},
	Husbandry: {EffectID: 39, Name: "Husbandry", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeSpeed, 1.1},
		{MultiplyAttribute, AnyUnit, ClassCavalryArcher, AttributeSpeed, 1.1},
	}},
	Forging: {EffectID: 67, Name: "Forging", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
	}},
	IronCasting: {EffectID: 68, Name: "Iron casting", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
	}},
	BlastFurnace: {EffectID: 75, Name: "Blast Furnace", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 2)},
		{AddAttribute, AnyUnit, ClassInfantry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 2)},
	}},
	PlateBardingArmor: {EffectID: 80, Name: "Plate Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 2)},
	}},
	ScaleBardingArmor: {EffectID: 81, Name: "Scale Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	ChainBardingArmor: {EffectID: 82, Name: "Chain Barding Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
	}},
	AutoUpgradeScoutFeudalAge: {EffectID: 16, Name: "Auto upgrade Scout Feudal Age", Commands: []EffectCommand{
		{AddAttribute, ScoutCavalry, AnyClass, AttributeSpeed, 0.35},
		{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassMelee, 2)},
	}},
	Bloodlines: {EffectID: 450, Name: "Bloodlines", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeHP, 20},
		{AddAttribute, AnyUnit, ClassCavalryArcher, AttributeHP, 20},
	}},
	Farimba: {EffectID: 606, Name: "Malian UT", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 5)},
	}, IsUniqueTechnology: true},
	ManipurCavalry: {EffectID: 667, Name: "Burmese UT", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassArchers, 4)},
	}, IsUniqueTechnology: true},
	Stirrups: {EffectID: 722, Name: "Khmer UT", Commands: []EffectCommand{
		{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeReloadTime, 0.75},
	}, IsUniqueTechnology: true},
	SilkArmor: {EffectID: 724, Name: "Malay UT", Commands: []EffectCommand{
		{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, LightCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, Hussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, WingedHussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, LightCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, Hussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		{AddAttribute, WingedHussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
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
		{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 0.75},
	}, IsUniqueTechnology: true},
	AznauriCavalry: {EffectID: 936, Name: "Aznauri Cavalry", Commands: []EffectCommand{
		{SetAttribute, AnyUnit, ClassCavalry, AttributePopulation, -0.8},
	}, IsUniqueTechnology: true},
	MingGuangArmor: {EffectID: 1062, Name: "Ming-Kuang Armor", Commands: []EffectCommand{
		{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 4)},
	}, IsUniqueTechnology: true},
}

// CivilizationBonusEffects are the civilization bonuses and team bonuses,
// active in the Imperial Age after researching every available technology,
// by civilization.
// A ranking applies them after the technologies, in this order.
var CivilizationBonusEffects = map[CivilizationID][]Effect{
	Byzantines: {
		{EffectID: 283, Name: "C-Bonus, Skirms Pikes cost -25%", Commands: []EffectCommand{
			{MultiplyAttribute, Skirmisher, AnyClass, AttributeAllCosts, 0.75},
			{MultiplyAttribute, EliteSkirmisher, AnyClass, AttributeAllCosts, 0.75},
			{MultiplyAttribute, ImperialSkirmisher, AnyClass, AttributeAllCosts, 0.75},
		}},
	},
	Bengalis: {
		{EffectID: 888, Name: "C-Bonus, Cavalry vs Skirm", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSkirmishers, 2)},
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
	Dravidians: {
		{EffectID: 864, Name: "C-Bonus, Skirm and EA faster attack", Commands: []EffectCommand{
			{MultiplyAttribute, Skirmisher, AnyClass, AttributeReloadTime, 0.8},
			{MultiplyAttribute, EliteSkirmisher, AnyClass, AttributeReloadTime, 0.8},
			{MultiplyAttribute, ImperialSkirmisher, AnyClass, AttributeReloadTime, 0.8},
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
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassInfantry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCapitalShips, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassElephantUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor6, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor7, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor9, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor10, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassAllBuildings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor12, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStoneDefense, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassPredatorAnimals, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassArchers, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShips, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassRams, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassTrees, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUniqueUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSiegeWeapons, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassWallsAndGates, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassGunpowderUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassBoars, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMonks, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCastles, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSpearmen, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalryArchers, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShockInfantry, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCamelUnits, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor33, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMamelukes, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeroesKings, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeavySiege, 120)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSkirmishers, 120)},
		}},
		{EffectID: 873, Name: "C-Bonus, +30% Cav bonus damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassInfantry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCapitalShips, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassElephantUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor6, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor7, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor9, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor10, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassAllBuildings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor12, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStoneDefense, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassPredatorAnimals, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassArchers, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShips, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassRams, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassTrees, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUniqueUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSiegeWeapons, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassWallsAndGates, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassGunpowderUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassBoars, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMonks, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCastles, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSpearmen, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalryArchers, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShockInfantry, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCamelUnits, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor33, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMamelukes, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeroesKings, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeavySiege, 111)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSkirmishers, 111)},
		}},
		{EffectID: 884, Name: "C-Bonus, +40% Cav bonus damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassInfantry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCapitalShips, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassElephantUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor6, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor7, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor9, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor10, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassAllBuildings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor12, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStoneDefense, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassPredatorAnimals, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassArchers, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShips, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassRams, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassTrees, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUniqueUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSiegeWeapons, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassWallsAndGates, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassGunpowderUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassBoars, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMonks, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCastles, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSpearmen, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCavalryArchers, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassShockInfantry, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassCamelUnits, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassUnusedArmor33, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMamelukes, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeroesKings, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassHeavySiege, 105)},
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSkirmishers, 105)},
		}},
	},
	Hindustanis: {
		{EffectID: 2, Name: "Hindustanis Team Bonus", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, LightCavalry, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, Hussar, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 2)},
			{AddAttribute, WingedHussar, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 2)},
		}},
	},
	// Inca: each age bonus first undoes the previous one, e.g. * 1.174 (for 1 / 0.85) then * 0.8.
	Inca: {
		{EffectID: 156, Name: "C-Bonus, Military cost -15%", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 0.85},
		}},
		{EffectID: 158, Name: "C-Bonus, Military cost -20%", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 1.174},
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 0.8},
		}},
		{EffectID: 159, Name: "C-Bonus, Military cost -25%", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 1.25},
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 0.75},
		}},
		{EffectID: 160, Name: "C-Bonus, Military cost -30%", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 1.328},
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFoodCost, 0.7},
		}},
	},
	Jurchens: {
		{EffectID: 994, Name: "C-Bonus, cav & Fire Lancer +25% attack", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeReloadTime, 0.8},
		}},
		{EffectID: 1009, Name: "C-Bonus, No friendly damage", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeFriendlyFireDamage, 0.5},
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeFriendlyFireDamage, 0.5},
		}},
	},
	Khitans: {
		{EffectID: 67, Name: "C-Bonus, Double Forging", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 68, Name: "C-Bonus, Double Iron Casting", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 1011, Name: "C-Bonus, Skirms Spear Scout train 25% faster", Commands: []EffectCommand{
			{MultiplyAttribute, ScoutCavalry, AnyClass, AttributeTrainTime, 0.869565},
			{MultiplyAttribute, LightCavalry, AnyClass, AttributeTrainTime, 0.869565},
			{MultiplyAttribute, Hussar, AnyClass, AttributeTrainTime, 0.869565},
		}},
	},
	Koreans: {
		{EffectID: 561, Name: "C-Bonus, Units cost less wood", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeWoodCost, 0.5},
		}},
	},
	Lithuanians: {
		{EffectID: 745, Name: "C-Bonus, Skirms Pikes 10% faster", Commands: []EffectCommand{
			{MultiplyAttribute, Skirmisher, AnyClass, AttributeSpeed, 1.1},
			{MultiplyAttribute, EliteSkirmisher, AnyClass, AttributeSpeed, 1.1},
			{MultiplyAttribute, ImperialSkirmisher, AnyClass, AttributeSpeed, 1.1},
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
	Mapuche: {
		{EffectID: 1383, TechnologyID: 1383, Name: "C-Bonus, Skirm Spear +5 HP Feudal", Commands: mapucheSkirmisherHP},
		{EffectID: 1383, TechnologyID: 1384, Name: "C-Bonus, Skirm Spear +5 HP Castle", Commands: mapucheSkirmisherHP},
		{EffectID: 1383, TechnologyID: 1385, Name: "C-Bonus, Skirm Spear +5 HP Imperial", Commands: mapucheSkirmisherHP},
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
	Muisca: {
		{EffectID: 1369, TechnologyID: 1369, Name: "C-Bonus, Arch + Champi +1 armor Age2", Commands: muiscaArcherMeleeArmor},
		{EffectID: 1369, TechnologyID: 1370, Name: "C-Bonus, Arch + Champi +1 armor Age3", Commands: muiscaArcherMeleeArmor},
		{EffectID: 1369, TechnologyID: 1371, Name: "C-Bonus, Arch + Champi +1 armor Age4", Commands: muiscaArcherMeleeArmor},
	},
	Poles: {
		{EffectID: 802, Name: "Poles Team Bonus", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassArchers, 1)},
			{AddAttribute, LightCavalry, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassArchers, 1)},
			{AddAttribute, Hussar, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassArchers, 1)},
			{AddAttribute, WingedHussar, AnyClass, AttributeAttack, encodeClassAmount(ArmorClassArchers, 1)},
		}},
	},
	Portuguese: {
		{EffectID: 33, Name: "C-Bonus, All units -20% gold", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassCavalry, AttributeGoldCost, 0.8},
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeGoldCost, 0.8},
		}},
	},
	Saracens: {
		{EffectID: 409, Name: "Saracens Team Bonus", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassArcher, AttributeAttack, encodeClassAmount(ArmorClassStandardBuildings, 3)},
		}},
	},
	Sicilians: {
		{EffectID: 796, Name: "C-Bonus, Bonus damage resistance", Commands: []EffectCommand{
			{SetAttribute, AnyUnit, ClassCavalry, AttributeBonusDamageResistance, 0.4},
			{SetAttribute, AnyUnit, ClassArcher, AttributeBonusDamageResistance, 0.4},
		}},
	},
	Teutons: {
		{EffectID: 333, Name: "C-Bonus, Inf Cav +1 armor Age3", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		}},
		{EffectID: 334, Name: "C-Bonus, Inf Cav +1 armor Age4", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
		}},
	},
	Turks: {
		{EffectID: 509, Name: "C-Bonus, Light cavalry +1P armor", Commands: []EffectCommand{
			{AddAttribute, ScoutCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
			{AddAttribute, LightCavalry, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
			{AddAttribute, Hussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
			{AddAttribute, WingedHussar, AnyClass, AttributeArmor, encodeClassAmount(ArmorClassPierce, 1)},
		}},
	},
	Varangians: {
		// the game describes it as "Bloodlines and Caravan effects +50%": +10 HP on top of +20,
		// the .dat name "+65%" is outdated
		{EffectID: 1489, Name: "C-Bonus, Bloodlines +65% more effective", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeHP, 10},
		}},
	},
	Vietnamese: {
		{EffectID: 672, Name: "C-Bonus, Archers +20% HP", Commands: []EffectCommand{
			{MultiplyAttribute, AnyUnit, ClassArcher, AttributeHP, 1.2},
		}},
	},
	Wei: {
		{EffectID: 1027, Name: "Wei Team Bonus", Commands: []EffectCommand{
			{AddAttribute, AnyUnit, ClassCavalry, AttributeAttack, encodeClassAmount(ArmorClassSiegeWeapons, 2)},
		}},
	},
}

// mapucheSkirmisherHP is the effect of the Mapuche "Skirm Spear +5 HP" bonus, triggered in each age.
var mapucheSkirmisherHP = []EffectCommand{
	{AddAttribute, EliteSkirmisher, AnyClass, AttributeHP, 5},
	{AddAttribute, Skirmisher, AnyClass, AttributeHP, 5},
	{AddAttribute, ImperialSkirmisher, AnyClass, AttributeHP, 5},
}

// muiscaArcherMeleeArmor is the effect of the Muisca "Arch + Champi +1 armor" bonus, triggered in each age.
var muiscaArcherMeleeArmor = []EffectCommand{
	{AddAttribute, AnyUnit, ClassArcher, AttributeArmor, encodeClassAmount(ArmorClassMelee, 1)},
}
