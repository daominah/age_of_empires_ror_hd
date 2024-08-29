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

	case Yamato:
		c.Name = "Yamato"

	default:
		return nil, fmt.Errorf("civID %v: %w", civID, ErrInvalidCivID)
	}
	return c, nil
}
