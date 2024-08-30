package aoego

import (
	_ "embed"
	"errors"
	"testing"
)

func TestStep_String_NewStep(t *testing.T) {
	for _, c := range []struct {
		step    Step
		stepStr string
	}{
		{
			step: Step{
				Action:       Train,
				UnitOrTechID: Hoplite,
				Quantity:     4,
			},
			stepStr: "U93       Soldier-Phal1        4      0",
		},
		{
			step: Step{
				Action:       TrainLimit,
				UnitOrTechID: Scout,
				Quantity:     1,
				LimitRebuild: 2,
			},
			stepStr: "T299      Soldier-Scout        1      101       2",
		},
		{
			step: Step{
				Action:       Build,
				UnitOrTechID: GovernmentCenter,
				Quantity:     1,
			},
			stepStr: "B82       Government_Center    1      -1",
		},
		{
			step: Step{
				Action:       BuildLimit,
				UnitOrTechID: Tower,
				Quantity:     2,
				LimitRebuild: 1,
			},
			stepStr: "A79       Watch_Tower          2      -1        1",
		},
		{
			step: Step{
				Action:       ResearchCritical,
				UnitOrTechID: BronzeAge,
				Quantity:     1,
			},
			stepStr: "C102    Bronze_Age             1      109",
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: Wheel,
				Quantity:     1,
			},
			stepStr: "R28     Wheel                  1      84",
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: LeatherArmorInfantry,
				Quantity:     1,
			},
			stepStr: "R40     Leather_Armor_-_Soldie 1      103",
		},
	} {
		stepStr, err := c.step.String()
		if err != nil {
			t.Errorf("error step.String(%+v): %v", c.step, err)
			continue
		}
		if stepStr != c.stepStr {
			t.Errorf("error step.String(%+v) got:\n%v, but want:\n%v", c.step, stepStr, c.stepStr)
		}

		parsedStep, err := NewStep(c.stepStr)
		if err != nil {
			t.Errorf("error NewStep(%v): %v", c.stepStr, err)
			continue
		}
		if !c.step.CheckEqual(*parsedStep) {
			t.Errorf("error NewStep(%v): got: %+v, but want: %+v", c.stepStr, parsedStep, c.step)
		}
	}
}

func TestNewStep_Weird(t *testing.T) {
	// `Default.ai` has a weird exception internal name with space
	step, err := NewStep("R125    Armored Elephants      1      101")
	if err != nil {
		t.Fatalf("error NewStep weird Armored Elephants: %v", err)
	}
	want := Step{Action: Research, UnitOrTechID: ArmoredElephant, Quantity: 1}
	if !step.CheckEqual(want) {
		t.Errorf("error NewStep weird Armored Elephants: got: %+v, but want: %+v", step, want)
	}
}

func TestNewStep_Error(t *testing.T) {
	for _, c := range []struct {
		line    string
		wantErr error
	}{
		{
			line:    "",
			wantErr: ErrEmptyLine,
		},
		{
			line:    "  ",
			wantErr: ErrEmptyLine,
		},
		{
			line:    "// Assyrian bonus:",
			wantErr: ErrEmptyLine,
		},
		{
			line:    "B109      Town_Center1         1      -1  // pre-built",
			wantErr: nil,
		},
		{
			line:    "B109      Town_Center1         // pre-built 1      -1  ",
			wantErr: ErrMissingStepFields,
		},
	} {
		_, err := NewStep(c.line)
		if !errors.Is(err, c.wantErr) {
			t.Errorf("error NewStep(%v): got: %v, but want: %v", c.line, err, c.wantErr)
		}
	}
}

func TestStep_StringError(t *testing.T) {
	for _, c := range []struct {
		step    Step
		wantErr error
	}{
		{
			step: Step{
				Action:       TrainLimit,
				UnitOrTechID: Cavalry,
				Quantity:     1,
				LimitRebuild: 2,
			},
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: Wheel,
				Quantity:     2,
			},
			wantErr: ErrResearchQuantity,
		},
		{
			step: Step{
				Action:       Build,
				UnitOrTechID: NullUnit,
				Quantity:     1,
			},
			wantErr: ErrUnitIDNotFound,
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: TechID(123456789),
				Quantity:     1,
			},
			wantErr: ErrTechIDNotFound,
		},
	} {
		_, err := c.step.String()
		if !errors.Is(err, c.wantErr) {
			t.Errorf("error step.String(%+v): got: %v, but want: %v", c.step, err, c.wantErr)
		}
	}
}

func TestNewBuildOrder_RequiredTechs(t *testing.T) {
	empire, err := NewEmpireDeveloping(WithCivilization(Sumerian))
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	prepare, errs := NewBuildOrder(`
		B109      Town_Center1         1      -1
U83       Man                  10      109
B12       Barracks1            1      -1
C101    Tool_Age               1      109
// R16     Watch_Tower            1      68
// R63     Axe                    1      12
B101      Stable1              1      -1
B84       Market1              1      -1
C102    Bronze_Age             1      109`)
	if len(errs) != 0 {
		t.Fatalf("error prepare NewBuildOrder: %v", errs)
	}
	for _, s := range prepare {
		err := empire.Do(s)
		if err != nil {
			t.Fatalf("error prepare Do Step(%+v): %v", s, err)
		}
	}
	t.Logf("empire prepared: %v", empire.Summary())
	if empire.Combatants[Villager] != 10 {
		t.Errorf("error prepare Villager: got: %v, but want: 10", empire.Combatants[Villager])
	}
	if empire.Buildings[TownCenter] != 1 {
		t.Errorf("error prepare TownCenter: got: %v, but want: 1", empire.Buildings[TownCenter])
	}
	if empire.Buildings[Granary] != 1 { // default AI personality auto build Granary
		t.Errorf("error prepare Granary: got: %v, but want: 1", empire.Buildings[Granary])
	}
	if !empire.Techs[ToolAge] {
		t.Errorf("error prepare ToolAge should be researched")
	}
	if !empire.Techs[StableBuilt] {
		t.Errorf("error prepare StableBuilt should be researched")
	}
	if !empire.Techs[EnableAcademy] {
		t.Errorf("error prepare EnableAcademy should be researched")
	}
	if empire.Techs[EnableSiegeWorkshop] {
		t.Errorf("error prepare EnableSiegeWorkshop should not be researched")
	}

	for _, c := range []struct {
		line    string
		wantErr error
	}{
		{line: "B79       Watch_Tower          1      -1", wantErr: ErrMissingRequireTechs},
		{line: "R12     Sentry_Tower           1      68", wantErr: ErrMissingRequireTechs},
		{line: "R64     Short_Sword            1      12", wantErr: ErrMissingRequireTechs},
		{line: "T75       Soldier-Inf3         5      12        0", wantErr: ErrMissingRequireTechs},
		{line: "T37       Soldier-Cavalry1     2      101       0", wantErr: ErrUnitDisabledByCiv},
		{line: "T37       Soldier-Cavalry1     2      101       1", wantErr: ErrTechDisabledByCiv},
		{line: "B49       Siege_Workshop       4      -1", wantErr: ErrMissingRequireTechs},
		{line: "R56     Improved_bow           1      87", wantErr: ErrTechDisabledByCiv},
	} {
		step, err := NewStep(c.line)
		if err != nil {
			t.Errorf("err bad test code, NewStep(%v): %v", c.line, err)
			continue
		}
		err = empire.Do(*step)
		if !errors.Is(err, c.wantErr) {
			t.Errorf("error DoStep(%v): got: %v, but want: %v", step, err, c.wantErr)
		}
	}

	trainChariot := Step{Action: Train, UnitOrTechID: Chariot, Quantity: 6}
	err = empire.Do(trainChariot)
	if !errors.Is(err, ErrMissingRequireTechs) {
		t.Errorf("error Train Chariot: got: %v, but want: %v", err, ErrMissingRequireTechs)
	}
	err = empire.Do(Step{Action: Research, UnitOrTechID: Wheel})
	if err != nil {
		t.Errorf("error Research Wheel: %v", err)
	}
	err = empire.Do(trainChariot)
	if err != nil {
		t.Errorf("error Train Chariot after Wheel: %v", err)
	}

	trainHoplite := Step{Action: Train, UnitOrTechID: Hoplite, Quantity: 3}
	err = empire.Do(trainHoplite)
	if !errors.Is(err, ErrLocationNotBuilt) {
		t.Errorf("error Train Hoplite: got: %v, but want: %v", err, ErrLocationNotBuilt)
	}
	err = empire.Do(Step{Action: Build, UnitOrTechID: Academy, Quantity: 4})
	if err != nil {
		t.Errorf("error Build Academy: %v", err)
	}
	err = empire.Do(trainHoplite)
	if err != nil {
		t.Errorf("error Train Hoplite after Academy: %v", err)
	}
}

//go:embed Default.ai
var testDefaultAI string

func TestNewStep_DefaultAI(t *testing.T) {
	if len(testDefaultAI) == 0 {
		t.Fatalf("error testDefaultAI: empty")
	}
	steps, errs := NewBuildOrder(testDefaultAI)
	if len(errs) > 0 {
		t.Logf("errors in NewBuildOrder:")
		for _, err := range errs {
			t.Error(err)
		}
	}
	empire, err := NewEmpireDeveloping()
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	for _, step := range steps {
		err := empire.Do(step)
		if err != nil {
			t.Errorf("error DoStep(%v): %v", step, err)
		}
	}

	t.Logf("spent: %+v", empire.Spent)
	t.Logf("population: %+v", empire.CountPopulation())
	t.Logf("buildings: %+v", empire.Buildings)
	t.Logf("techs count: %v", empire.TechnologyCount)
}

func TestTechnology_GetName(t *testing.T) {
	for k, v := range AllTechs {
		name := v.NameInGame
		if name == "" {
			t.Errorf("error tech empty name: id %v", k)
		}
	}
}
