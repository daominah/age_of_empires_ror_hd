package logic

// Criterion sorts the rows, in the configured order.
type Criterion struct {
	Name           string
	Value          func(row Row) float64
	IsHigherBetter bool
}

// NewColumnFromCriterion shows the value of a criterion, e.g. the tank vs a Castle.
func NewColumnFromCriterion(criterion Criterion) Column {
	return Column{Name: criterion.Name, Value: func(row Row) string { return formatNumber(criterion.Value(row)) }}
}

// Criteria on the unit stats, lower is better for cost and reload time.
var (
	CriterionCost = Criterion{Name: "Cost", Value: func(r Row) float64 { return r.Unit.Cost.Sum() }}
	// CriterionCostBeforeUniqueTechnologies ranks a cost saving from a unique technology
	// below the same cost for free, e.g. Magyars and Berbers above Gurjaras (Kshatriyas).
	CriterionCostBeforeUniqueTechnologies = Criterion{
		Name:  "Cost before unique tech",
		Value: func(r Row) float64 { return r.UnitWithoutUniqueTechnologies.Cost.Sum() },
	}
	CriterionAttack     = Criterion{Name: "Attack", Value: func(r Row) float64 { return getAttack(r.Unit) }, IsHigherBetter: true}
	CriterionReloadTime = Criterion{Name: "Reload time", Value: func(r Row) float64 { return r.Unit.ReloadTime }}
	CriterionRange      = Criterion{Name: "Range", Value: func(r Row) float64 { return r.Unit.Range }, IsHigherBetter: true}
)

// NewCriterionTankVs is the hits of the attacker to kill the row unit, more is better.
func NewCriterionTankVs(attacker Attacker) Criterion {
	return Criterion{
		Name:           "Tank vs " + attacker.Name,
		Value:          func(r Row) float64 { return float64(CalculateHitsToKill(attacker, r.Unit)) },
		IsHigherBetter: true,
	}
}

// NewCriterionTankVsBeforeUniqueTechnologies is NewCriterionTankVs on the unit without its unique technologies,
// to rank a strength from a unique technology below the same strength for free,
// e.g. Turks above Tatars (Silk Armor).
func NewCriterionTankVsBeforeUniqueTechnologies(attacker Attacker) Criterion {
	return Criterion{
		Name:           "Tank vs " + attacker.Name + " before unique tech",
		Value:          func(r Row) float64 { return float64(CalculateHitsToKill(attacker, r.UnitWithoutUniqueTechnologies)) },
		IsHigherBetter: true,
	}
}

// NewCriterionHitsToKillTarget is the hits of the row unit to kill the target, fewer is better,
// e.g. a skirmisher shooting an Arbalester.
func NewCriterionHitsToKillTarget(target Unit) Criterion {
	return Criterion{
		Name: "Hits to kill " + target.Name,
		Value: func(r Row) float64 {
			return float64(CalculateHitsToKill(Attacker{Name: r.Unit.Name, Attacks: r.Unit.Attacks}, target))
		},
	}
}
