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
	// arbalester6Plus4 is a fully upgraded Arbalester: 6 + Fletching, Bodkin Arrow, Bracer.
	arbalester6Plus4 = logic.Attacker{Name: "Arbalester 6+4", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 10}}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// Reorder or comment out lines to rank differently,
	// e.g. put arrowsToKill(arbalester6Plus4) first.
	criteria = []logic.Criterion{
		arrowsToKill(castle11Plus4),
		arrowsToKill(arbalester6Plus4),
		// cost saving from a unique technology only breaks ties of the same cost for free,
		// e.g. Magyars and Berbers above Gurjaras (Kshatriyas)
		{Name: "Cost before unique tech", Value: func(r logic.Row) float64 { return r.UnitWithoutUniqueTechnologies.Cost.Total() }, IsHigherBetter: false},
		{Name: "Cost", Value: func(r logic.Row) float64 { return r.Unit.Cost.Total() }, IsHigherBetter: false},
		{Name: "Attack", Value: func(r logic.Row) float64 { return r.Unit.Attacks[logic.ArmorClassMelee] }, IsHigherBetter: true},
		arrowsToKill(castle11Plus3),
		arrowsToKill(townCenter5Plus4),
		// last, a strength from a unique technology only breaks the remaining ties,
		// e.g. Turks above Tatars (Silk Armor)
		arrowsToKillBeforeUniqueTechnologies(castle11Plus4),
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		logic.CriterionColumn(arrowsToKill(castle11Plus4)),
		logic.CriterionColumn(arrowsToKill(townCenter5Plus4)),
		logic.CriterionColumn(arrowsToKill(arbalester6Plus4)),
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
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	rows := logic.FullyUpgradedRows(tree, scoutLine)
	for _, line := range extraLines {
		rows = append(rows, logic.FullyUpgradedRows(tree, line)...)
	}
	logic.SortRows(rows, criteria)
	table := logic.MarkdownTable(rows, columns)
	fmt.Print(table)

	outputPath := filepath.Join(sourceDir(), outputFileName)
	if err := os.WriteFile(outputPath, []byte(table), 0o644); err != nil {
		log.Fatalf("error os.WriteFile: %v", err)
	}
	log.Printf("wrote %v", outputPath)
}

// sourceDir is the directory of this source file,
// so the output lands next to it whatever the working directory of "go run".
func sourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

// arrowsToKillBeforeUniqueTechnologies is arrowsToKill on the unit without its unique technologies.
func arrowsToKillBeforeUniqueTechnologies(attacker logic.Attacker) logic.Criterion {
	return logic.Criterion{
		Name:           "Arrows " + attacker.Name + " before unique tech",
		Value:          func(r logic.Row) float64 { return float64(logic.HitsToKill(attacker, r.UnitWithoutUniqueTechnologies)) },
		IsHigherBetter: true,
	}
}

// arrowsToKill is a criterion: hits of the attacker to kill the unit, more is better.
func arrowsToKill(attacker logic.Attacker) logic.Criterion {
	return logic.Criterion{
		Name:           "Arrows " + attacker.Name,
		Value:          func(r logic.Row) float64 { return float64(logic.HitsToKill(attacker, r.Unit)) },
		IsHigherBetter: true,
	}
}
