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
		stepStr, err := c.step.Marshal()
		if err != nil {
			t.Errorf("error step.Marshal(%+v): %v", c.step, err)
			continue
		}
		if stepStr != c.stepStr {
			t.Errorf("error step.Marshal(%+v) got:\n%v, but want:\n%v", c.step, stepStr, c.stepStr)
		}

		parsedStep, err := NewStep(c.stepStr)
		if err != nil {
			t.Errorf("error NewStep(%v): %v", c.stepStr, err)
			continue
		}
		if !c.step.checkEqual(*parsedStep) {
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
	if !step.checkEqual(want) {
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
		_, err := c.step.Marshal()
		if !errors.Is(err, c.wantErr) {
			t.Errorf("error step.Marshal(%+v): got: %v, but want: %v", c.step, err, c.wantErr)
		}
	}
}

func TestNewStrategy_RequiredTechs(t *testing.T) {
	empire, err := NewEmpireDeveloping(WithCivilization(Sumerian))
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	prepare, errs := NewStrategy(`
		B109      Town_Center1         1      -1
U83       Man                  10      109
B12       Barracks1            1      -1
C101    Tool_Age               1      109
// R16     Watch_Tower            1      68
// R63     Axe                    1      12
B101      Stable1              1      -1
R46     Toolworking            1      103
B84       Market1              1      -1
C102    Bronze_Age             1      109`)
	if len(errs) != 0 {
		t.Fatalf("error prepare NewStrategy: %v", errs)
	}
	for _, s := range prepare {
		err := empire.Do(s)
		if err != nil {
			t.Fatalf("error prepare Do Step(%+v): %v", s, err)
		}
	}
	// t.Logf("empire prepared: %v", empire.Summary())
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
		{line: "T37       Soldier-Cavalry1     2      101       1", wantErr: ErrUnitDisabledByCiv},
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

func TestTechnology_GetName(t *testing.T) {
	for k, v := range AllTechs {
		name := v.NameInGame
		if name == "" {
			t.Errorf("error tech empty name: id %v", k)
		}
	}
}

func TestAutoBuildHouse(t *testing.T) {
	//for i := 40.0; i > -40; i-- {
	//	nHouses := calcNumberHousesAutoBuild(i)
	//	t.Logf("i: %v, nHouses: %v", i, nHouses)
	//}

	empire, err := NewEmpireDeveloping()
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	err = empire.Do(Step{Action: Build, UnitOrTechID: TownCenter, Quantity: 1})
	if err != nil {
		t.Fatalf("error Build TownCenter: %v", err)
	}
	err = empire.Do(Step{Action: Build, UnitOrTechID: Villager, Quantity: 10})
	if err != nil {
		t.Fatalf("error Build Villager: %v", err)
	}
	if empire.Buildings[House] != 5 {
		t.Errorf("should auto build House here, but nHouse: %v", empire.Buildings[House])
	}
	err = empire.Do(Step{Action: Build, UnitOrTechID: Villager, Quantity: 14})
	if err != nil {
		t.Fatalf("error Build Villager: %v", err)
	}
	if empire.Buildings[House] != 10 {
		t.Errorf("should auto build House here, but nHouse: %v", empire.Buildings[House])
	}
	// t.Logf("empire: %v", empire.Summary())
}

//go:embed Default.ai
var testDefaultAI string

func TestStrategy_DefaultAI(t *testing.T) {
	if len(testDefaultAI) == 0 {
		t.Fatalf("error testDefaultAI: empty")
	}
	steps, errs := NewStrategy(testDefaultAI)
	if len(errs) > 0 {
		t.Logf("errors in NewStrategy:")
		for _, err := range errs {
			t.Error(err)
		}
	}
	empire, err := NewEmpireDeveloping()
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	for i, step := range steps {
		err := empire.Do(step)
		if err != nil {
			if errors.Is(err, ErrMissingRequireTechs) && step.UnitOrTechID == Medicine {
				continue // maybe Microsoft wants Medicine is a Bronze tech, but it is Iron tech
			}
			t.Errorf("error i %v DoStep(%v): %v", i, step, err)
		}
	}

	// t.Logf("summary: %v", empire.Summary())

	//for _, tech := range AllNormalTechs {
	//	if !empire.Techs[tech.ID] {
	//		t.Logf("tech not in the file: %v", tech.NameInGame)
	//		// Medicine, Sacrifice, Writing, SmallWall, MediumWall, FortifiedWall
	//	}
	//}
}

func TestNewStrategy_Macedonian(t *testing.T) {
	buildOrder := `
B109      Town_Center1         1      -1
U83       Man                  10      109
B12       Barracks1            1      -1


C101    Tool_Age               1      109
R16     Watch_Tower            1      68
B101      Stable1              1      -1
R46     Toolworking            1      103
B84       Market1              1      -1
B79       Watch_Tower          1      -1
T83       Man                  2      109       0
T299      Soldier-Scout        1      101       0
// 12 villager, 1 scout


C102    Bronze_Age             1      109
R107    Wood_Working           1      84
R40     Leather_Armor_-_Soldie 1      103
B79       Watch_Tower          2      -1

T37       Soldier-Cavalry1     1      101       1
R32     Artisanship            1      84
R12     Sentry_Tower           1      68
R43     Scale_Armor_-_Soldiers 1      103
T83       Man                  2      109       0
B0        Academy              1      -1
B82       Government_Center    1      -1
B0        Academy              1      -1
B0        Academy              1      -1

T37       Soldier-Cavalry1     1      101       0
T83       Man                  2      109       0
T93       Soldier-Phal1        1      0         0

// 16 villager, 2 cavalry, 1 hoplite, 3 tower
// spent Bronze:

C103    Iron_Age               1      109
R112    Architecture           1      82
R47     Bronze_Shield          1      103
T37       Soldier-Cavalry1     1      101       0
T93       Soldier-Phal1        3      0         0
B79       Watch_Tower          2      -1
B87       Range1               1      -1
B0        Academy              2      -1
R51     Metal_Working          1      103
T93       Soldier-Phal1        3      0         0
B79       Watch_Tower          2       -1
B49       Siege_Workshop       2      -1
T93       Soldier-Phal1        3      0         0

C73     Phalanx                1      0
C113    Aristocracy            1      82
R48     Chain_Mail_-_Soldiers  1      103
R15     Guard_Tower            1      68
T93       Soldier-Phal1        4      0         0
T11       Soldier-Ballista     2      49        0
T83       Man                  2      109       0
B79       Watch_Tower          2      -1

T93       Soldier-Phal1        4      0         0
R37     Alchemy                1      82
T83       Man                  1      109       0
T11       Soldier-Ballista     2      49        0
B49       Siege_Workshop       2      -1
U93       Soldier-Phal1        4      0
T83       Man                  1      109       0
B79       Watch_Tower          2      -1

R117    Iron_Shield            1      103
R106    Ballistics             1      82
R79     Centurion              1      0
U93       Soldier-Phal1        4      0
T83       Man                  2      109       0
B49       Siege_Workshop       2      -1
B79       Watch_Tower          2      -1

U11       Soldier-Ballista     4      49
R122    Tower_Shield           1      103
// R114    Writing                1      82
U93       Soldier-Phal1        4      0
T83       Man                  1      109       0
B79       Watch_Tower          1      -1

R52     Metallurgy             1      103
U11       Soldier-Ballista     6      49
U93       Soldier-Phal1        4      0
T83       Man                  1      109       0
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B109      Town_Center1         1      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B79       Watch_Tower          1      -1

U11       Soldier-Ballista     6      49
B50       Farm                 2      -1
B109      Town_Center1         1      -1
B79       Watch_Tower          2      -1

// spent army:`
	steps, errs := NewStrategy(buildOrder)
	if len(errs) > 0 {
		t.Logf("errors in NewStrategy:")
		for _, err := range errs {
			t.Error(err)
		}
	}
	empire, err := NewEmpireDeveloping(WithCivilization(Macedonian))
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	if empire.FreeUnits[TownCenter] != 1 || empire.FreeUnits[Villager] != 3 {
		t.Errorf("error free units got: %+v, but want 1 TownCenter 3 Villager", empire.FreeUnits)
	}

	prevSpent := 0
	for i, step := range steps {
		prevSpent = int(empire.Spent.Food) / 1000
		err := empire.Do(step)
		if err != nil {
			t.Errorf("error i %v DoStep(%v): %v", i, step, err)
		}

		spent := int(empire.Spent.Food) / 1000
		if spent > prevSpent {
			// t.Logf("i %v DoStep(%v): spent: %+v", i, step, empire.Spent.Food)
		}
		prevSpent = spent
	}
	_ = prevSpent

	// t.Logf("summary: %v", empire.Summary())

	if empire.UnitStats[Ballista].Cost.Wood != 50 || empire.UnitStats[Ballista].Cost.Gold != 40 {
		t.Errorf("error Ballista cost: %v", empire.UnitStats[Ballista].Cost)
	}
	if empire.UnitStats[House].Cost.Wood != 30 {
		t.Errorf("error House cost: %v", empire.UnitStats[House].Cost)
	}
	if empire.FreeUnits[TownCenter] != 0 || empire.FreeUnits[Villager] != 0 {
		t.Errorf("error free units got: %+v, but want 0", empire.FreeUnits)
	}
	if empire.Combatants[Hoplite] != 34 {
		t.Errorf("error Hoplite: got: %v, but want 34", empire.Combatants[Hoplite])
	}
	wantSpent := Cost{Wood: 9190, Food: 10835, Gold: 7430, Stone: 3750}
	if !empire.Spent.CheckEqual(wantSpent) {
		t.Errorf("error spent: got: %+v, but want %+v", empire.Spent, wantSpent)
	}
}

func TestUnitOrTechID_GetAge(t *testing.T) {
	for _, c := range []struct {
		id  UnitOrTechID
		age TechID
	}{
		{id: TownCenter, age: StoneAge},
		{id: House, age: StoneAge},
		{id: Villager, age: StoneAge},
		{id: Clubman, age: StoneAge},
		{id: GranaryBuilt, age: StoneAge},
		{id: StoneAge, age: StoneAge},
		{id: Market, age: ToolAge},
		{id: Scout, age: ToolAge},
		{id: Bowman, age: ToolAge},
		{id: ArcheryRange, age: ToolAge},
		{id: EnableSlinger, age: ToolAge},
		{id: Woodworking, age: ToolAge},
		{id: Camel, age: BronzeAge},
		{id: Temple, age: BronzeAge},
		{id: GovernmentCenter, age: BronzeAge},
		{id: Hoplite, age: BronzeAge},
		{id: Wheel, age: BronzeAge},
		{id: EnableCavalry, age: BronzeAge},
		{id: HorseArcher, age: IronAge},
		{id: Elephant, age: IronAge},
		{id: Ballista, age: IronAge},
		{id: Wonder, age: IronAge},
		{id: Centurion, age: IronAge},
		{id: Alchemy, age: IronAge},
	} {
		if age := c.id.GetAge(); age != c.age {
			t.Errorf("error %v.GetAge(): got: %v, but want: %v", c.id, age, c.age)
		}
	}
}

func TestGuessCivilization(t *testing.T) {
	for _, c := range []struct {
		fileName string
		want     CivilizationID
	}{
		{fileName: "Assyria_Archer.ai", want: Assyrian},
		{fileName: "Babylon_Tower_Priest.ai.ai", want: Babylonian},
		{fileName: "Minoa_Composite_Bowmen.ai", want: Minoan},
		{fileName: "Rome Legion.ai", want: Roman},
	} {
		if guess := GuessCivilization(c.fileName); guess != c.want {
			t.Errorf("error GuessCivilization(%v): got: %v, but want: %v", c.fileName, guess, c.want)
		}
	}
}

func TestSpentRoman(t *testing.T) {
	empire, err := NewEmpireDeveloping(WithCivilization(Roman))
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	strategy := []Step{
		{Action: Build, UnitOrTechID: TownCenter, Quantity: 1},
		{Action: Build, UnitOrTechID: Villager, Quantity: 10},
		{Action: Build, UnitOrTechID: Barracks, Quantity: 1},
		{Action: Research, UnitOrTechID: ToolAge, Quantity: 1},
		{Action: Build, UnitOrTechID: Barracks, Quantity: 5},
		{Action: Research, UnitOrTechID: WatchTower, Quantity: 1},
		{Action: Research, UnitOrTechID: Axe, Quantity: 1},
		{Action: Build, UnitOrTechID: Stable, Quantity: 1},
		{Action: Research, UnitOrTechID: Toolworking, Quantity: 1},
		{Action: Build, UnitOrTechID: Market, Quantity: 1},
		{Action: Build, UnitOrTechID: Tower, Quantity: 1},
		{Action: Train, UnitOrTechID: Villager, Quantity: 2},
		{Action: Research, UnitOrTechID: BronzeAge, Quantity: 1},
		{Action: Research, UnitOrTechID: Woodworking, Quantity: 1},
	}
	for i, step := range strategy {
		err := empire.Do(step)
		if err != nil {
			t.Fatalf("error i %v DoStep(%v): %v", i, step, err)
		}
		if i == 0 {
			if want := (120 + 120) * 0.85; empire.Spent.Wood != want {
				t.Errorf("error i %v spent Wood: got: %v, but want %v", i, empire.Spent.Wood, want)
			}
		} else if i == 1 {
			if want := (120 + 120 + 30*5) * 0.85; empire.Spent.Wood != want {
				t.Errorf("error i %v spent Wood: got: %v, but want %v", i, empire.Spent.Wood, want)
			}
		} else if i == 4 {
			if want := (120 + 120 + 30*5 + 125*6) * 0.85; empire.Spent.Wood != want {
				t.Errorf("error i %v spent Wood: got: %v, but want %v", i, empire.Spent.Wood, want)
			}
		}
	}
	want := Cost{
		Wood:  (120+120+30*5+125*6+150+150)*0.85 + 75,
		Food:  9*50 + 500 + 800 + 100 + 120 + 100 + 50,
		Gold:  0,
		Stone: 75,
	}
	if !empire.Spent.CheckEqual(want) {
		t.Errorf("error spent: got: %+v, but want %+v", empire.Spent, want)
	}
}

func TestEmpireDeveloping_BeautyMainArmy(t *testing.T) {
	empire, err := NewEmpireDeveloping(WithCivilization(Yamato))
	if err != nil {
		t.Fatalf("error NewEmpireDeveloping: %v", err)
	}
	empire.Combatants = map[UnitID]int{
		Villager:    30,
		Scout:       1,
		Cavalry:     48,
		Priest:      12,
		HorseArcher: 8,
	}
	empire.Buildings = map[UnitID]int{
		Barracks: 1, Granary: 1, House: 30, StoragePit: 1, TownCenter: 4,
		Farm: 16, Market: 1, ArcheryRange: 4, Stable: 7, Tower: 16,
		GovernmentCenter: 1, Temple: 2,
	}
	CataphractEffect126(empire)
	TempleBuiltEffect17(empire)
	EnableHorseArcherEffect60(empire)

	got := empire.beautyMainArmy()
	want := "48 cataphract, 12 priest, 8 horse archer (7 stable, 2 temple, 4 range)"
	if got != want {
		t.Errorf("error beautyMainArmy: got: %v, want: %v", got, want)
	}
}
