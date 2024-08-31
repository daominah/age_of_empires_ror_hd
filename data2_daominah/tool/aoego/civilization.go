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
			ImprovedBow:   true,
			Metallurgy:    true,
			Astrology:     true,
			EnableCavalry: true,

			IronShield:     true,
			Craftsmanship:  true,
			Coinage:        true,
			Afterlife:      true,
			Monotheism:     true,
			Fanaticism:     true,
			Zealotry:       true,
			EnableBallista: true,

			HeavyTransport:  true,
			CatapultTrireme: true,
		}
		c.DisabledUnits = map[UnitID]bool{
			Cavalry:  true,
			Ballista: true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// Siege units attack speed +43%.
			},
			func(e *EmpireDeveloping) {
				// Villager HP +15 (so 40 instead of 25).
			},
			func(e *EmpireDeveloping) {
				// Farm food +250 (starting at 500 instead of 250).
			},
		}

	case Yamato:
		c.Name = "Yamato"
		c.DisabledTechs = map[TechID]bool{
			EnableChariotArcher:  true,
			EnableChariot:        true,
			EnableCamel:          true,
			Broadsword:           true,
			Astrology:            true,
			Mysticism:            true,
			Medicine:             true,
			Monotheism:           true,
			Fanaticism:           true,
			Zealotry:             true,
			Sacrifice:            true,
			GuardTower:           true,
			EnableBallista:       true,
			Catapult:             true,
			FortifiedWall:        true,
			EnableElephantArcher: true,
			EnableWarElephant:    true,
			EnableFireBoat:       true,
		}
		c.Bonuses = []EffectFunc{
			func(e *EmpireDeveloping) {
				// Villager move speed +18%
			},
			func(e *EmpireDeveloping) {
				// Mounted units cost -25%
				e.UnitStats[Scout].Cost.Multiply(0.75)
				e.UnitStats[Cavalry].Cost.Multiply(0.75)
				e.UnitStats[HorseArcher].Cost.Multiply(0.75)
			},
			func(e *EmpireDeveloping) {
				// Ships HP +30%
			},
		}

	default:
		return nil, fmt.Errorf("civID %v: %w", civID, ErrInvalidCivID)
	}
	return c, nil
}
