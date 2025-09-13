// Package aoego is used to define the data structure in the game
// Age of Empires: The Rise of Rome.
// This can be used to validate and calculate the cost of a strategy (".ai" file),
// which defines computer player's order to build and train units.
package aoego

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrEmptyLine               = errors.New("line is empty or a comment")
	ErrMissingStepFields       = errors.New("not enough fields in a step")
	ErrRedundantFields         = errors.New("too many fields in a step")
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
	ErrUnitDisabledByCiv   = errors.New("unit disabled by civ")
	ErrTechDisabledByCiv   = errors.New("tech disabled by civ")
	ErrMissingRequireTechs = errors.New("missing required techs")
	ErrTechResearched      = errors.New("technology is already researched")
	ErrExceedPopLimit      = errors.New("build units will exceed population limit, build more houses first")
)

type ErrorWithLineNo struct {
	LineNo        int // line number in the original file, that caused the error
	Err           error
	IsJustWarning bool // if true, the error can be optionally ignored
}

func (e ErrorWithLineNo) Error() string {
	if !e.IsJustWarning {
		return fmt.Sprintf("error at line %-3v: %v", e.LineNo, e.Err)
	} else {
		return fmt.Sprintf("warn  at line %-3v: %v", e.LineNo, e.Err)
	}
}

type SortByLineNo []ErrorWithLineNo

func (a SortByLineNo) Len() int           { return len(a) }
func (a SortByLineNo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByLineNo) Less(i, j int) bool { return a[i].LineNo < a[j].LineNo }

// Strategy defines order to build and train units.
type Strategy []Step

// NewStrategy creates a Strategy from a ".ai" file format,
// this parses the file line by line, returns error for the first invalid line if any.
func NewStrategy(aiFileData string) ([]Step, []ErrorWithLineNo) {
	aiFileData = strings.ReplaceAll(aiFileData, "\r\n", "\n")
	lines := strings.Split(aiFileData, "\n")
	var strategy []Step
	var errs []ErrorWithLineNo
	for i, line := range lines {
		step, err := NewStep(line)
		if err != nil {
			if !errors.Is(err, ErrEmptyLine) {
				errs = append(errs, ErrorWithLineNo{
					LineNo: i + 1,
					Err:    fmt.Errorf("%-60v: %w", line, err),
				})
			}
			continue
		}
		step.OriginLineNo = i + 1
		strategy = append(strategy, *step)
	}
	return strategy, errs
}

func (b Strategy) Marshal() (string, error) {
	var lines []string
	for i, step := range b {
		line, err := step.Marshal()
		if err != nil {
			return "", fmt.Errorf("line %v step.Marshal: %w", i+1, err)
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
	OriginStr    string
	OriginLineNo int // line number in the original file, not handled in func NewStep
}

// NewStep parses a line in ".ai" file format to a Step object.
// e.g "T299      Soldier-Scout        1      101       2" means
// train 1 Scout at Stable, if killed, retrain max 2 times.
// This func is the inverse function of Step.Marshal().
func NewStep(line string) (*Step, error) {
	originalLine := line
	commentBegin := strings.Index(line, `//`)
	if commentBegin != -1 {
		if strings.HasPrefix(line, `// spent`) {
			return &Step{Action: PrintSummary, OriginStr: originalLine}, nil
		}
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
	s := &Step{Action: Action(words[0][:1]), OriginStr: originalLine}
	unitOrTech, err := s.determineUnitOrTech(unitOrTechID)
	if err != nil {
		return nil, fmt.Errorf("determineUnitOrTech(%v): %w", unitOrTechID, err)
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
			return nil, fmt.Errorf("missing limit times for action %v: %w", s.Action, ErrMissingStepFields)
		}
		s.LimitRebuild, err = strconv.Atoi(words[4])
		if err != nil {
			return nil, fmt.Errorf("%v: %w", words[4], ErrLimitRebuildNotInt)
		}
	} else {
		if len(words) > 4 {
			return nil, fmt.Errorf("too many fields for action %v: %w", s.Action, ErrRedundantFields)
		}
	}
	return s, nil
}

// Marshal returns a string representation of a Step in ".ai" file format,
// e.g "T299      Soldier-Scout        1      101       2" means
// train 1 Scout at Stable, if killed, retrain max 2 times.
// This func is the inverse function of NewStep().
func (s Step) Marshal() (string, error) {
	if s.Action == PrintSummary {
		return "// spent:", nil
	}
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

// Strings is for debugging purpose.
func (s Step) String() string {
	if s.Action == PrintSummary {
		return "// spent:"
	}
	target, err := s.determineUnitOrTech(s.UnitOrTechID.IntID())
	if err != nil { // unlikely to happen
		return fmt.Sprintf("%#v", s)
	}
	return fmt.Sprintf("%v %v%v %v", target.GetName(), s.Action, s.UnitOrTechID, s.Quantity)
}

func (s Step) checkEqual(other Step) bool {
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

	PrintSummary Action = "P" // special action treats a line has prefix "// spent" as a Step
)

// UnitOrTechID is a UnitID or TechID
type UnitOrTechID interface {
	IntID() int
	GetNameInGame() string // name with spaces, easier to read
	GetAge() TechID
}

type UnitOrTech interface {
	IsUnit() bool
	GetID() UnitOrTechID
	GetName() string     // name without spaces, e.g. "Town_Center1"
	GetFullName() string // e.g. "Town Center(B109)", for print debug purpose
	GetLocation() UnitID
	GetCost() *Cost // clone value of the cost, for easily chain Multi without changing the original cost
}

// EmpireDeveloping represents state of an empire at a moment,
// This can be used to store state of a running Strategy.
// MUST be initialized by func NewEmpireDeveloping.
type EmpireDeveloping struct {
	Civilization     Civilization
	IsAutoBuildHouse bool             // build 5 houses if close to population limit
	UnitStats        map[UnitID]*Unit // excluding the disabled units of the civilization

	EnabledUnits    map[UnitID]bool // example research Wheel add Chariot and Chariot Archer to this map
	Combatants      map[UnitID]int  // trained units are not buildings
	Buildings       map[UnitID]int  // built buildings
	Techs           map[TechID]bool // researched technologies, including auto-researched
	TechnologyCount int             // only count the techs that are not auto-researched
	Spent           *Cost
	// FreeUnits for example 1st Town Center and 3 Villagers when every game starts,
	// if there are already units that are also in the Strategy,
	// these units will be counted as already created and will not be rebuilt,
	// the cost of these units will be subtracted from Spent.
	FreeUnits map[UnitID]int

	// for testing Spent, allow to build or research without checking required
	// techs or location
	IsIgnoreRequiredTechOrBuilding bool
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
		FreeUnits:        map[UnitID]int{TownCenter: 1, Villager: 3},

		EnabledUnits: map[UnitID]bool{
			TownCenter: true, Villager: true, House: true,
			Granary: true, StoragePit: true, Barracks: true, Dock: true,
		},
		Combatants: map[UnitID]int{},
		Buildings:  map[UnitID]int{Granary: 1, StoragePit: 1}, // AI will always auto build Granary and Storage Pit
		Techs:      map[TechID]bool{GranaryBuilt: true, StoragePitBuilt: true},
		Spent:      &Cost{},
	}
	for _, option := range options {
		option(e)
	}

	disabledUnits := make(map[UnitID]bool)
	for unitID, techID := range UnitEnabledByTechs {
		if e.Civilization.DisabledTechs[techID] {
			disabledUnits[unitID] = true
		}
	}
	for unitID := range AllUnits {
		if disabledUnits[unitID] {
			continue
		}
		u, err := NewUnit(unitID)
		if err != nil { // unlikely to happen
			continue
		}
		e.UnitStats[unitID] = u
	}
	for _, v := range e.Civilization.Bonuses {
		v(e)
	}
	if e.UnitStats[Granary] != nil && e.UnitStats[StoragePit] != nil {
		granaries := e.UnitStats[Granary].GetCost().Multiply(float64(e.Buildings[Granary]))
		storagePits := e.UnitStats[StoragePit].GetCost().Multiply(float64(e.Buildings[StoragePit]))
		e.Spent.Add(*granaries).Add(*storagePits)
	}
	return e, nil
}

// Do tries to execute a Step (probably from a Strategy),
// it will return error if the Step is invalid, e.g. technology is not available
// for the civilization or not satisfied the requirement,
// unit's location is not built yet, ...
func (e *EmpireDeveloping) Do(s Step) error {
	if s.Action == PrintSummary {
		return nil
	}
	target, err := s.determineUnitOrTech(s.UnitOrTechID.IntID())
	if err != nil {
		return fmt.Errorf("determineUnitOrTech: %w", err)
	}
	switch v := target.(type) {
	case *Unit:
		return e.build(v.ID, s.Quantity)
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
func (e *EmpireDeveloping) build(unitID UnitID, quantity ...int) error {
	n := 1 // number of units to build
	if len(quantity) > 0 {
		n = quantity[0]
	}
	civUnit, found := e.UnitStats[unitID] // get unit stats applied civilization bonuses
	if !found {
		return fmt.Errorf("%w: %v is disabled for %v", ErrUnitDisabledByCiv, unitID.GetNameInGame(), e.Civilization.Name)
	}

	if !e.EnabledUnits[unitID] { // cannot build this unit: return why
		u, err := NewUnit(unitID)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrUnitIDNotFound, unitID)
		}
		enableTechID, found := UnitEnabledByTechs[unitID]
		if !found {
			if u.Location != NullUnit && e.Buildings[u.Location] <= 0 {
				return fmt.Errorf("%w: need %v(%v) first to train %v", ErrLocationNotBuilt, u.Location.GetNameInGame(), u.Location.ActionID(), u.GetFullName())
			}
			return fmt.Errorf("%w for %v, missing key %v in UnitEnabledByTechs", ErrMissingRequireTechs, u.NameInGame, u.ID)
		}
		enableTech, found := AllTechs[enableTechID]
		if !found { // unlikely to happen
			return fmt.Errorf("%w for %v, missing key %v in AllTechs", ErrMissingRequireTechs, u.NameInGame, enableTechID)
		}
		if !CheckIsAutoTech(enableTechID) {
			return fmt.Errorf("%w for %v, need to research %v first", ErrMissingRequireTechs, u.GetFullName(), enableTech.GetFullName())
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
			return fmt.Errorf("not miss required techs but %v is still disabled, check effect funcs: %v", u.GetFullName(), enableTech.GetEffectsName())
		}
		return fmt.Errorf("%w for %v: at least %v of %v", ErrMissingRequireTechs, u.NameInGame, missingTechsCount, strings.Join(missingTechs, ", "))
	}
	if civUnit.Location != NullUnit && e.Buildings[civUnit.Location] <= 0 {
		return fmt.Errorf("%w: need %v(%v) first to train %v(%v)", ErrLocationNotBuilt, civUnit.Location.GetNameInGame(), civUnit.Location.ActionID(), civUnit.NameInGame, unitID.ActionID())
	}

	if e.FreeUnits[unitID] > 0 {
		freeCountThisStep := min(n, e.FreeUnits[unitID])
		e.Spent.Add(*civUnit.GetCost().Multiply(-float64(freeCountThisStep)))
		e.FreeUnits[unitID] -= freeCountThisStep
	}

	if civUnit.IsBuilding {
		if e.Buildings[unitID] == 0 && civUnit.InitiateTech != NullTech {
			e.Techs[civUnit.InitiateTech] = true
			for _, effect := range AllTechs[civUnit.InitiateTech].Effects {
				effect(e)
			}
			e.refreshAutoTechs()
		}
		e.Spent.Add(*(civUnit.GetCost().Multiply(float64(n))))
		e.Buildings[unitID] += n
	} else {
		now := e.CountPopulation()
		need := float64(n) * civUnit.Population
		goodLimit := now + need + 10
		if e.IsAutoBuildHouse {
			for e.CountPopulationLimit() < goodLimit {
				e.buildHouse(5)
			}
		}
		if popCap := e.CountPopulationLimit(); popCap < now+need {
			return fmt.Errorf("%w: now %.0f/%.0f, need %.0f", ErrExceedPopLimit, now, popCap, need)
		}
		e.Spent.Add(*(civUnit.GetCost().Multiply(float64(n))))
		e.Combatants[unitID] += n
	}
	return nil
}

func (e *EmpireDeveloping) buildHouse(n int) {
	house := e.UnitStats[House]
	e.Spent.Add(*(house.GetCost().Multiply(float64(n))))
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
	e.refreshAutoTechs()
	return nil
}

func (e *EmpireDeveloping) refreshAutoTechs() {
	for autoID, autoTech := range AllAutoTechs {
		if e.Civilization.DisabledTechs[autoID] {
			continue
		}
		if e.Techs[autoID] {
			continue
		}
		nSatisfied := 0
		for _, require := range autoTech.RequiredTechs {
			if e.Techs[require] {
				nSatisfied += 1
			}
		}
		if nSatisfied >= autoTech.MinRequiredTechs {
			e.Techs[autoID] = true
			for _, effect := range autoTech.Effects {
				effect(e)
			}
		}
	}
}

func (e *EmpireDeveloping) Summary() string {
	var lines []string
	lines = append(lines, "//")
	lines = append(lines, "// "+e.beautyMainArmy())
	economy := fmt.Sprintf("// %v %v", e.Combatants[Villager], strings.ToLower(e.UnitStats[Villager].NameInGame))
	if e.Buildings[Farm] > 0 {
		economy += fmt.Sprintf(", %v farm", e.Buildings[Farm])
	}
	if e.Buildings[Tower] > 0 {
		economy += fmt.Sprintf(", %v tower", e.Buildings[Tower])
	}
	if e.Buildings[Wonder] > 0 {
		economy += fmt.Sprintf(", %v wonder", e.Buildings[Wonder])
	}
	economy += fmt.Sprintf(" (pop %.0f, tech %v)", e.CountPopulation(), e.TechnologyCount)
	lines = append(lines, economy)
	lines = append(lines, fmt.Sprintf("// spent: %+v", e.Spent))
	lines = append(lines, "//")
	lines = append(lines, fmt.Sprintf("// civilization: %v", e.Civilization.Name))
	lines = append(lines, fmt.Sprintf("// combatants: %v", beautyUnits(e.Combatants)))
	lines = append(lines, fmt.Sprintf("// buildings: %v", beautyUnits(e.Buildings)))
	lines = append(lines, fmt.Sprintf("// techs researched: %+v", beautyTechs(e.Techs)))
	return "\n" + strings.Join(lines, "\n") + "\n"
}

func (e *EmpireDeveloping) beautyMainArmy() string {
	type Pair struct {
		UnitID UnitID
		Count  int
	}
	var a []Pair
	for unit, count := range e.Combatants {
		if unit.GetAge() < BronzeAge {
			continue
		}
		a = append(a, Pair{UnitID: unit, Count: count})
	}
	sort.Slice(a, func(i, j int) bool {
		if a[i].Count == a[j].Count {
			return a[i].UnitID.GetAge() > a[j].UnitID.GetAge()
		}
		return a[i].Count > a[j].Count
	})
	var c strings.Builder // main combatants
	var b strings.Builder // main combatants location
	sameLocations := make(map[UnitID]bool)
	for i, pair := range a {
		upgradedUnit, found := e.UnitStats[pair.UnitID]
		if !found {
			continue
		}
		c.WriteString(fmt.Sprintf("%v %v", pair.Count, strings.ToLower(upgradedUnit.NameInGame)))
		if i != len(a)-1 {
			c.WriteString(", ")
		}
		if !sameLocations[upgradedUnit.Location] {
			sameLocations[upgradedUnit.Location] = true
			b.WriteString(fmt.Sprintf("%v %v", e.Buildings[upgradedUnit.Location], strings.ToLower(upgradedUnit.Location.GetNameInGame())))
			b.WriteString(", ")
		}
	}
	return fmt.Sprintf("%v (%v)", c.String(), strings.TrimSuffix(b.String(), ", "))
}

func beautyUnits(m map[UnitID]int) string {
	var a SortByAgeLocationName
	for unitID := range m {
		u, err := NewUnit(unitID)
		if err == nil {
			a = append(a, u)
		}
	}
	sort.Sort(a)
	var ret strings.Builder
	for i, v := range a {
		unitID := UnitID(v.GetID().IntID())
		ret.WriteString(fmt.Sprintf("%v %v", m[unitID], v.GetFullName()))
		if i < len(a)-1 && a[i+1].GetID().GetAge() != v.GetID().GetAge() {
			ret.WriteString(",\n   ")
		} else {
			ret.WriteString(", ")
		}
	}
	return ret.String()
}

func beautyTechs(m map[TechID]bool) string {
	var a SortByAgeLocationName
	for techID, researched := range m {
		if !researched {
			continue
		}
		if CheckIsAutoTech(techID) || CheckIsBuiltTech(techID) {
			continue
		}
		t, found := AllTechs[techID]
		if found {
			a = append(a, t)
		} else { // should not happen
			println("missing key in AllTechs: ", techID)
		}
	}
	sort.Sort(a)
	var ret strings.Builder
	for i, v := range a {
		ret.WriteString(fmt.Sprintf("%v", v.GetFullName()))
		if i < len(a)-1 && a[i+1].GetID().GetAge() != v.GetID().GetAge() {
			ret.WriteString(",\n   ")
		} else if i == len(a)-1 {
			ret.WriteString(".")
		} else {
			ret.WriteString(", ")
		}
	}
	return ret.String()
}

type SortByAgeLocationName []UnitOrTech

func (a SortByAgeLocationName) Len() int      { return len(a) }
func (a SortByAgeLocationName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortByAgeLocationName) Less(i, j int) bool {
	age1 := a[i].GetID().GetAge()
	age2 := a[j].GetID().GetAge()
	if age1 != age2 {
		return age1 < age2
	}
	location1, location2 := a[i].GetLocation(), a[j].GetLocation()
	if location1 != location2 {
		locationAge1, locationAge2 := location1.GetAge(), location2.GetAge()
		if locationAge1 != locationAge2 {
			return locationAge1 < locationAge2
		}
		return location1 > location2
	}
	return a[i].GetFullName() < a[j].GetFullName()
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
		e.Spent = &Cost{}
		e.FreeUnits = map[UnitID]int{}
	}
}

func WithDisableAutoBuildHouse() EmpireOption {
	return func(e *EmpireDeveloping) {
		e.IsAutoBuildHouse = false
	}
}

func WithIgnoreRequiredTechOrBuilding() EmpireOption {
	return func(e *EmpireDeveloping) {
		e.IsIgnoreRequiredTechOrBuilding = true
	}
}
