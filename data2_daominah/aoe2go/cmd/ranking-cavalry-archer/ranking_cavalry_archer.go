// Command ranking-cavalry-archer prints the AoE2 DE civilizations ranked by their fully upgraded cavalry archer line unit,
// as a Markdown table, one row per civilization, plus extra rows for unique cavalry archers.
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
const outputFileName = "ranking_cavalry_archer.md"

// Ranking config (edit before running):
var (
	cavalryArcherLine = logic.Line{
		UnitIDs: []logic.UnitID{logic.CavalryArcher, logic.HeavyCavalryArcher},
		Class:   logic.ClassCavalryArcher,
	}
	// extraLines add a row for civilizations with a unit that plays like the cavalry archer line,
	// so they can have 2 rows, e.g. Mongols Heavy Cavalry Archer and Elite Mangudai.
	// The Genitour and the Elephant Archer are left out: they play too differently.
	extraLines = []logic.Line{
		{UnitIDs: []logic.UnitID{logic.MountedCrossbowman, logic.HeavyMountedCrossbowman}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.Mangudai, logic.EliteMangudai}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.CamelArcher, logic.EliteCamelArcher}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.XianbeiRaider}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.Kipchak, logic.EliteKipchak}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.RathaRanged, logic.EliteRathaRanged}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.BolasRider, logic.EliteBolasRider}, Class: logic.ClassCavalryArcher},
		{UnitIDs: []logic.UnitID{logic.WarWagon, logic.EliteWarWagon}, Class: logic.ClassCavalryArcher},
	}

	// arbalester6Plus4 is a fully upgraded Arbalester attack: 6 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	arbalester6Plus4 = tree.LoadAttackerUpgraded(logic.Arbalester)
	// arbalesterTarget is a fully upgraded Arbalester shot by the cavalry archer:
	// 40 HP, pierce armor 0 + Padded, Leather and Ring Archer Armor.
	arbalesterTarget = tree.LoadUnitUpgraded(logic.Arbalester)

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// Reorder or comment out lines to rank differently.
	criteria = []logic.Criterion{
		logic.CriterionRange,
		logic.NewCriterionTankVs(arbalester6Plus4),
		logic.NewCriterionHitsToKillTarget(arbalesterTarget),
		logic.CriterionReloadTime,
		// cost saving from a unique technology only breaks ties of the same cost for free
		logic.CriterionCostBeforeUniqueTechnologies,
		logic.CriterionCost,
		// last, a strength from a unique technology only breaks the remaining ties
		logic.NewCriterionTankVsBeforeUniqueTechnologies(arbalester6Plus4),
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		// the basic stats first, then the stats that matter for this line
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnAttack,
		logic.ColumnReloadTime,
		logic.ColumnCost,
		logic.ColumnRange,
		logic.NewColumnFromCriterion(logic.NewCriterionTankVs(arbalester6Plus4)),
		logic.NewColumnFromCriterion(logic.NewCriterionHitsToKillTarget(arbalesterTarget)),
		logic.NewColumnBonusAttack("Bonus vs Spearman", logic.ArmorClassSpearmen),
		logic.ColumnSpeed,
		// the bonus attack vs Spearmen is common to the line (Parthian Tactics), so it is not noted
		logic.NewColumnNoteWithBonusAttack(logic.ArmorClassSpearmen),
	}
)

func main() {
	rows := logic.LoadRowsUpgraded(tree, cavalryArcherLine)
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
// e.g. tree.LoadAttackerUpgraded(logic.Arbalester).
var tree = parseTechTreeOrExit()

// parseTechTreeOrExit parses the embedded game data, it exits on error.
func parseTechTreeOrExit() *logic.TechTree {
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	return tree
}
