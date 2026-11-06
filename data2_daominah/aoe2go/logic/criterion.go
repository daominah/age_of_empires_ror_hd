package logic

import "math"

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

// NewColumnFromCriterionNamed is NewColumnFromCriterion with a shorter name for the header,
// e.g. "DP2Svs Scout" for the damage vs a fully upgraded Hussar.
func NewColumnFromCriterionNamed(name string, criterion Criterion) Column {
	column := NewColumnFromCriterion(criterion)
	column.Name = name
	return column
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
	// CriterionFoodCost is the food alone, for a line where the food is the scarce resource, e.g. the skirmisher.
	CriterionFoodCost   = Criterion{Name: "Food cost", Value: func(r Row) float64 { return r.Unit.Cost.Food }}
	CriterionAttack     = Criterion{Name: "Attack", Value: func(r Row) float64 { return getAttack(r.Unit) }, IsHigherBetter: true}
	CriterionReloadTime = Criterion{Name: "Reload time", Value: func(r Row) float64 { return r.Unit.ReloadTime }}
	CriterionRange      = Criterion{Name: "Range", Value: func(r Row) float64 { return r.Unit.Range }, IsHigherBetter: true}
	CriterionSpeed      = Criterion{Name: "Speed", Value: func(r Row) float64 { return r.Unit.Speed }, IsHigherBetter: true}
)

// NewCriterionTankVs is the whole hits of the attacker to kill the row unit, more is better,
// as in the game: a few HP or armor that do not save a hit do not count.
// It combines the HP and the armor, e.g. Bloodlines +20 HP and Parthian Tactics +2 pierce armor.
func NewCriterionTankVs(attacker Unit) Criterion {
	return Criterion{
		Name:           "Tank vs " + attacker.Name,
		Value:          func(r Row) float64 { return float64(CalculateHitsToKill(attacker, r.Unit)) },
		IsHigherBetter: true,
	}
}

// NewCriterionTankVsBeforeUniqueTechnologies is NewCriterionTankVs on the unit without its unique technologies,
// to rank a strength from a unique technology below the same strength for free,
// e.g. Turks above Tatars (Silk Armor).
func NewCriterionTankVsBeforeUniqueTechnologies(attacker Unit) Criterion {
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
			return float64(CalculateHitsToKill(r.Unit, target))
		},
	}
}

// NewCriterionBonusAttackVs is the bonus attack of the row unit against the armor class, more is better,
// e.g. the Saracens Elite Skirmisher +3 vs Standard Buildings.
func NewCriterionBonusAttackVs(class ArmorClass) Criterion {
	return Criterion{
		Name:           "Bonus attack vs " + class.String(),
		Value:          func(r Row) float64 { return r.Unit.Attacks[class] },
		IsHigherBetter: true,
	}
}

// CriterionBonusAttackSum is the sum of the bonus attacks of the row unit, all classes but melee and pierce,
// more is better, a rough tie-breaker, e.g. the Mongols Elite Mangudai +5 vs Siege Weapons.
var CriterionBonusAttackSum = Criterion{
	Name: "Bonus attack sum",
	Value: func(r Row) float64 {
		sum := 0.0
		for class, amount := range r.Unit.Attacks {
			if class != ArmorClassMelee && class != ArmorClassPierce {
				sum += amount
			}
		}
		return sum
	},
	IsHigherBetter: true,
}

// CriterionTraitWeightSum is the sum of the weights of the hand written traits of the row unit, more is better,
// a rough tie-breaker as every trait is a strength, e.g. the Sicilians bonus damage taken -40%.
var CriterionTraitWeightSum = Criterion{
	Name: "Trait weight sum",
	Value: func(r Row) float64 {
		sum := 0.0
		for _, trait := range Traits {
			if trait.CivilizationID == r.CivilizationID && trait.CheckHasUnit(r.Unit.ID) {
				sum += trait.Weight
			}
		}
		return sum
	},
	IsHigherBetter: true,
}

// CriterionHasCheaperTechnologies is 1 if the civilization researches some upgrades of the row unit cheaper,
// more is better, a rough tie-breaker that leaves out how much cheaper, e.g. the Spanish Blacksmith without gold.
var CriterionHasCheaperTechnologies = Criterion{
	Name: "Has cheaper technologies",
	Value: func(r Row) float64 {
		if FormatTechnologyCostNote(r) == "" {
			return 0
		}
		return 1
	},
	IsHigherBetter: true,
}

// NewCriterionDamagePerThreeSecondsVs is NewCriterionDamagePerTwoSecondsVs in 3 seconds,
// for a line that reloads in 3 seconds, so the value is the damage of one shot, e.g. a skirmisher.
func NewCriterionDamagePerThreeSecondsVs(target Unit) Criterion {
	return Criterion{
		// DP3S is short for damage per 3 seconds, the column reads e.g. "DP3SvsArbalester"
		Name:           "DP3Svs " + target.Name,
		Value:          func(r Row) float64 { return CalculateDamagePerTwoSeconds(r.Unit, target) * 3 / 2 },
		IsHigherBetter: true,
	}
}

// NewCriterionDamagePerTwoSecondsVs is the damage of the row unit on the target in 2 seconds, more is better,
// e.g. a Hussar vs an Arbalester. Unlike the attack, it includes the reload time, the armor of the target,
// the extra arrows and the charged attack.
func NewCriterionDamagePerTwoSecondsVs(target Unit) Criterion {
	return Criterion{
		// DP2S is short for damage per 2 seconds, the column reads e.g. "DP2SvsHussar"
		Name:           "DP2Svs " + target.Name,
		Value:          func(r Row) float64 { return CalculateDamagePerTwoSeconds(r.Unit, target) },
		IsHigherBetter: true,
	}
}

// NewCriterionDamagePerTwoSecondsPer100GoldVs is NewCriterionDamagePerTwoSecondsVs per 100 gold of the row unit cost,
// more is better, e.g. the Xianbei Raider 35 gold.
// Per 100 gold keeps it on the scale of the damage, e.g. a Heavy Cavalry Archer vs a Hussar: 5.56 and 9.26.
func NewCriterionDamagePerTwoSecondsPer100GoldVs(target Unit) Criterion {
	return Criterion{
		Name: "DP2Svs " + target.Name + " per 100 gold",
		Value: func(r Row) float64 {
			// a unit without gold counts as 1 gold
			return CalculateDamagePerTwoSeconds(r.Unit, target) * 100 / math.Max(1, r.Unit.Cost.Gold)
		},
		IsHigherBetter: true,
	}
}

// NewCriterionDiscountUniqueTechnologies counts only a share of what the unique technologies add to the criterion,
// as a stat locked behind a unique technology is a downside: it costs a Castle and a lot of resources,
// e.g. 0.75 counts 3/4 of the Aztecs Atlatl +1 attack +1 range.
func NewCriterionDiscountUniqueTechnologies(criterion Criterion, uniqueTechnologyShare float64) Criterion {
	return Criterion{
		Name: criterion.Name,
		Value: func(r Row) float64 {
			rowWithoutUniqueTechnologies := r
			rowWithoutUniqueTechnologies.Unit = r.UnitWithoutUniqueTechnologies
			without := criterion.Value(rowWithoutUniqueTechnologies)
			return without + uniqueTechnologyShare*(criterion.Value(r)-without)
		},
		IsHigherBetter: criterion.IsHigherBetter,
	}
}

// NewCriterionWithExtraScores adds a hand written score to the criterion for some civilizations,
// for a strength the stats leave out, e.g. the Koreans free archer armors.
func NewCriterionWithExtraScores(criterion Criterion, extraScores map[CivilizationID]float64) Criterion {
	return Criterion{
		Name:           criterion.Name,
		Value:          func(r Row) float64 { return criterion.Value(r) + extraScores[r.CivilizationID] },
		IsHigherBetter: criterion.IsHigherBetter,
	}
}

// WeightedCriterion is a criterion in a NewCriterionWeightedSum.
type WeightedCriterion struct {
	Criterion Criterion
	// Weight is the score of the criterion value 100% better than the baseline,
	// e.g. 7 for the range: a range of 8 vs 7 is 14% better, so it adds 1.
	// 0 is the default: 1 for the last criterion, 4 times the next one for each criterion above it,
	// so a criterion added without tuning ranks by its place in the list.
	Weight float64
	// WeightIfWorse replaces the Weight when the row is worse than the baseline, 0 to keep the Weight,
	// e.g. a range of 6 hurts a cavalry archer more than a range of 8 helps it.
	WeightIfWorse float64
	// Limit caps how much better or worse than the baseline the criterion counts, 0 for no cap,
	// e.g. 0.5 counts the Elite War Wagon tank +317% vs a Heavy Cavalry Archer as +50%,
	// so one stat far from the baseline does not decide the rank alone.
	Limit float64
}

// NewCriterionWeightedSum is the sum of the criteria values compared to the baseline row, more is better.
// Each value is how much better the row is than the baseline, e.g. +14% for a range of 8 vs 7,
// -14% for a reload time of 2 vs 1.75 as lower is better, times its weight.
// Comparing to a baseline puts all the criteria on the same scale, whatever their unit,
// e.g. a standard fully upgraded Heavy Cavalry Archer for the cavalry archer line.
// A criterion of 0 for the baseline counts its difference instead, as no ratio exists.
func NewCriterionWeightedSum(name string, baseline Row, terms []WeightedCriterion) Criterion {
	weights := make([]float64, len(terms))
	defaultWeight := 1.0
	for i := len(terms) - 1; i >= 0; i-- {
		weights[i] = terms[i].Weight
		if weights[i] == 0 {
			weights[i] = defaultWeight
		}
		defaultWeight *= 4
	}
	return Criterion{
		Name: name,
		Value: func(r Row) float64 {
			score := 0.0
			for i, term := range terms {
				value, reference := term.Criterion.Value(r), term.Criterion.Value(baseline)
				better := value - reference
				if reference != 0 {
					better = better / math.Abs(reference)
				}
				if !term.Criterion.IsHigherBetter {
					better = -better
				}
				if term.Limit > 0 {
					better = math.Max(-term.Limit, math.Min(term.Limit, better))
				}
				weight := weights[i]
				if better < 0 && term.WeightIfWorse != 0 {
					weight = term.WeightIfWorse
				}
				score += weight * better
			}
			return score
		},
		IsHigherBetter: true,
	}
}
