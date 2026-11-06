package logic

import (
	"math"
	"reflect"
	"testing"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
)

func TestLoadUnitBase(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	for _, c := range []struct {
		id   UnitID
		want Unit
	}{
		// the scout line, with the Feudal Age automatic +2 attack and +0.35 speed
		{ScoutCavalry, Unit{
			ID:    ScoutCavalry,
			Name:  "Scout Cavalry",
			Class: ClassScoutCavalry,
			HP:    45,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        5,
				ArmorClassMonks:        6,
				ArmorClassAllBuildings: 0, ArmorClassArchers: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 2,
				ArmorClassMelee:  0, ArmorClassCavalry: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 2,
			Speed:      1.55,
			TrainTime:  30,
			Cost:       Cost{Food: 80},
		}},
		// cavalry
		{Knight, Unit{
			ID:    Knight,
			Name:  "Knight",
			Class: ClassCavalry,
			HP:    100,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        10,
				ArmorClassAllBuildings: 0, ArmorClassArchers: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  2,
				ArmorClassMelee:   2,
				ArmorClassCavalry: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 1.8,
			Speed:      1.35,
			TrainTime:  30,
			Cost:       Cost{Food: 60, Gold: 75},
		}},
		// cavalry, melee attack at range 3, bonus against cavalry
		{Mameluke, Unit{
			ID:    Mameluke,
			Name:  "Mameluke",
			Class: ClassCavalry,
			HP:    80,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        8,
				ArmorClassCavalry:      9,
				ArmorClassAllBuildings: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 0, ArmorClassMelee: 0, ArmorClassUniqueUnits: 0, ArmorClassCamelUnits: 0, ArmorClassUnusedArmor31: 0, ArmorClassMamelukes: 0, ArmorClassRoyalHeirs: -3,
			},
			Range:      3,
			ReloadTime: 2,
			Speed:      1.4,
			TrainTime:  23,
			Cost:       Cost{Food: 55, Gold: 85},
		}},
		// cavalry, melee attack at range 1
		{SteppeLancer, Unit{
			ID:    SteppeLancer,
			Name:  "Steppe Lancer",
			Class: ClassCavalry,
			HP:    60,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        9,
				ArmorClassAllBuildings: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 1,
				ArmorClassMelee:  0, ArmorClassCavalry: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      1,
			ReloadTime: 2,
			Speed:      1.45,
			TrainTime:  24,
			Cost:       Cost{Food: 70, Gold: 40},
		}},
		// foot archer, bonus against spearmen
		{Archer, Unit{
			ID:    Archer,
			Name:  "Archer",
			Class: ClassArcher,
			HP:    30,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:       4,
				ArmorClassSpearmen:     3,
				ArmorClassStoneDefense: 0, ArmorClassRams: 0, ArmorClassStandardBuildings: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 0, ArmorClassMelee: 0, ArmorClassArchers: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      4,
			ReloadTime: 2,
			Speed:      0.96,
			TrainTime:  35,
			Cost:       Cost{Wood: 25, Gold: 45},
		}},
		// foot archer, bonus against archers, the Cavalry Archers bonus starts at 0 before the Elite Skirmisher
		{Skirmisher, Unit{
			ID:    Skirmisher,
			Name:  "Skirmisher",
			Class: ClassArcher,
			HP:    30,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:    2,
				ArmorClassArchers:   3,
				ArmorClassSpearmen:  3,
				ArmorClassMamelukes: 2,
				ArmorClassRams:      0, ArmorClassStandardBuildings: 0, ArmorClassCavalryArchers: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 3,
				ArmorClassMelee:  0, ArmorClassArchers: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0,
			},
			Range:      4,
			MinRange:   1,
			ReloadTime: 3,
			Speed:      0.96,
			TrainTime:  26,
			Cost:       Cost{Food: 25, Wood: 35},
		}},
		// unique foot archer, high pierce armor
		{RattanArcher, Unit{
			ID:    RattanArcher,
			Name:  "Rattan Archer",
			Class: ClassArcher,
			HP:    35,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:   6,
				ArmorClassSpearmen: 2,
				ArmorClassInfantry: 0, ArmorClassRams: 0, ArmorClassStandardBuildings: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 4,
				ArmorClassMelee:  0, ArmorClassArchers: 0, ArmorClassUniqueUnits: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      4,
			ReloadTime: 2,
			Speed:      1.1,
			TrainTime:  16,
			Cost:       Cost{Wood: 50, Gold: 45},
		}},
		// cavalry archer
		{CavalryArcher, Unit{
			ID:    CavalryArcher,
			Name:  "Cavalry Archer",
			Class: ClassCavalryArcher,
			HP:    50,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:   6,
				ArmorClassSpearmen: 2,
				ArmorClassArchers:  0, ArmorClassRams: 0, ArmorClassStandardBuildings: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce: 0, ArmorClassMelee: 0, ArmorClassCavalry: 0, ArmorClassArchers: 0, ArmorClassCavalryArchers: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      4,
			ReloadTime: 2,
			Speed:      1.4,
			TrainTime:  37,
			Cost:       Cost{Wood: 40, Gold: 60},
		}},
		// infantry
		{ManAtArms, Unit{
			ID:    ManAtArms,
			Name:  "Man-at-Arms",
			Class: ClassInfantry,
			HP:    45,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:             6,
				ArmorClassStandardBuildings: 2,
				ArmorClassShockInfantry:     3,
				ArmorClassCavalry:           0, ArmorClassArchers: 0, ArmorClassCamelUnits: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:   1,
				ArmorClassInfantry: 0, ArmorClassMelee: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 2,
			Speed:      0.96,
			TrainTime:  21,
			Cost:       Cost{Food: 50, Gold: 20},
		}},
		// infantry, melee attack at range 5
		{ThrowingAxeman, Unit{
			ID:    ThrowingAxeman,
			Name:  "Throwing Axeman",
			Class: ClassInfantry,
			HP:    60,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:             7,
				ArmorClassStandardBuildings: 1,
				ArmorClassShockInfantry:     1,
				ArmorClassCavalry:           0, ArmorClassArchers: 0, ArmorClassCamelUnits: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassInfantry: 0, ArmorClassPierce: 0, ArmorClassMelee: 0, ArmorClassUniqueUnits: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      5,
			ReloadTime: 2,
			Speed:      1,
			TrainTime:  13,
			Cost:       Cost{Food: 55, Gold: 25},
		}},
	} {
		// WHEN the unit is loaded before any technology or civilization bonus
		got := roundUnit(tree.LoadUnitBase(c.id))

		// THEN every stat is the game data, an armor class missing from the map is not hit at all
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%v:\n got  %+v\n want %+v", c.want.Name, got, c.want)
		}
	}
}

// roundUnit rounds the float stats to 6 decimals,
// so a value from a multiplied effect compares equal, e.g. speed 1.2 + 0.35 is 1.5499999999999998.
func roundUnit(u Unit) Unit {
	round := func(v float64) float64 { return math.Round(v*1e6) / 1e6 }
	u = u.clone()
	u.HP, u.Range, u.ReloadTime, u.Speed, u.TrainTime = round(u.HP), round(u.Range), round(u.ReloadTime), round(u.Speed), round(u.TrainTime)
	u.Cost = Cost{Food: round(u.Cost.Food), Wood: round(u.Cost.Wood), Gold: round(u.Cost.Gold), Stone: round(u.Cost.Stone)}
	for class, v := range u.Attacks {
		u.Attacks[class] = round(v)
	}
	for class, v := range u.Armors {
		u.Armors[class] = round(v)
	}
	return u
}

func TestLoadUnitUpgraded(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	for _, c := range []struct {
		id   UnitID
		want Unit
	}{
		// Fletching, Bodkin Arrow, Bracer, Chemistry +4 attack +3 range, archer armors and Parthian Tactics +4/+6,
		// Bloodlines +20 HP, Husbandry +10% speed, Thumb Ring 10% faster attack
		{HeavyCavalryArcher, Unit{
			ID:    HeavyCavalryArcher,
			Name:  "Heavy Cavalry Archer",
			Class: ClassCavalryArcher,
			HP:    80,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:   11,
				ArmorClassSpearmen: 6,
				ArmorClassArchers:  0, ArmorClassRams: 0, ArmorClassStandardBuildings: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  6,
				ArmorClassMelee:   5,
				ArmorClassCavalry: 0, ArmorClassArchers: 0, ArmorClassCavalryArchers: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      7,
			ReloadTime: 1.8,
			Speed:      1.54,
			TrainTime:  30,
			Cost:       Cost{Wood: 40, Gold: 60},
		}},
		// a unique unit gets the archer upgrades too: +4 attack +3 range, archer armors +3/+4,
		// Thumb Ring 15% faster attack,
		// only the Vietnamese train it, so it gets their archer HP +20% bonus too
		{EliteRattanArcher, Unit{
			ID:    EliteRattanArcher,
			Name:  "Elite Rattan Archer",
			Class: ClassArcher,
			HP:    48,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:   11,
				ArmorClassSpearmen: 2,
				ArmorClassInfantry: 0, ArmorClassRams: 0, ArmorClassStandardBuildings: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  10,
				ArmorClassMelee:   3,
				ArmorClassArchers: 0, ArmorClassUniqueUnits: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      8,
			ReloadTime: 1.7,
			Speed:      1.1,
			TrainTime:  16,
			Cost:       Cost{Wood: 50, Gold: 45},
		}},
		// +4 attack +3 range, archer armors +3/+4, Thumb Ring does not change the skirmisher reload time
		{EliteSkirmisher, Unit{
			ID:    EliteSkirmisher,
			Name:  "Elite Skirmisher",
			Class: ClassArcher,
			HP:    35,
			Attacks: map[ArmorClass]float64{
				ArmorClassPierce:         7,
				ArmorClassArchers:        4,
				ArmorClassSpearmen:       4,
				ArmorClassCavalryArchers: 2,
				ArmorClassMamelukes:      2,
				ArmorClassRams:           0, ArmorClassStandardBuildings: 0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  8,
				ArmorClassMelee:   3,
				ArmorClassArchers: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0,
			},
			Range:      8,
			MinRange:   1,
			ReloadTime: 3,
			Speed:      0.96,
			TrainTime:  22,
			Cost:       Cost{Food: 25, Wood: 35},
		}},
		// Forging, Iron Casting, Blast Furnace +4 attack, the mail armors +3/+4, Squires +10% speed, Arson +2 against buildings
		{Halberdier, Unit{
			ID:    Halberdier,
			Name:  "Halberdier",
			Class: ClassInfantry,
			HP:    60,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:             10,
				ArmorClassElephantUnits:     28,
				ArmorClassCavalry:           32,
				ArmorClassShips:             17,
				ArmorClassStandardBuildings: 3,
				ArmorClassShockInfantry:     1,
				ArmorClassCamelUnits:        26,
				ArmorClassMamelukes:         7,
				ArmorClassArchers:           0,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:   4,
				ArmorClassMelee:    3,
				ArmorClassInfantry: 0, ArmorClassSpearmen: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 3,
			Speed:      1.1,
			TrainTime:  22,
			Cost:       Cost{Food: 35, Wood: 25},
		}},
		// Forging line +4 attack, the barding armors +3/+4, Bloodlines +20 HP, Husbandry +10% speed
		{Hussar, Unit{
			ID:    Hussar,
			Name:  "Hussar",
			Class: ClassCavalry,
			HP:    95,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        11,
				ArmorClassMonks:        12,
				ArmorClassAllBuildings: 0, ArmorClassArchers: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  6,
				ArmorClassMelee:   3,
				ArmorClassCavalry: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 1.9,
			Speed:      1.65,
			TrainTime:  30,
			Cost:       Cost{Food: 80},
		}},
		// the same cavalry upgrades
		{Paladin, Unit{
			ID:    Paladin,
			Name:  "Paladin",
			Class: ClassCavalry,
			HP:    180,
			Attacks: map[ArmorClass]float64{
				ArmorClassMelee:        18,
				ArmorClassAllBuildings: 0, ArmorClassArchers: 0, ArmorClassSiegeWeapons: 0, ArmorClassStandardBuildings: 0, ArmorClassUnusedArmor31: 0, ArmorClassSkirmishers: 0, ArmorClassRoyalHeirs: -3,
			},
			Armors: map[ArmorClass]float64{
				ArmorClassPierce:  7,
				ArmorClassMelee:   5,
				ArmorClassCavalry: 0, ArmorClassUnusedArmor31: 0,
			},
			Range:      0,
			ReloadTime: 1.9,
			Speed:      1.485,
			TrainTime:  30,
			Cost:       Cost{Food: 60, Gold: 75},
		}},
	} {
		// WHEN the unit gets every technology any civilization can research
		got := roundUnit(tree.LoadUnitUpgraded(c.id))

		// THEN every stat includes the upgrades, without civilization bonuses or unique technologies,
		// except for a unit only one civilization trains
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("%v:\n got  %+v\n want %+v", c.want.Name, got, c.want)
		}
	}
}
