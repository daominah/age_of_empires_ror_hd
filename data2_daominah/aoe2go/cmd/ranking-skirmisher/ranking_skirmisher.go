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

	// arbalester6Plus4 is a fully upgraded Arbalester attack: 6 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	arbalester6Plus4 = tree.LoadAttackerUpgraded(logic.Arbalester)
	// arbalesterTarget is a fully upgraded Arbalester shot by the skirmisher:
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
		logic.ColumnRange,
		logic.NewColumnFromCriterion(logic.NewCriterionTankVs(arbalester6Plus4)),
		logic.NewColumnFromCriterion(logic.NewCriterionHitsToKillTarget(arbalesterTarget)),
		logic.ColumnCost,
		logic.ColumnAttack,
		logic.NewColumnBonusAttack("Bonus vs Archer, CA", logic.ArmorClassArchers, logic.ArmorClassCavalryArchers),
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnReloadTime,
		logic.ColumnSpeed,
		logic.ColumnNote,
	}
)

func main() {
	rows := logic.LoadRowsUpgraded(tree, skirmisherLine)
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
