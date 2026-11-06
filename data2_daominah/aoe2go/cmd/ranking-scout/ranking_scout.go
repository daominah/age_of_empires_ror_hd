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

	// castle11Plus4 is a fully upgraded Castle arrow:
	// 11 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	castle11Plus4 = logic.Attacker{Name: "Castle 11+4", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 15}}
	// castle11Plus3 is a Castle arrow without Bracer.
	castle11Plus3 = logic.Attacker{Name: "Castle 11+3", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 14}}
	// townCenter5Plus4 is a fully upgraded Town Center arrow:
	// 5 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	townCenter5Plus4 = logic.Attacker{Name: "Town Center 5+4", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 9}}
	// arbalester6Plus4 is a fully upgraded Arbalester: 6 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	arbalester6Plus4 = tree.LoadAttackerUpgraded(logic.Arbalester)

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// Reorder or comment out lines to rank differently,
	// e.g. put logic.NewCriterionTankVs(arbalester6Plus4) first.
	criteria = []logic.Criterion{
		logic.NewCriterionTankVs(castle11Plus4),
		logic.NewCriterionTankVs(arbalester6Plus4),
		// cost saving from a unique technology only breaks ties of the same cost for free,
		// e.g. Magyars and Berbers above Gurjaras (Kshatriyas)
		logic.CriterionCostBeforeUniqueTechnologies,
		logic.CriterionCost,
		logic.CriterionAttack,
		logic.NewCriterionTankVs(castle11Plus3),
		logic.NewCriterionTankVs(townCenter5Plus4),
		// last, a strength from a unique technology only breaks the remaining ties,
		// e.g. Turks above Tatars (Silk Armor)
		logic.NewCriterionTankVsBeforeUniqueTechnologies(castle11Plus4),
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		logic.NewColumnFromCriterion(logic.NewCriterionTankVs(castle11Plus4)),
		logic.NewColumnFromCriterion(logic.NewCriterionTankVs(townCenter5Plus4)),
		logic.NewColumnFromCriterion(logic.NewCriterionTankVs(arbalester6Plus4)),
		logic.ColumnCost,
		logic.ColumnAttack,
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnReloadTime,
		logic.ColumnSpeed,
		logic.ColumnNote,
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
