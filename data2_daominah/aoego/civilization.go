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
		DisabledTechs: make(map[TechID]bool),
	}
	switch civID {
	case FullTechTree:
		c.Name, c.Name2 = "FullTechTree", "FullTechTree"
		return c, nil

	case Assyrian:
		c.DisabledTechs = map[TechID]bool{
			EnableSlinger:        true,
			ImprovedBow:          true,
			EnableElephantArcher: true,
			EnableWarElephant:    true,
			Nobility:             true,
			Architecture:         true,
			Aristocracy:          true,
			Alchemy:              true,
			Engineering:          true,
			ChainMailInfantry:    true,
			ChainMailArchers:     true,
			ChainMailCavalry:     true,
			BronzeShield:         true,
			Phalanx:              true,
			CatapultTrireme:      true,
			HeavyTransport:       true,
		}
		c.Name, c.Name2 = "Assyrian", "Assyria"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Archers attack reload to 1.1
			},
			func(e *EmpireDeveloping) {
				// * Villager move speed +0.2
			},
		}

	case Babylonian:
		c.Name, c.Name2 = "Babylonian", "Babylon"
		c.DisabledTechs = map[TechID]bool{
			EnableElephantArcher: true,
			HeavyCalvary:         true,
			EnableWarElephant:    true,
			Metallurgy:           true,
			ChainMailInfantry:    true,
			ChainMailArchers:     true,
			ChainMailCavalry:     true,
			IronShield:           true,
			EnableBallista:       true,
			Phalanx:              true,
			Trireme:              true,
			CatapultTrireme:      true,
			HeavyTransport:       true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Priest's rejuvenation +0.75
			},
			func(e *EmpireDeveloping) {
				// * Stone Miner work rate +0.2 and capacity +3 (stated 30%).
			},
			func(e *EmpireDeveloping) {
				// * Tower and Wall HP x2.
			},
		}

	case Carthaginian:
		c.Name, c.Name2 = "Carthaginian", "Carthage"
		c.DisabledTechs = map[TechID]bool{
			CompositeBow:        true,
			EnableChariotArcher: true,
			EnableChariot:       true,
			Metallurgy:          true,
			ChainMailInfantry:   true,
			ChainMailArchers:    true,
			ChainMailCavalry:    true,
			Siegecraft:          true,
			Astrology:           true,
			Monotheism:          true,
			Fanaticism:          true,
			Catapult:            true,
			FortifiedWall:       true,
			CatapultTrireme:     true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Elephant and Academy units HP +25%.
			},
			func(e *EmpireDeveloping) {
				// * Light Transport move speed +25%, Heavy Transport move speed +43%
			},
			func(e *EmpireDeveloping) {
				// * Fire Galley attack +6 (24+12 instead of 24+6).
			},
		}

	case Choson:
		c.Name, c.Name2 = "Choson", "Choson"
		c.DisabledTechs = map[TechID]bool{
			CompositeBow:         true,
			EnableChariotArcher:  true,
			EnableElephantArcher: true,
			EnableCamel:          true,
			EnableChariot:        true,
			EnableWarElephant:    true,
			Nobility:             true,
			Aristocracy:          true,
			Alchemy:              true,
			Engineering:          true,
			ChainMailInfantry:    true,
			ChainMailArchers:     true,
			ChainMailCavalry:     true,
			IronShield:           true,
			Phalanx:              true,
			Catapult:             true,
			CatapultTrireme:      true,
			EnableFireBoat:       true,
			HeavyTransport:       true,
		}
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
		c.DisabledTechs = map[TechID]bool{
			EnableHorseArcher: true,
			EnableCavalry:     true,
			BronzeShield:      true,
			Coinage:           true,
			Siegecraft:        true,
			Catapult:          true,
			EnableBallista:    true,
			Phalanx:           true,
			LongSword:         true,
			EnableFireBoat:    true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Chariots HP +33%
			},
			func(e *EmpireDeveloping) {
				// * Priest range +3 (10+6 instead of 10+3).
			},
			func(e *EmpireDeveloping) {
				// * Gold Miner work rate +44% and capacity +2 (stated +20%).
			},
		}

	case Greek:
		c.DisabledTechs = map[TechID]bool{
			ImprovedBow:          true,
			EnableChariotArcher:  true,
			EnableHorseArcher:    true,
			EnableElephantArcher: true,
			EnableChariot:        true,
			EnableCamel:          true,
			EnableWarElephant:    true,
			Metallurgy:           true,
			Monotheism:           true,
			Zealotry:             true,
			Sacrifice:            true,
			BroadSword:           true,
			EnableFireBoat:       true,
		}
		c.Name, c.Name2 = "Greek", "Greek"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Academy  +0.3 tiles/s
			},
			func(e *EmpireDeveloping) {
				// * Warships move speed +17% (stated +30%).
			},
		}

	case Hittite:
		c.DisabledTechs = map[TechID]bool{
			EnableSlinger:   true,
			ImprovedBow:     true,
			HeavyCalvary:    true,
			Mysticism:       true,
			Polytheism:      true,
			Medicine:        true,
			Afterlife:       true,
			Monotheism:      true,
			Fanaticism:      true,
			Zealotry:        true,
			Sacrifice:       true,
			EnableBallista:  true,
			LongSword:       true,
			Trireme:         true,
			CatapultTrireme: true,
			FishingShip:     true,
			HeavyTransport:  true,
		}
		c.Name, c.Name2 = "Hittite", "Hittite"
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Archers attack +1.
			},
			func(e *EmpireDeveloping) {
				// * Siege units HP x2.
			},
			func(e *EmpireDeveloping) {
				// * Warships range +4.
			},
		}

	case Macedonian:
		c.Name, c.Name2 = "Macedonian", "Macedon"
		c.DisabledTechs = map[TechID]bool{
			Wheel:                true,
			EnableChariotArcher:  true,
			EnableElephantArcher: true,
			EnableChariot:        true,
			EnableCamel:          true,
			Nobility:             true,
			Engineering:          true,
			Craftsmanship:        true,
			Siegecraft:           true,
			EnableTemple:         true,
			TempleBuilt:          true,
			Astrology:            true,
			Mysticism:            true,
			Polytheism:           true,
			Medicine:             true,
			Afterlife:            true,
			Monotheism:           true,
			Fanaticism:           true,
			Zealotry:             true,
			Sacrifice:            true,
			Catapult:             true,
			LongSword:            true,
			FortifiedWall:        true,
			EnableFireBoat:       true,
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
		c.DisabledTechs = map[TechID]bool{
			EnableChariotArcher:  true,
			EnableHorseArcher:    true,
			EnableElephantArcher: true,
			EnableChariot:        true,
			HeavyCalvary:         true,
			EnableWarElephant:    true,
			Astrology:            true,
			Mysticism:            true,
			Afterlife:            true,
			Monotheism:           true,
			Fanaticism:           true,
			Zealotry:             true,
			Sacrifice:            true,
			GuardTower:           true,
			FortifiedWall:        true,
			EnableFireBoat:       true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Composite Bowman range +2.
			},
			func(e *EmpireDeveloping) {
				// * Farm food +60 (starting at 310 instead of 250).
			},
			func(e *EmpireDeveloping) {
				// * Ships cost -30%.
			},
		}

	case Palmyran:
		c.Name, c.Name2 = "Palmyran", "Palmyra"
		c.DisabledTechs = map[TechID]bool{
			EnableElephantArcher: true,
			Metallurgy:           true,
			Logistics:            true,
			Aristocracy:          true,
			Engineering:          true,
			Craftsmanship:        true,
			Coinage:              true,
			TowerShield:          true,
			Mysticism:            true,
			Polytheism:           true,
			Monotheism:           true,
			Medicine:             true,
			Sacrifice:            true,
			LongSword:            true,
			Plow:                 true,
			CatapultTrireme:      true,
			HeavyTransport:       true,
		}
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
		c.DisabledTechs = map[TechID]bool{
			EnableChariotArcher: true,
			EnableChariot:       true,
			Aristocracy:         true,
			Ballistics:          true,
			Wheel:               true,
			Plow:                true,
			Artisanship:         true,
			Coinage:             true,
			Siegecraft:          true,
			EnableAcademy:       true,
			EnableBallista:      true,
			EnableFireBoat:      true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// * Hunter work rate +66% and capacity +3 (stated +30%).
			},
			func(e *EmpireDeveloping) {
				// * Elephants move speed +56% (stated +50%).
			},
			func(e *EmpireDeveloping) {
				// * Trireme attack speed +38% (stated 50%).
			},
		}

	case Phoenician:
		c.Name, c.Name2 = "Phoenician", "Phoenicia"
		c.DisabledTechs = map[TechID]bool{
			EnableHorseArcher: true,
			HeavyCalvary:      true,
			Architecture:      true,
			Siegecraft:        true,
			Catapult:          true,
			EnableBallista:    true,
			Metallurgy:        true,
			ChainMailInfantry: true,
			ChainMailArchers:  true,
			ChainMailCavalry:  true,
			EnableFireBoat:    true,
		}
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
		c.DisabledTechs = map[TechID]bool{
			CompositeBow:         true,
			EnableChariotArcher:  true,
			EnableHorseArcher:    true,
			EnableElephantArcher: true,
			HeavyCalvary:         true,
			EnableWarElephant:    true,
			EnableCamel:          true,
			Alchemy:              true,
			Astrology:            true,
			Afterlife:            true,
			GuardTower:           true,
			Irrigation:           true,
			EnableFireBoat:       true,
		}
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
				// * Tower cost -50%.
				e.UnitStats[Tower].Cost.Multiply(0.5)
			},
			func(e *EmpireDeveloping) {
				// * Swordsmen attack speed +50% (stated 33%, they mean attack reload time).
			},
		}

	case Shang:
		c.Name, c.Name2 = "Shang", "Shang"
		c.DisabledTechs = map[TechID]bool{
			EnableElephantArcher: true,
			EnableWarElephant:    true,
			Aristocracy:          true,
			Ballistics:           true,
			Alchemy:              true,
			Engineering:          true,
			Coinage:              true,
			Siegecraft:           true,
			Phalanx:              true,
			LongSword:            true,
			Trireme:              true,
			CatapultTrireme:      true,
			HeavyTransport:       true,
		}
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
			ImprovedBow:     true,
			EnableCavalry:   true,
			Craftsmanship:   true,
			Metallurgy:      true,
			IronShield:      true,
			Coinage:         true,
			Astrology:       true,
			Afterlife:       true,
			Monotheism:      true,
			Fanaticism:      true,
			Zealotry:        true,
			EnableBallista:  true,
			CatapultTrireme: true,
			HeavyTransport:  true,
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
			EnableChariotArcher:  true,
			EnableElephantArcher: true,
			EnableChariot:        true,
			EnableWarElephant:    true,
			EnableCamel:          true,
			Astrology:            true,
			Mysticism:            true,
			Medicine:             true,
			Monotheism:           true,
			Fanaticism:           true,
			Zealotry:             true,
			Sacrifice:            true,
			Catapult:             true,
			EnableBallista:       true,
			BroadSword:           true,
			GuardTower:           true,
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
