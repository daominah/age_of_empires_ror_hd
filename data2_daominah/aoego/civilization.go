package aoego

import (
	"errors"
	"fmt"
)

var ErrInvalidCivID = errors.New("invalid civilization ID, check civilization enum list")

// Civilization must be initialized with func NewCivilization
type Civilization struct {
	ID            CivilizationID
	Name          string
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
		c.Name = "FullTechTree"
		return c, nil

	case Assyrian:
		c.Name = "Assyrian"

	case Babylonian:
		c.Name = "Babylonian"

	case Carthaginian:
		c.Name = "Carthaginian"

	case Choson:
		c.Name = "Choson"

	case Egyptian:
		c.Name = "Egyptian"

	case Greek:
		c.Name = "Greek"

	case Hittite:
		c.Name = "Hittite"

	case Macedonian:
		c.Name = "Macedonian"
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
		c.Name = "Minoan"

	case Palmyran:
		c.Name = "Palmyran"

	case Persian:
		c.Name = "Persian"

	case Phoenician:
		c.Name = "Phoenician"

	case Roman:
		c.Name = "Roman"

	case Shang:
		c.Name = "Shang"

	case Sumerian:
		c.Name = "Sumerian"
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
		c.Name = "Yamato"
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
