// Package aoego is used to define the data structure in the game
// Age of Empires: The Rise of Rome.
// This can be used to validate and calculate the cost of
// build order (".ai" file), which defines computer player's strategy.
package aoego

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrNotImplemented = errors.New("not implemented")

	ErrEmptyLine               = errors.New("line is empty or a comment")
	ErrMissingStepFields       = errors.New("not enough fields in a step")
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

	ErrLocationNotBuilt    = errors.New("location is not built yet")
	ErrUnitDisabledByCiv   = errors.New("unit is disabled by civilization")
	ErrTechDisabledByCiv   = errors.New("technology is disabled by civilization")
	ErrMissingRequireTechs = errors.New("missing required techs")
	ErrTechResearched      = errors.New("technology is already researched")
)

type BuildOrder []Step

// NewBuildOrder creates a BuildOrder from a ".ai" file format,
// this parses the file line by line, returns error for the first invalid line if any.
func NewBuildOrder(aiFileData string) ([]Step, []error) {
	aiFileData = strings.ReplaceAll(aiFileData, "\r\n", "\n")
	lines := strings.Split(aiFileData, "\n")
	var buildOrder []Step
	var errs []error
	for i, line := range lines {
		step, err := NewStep(line)
		if err != nil {
			if !errors.Is(err, ErrEmptyLine) {
				errs = append(errs, fmt.Errorf("line %v: err: %w, line: %v", i+1, err, line))
			}
			continue
		}
		buildOrder = append(buildOrder, *step)
	}
	return buildOrder, errs
}

func (b BuildOrder) Marshal() (string, error) {
	var lines []string
	for i, step := range b {
		line, err := step.String()
		if err != nil {
			return "", fmt.Errorf("line %v step.String: %w", i+1, err)
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n"), nil
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
	commentBegin := strings.Index(line, `//`)
	if commentBegin != -1 {
		line = line[:commentBegin]
	}
	line = strings.TrimSpace(line)

	// workaround for the exceptional name with space  in `Default.ai` file
	line = strings.ReplaceAll(line, "Armored Elephants", "Armored_Elephants")

	words := strings.Fields(line)
	if len(words) == 0 {
		return nil, ErrEmptyLine
	}
	if len(words) < 4 {
		return nil, ErrMissingStepFields
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
			return nil, fmt.Errorf("missing limit times: %w", ErrMissingStepFields)
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
	internalName = strings.ReplaceAll(internalName, " ", "_")
	if target.IsUnit() {
		internalName = "  " + internalName
	} else { // research quantity should always be 1
		if s.Quantity != 1 {
			return "", ErrResearchQuantity
		}
	}
	if len(internalName) > 22 {
		internalName = internalName[:22]
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
type UnitOrTechID interface {
	IntID() int
	GetNameInGame() string // name with spaces, easier to read
}

type UnitOrTech interface {
	IsUnit() bool
	GetID() UnitOrTechID
	GetName() string // name without spaces
	GetLocation() UnitID
}

// EmpireDeveloping represents state of an empire at a moment,
// This can be used to store state of a running BuildOrder.
// MUST be initialized by func NewEmpireDeveloping.
type EmpireDeveloping struct {
	Civilization     Civilization
	IsAutoBuildHouse bool             // build 5 houses if close to population limit
	UnitStats        map[UnitID]*Unit // excluding the disabled units of the civilization
	EnabledUnits     map[UnitID]bool  // example research Wheel add Chariot and Chariot Archer to this map
	Combatants       map[UnitID]int   // trained units are not buildings
	Buildings        map[UnitID]int   // built buildings
	Techs            map[TechID]bool  // researched technologies, including auto-researched
	TechnologyCount  int              // only count the techs that are not auto-researched
	Spent            Cost
}

func NewEmpireDeveloping(options ...EmpireOption) (*EmpireDeveloping, error) {
	fullTechTreeCiv, err := NewCivilization(FullTechTree)
	if err != nil { // unlikely to happen
		return nil, fmt.Errorf("NewCivilization FullTechTree: %w", err)
	}
	e := &EmpireDeveloping{
		Civilization:     *fullTechTreeCiv,
		IsAutoBuildHouse: true,
		UnitStats:        make(map[UnitID]*Unit),
		EnabledUnits: map[UnitID]bool{
			TownCenter: true, Villager: true, House: true,
			Granary: true, StoragePit: true, Barracks: true, Dock: true,
		},
		Combatants: map[UnitID]int{},
		Buildings:  map[UnitID]int{Granary: 1, StoragePit: 1}, // AI will always build these 2 buildings without instruction
		Techs:      map[TechID]bool{GranaryBuilt: true, StoragePitBuilt: true},
		Spent:      Cost{Wood: 240},
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
		return e.build(*v, s.Quantity)
	case *Technology:
		return e.research(*v)
	default:
		return fmt.Errorf("invalid target type: %T", target)
	}
}

func (e *EmpireDeveloping) CountPopulation() float64 {
	pop := float64(0)
	for unitID, count := range e.Combatants {
		unit, found := e.UnitStats[unitID]
		if !found { // unlikely to happen
			continue
		}
		pop += float64(count) * unit.Population
	}
	return pop
}

func (e *EmpireDeveloping) CountPopulationLimit() float64 {
	popLimit := float64(0)
	for unitID, count := range e.Buildings {
		if unitID == House || unitID == TownCenter {
			popLimit += 4 * float64(count)
		}
	}
	return popLimit
}

// Build a combatant(s) or a building(s). if quantity is not provided, build 1.
func (e *EmpireDeveloping) build(u Unit, quantity ...int) error {
	n := 1 // number of units to build
	if len(quantity) > 0 {
		n = quantity[0]
	}
	if _, found := e.UnitStats[u.ID]; !found {
		return fmt.Errorf("%w: %v is disabled for %v", ErrUnitDisabledByCiv, u.NameInGame, e.Civilization.Name)
	}

	// cannot build this unit: return why

	if !e.EnabledUnits[u.ID] {
		enableTechID, found := UnitEnabledByTechs[u.ID]
		if !found { // unlikely to happen
			return fmt.Errorf("%w for %v, missing key %v in UnitEnabledByTechs", ErrMissingRequireTechs, u.NameInGame, u.ID)
		}
		enableTech, found := AllTechs[enableTechID]
		if !found { // unlikely to happen
			return fmt.Errorf("%w for %v, missing key %v in AllTechs", ErrMissingRequireTechs, u.NameInGame, enableTechID)
		}
		missingTechsCount := enableTech.MinRequiredTechs
		var missingTechs []string
		for _, techID := range enableTech.RequiredTechs {
			if e.Techs[techID] {
				missingTechsCount--
				continue
			}
			missingTechs = append(missingTechs, techID.GetNameInGame())
		}
		if missingTechsCount == 0 { // unlikely to happen
			return fmt.Errorf("not miss required techs but %v(%v) is still disabled, check effect funcs: %v", u.NameInGame, u.ID.ActionID(), enableTech.GetEffectsName())
		}
		return fmt.Errorf("%w for %v: at least %v of %v", ErrMissingRequireTechs, u.NameInGame, missingTechsCount, strings.Join(missingTechs, ", "))
	}
	if u.Location != NullUnit {
		if e.Buildings[u.Location] <= 0 {
			return fmt.Errorf("%w: need %v first to train %v", ErrLocationNotBuilt, u.Location, u.NameInGame)
		}
	}

	// can build this unit: update the empire

	e.Spent.Add(*(u.Cost.Multiply(float64(n))))
	if u.IsBuilding {
		if e.Buildings[u.ID] == 0 && u.InitiateTech != NullTech {
			e.Techs[u.InitiateTech] = true
			for _, effect := range AllTechs[u.InitiateTech].Effects {
				effect(e)
			}
			e.refreshAutoTechs()
		}
		e.Buildings[u.ID] += n
	} else {
		e.Combatants[u.ID] += n
		if e.IsAutoBuildHouse {
			remainingHousesBlock := (e.CountPopulationLimit() - e.CountPopulation()) / 20.0
			if remainingHousesBlock <= 0.5 {
				e.buildHouse(5 * int(math.Ceil(1-remainingHousesBlock)))
			}
		}
	}
	return nil
}

func (e *EmpireDeveloping) buildHouse(n int) {
	house := AllUnits[House]
	e.Spent.Add(*(house.Cost.Multiply(float64(n))))
	e.Buildings[House] += n
}

func (e *EmpireDeveloping) research(t Technology) error {
	if e.Civilization.DisabledTechs[t.ID] {
		return fmt.Errorf("%w: %v is disabled for %v", ErrTechDisabledByCiv, t.NameInGame, e.Civilization.Name)
	}
	if e.Techs[t.ID] {
		return fmt.Errorf("%w: %v", ErrTechResearched, t.NameInGame)
	}
	missingTechsCount := t.MinRequiredTechs
	var missingTechs []string
	for _, techID := range t.RequiredTechs {
		if e.Techs[techID] {
			missingTechsCount--
			continue
		}
		missingTechs = append(missingTechs, techID.GetNameInGame())
	}
	if missingTechsCount > 0 {
		return fmt.Errorf("%w for %v: at least %v of %v", ErrMissingRequireTechs, t.NameInGame, missingTechsCount, strings.Join(missingTechs, ", "))
	}

	e.Spent.Add(t.Cost)
	e.Techs[t.ID] = true
	if !t.Cost.IsZero() {
		e.TechnologyCount++
	}
	for _, effect := range t.Effects {
		effect(e)
	}
	return nil
}

func (e *EmpireDeveloping) refreshAutoTechs() {
	for autoID, autoTech := range AllAutoTechs {
		if e.Techs[autoID] {
			continue
		}
		isSatisfied := true
		for _, require := range autoTech.RequiredTechs {
			if !e.Techs[require] {
				isSatisfied = false
				break
			}
		}
		if isSatisfied {
			e.Techs[autoID] = true
			for _, effect := range autoTech.Effects {
				effect(e)
			}
		}
	}
}

func (e *EmpireDeveloping) Summary() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("civilization: %v", e.Civilization.Name))
	lines = append(lines, fmt.Sprintf("spent: %+v", e.Spent))
	lines = append(lines, fmt.Sprintf("population: %.0f/%.0f", e.CountPopulation(), e.CountPopulationLimit()))
	lines = append(lines, fmt.Sprintf("buildings: %v", beautyUnits(e.Buildings)))
	lines = append(lines, fmt.Sprintf("combatants: %v", beautyUnits(e.Combatants)))
	lines = append(lines, fmt.Sprintf("techs count: %v", e.TechnologyCount))
	lines = append(lines, fmt.Sprintf("techs researched: %+v", beautyTechs(e.Techs)))
	return strings.Join(lines, "\n")
}

func beautyUnits(m map[UnitID]int) string {
	var keyValues []string
	for unitID, count := range m {
		keyValues = append(keyValues, fmt.Sprintf("%v(%v): %v", unitID.GetNameInGame(), unitID.ActionID(), count))
	}
	return strings.Join(keyValues, ", ")
}

func beautyTechs(m map[TechID]bool) string {
	var words []string
	for techID, researched := range m {
		if researched {
			words = append(words, fmt.Sprintf("%v(R%v)", techID.GetNameInGame(), techID))
		}
	}
	return strings.Join(words, ", ")
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
		e.Buildings = map[UnitID]int{}
		e.Techs = map[TechID]bool{}
		e.Spent = Cost{}
	}
}

func WithDisableAutoBuildHouse() EmpireOption {
	return func(e *EmpireDeveloping) {
		e.IsAutoBuildHouse = false
	}
}
