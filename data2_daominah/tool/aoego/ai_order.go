// Package aoego is used to define the data structure in the game
// Age of Empires: The Rise of Rome.
// This can be used to validate and calculate the cost of
// build order (".ai" file), which defines computer player's strategy.
package aoego

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNotImplemented = errors.New("not implemented")

	ErrNotEnoughStepFields     = errors.New("not enough fields in a step")
	ErrInvalidAction           = errors.New("invalid action, check action enum list")
	ErrInvalidActionOrTargetID = errors.New("invalid action and targetID")
	ErrTargetIDNotInt          = errors.New("invalid unitID or techID, must be an integer")
	ErrQuantityNotInt          = errors.New("invalid quantity, must be an integer")
	ErrLocationNotInt          = errors.New("invalid location, must be an integer")
	ErrLocationNotMatch        = errors.New("location not match with the unit or tech")
	ErrLimitRebuildNotInt      = errors.New("invalid limit rebuild, must be an integer")

	ErrUnitIDNotFound   = errors.New("unitID not found")
	ErrTechIDNotFound   = errors.New("techID not found")
	ErrResearchQuantity = errors.New("research quantity must be 1")

	ErrUnitDisabledByCiv    = errors.New("unit is disabled by civilization")
	ErrUnitNotAvailableYet  = errors.New("unit is not available yet")
	ErrUnitLocationNotBuilt = errors.New("location is not built yet")
	ErrTechDisabledByCiv    = errors.New("technology is disabled by civilization")
	ErrTechNotAvailableYet  = errors.New("technology is not available yet")
)

// NewBuildOrder creates a BuildOrder from a ".ai" file format,
// this parses the file line by line, returns error for the first invalid line if any.
func NewBuildOrder(aiFileData string) ([]Step, error) {
	aiFileData = strings.ReplaceAll(aiFileData, "\r\n", "\n")
	lines := strings.Split(aiFileData, "\n")
	var buildOrder []Step
	for i, line := range lines {
		step, err := NewStep(line)
		if err != nil {
			return nil, fmt.Errorf("line %v: %w", i+1, err)
		}
		buildOrder = append(buildOrder, *step)
	}
	return buildOrder, nil
}

type Step struct {
	Action       Action
	UnitOrTechID UnitOrTechID
	Quantity     int
	LimitRebuild int // number times retrain if unit is destroyed, only meaningful if Action is BuildLimit or TrainLimit
}

// NewStep parses a line in ".ai" file format to a Step object.
// e.g "T299      Soldier-Scout        1      101       2" means
// train 1 Scout at Stable, if killed, retrain max 2 times.
// This func is the inverse function of Step.String().
func NewStep(line string) (*Step, error) {
	// workaround for the exceptional name with space  in `Default.ai` file
	line = strings.ReplaceAll(line, "Armored Elephants", "Armored_Elephants")

	words := strings.Fields(line)
	if len(words) < 4 {
		return nil, ErrNotEnoughStepFields
	}
	if len(words[0]) < 2 {
		return nil, fmt.Errorf("%v: %w", words[0], ErrInvalidActionOrTargetID)
	}
	unitOrTechIDStr := words[0][1:]
	unitOrTechID, err := strconv.Atoi(unitOrTechIDStr)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", unitOrTechIDStr, ErrTargetIDNotInt)
	}
	s := &Step{Action: Action(words[0][:1])}
	unitOrTech, err := s.determineUnitOrTech(unitOrTechID)
	if err != nil {
		return nil, fmt.Errorf("determineUnitOrTech: %w", err)
	}
	s.UnitOrTechID = unitOrTech.GetID()
	s.Quantity, err = strconv.Atoi(words[2])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", words[2], ErrQuantityNotInt)
	}
	locationID, err := strconv.Atoi(words[3])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", words[3], ErrLocationNotInt)
	}
	if unitOrTech.GetLocation().IntID() != locationID {
		return nil, ErrLocationNotMatch
	}
	if s.Action == BuildLimit || s.Action == TrainLimit {
		if len(words) < 5 {
			return nil, fmt.Errorf("missing limit times: %w", ErrNotEnoughStepFields)
		}
		s.LimitRebuild, err = strconv.Atoi(words[4])
		if err != nil {
			return nil, fmt.Errorf("%v: %w", words[4], ErrLimitRebuildNotInt)
		}
	}
	return s, nil
}

// String returns a string representation of a Step in ".ai" file format,
// e.g "T299      Soldier-Scout        1      101       2" means
// train 1 Scout at Stable, if killed, retrain max 2 times.
// This func is the inverse function of NewStep().
func (s Step) String() (string, error) {
	target, err := s.determineUnitOrTech(s.UnitOrTechID.IntID())
	if err != nil {
		return "", fmt.Errorf("determineUnitOrTech: %w", err)
	}
	internalName := target.GetName()
	if target.IsUnit() {
		internalName = "  " + internalName
	} else { // research quantity should always be 1
		if s.Quantity != 1 {
			return "", ErrResearchQuantity
		}
	}
	line := fmt.Sprintf("%v%-7v%-23v%-7v%-10v",
		s.Action, s.UnitOrTechID, internalName, s.Quantity, target.GetLocation())
	if s.Action == BuildLimit || s.Action == TrainLimit {
		line += fmt.Sprintf("%v", s.LimitRebuild)
	} else {
		line = strings.TrimSpace(line)
	}
	return line, nil
}

func (s Step) CheckEqual(other Step) bool {
	if s.Action != other.Action {
		return false
	}
	if s.UnitOrTechID.IntID() != other.UnitOrTechID.IntID() {
		return false
	}
	if s.Quantity != other.Quantity {
		return false
	}
	if s.LimitRebuild != other.LimitRebuild {
		return false
	}
	return true
}

// determineUnitOrTech use the Action to determine if the input arg is a UnitID or a TechID,
// this func does not use Step.UnitOrTechID so the field can be nil when calling this func.
func (s Step) determineUnitOrTech(unitOrTechID int) (UnitOrTech, error) {
	if s.Action == Research || s.Action == ResearchCritical {
		t, err := NewTechnology(TechID(unitOrTechID))
		if err != nil {
			return nil, ErrTechIDNotFound
		}
		return t, nil
	}
	if s.Action == Build || s.Action == BuildLimit ||
		s.Action == Train || s.Action == TrainLimit {
		u, err := NewUnit(UnitID(unitOrTechID))
		if err != nil {
			return nil, ErrUnitIDNotFound
		}
		return u, nil
	}
	return nil, ErrInvalidAction
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

// UnitOrTechID is a UnitID or TechID
type UnitOrTechID interface{ IntID() int }

type UnitOrTech interface {
	IsUnit() bool
	GetID() UnitOrTechID
	GetName() string // name without spaces
	GetLocation() UnitID
	GetCost() Cost
}

// EmpireDeveloping represents state of an empire at a moment,
// This can be used to store state of a running BuildOrder.
// MUST be initialized by func NewEmpireDeveloping.
type EmpireDeveloping struct {
	Civilization    Civilization
	AutoBuildHouse  int // if reach population limit, automatically build this number of houses
	Spent           Cost
	UnitStats       map[UnitID]*Unit // excluding the disabled units of the civilization
	EnabledUnits    map[UnitID]bool
	Combatants      map[UnitID]int  // trained units are not buildings
	Buildings       map[UnitID]int  // built buildings
	Techs           map[TechID]bool // researched technologies, including auto-researched
	TechnologyCount int             // only count the techs that are not auto-researched
}

func NewEmpireDeveloping(options ...EmpireOption) (*EmpireDeveloping, error) {
	fullTechTreeCiv, err := NewCivilization(FullTechTree)
	if err != nil { // unlikely to happen
		return nil, fmt.Errorf("NewCivilization FullTechTree: %w", err)
	}
	e := &EmpireDeveloping{
		Civilization:   *fullTechTreeCiv,
		AutoBuildHouse: 5,
		UnitStats:      make(map[UnitID]*Unit),
		Combatants:     map[UnitID]int{Villager: 3},
		Buildings:      map[UnitID]int{TownCenter: 1},
		Techs:          make(map[TechID]bool),
	}

	for _, option := range options {
		option(e)
	}

	for unitID := range AllUnits {
		if e.Civilization.DisabledUnits[unitID] {
			continue
		}
		u, err := NewUnit(unitID)
		if err != nil { // unlikely to happen
			continue
		}
		e.UnitStats[unitID] = u
	}

	for _, v := range e.Civilization.Bonuses {
		v(e) // TODO: implement Civilization Bonuses
	}
	return e, nil
}

// Do tries to execute a Step (probably from a BuildOrder),
// it will return error if the Step is invalid, e.g. technology is not available
// for the civilization or not satisfied the requirement,
// unit's location is not built yet, ...
func (e *EmpireDeveloping) Do(s Step) error {
	target, err := s.determineUnitOrTech(s.UnitOrTechID.IntID())
	if err != nil {
		return fmt.Errorf("determineUnitOrTech: %w", err)
	}
	switch v := target.(type) {
	case *Unit:
		if _, found := e.UnitStats[v.ID]; !found {
			return fmt.Errorf("%w: %v is disabled for civilization %v", ErrUnitDisabledByCiv, v.NameInGame, e.Civilization.Name)
		}
		if v.Location != NullUnit {
			if !(e.Buildings[v.Location] > 0) {
				return fmt.Errorf("%w: need %v first to train %v", ErrUnitLocationNotBuilt, v.Location, v.NameInGame)
			}
		}

	case *Technology:
		return ErrNotImplemented
	default:
		return fmt.Errorf("invalid target type: %T", target)
	}
	return ErrNotImplemented
}

func (e *EmpireDeveloping) CountPopulation() int {
	pop := float64(0)
	for unitID, count := range e.Combatants {
		unit, found := e.UnitStats[unitID]
		if !found { // unlikely to happen
			continue
		}
		pop += float64(count) * unit.Population
	}
	return int(pop)
}

func (e *EmpireDeveloping) build(unitID UnitID, quantity int) error {
	//unit, found := e.UnitStats[unitID]
	//if !found {
	//	return fmt.Errorf("%w: %v is disabled for civilization %v", ErrUnitDisabledByCiv, UnitName(unitID), e.Civilization.Name)
	//}
	return ErrNotImplemented
}

func (e *EmpireDeveloping) autoBuildHouse() error {
	if e.AutoBuildHouse <= 0 {
		return nil
	}
	return ErrNotImplemented
}

// EmpireOption can set civilization and choose to
type EmpireOption func(*EmpireDeveloping)

func WithCivilization(civID CivilizationID) EmpireOption {
	return func(e *EmpireDeveloping) {
		tmp, err := NewCivilization(civID)
		if err != nil {
			return
		}
		e.Civilization = *tmp
	}
}

func WithNoUnit() EmpireOption {
	return func(e *EmpireDeveloping) {
		e.Buildings[TownCenter] = 0
		e.Combatants[Villager] = 0
	}
}

func WithDisableAutoBuildHouse() EmpireOption {
	return func(e *EmpireDeveloping) {
		e.AutoBuildHouse = 0
	}
}
