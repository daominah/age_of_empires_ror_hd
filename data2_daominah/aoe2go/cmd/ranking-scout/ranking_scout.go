// Command ranking-scout prints the AoE2 DE civilizations ranked by their fully upgraded scout line unit,
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
const outputFileName = "ranking_scout.md"

// Ranking config (edit before running):
var (
	scoutLine = logic.Line{
		UnitIDs: []logic.UnitID{logic.ScoutCavalry, logic.LightCavalry, logic.Hussar, logic.WingedHussar},
		Class:   logic.ClassCavalry,
	}
	// extraLines add a row for civilizations with a unit that plays like the scout line,
	// so they can have 2 rows, e.g. Magyars Hussar and Elite Magyar Huszar.
	extraLines = []logic.Line{
		{UnitIDs: []logic.UnitID{logic.MagyarHuszar, logic.EliteMagyarHuszar}, Class: logic.ClassCavalry},
	}

	// castle15 is a fully upgraded Castle arrow: 11 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	castle15 = logic.Unit{Name: "Castle 15", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 15}}
	// townCenter9 is a fully upgraded Town Center arrow: 5 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	townCenter9 = logic.Unit{Name: "Town Center 9", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 9}}
	// arbalester is a fully upgraded Arbalester: attack 6+4.
	arbalester = tree.LoadUnitUpgraded(logic.Arbalester)
	// heavyCavalryArcher is a fully upgraded Heavy Cavalry Archer: attack 7+4.
	heavyCavalryArcher = tree.LoadUnitUpgraded(logic.HeavyCavalryArcher)
	// hussar is a fully upgraded Hussar, for the tank in the Hussar mirror: attack 7+4.
	hussar = tree.LoadUnitUpgraded(logic.Hussar)
	// baseline is a standard fully upgraded Hussar, the 100% of each score weight.
	baseline = logic.NewRowFromUnit(hussar)
	// scoreWeights weight the arrows that hit a scout most: an Arbalester, a Castle, a Heavy Cavalry Archer
	// and a Town Center, all fully upgraded, in whole hits to kill.
	scoreWeights = []logic.WeightedCriterion{
		{Criterion: logic.NewCriterionTankVs(arbalester), Weight: 4},
		{Criterion: logic.NewCriterionTankVs(castle15), Weight: 3},
		{Criterion: logic.NewCriterionTankVs(heavyCavalryArcher), Weight: 2},
		{Criterion: logic.NewCriterionTankVs(townCenter9), Weight: 1},
	}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// Reorder or comment out lines to rank differently,
	// e.g. put logic.CriterionCost first.
	criteria = []logic.Criterion{
		logic.NewCriterionWeightedSum("Tank vs arrows", baseline, scoreWeights),
		// cost saving from a unique technology only breaks ties of the same cost for free,
		// e.g. Magyars and Berbers above Gurjaras (Kshatriyas)
		logic.CriterionCostBeforeUniqueTechnologies,
		logic.CriterionCost,
		logic.CriterionAttack,
		// the strengths the arrow tank leaves out, roughly ordered:
		// melee armor in the Hussar mirror (Teutons), the speed (Cumans), the bonus attacks,
		// then the hand written traits (regeneration, bonus damage taken, trample)
		logic.NewCriterionTankVs(hussar),
		logic.CriterionSpeed,
		logic.CriterionBonusAttackSum,
		logic.CriterionTraitWeightSum,
		// last, a strength from a unique technology only breaks the remaining ties,
		// e.g. Turks above Tatars (Silk Armor)
		logic.NewCriterionTankVsBeforeUniqueTechnologies(arbalester),
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		// the basic stats, add the stats that matter for this line later
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnAttack,
		logic.ColumnReloadTime,
		logic.ColumnCost,
		logic.ColumnSpeed,
		// the bonus attack vs Monks is common to the line, so it is not noted
		logic.NewColumnNoteWithBonusAttack(logic.ArmorClassMonks),
	}
)

func main() {
	rows := logic.LoadRowsUpgraded(tree, scoutLine)
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
