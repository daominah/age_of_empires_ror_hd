package aoego

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidCivID = errors.New("invalid civilization ID, check civilization enum list")

// Civilization must be initialized with func NewCivilization
type Civilization struct {
	ID            CivilizationID
	Name          string
	Name2         string
	DisabledUnits map[UnitID]bool
	DisabledTechs map[TechID]bool
	Bonuses       []EffectFunc
}

// CivilizationID is enum
type CivilizationID int

// CivilizationID enum
const (
	Assyrian     CivilizationID = 81
	Babylonian   CivilizationID = 82
	Carthaginian CivilizationID = 205
	Choson       CivilizationID = 91
	Egyptian     CivilizationID = 83
	Greek        CivilizationID = 84
	Hittite      CivilizationID = 85
	Macedonian   CivilizationID = 206
	Minoan       CivilizationID = 86
	Palmyran     CivilizationID = 207
	Persian      CivilizationID = 87
	Phoenician   CivilizationID = 88
	Roman        CivilizationID = 208
	Shang        CivilizationID = 89
	Sumerian     CivilizationID = 90
	Yamato       CivilizationID = 92

	FullTechTree CivilizationID = 0
)

func NewCivilization(civID CivilizationID) (*Civilization, error) {
	c := &Civilization{
		ID:            civID,
		DisabledUnits: make(map[UnitID]bool),
		DisabledTechs: make(map[TechID]bool),
	}
	switch civID {
	case FullTechTree:
		c.Name, c.Name2 = "FullTechTree", "FullTechTree"
		return c, nil

	case Assyrian:
		c.Name, c.Name2 = "Assyrian", "Assyria"

	case Babylonian:
		c.Name, c.Name2 = "Babylonian", "Babylon"

	case Carthaginian:
		c.Name, c.Name2 = "Carthaginian", "Carthage"

	case Choson:
		c.Name, c.Name2 = "Choson", "Choson"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Priest cost -32% (stated -30%): 85 gold instead of 125.
				e.UnitStats[Priest].Cost.Multiply(0.68)
			},
			func(e *EmpireDeveloping) {
				// * Iron Age Swordsmen HP +80.
			},
			func(e *EmpireDeveloping) {
				// * Towers range +2.
			},
		}

	case Egyptian:
		c.Name, c.Name2 = "Egyptian", "Egypt"

	case Greek:
		c.Name, c.Name2 = "Greek", "Greek"

	case Hittite:
		c.Name, c.Name2 = "Hittite", "Hittite"

	case Macedonian:
		c.Name, c.Name2 = "Macedonian", "Macedon"
		c.DisabledTechs = map[TechID]bool{
			Wheel:                true,
			EnableChariotArcher:  true,
			EnableElephantArcher: true,
			EnableChariot:        true,
			ScytheChariot:        true,
			EnableCamel:          true,
			Nobility:             true,
			LongSword:            true,
			Legion:               true,

			Engineering:     true,
			Siegecraft:      true,
			Craftsmanship:   true,
			Helepolis:       true,
			Catapult:        true,
			MassiveCatapult: true,
			EnableFireBoat:  true,

			EnableTemple: true,
			TempleBuilt:  true,
			Astrology:    true,
			Mysticism:    true,
			Polytheism:   true,
			Afterlife:    true,
			Monotheism:   true,
			Fanaticism:   true,
			Zealotry:     true,
			Sacrifice:    true,

			FortifiedWall: true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Academy units pierce armor +2.
			},
			func(e *EmpireDeveloping) {
				// * Siege units cost -50%.
				e.UnitStats[StoneThrower].Cost.Multiply(0.5)
				e.UnitStats[Ballista].Cost.Multiply(0.5)
			},
			func(e *EmpireDeveloping) {
				// * Melee units sight +2.
			},
			func(e *EmpireDeveloping) {
				// * All units are 4 times more resistant to conversion.
			},
		}

	case Minoan:
		c.Name, c.Name2 = "Minoan", "Minoa"

	case Palmyran:
		c.Name, c.Name2 = "Palmyran", "Palmyra"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Forager, Hunter, Gold Miner, Stone Miner work rate +44%; Woodcutter +36%,
				//	(Farmer and Builder work rate are normal, stated Villager work rate +25%).
				// * Villager armor +1 and pierce armor +1 (in game only show armor +1).
				// * Villager cost +50% (so 75 food instead of 50).
				e.UnitStats[Villager].Cost.Multiply(1.5)
			},
			func(e *EmpireDeveloping) {
				// * Camel Rider move speed +25% (so as fast as Heavy Horse Archer).
			},
			func(e *EmpireDeveloping) {
				// * Free tribute (instead of 25% taxed, other civilizations lost 125 resource to give 100).
			},
		}

	case Persian:
		c.Name, c.Name2 = "Persian", "Persia"

	case Phoenician:
		c.Name, c.Name2 = "Phoenician", "Phoenicia"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Woodcutter work rate +36% and capacity +3 (stated +30%).
			},
			func(e *EmpireDeveloping) {
				// * Elephants cost -25%.
				e.UnitStats[Elephant].Cost.Multiply(0.75)
				e.UnitStats[ElephantArcher].Cost.Multiply(0.75)
			},
			func(e *EmpireDeveloping) {
				// * Catapult Trireme and Juggernaught attack speed +72%.
			},
		}

	case Roman:
		c.Name, c.Name2 = "Roman", "Rome"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Buildings cost -15% (except Tower, Wall, Wonder).
				for building := range AllBuildings {
					if building == Tower || building == Wall || building == Wonder {
						continue
					}
					e.UnitStats[building].Cost.Multiply(0.85)
				}
			},
			func(e *EmpireDeveloping) {
				//   Tower cost -50%.
				e.UnitStats[Tower].Cost.Multiply(0.5)
			},
			func(e *EmpireDeveloping) {
				//* Swordsmen attack speed +50% (stated 33%, they mean attack reload time).
			},
		}

	case Shang:
		c.Name, c.Name2 = "Shang", "Shang"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Villager cost 35 food instead of 50 (so Villager cost -30%).
				e.UnitStats[Villager].Cost.Food = 35
			},
			func(e *EmpireDeveloping) {
				// * Wall HP x2.
			},
		}

	case Sumerian:
		c.Name, c.Name2 = "Sumerian", "Sumeria"
		c.DisabledTechs = map[TechID]bool{
			EnableCavalry: true,
			ImprovedBow:   true,
			Astrology:     true,

			Metallurgy:     true,
			IronShield:     true,
			Craftsmanship:  true,
			Coinage:        true,
			EnableBallista: true,

			Afterlife:  true,
			Monotheism: true,
			Fanaticism: true,
			Zealotry:   true,

			HeavyTransport:  true,
			CatapultTrireme: true,
		}
		c.DisabledUnits = map[UnitID]bool{
			Cavalry:  true,
			Ballista: true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Siege units attack speed +43%.
			},
			func(e *EmpireDeveloping) {
				// * Villager HP +15 (so 40 instead of 25).
			},
			func(e *EmpireDeveloping) {
				// * Farm food +250 (starting at 500 instead of 250).
			},
		}

	case Yamato:
		c.Name, c.Name2 = "Yamato", "Yamato"
		c.DisabledTechs = map[TechID]bool{
			EnableChariotArcher: true,
			EnableChariot:       true,
			EnableCamel:         true,
			Broadsword:          true,
			Astrology:           true,
			Mysticism:           true,

			Catapult:       true,
			EnableBallista: true,

			Medicine:   true,
			Monotheism: true,
			Fanaticism: true,
			Zealotry:   true,
			Sacrifice:  true,

			GuardTower:           true,
			EnableElephantArcher: true,
			EnableWarElephant:    true,
			FortifiedWall:        true,
			EnableFireBoat:       true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Villager move speed +18%
			},
			func(e *EmpireDeveloping) {
				// * Mounted units cost -25%
				e.UnitStats[Scout].Cost.Multiply(0.75)
				e.UnitStats[Cavalry].Cost.Multiply(0.75)
				e.UnitStats[HorseArcher].Cost.Multiply(0.75)
			},
			func(e *EmpireDeveloping) {
				// * Ships HP +30%
			},
		}

	default:
		return nil, fmt.Errorf("civID %v: %w", civID, ErrInvalidCivID)
	}
	return c, nil
}

var AllCivilizations []Civilization

func init() {
	for _, civID := range []CivilizationID{
		Assyrian, Egyptian, Sumerian,
		Babylonian, Hittite, Persian,
		Carthaginian, Macedonian, Palmyran, Roman,
		Choson, Shang, Yamato,
		Greek, Minoan, Phoenician,
	} {
		civilization, err := NewCivilization(civID)
		if err != nil {
			panic(fmt.Errorf("error NewCivilization(%v): %w", civID, err))
		}
		AllCivilizations = append(AllCivilizations, *civilization)
	}
}

func GuessCivilization(fileName string) CivilizationID {
	for _, civ := range AllCivilizations {
		if strings.Contains(fileName, civ.Name) ||
			strings.Contains(fileName, civ.Name2) {
			return civ.ID
		}
	}
	return FullTechTree
}
