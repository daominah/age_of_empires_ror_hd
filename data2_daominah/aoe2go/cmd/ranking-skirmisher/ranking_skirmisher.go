// Command ranking-skirmisher prints the AoE2 DE civilizations ranked by their fully upgraded skirmisher line unit,
// as a Markdown table, one row per civilization.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/logic"
)

// outputFileName is the Markdown table written next to this source file.
const outputFileName = "ranking_skirmisher.md"

// Ranking config (edit before running):
var (
	skirmisherLine = logic.Line{
		UnitIDs: []logic.UnitID{logic.Skirmisher, logic.EliteSkirmisher, logic.ImperialSkirmisher},
		Class:   logic.ClassArcher,
	}
	// extraLines add a row for civilizations with a unit that plays like the skirmisher line,
	// so they can have 2 rows, e.g. Berbers Elite Skirmisher and Elite Genitour.
	extraLines = []logic.Line{
		{UnitIDs: []logic.UnitID{logic.Genitour, logic.EliteGenitour}, Class: logic.ClassCavalryArcher},
	}

	// arbalester is a fully upgraded Arbalester, both the target of the skirmisher and the attacker for its tank:
	// attack 6+4, 40 HP, pierce armor 0+4.
	arbalester = tree.LoadUnitUpgraded(logic.Arbalester)
	// hussar is a fully upgraded Hussar, both a target of the skirmisher and the attacker for its tank vs melee:
	// 95 HP, pierce armor 6, attack 7+4.
	hussar = tree.LoadUnitUpgraded(logic.Hussar)
	// heavyCavalryArcher is a fully upgraded Heavy Cavalry Archer, the other attacker for the tank: attack 7+4.
	heavyCavalryArcher = tree.LoadUnitUpgraded(logic.HeavyCavalryArcher)

	// baseline is a standard fully upgraded Elite Skirmisher, the 100% of each score weight:
	// attack 3+4, range 5+3, reload 3, 35 HP, pierce armor 4+4, 25f 35w.
	baseline = logic.NewRowFromUnit(tree.LoadUnitUpgraded(logic.EliteSkirmisher))

	// scoreWeights are the stats of the score, each weight is the score of that stat 100% better than the baseline.
	// A weight of 0 is 1 for the last stat and 4 times the next one for each stat above it.
	// A skirmisher has to tank and kill archers. An attack or a discount counts more than
	// a speed or a range above 8, e.g. +1 attack (+14% damage) adds 0.86, a 25% discount 0.75 (1.5 with the food),
	// +0.1 speed (+10%) 0.62 and +1 range (+12.5%) 0.5.
	scoreWeights = []logic.WeightedCriterion{
		{Criterion: logic.NewCriterionDamagePerThreeSecondsVs(arbalester), Weight: 6},
		// the tank vs the 2 archers a skirmisher counters, e.g. the missing Ring Archer Armor costs 0.5 vs the Arbalester
		// (2 damage instead of 1) and 0.42 vs the Heavy Cavalry Archer (5 damage instead of 3)
		{Criterion: logic.NewCriterionTankVs(arbalester), Weight: 1},
		{Criterion: logic.NewCriterionTankVs(heavyCavalryArcher), Weight: 1},
		// vs a raiding Hussar (pierce armor 6) a standard skirmisher deals the minimum 1 damage,
		// so +1 attack doubles it, a small weight as killing a Hussar is not its job
		{Criterion: logic.NewCriterionDamagePerThreeSecondsVs(hussar), Weight: 0.5},
		// a range above 8 helps a little, a tile less costs 2: Bracer (+1 range, +1 attack) counts more than an armor upgrade,
		// e.g. the Bulgarians without Ring Archer Armor rank above the Celts without Bracer
		{Criterion: logic.CriterionRange, Weight: 4, WeightIfWorse: 16},
		{Criterion: logic.CriterionSpeed, Weight: 6},
		{Criterion: logic.CriterionCost, Weight: 3},
		// the food on top of the total cost, as the food is the scarce resource, e.g. the Inca 17f
		{Criterion: logic.CriterionFoodCost, Weight: 3},
	}

	// extraScores are the strengths and weaknesses the score leaves out:
	//   - Koreans: the free archer armors (Padded, Leather and Ring Archer Armor, 1150 resources)
	//     rank them above the Byzantines.
	//   - Inca: the food discount grows by age, 20f in the Feudal Age, 17f only in the Imperial Age,
	//     while the food is worth more early.
	//   - Saxons: -5% cost per Town Center or Castle, up to -20%, about -10% in a real game,
	//     so they rank just below the Britons.
	extraScores = map[logic.CivilizationID]float64{logic.Koreans: 0.75, logic.Inca: -0.4, logic.Saxons: 0.3}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// The score decides, the others only break exact ties with what the score leaves out.
	criteria = []logic.Criterion{
		// a stat locked behind a unique technology counts 75%, e.g. the Aztecs Atlatl and the Maya Hul'che Javelineers
		logic.NewCriterionWithExtraScores(
			logic.NewCriterionDiscountUniqueTechnologies(logic.NewCriterionWeightedSum("Score", baseline, scoreWeights), 0.75),
			extraScores),
		// a bonus attack on top of the line, e.g. Saracens above the standard Elite Skirmisher
		logic.NewCriterionBonusAttackVs(logic.ArmorClassStandardBuildings),
		// the strengths the score leaves out, roughly ordered:
		// melee armor vs a raiding Hussar (Muisca), the other bonus attacks,
		// then the hand written traits by weight (Tupi, Saxons, Khitans, see logic.Traits),
		// then the cheaper upgrades (Spanish, Italians)
		logic.NewCriterionTankVs(hussar),
		logic.CriterionBonusAttackSum,
		logic.CriterionTraitWeightSum,
		logic.CriterionHasCheaperTechnologies,
		// cost saving from a unique technology only breaks ties of the same cost for free
		logic.CriterionCostBeforeUniqueTechnologies,
		// last, a strength from a unique technology only breaks the remaining ties
		logic.NewCriterionTankVsBeforeUniqueTechnologies(arbalester),
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	// The same columns as the cavalry archer ranking, both ranged units.
	columns = []logic.Column{
		// a range unit shows its offense first
		logic.ColumnRange,
		logic.ColumnAttack,
		logic.ColumnReloadTime,
		logic.NewColumnFromCriterionNamed("DP3Svs Archer", logic.NewCriterionDamagePerThreeSecondsVs(arbalester)),
		logic.NewColumnFromCriterionNamed("DP3Svs Scout", logic.NewCriterionDamagePerThreeSecondsVs(hussar)),
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnCost,
		logic.ColumnSpeed,
		// the bonus attacks common to the line are not noted
		logic.NewColumnNoteWithBonusAttack(logic.ArmorClassArchers, logic.ArmorClassCavalryArchers, logic.ArmorClassSpearmen, logic.ArmorClassMamelukes),
	}
)

func main() {
	rows := logic.LoadRowsUpgraded(tree, skirmisherLine)
	for _, line := range extraLines {
		rows = append(rows, logic.LoadRowsUpgraded(tree, line)...)
	}
	logic.SortRows(rows, criteria)
	table := logic.RenderMarkdownTable(rows, columns)
	fmt.Print(table)

	outputPath := filepath.Join(findSourceDir(), outputFileName)
	if err := os.WriteFile(outputPath, []byte(table), 0o644); err != nil {
		log.Fatalf("error os.WriteFile: %v", err)
	}
	log.Printf("wrote %v", outputPath)
}

// findSourceDir is the directory of this source file,
// so the output lands next to it whatever the working directory of "go run".
func findSourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

// tree is the embedded game data, parsed before the config that uses it,
// e.g. tree.LoadUnitUpgraded(logic.Arbalester).
var tree = parseTechTreeOrExit()

// parseTechTreeOrExit parses the embedded game data, it exits on error.
func parseTechTreeOrExit() *logic.TechTree {
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	return tree
}
