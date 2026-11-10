// Command building-bonus writes a report of the units with an attack against the building armor classes,
// and how Masonry and Architecture change the damage they deal to a building,
// as a Markdown file next to its source file.
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/logic"
)

// outputFileName is the Markdown report written next to this source file.
const outputFileName = "building_bonus.md"

// Report config (edit before running):
var (
	// buildingClasses are the armor classes only buildings have,
	// so an attack against them is a bonus attack against buildings.
	buildingClasses = []logic.ArmorClassLabel{
		{Class: logic.ArmorClassAllBuildings, Label: "All Buildings"},
		{Class: logic.ArmorClassStandardBuildings, Label: "Standard Buildings"},
		{Class: logic.ArmorClassStoneDefense, Label: "Stone Defense"},
		{Class: logic.ArmorClassWallsAndGates, Label: "Walls and Gates"},
		{Class: logic.ArmorClassCastles, Label: "Castles"},
	}
	// warshipClass is the unit class of the ships, left out: the report is about the land units.
	warshipClass logic.UnitClass = 22

	// armorBuildings are shown with their HP and armors, to see which classes each kind of building has.
	// A building that the age technologies replace is its Imperial Age unit, as the attackers are fully upgraded:
	// the Dark Age unit (e.g. TownCenter, 3 melee and 5 pierce armor) has less armor than the Imperial Age one (6 and 8).
	// The walls are the control: Masonry and Architecture do not change them.
	armorBuildings = []logic.UnitID{
		logic.BarracksImperialAge, logic.TownCenterImperialAge, logic.Castle, logic.MonasteryImperialAge,
		logic.Keep, logic.StoneWall, logic.FortifiedWall,
	}
	// targetBuildings are attacked in the damage tables: the Town Center, the Castle and the Monastery,
	// the buildings that a raid or a siege goes for.
	targetBuildings = []logic.UnitID{logic.TownCenterImperialAge, logic.Castle, logic.MonasteryImperialAge}

	// steps are the states of the target building compared in the damage tables.
	// Architecture requires Masonry, so its step has both.
	steps = []step{
		{Name: "None"},
		{Name: "Masonry", TechnologyIDs: []logic.TechID{logic.Masonry}},
		{Name: "Architecture", TechnologyIDs: []logic.TechID{logic.Masonry, logic.Architecture}},
	}
)

func main() {
	attackers := listAttackers()
	var b strings.Builder
	b.WriteString(reportIntro)
	b.WriteString("\n## Building armor classes\n\n")
	b.WriteString(armorClassIntro)
	b.WriteString("\n")
	b.WriteString(logic.RenderTable(buildArmorTable()))
	b.WriteString("\n## Units with an attack against buildings\n\n")
	b.WriteString(attackersIntro)
	b.WriteString("\n")
	b.WriteString(logic.RenderTable(buildAttackerTable(attackers)))
	for _, id := range targetBuildings {
		b.WriteString(fmt.Sprintf("\n## Damage to a %v\n\n", tree.Units[id].Name))
		b.WriteString(damageIntro)
		b.WriteString("\n")
		b.WriteString(logic.RenderTable(buildDamageTable(attackers, id)))
	}

	outputPath := filepath.Join(findSourceDir(), outputFileName)
	if err := os.WriteFile(outputPath, []byte(b.String()), 0o644); err != nil {
		log.Fatalf("error os.WriteFile: %v", err)
	}
	log.Printf("wrote %v", outputPath)
}

// listAttackers is the fully upgraded land units that any civilization trains,
// with an attack above 0 against a building class, once per name and attacks,
// sorted by the strongest building attack.
// A unit repeats when several civilizations have their own copy, e.g. the Huskarl.
func listAttackers() []logic.Unit {
	trained := map[logic.UnitID]bool{}
	for _, civilization := range tree.Civilizations {
		for id := range civilization.UnitIDs {
			trained[id] = true
		}
	}
	ids := make([]logic.UnitID, 0, len(trained))
	for id := range trained {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	var result []logic.Unit
	seen := map[string]bool{}
	for _, id := range ids {
		if tree.Units[id].Class == warshipClass {
			continue
		}
		unit := tree.LoadUnitUpgraded(id)
		if getStrongestBuildingAttack(unit) <= 0 {
			continue
		}
		key := fmt.Sprint(unit.Name, unit.Attacks, unit.ExtraProjectiles)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, unit)
	}
	slices.SortStableFunc(result, func(a, b logic.Unit) int {
		return -compareFloat(getStrongestBuildingAttack(a), getStrongestBuildingAttack(b))
	})
	return result
}

// getStrongestBuildingAttack is the largest attack of the unit against the building classes.
func getStrongestBuildingAttack(unit logic.Unit) float64 {
	strongest := 0.0
	for _, class := range buildingClasses {
		strongest = max(strongest, unit.Attacks[class.Class])
	}
	return strongest
}

// buildArmorTable is the HP and armors of armorBuildings, "None → Masonry → Architecture" for a value that changes.
func buildArmorTable() [][]string {
	header := []string{"Building", "HP", "Melee", "Pierce"}
	for _, class := range buildingClasses {
		header = append(header, class.Label)
	}
	table := [][]string{header}
	for _, id := range armorBuildings {
		stages := loadStages(id)
		row := []string{stages[0].Name}
		row = append(row, formatStages(stages, func(u logic.Unit) (float64, bool) { return u.HP, true }))
		row = append(row, formatStages(stages, func(u logic.Unit) (float64, bool) {
			amount, isFound := u.Armors[logic.ArmorClassMelee]
			return amount, isFound
		}))
		row = append(row, formatStages(stages, func(u logic.Unit) (float64, bool) {
			amount, isFound := u.Armors[logic.ArmorClassPierce]
			return amount, isFound
		}))
		for _, class := range buildingClasses {
			row = append(row, formatStages(stages, func(u logic.Unit) (float64, bool) {
				amount, isFound := u.Armors[class.Class]
				return amount, isFound
			}))
		}
		table = append(table, row)
	}
	return table
}

// buildAttackerTable is the attack of each unit against each building class, empty when 0.
func buildAttackerTable(attackers []logic.Unit) [][]string {
	header := []string{"Unit"}
	for _, class := range buildingClasses {
		header = append(header, class.Label)
	}
	table := [][]string{header}
	for _, unit := range attackers {
		row := []string{unit.Name}
		for _, class := range buildingClasses {
			amount := unit.Attacks[class.Class]
			if amount <= 0 {
				row = append(row, "")
				continue
			}
			row = append(row, formatNumber(amount))
		}
		table = append(table, row)
	}
	return table
}

// buildDamageTable is the damage per hit and the hits to kill of each attacker against the building,
// at each step, sorted by how much the hits grow from the first to the last step.
func buildDamageTable(attackers []logic.Unit, targetID logic.UnitID) [][]string {
	targets := loadStages(targetID)
	header := []string{"Unit"}
	for _, t := range targets {
		header = append(header, "Damage "+t.StepName)
	}
	for _, t := range targets {
		header = append(header, "Hits "+t.StepName)
	}
	header = append(header, "Hits ratio")
	type damageRow struct {
		cells []string
		ratio float64
		// firstDamage breaks the ties of the ratio
		firstDamage float64
	}
	var rows []damageRow
	for _, attacker := range attackers {
		row := []string{attacker.Name}
		for _, t := range targets {
			row = append(row, formatNumber(logic.CalculateDamagePerHit(attacker, t.Unit)))
		}
		var hits []int
		for _, t := range targets {
			hits = append(hits, logic.CalculateHitsToKill(attacker, t.Unit))
			row = append(row, strconv.Itoa(hits[len(hits)-1]))
		}
		ratio := float64(hits[len(hits)-1]) / float64(hits[0])
		row = append(row, "x"+formatNumber(ratio))
		rows = append(rows, damageRow{row, ratio, logic.CalculateDamagePerHit(attacker, targets[0].Unit)})
	}
	slices.SortStableFunc(rows, func(a, b damageRow) int {
		if c := compareFloat(b.ratio, a.ratio); c != 0 {
			return c
		}
		return compareFloat(b.firstDamage, a.firstDamage)
	})
	table := [][]string{header}
	for _, r := range rows {
		table = append(table, r.cells)
	}
	return table
}

// stage is a building at one step.
type stage struct {
	Name     string // the building name
	StepName string
	Unit     logic.Unit
}

// loadStages is the building at each of the steps.
func loadStages(id logic.UnitID) []stage {
	var result []stage
	for _, s := range steps {
		unit := tree.LoadUnitWithTechnologies(id, s.TechnologyIDs...)
		result = append(result, stage{Name: unit.Name, StepName: s.Name, Unit: unit})
	}
	return result
}

// step is the technologies researched in a state of the building, in order.
type step struct {
	Name          string
	TechnologyIDs []logic.TechID
}

// formatStages shows a stat at each step joined by arrows when it changes, e.g. "0 → 3 → 6",
// or one value when it does not, empty if the building does not have the stat.
func formatStages(stages []stage, getStat func(u logic.Unit) (float64, bool)) string {
	var values []string
	isChanged := false
	for i, s := range stages {
		value, isFound := getStat(s.Unit)
		if !isFound {
			return ""
		}
		values = append(values, formatNumber(value))
		if i > 0 && values[i] != values[0] {
			isChanged = true
		}
	}
	if !isChanged {
		return values[0]
	}
	return strings.Join(values, " → ")
}

// formatNumber rounds to 2 decimals, without trailing zeros.
func formatNumber(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// compareFloat is -1, 0 or 1.
func compareFloat(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// findSourceDir is the directory of this source file,
// so the output lands next to it whatever the working directory of "go run".
func findSourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

// tree is the embedded game data, parsed before the config that uses it.
var tree = parseTechTreeOrExit()

// parseTechTreeOrExit parses the embedded game data, it exits on error.
func parseTechTreeOrExit() *logic.TechTree {
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	return tree
}

const reportIntro = `# Building bonus attacks, Masonry and Architecture

## Summary

- **Every unit that attacks a building loses damage**, even a unit with a bonus attack:
  each technology gives the building +1 melee armor, +1 pierce armor and +10% HP.
- **The Standard Buildings bonus is not reduced**, so most infantry (regional and unique included)
  keep that part of their hit, as every building has 0 armor in that class.
  The bonus is small (mostly 2 to 6) next to the melee part of the hit,
  so infantry still need more hits: the melee part loses 1 per technology and the HP grows 10%.
- **The Stone Defense, Walls and Gates, and Castles bonuses are not reduced either**,
  for the same reason (elephants, Bombard Cannon, Petard, Tarkan).
- **The All Buildings bonus is reduced by 3 per technology**, 6 in total, which is small for a siege weapon
  (+150 to +200) but wipes out a small one, e.g. the Elite Battle Elephant +7 is 1 vs the Town Center.
  The game text says "+3 building armor" without a class: it means All Buildings,
  not Standard Buildings (the infantry bonus class), which the technologies do not change.
- **The Rocket Carts lose the most**: each rocket loses 3 per technology,
  and a cart fires 8 or 10 rockets, so one hit loses 24 to 30 per technology.
  Their armor ignoring attack does not work on buildings.
- **The walls are not changed by the technologies.**

A hit deals the sum over the armor classes that both the attack and the target list,
each one attack minus armor, never below 0, and at least 1 in total.
So a bonus attack against a building only counts if the building lists that armor class.

Masonry and Architecture are University technologies.
Each one gives the buildings and the towers (not the walls) +10% HP,
+1 melee armor, +1 pierce armor and +3 armor of the class All Buildings.
Architecture requires Masonry, so the Architecture step has both researched.
The game text is "Buildings +10% HP, +1 melee/+1 pierce armor and +3 building armor",
and the data tells the class of the last one: the armor class All Buildings (11).

The buildings are the Imperial Age ones: the age technologies replace a building with a unit of more armor,
e.g. the Town Center has melee and pierce armor 3 and 5 in the Dark Age and 6 and 8 in the Imperial Age.

Generated by cmd/building-bonus from data/dat_original_to_json, the units have all the technologies
in logic.TechnologyEffects (Forging line, etc.), without civilization bonuses.
`

const armorClassIntro = `A building has several armor classes, and a unit attacks each of them with its own amount.
A value shows None, then Masonry, then Architecture when the technologies change it.

- **Melee and Pierce** are the normal armors, they cut a swordsman melee attack or an archer pierce attack.
  Masonry and Architecture raise each by 1.
- **All Buildings** is the bonus class of the siege weapons and the ships,
  e.g. the Siege Ram +200, the Mangonel +35, the Bombard Cannon +200.
  Masonry and Architecture raise it by 3 each.
  The base value is 0 for most buildings, 8 for the Castle and 16 or more for the walls.
- **Standard Buildings** is the bonus class of the infantry and a few other units.
  Every building has 0 and the technologies do not change it,
  so this bonus is never reduced.
- **Stone Defense** is the bonus class of the elephants and the Bombard Cannon.
  Only the towers, the walls, the Gate, the Harbor and the Fortified Church have it, always 0.
- **Walls and Gates** is the bonus class of the Petard, only the walls and the gates have it, always 0.
- **Castles** is the bonus class of the Petard and the Tarkan, only the Castle and the Krepost have it, always 0.
`

const attackersIntro = `The attack against each building class, after the Blacksmith attack technologies.
A blank is 0 or no attack in that class.
The units are sorted by the strongest building attack.
`

const damageIntro = `Damage is one hit with all the attack classes added, and hits is how many hits destroy the building.
The units are sorted by the hits ratio, the hits with Architecture divided by the hits without any technology,
so the units that the technologies hurt the most are first.
`
