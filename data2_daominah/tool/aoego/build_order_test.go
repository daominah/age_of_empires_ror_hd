package aoego

import (
	"errors"
	"testing"
)

func TestStep(t *testing.T) {
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

func TestStep_Weird(t *testing.T) {
	// should be "R125    Armored_Elephant       1      101", but `Default.ai` has this weird exception:
	step, err := NewStep("R125    Armored Elephants      1      101")
	if err != nil {
		t.Fatalf("error NewStep weird Armored Elephants: %v", err)
	}
	want := Step{Action: Research, UnitOrTechID: ArmoredElephant, Quantity: 1}
	if !step.CheckEqual(want) {
		t.Errorf("error NewStep weird Armored Elephants: got: %+v, but want: %+v", step, want)
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
