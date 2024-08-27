package aoego

import (
	"errors"
	"testing"
)

func TestStep_String(t *testing.T) {
	for _, c := range []struct {
		step Step
		want string
	}{
		{
			step: Step{
				Action:       Train,
				UnitOrTechID: HopliteID.ID(),
				Quantity:     4,
				Location:     AcademyID,
			},
			want: "U93       Soldier-Phal1        4      0",
		},
		{
			step: Step{
				Action:       TrainLimit,
				UnitOrTechID: ScoutID.ID(),
				Quantity:     1,
				Location:     StableID,
				LimitRebuild: 2,
			},
			want: "T299      Soldier-Scout        1      101       2",
		},
		{
			step: Step{
				Action:       Build,
				UnitOrTechID: GovernmentCenterID.ID(),
				Quantity:     1,
				Location:     NullUnitID,
			},
			want: "B82       Government_Center    1      -1",
		},
		{
			step: Step{
				Action:       BuildLimit,
				UnitOrTechID: TowerID.ID(),
				Quantity:     2,
				Location:     NullUnitID,
				LimitRebuild: 1,
			},
			want: "A79       Watch_Tower          2      -1        1",
		},
		{
			step: Step{
				Action:       ResearchCritical,
				UnitOrTechID: BronzeAgeID.ID(),
				Quantity:     1,
				Location:     TownCenterID,
			},
			want: "C102    Bronze_Age             1      109",
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: WheelID.ID(),
				Quantity:     1,
				Location:     MarketID,
			},
			want: "C102    Bronze_Age             1      109",
		},
	} {
		stepStr, err := c.step.String()
		if err != nil {
			t.Errorf("error step.String(%+v): %v", c.step, err)
			continue
		}
		if stepStr != c.want {
			t.Errorf("error step.String(%+v) got:\n%v, but want:\n%v", c.step, stepStr, c.want)
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
				UnitOrTechID: CavalryID.ID(),
				Quantity:     1,
				Location:     StableID,
				LimitRebuild: 2,
			},
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: WheelID.ID(),
				Quantity:     1,
			},
			wantErr: ErrLocationNotMatched,
		},
		{
			step: Step{
				Action:       Research,
				UnitOrTechID: WheelID.ID(),
				Quantity:     2,
			},
			wantErr: ErrResearchQuantity,
		},
		{
			step: Step{
				Action:       Build,
				UnitOrTechID: 123123123,
				Quantity:     1,
			},
			wantErr: ErrUnitIDNotFound,
		},
	} {
		_, err := c.step.String()
		if !errors.Is(err, c.wantErr) {
			t.Errorf("error step.String(%+v): got: %v, but want: %v", c.step, err, c.wantErr)
		}
	}
}
