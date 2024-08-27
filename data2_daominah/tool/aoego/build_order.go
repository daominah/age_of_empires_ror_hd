// Package aoego is used to define the data structure in the game
// Age of Empires: The Rise of Rome.
// This can be used to validate and calculate the cost of
// build order (".ai" file), which defines computer player's strategy.
package aoego

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotImplemented     = errors.New("not implemented")
	ErrLocationNotMatched = errors.New("parent building does not matched with unit or tech")
	ErrUnitIDNotFound     = errors.New("unitID not found")
	ErrTechIDNotFound     = errors.New("techID not found")
	ErrResearchQuantity   = errors.New("research quantity must be 1")
)

// BuildOrder is equivalent to a ".ai" file, each step is a line in the file.
type BuildOrder struct {
	Civilization CivilizationID
	Steps        []Step
}

type Step struct {
	Action       Action
	UnitOrTechID UnitOrTechID
	Quantity     int
	Location     UnitID // parent building that trains unit or researches technology
	LimitRebuild int    // number times retrain if unit is destroyed, only meaningful if Action is BuildLimit or TrainLimit
}

// String returns a string representation of a Step in ".ai" file format,
// e.g "T299      Soldier-Scout        1      101       2" means
// train 1 Scout at Stable, if killed, retrain max 2 times
func (s Step) String() (string, error) {
	target, err := s.determineTarget()
	if err != nil {
		return "", fmt.Errorf("determineTarget: %w", err)
	}
	if target.GetLocation() != s.Location {
		return "", ErrLocationNotMatched
	}
	internalName := target.GetNameInternal()
	if target.IsUnit() {
		internalName = "  " + internalName
	}
	line := fmt.Sprintf("%v%-7v%-23v%-7v%-10v",
		s.Action, s.UnitOrTechID, internalName, s.Quantity, s.Location)
	if s.Action == BuildLimit || s.Action == TrainLimit {
		line += fmt.Sprintf("%v", s.LimitRebuild)
	} else {
		line = strings.TrimSpace(line)
	}
	return line, nil
}

func (s Step) determineTarget() (UnitOrTech, error) {
	if s.Action == Research || s.Action == ResearchCritical {
		t := NewTechnology(TechID(s.UnitOrTechID))
		if t == nil {
			return nil, ErrTechIDNotFound
		}
		return t, nil
	}
	u := NewUnit(UnitID(s.UnitOrTechID))
	if u == nil {
		return nil, ErrUnitIDNotFound
	}
	return u, nil
}

// Action is enum: build or train or research
type Action string

// Action enum
const (
	Build            Action = "B" // will rebuild if destroyed
	BuildLimit       Action = "A" // will rebuild up to the limit times specified if destroyed
	Train            Action = "U" // will always be replaced if killed
	TrainLimit       Action = "T" // will be trained up to the limit times specified
	Research         Action = "R" // item will be researched if possible, can be skipped if failed too many attempts
	ResearchCritical Action = "C" // the build stuck until this research is done
)

type UnitOrTechID int // UnitOrTechID is a UnitID or TechID

type UnitOrTech interface {
	IsUnit() bool
	GetID() UnitOrTechID
	GetNameInternal() string
	GetLocation() UnitID
	GetCost() Storage
}

// EmpireDeveloping represents state of a player's empire at a moment.
// This can be used to store state of a running BuildOrder
type EmpireDeveloping struct {
	Civilization   CivilizationID
	TechsDisabled  map[TechID]bool
	IsFullTechTree bool // if true, Civilization and TechsDisabled don't matter

	Storage    Storage
	Combatants []Unit          // trained units are not buildings
	Buildings  []Unit          // built buildings
	Techs      map[TechID]bool // researched technologies
}

// Do tries to execute a Step (probably from a BuildOrder),
// it will return error if the Step is invalid, e.g. not enough resources,
// technology is not available, unit is not trained by the building, ...
func (e *EmpireDeveloping) Do(s Step) error {
	return ErrNotImplemented
}
