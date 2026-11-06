package logic

import (
	"fmt"
	"maps"
	"math"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Line is a unit line ranked across civilizations.
type Line struct {
	// UnitIDs from the first unit to the last upgrade,
	// e.g. Scout Cavalry, Light Cavalry, Hussar, Winged Hussar.
	// A civilization gets the last one available in its Technology Tree.
	UnitIDs []UnitID
	Class   UnitClass // unit class in the .dat, e.g. ClassCavalry
}

// Attacker is what hits the line units, with its fully upgraded attack.
type Attacker struct {
	Name    string
	Attacks map[ArmorClass]float64 // amount by armor class, e.g. ArmorClassPierce
}

// Row is a civilization's fully upgraded unit of the line.
type Row struct {
	CivilizationID CivilizationID
	Civilization   string // English name
	Unit           Unit
	BaseUnit       Unit   // the unit before any technology or civilization bonus
	Steps          []Step // the applied effects, in order
	// UnitWithoutUniqueTechnologies is the unit without its unique technologies,
	// to rank a strength that costs a unique technology below the same strength for free,
	// e.g. Tatars Silk Armor against Turks, Gurjaras Kshatriyas against Magyars.
	UnitWithoutUniqueTechnologies Unit
}

// Criterion sorts the rows, in the configured order.
type Criterion struct {
	Name           string
	Value          func(row Row) float64
	IsHigherBetter bool
}

// Column is a column of the ranking table, apart from the sorting criteria.
type Column struct {
	Name  string
	Value func(row Row) string
	// IsFootnote shows a long value as a numbered note below the table, e.g. "[1]" in the cell.
	IsFootnote bool
}

// CriterionColumn shows the value of a criterion, e.g. the arrows to kill.
func CriterionColumn(criterion Criterion) Column {
	return Column{Name: criterion.Name, Value: func(row Row) string { return formatNumber(criterion.Value(row)) }}
}

// Columns of the unit stats.
var (
	ColumnHP          = Column{Name: "HP", Value: func(row Row) string { return formatNumber(ShownHP(row.Unit)) }}
	ColumnMeleeArmor  = Column{Name: "Melee armor", Value: func(row Row) string { return BasePlusUpgrade(row, meleeArmor) }}
	ColumnPierceArmor = Column{Name: "Pierce armor", Value: func(row Row) string { return BasePlusUpgrade(row, pierceArmor) }}
	ColumnAttack      = Column{Name: "Attack", Value: func(row Row) string { return BasePlusUpgrade(row, meleeAttack) }}
	ColumnReloadTime  = Column{Name: "Reload time", Value: func(row Row) string { return formatNumber(row.Unit.ReloadTime) }}
	ColumnSpeed       = Column{Name: "Speed", Value: func(row Row) string { return formatNumber(row.Unit.Speed) }}
	ColumnCost        = Column{Name: "Cost", Value: func(row Row) string { return formatCost(row.Unit.Cost) }}
	ColumnNote        = Column{Name: "Note", Value: Note, IsFootnote: true}
)

// Note is CostNote then AgeNote.
func Note(row Row) string {
	var notes []string
	for _, note := range []string{CostNote(row), AgeNote(row)} {
		if note != "" {
			notes = append(notes, note)
		}
	}
	return strings.Join(notes, "; ")
}

// AgeNote shows a stat that civilization bonuses improve in several ages,
// with the value after the bonuses of each age, e.g. Berbers "cost 68f in Castle Age, 64f in Imperial Age".
// The table shows the last one.
func AgeNote(row Row) string {
	var notes []string
	for _, stat := range []struct {
		name  string
		value func(u Unit) string
	}{
		{"cost", func(u Unit) string { return formatCost(u.Cost) }},
		{"HP", func(u Unit) string { return formatNumber(ShownHP(u)) }},
		{"attack", func(u Unit) string { return formatNumber(meleeAttack(u)) }},
		{"melee armor", func(u Unit) string { return formatNumber(meleeArmor(u)) }},
		{"pierce armor", func(u Unit) string { return formatNumber(pierceArmor(u)) }},
		{"reload time", func(u Unit) string { return formatNumber(u.ReloadTime) }},
		{"speed", func(u Unit) string { return formatNumber(u.Speed) }},
	} {
		var ages []TechID
		valueByAge := map[TechID]string{}
		for _, step := range row.Steps {
			age := step.Technology.Age()
			if step.Source != CivilizationBonuses || age == 0 || stat.value(step.Before) == stat.value(step.After) {
				continue
			}
			if _, isFound := valueByAge[age]; !isFound {
				ages = append(ages, age)
			}
			valueByAge[age] = stat.value(step.After)
		}
		if len(ages) < 2 {
			continue
		}
		var parts []string
		for _, age := range ages {
			parts = append(parts, valueByAge[age]+" in "+AgeNames[age])
		}
		notes = append(notes, stat.name+" "+strings.Join(parts, ", "))
	}
	return strings.Join(notes, "; ")
}

// CostNote explains a cost changed by a technology to research, empty if none,
// e.g. "base cost 35f 45g, 80f after Corvinian Army (200f 300g)".
// Civilization bonuses always apply, so they need no note, e.g. Magyars Hussar 68f.
func CostNote(row Row) string {
	var changes []string
	for _, step := range row.Steps {
		if step.Source == CivilizationBonuses || step.Before.Cost == step.After.Cost {
			continue
		}
		changes = append(changes, fmt.Sprintf("%v after %v (%v)", formatCost(step.After.Cost), step.Name, formatCost(step.Technology.Cost)))
	}
	if len(changes) == 0 {
		return ""
	}
	return "base cost " + formatCost(row.BaseUnit.Cost) + ", " + strings.Join(changes, ", ")
}

func meleeArmor(u Unit) float64  { return u.Armors[ArmorClassMelee] }
func pierceArmor(u Unit) float64 { return u.Armors[ArmorClassPierce] }
func meleeAttack(u Unit) float64 { return u.Attacks[ArmorClassMelee] }

// Step is an effect applied to a row unit, with the unit before and after it.
type Step struct {
	Name       string     // the technology name, or the effect name of a civilization bonus
	Technology Technology // the researched technology, empty for a civilization bonus
	Source     UpgradeSource
	Before     Unit
	After      Unit
}

// UpgradeSource groups the effects shown as one term of a base+upgrade stat.
type UpgradeSource int

// UpgradeSource values, in the order of their term.
const (
	Technologies UpgradeSource = iota
	UniqueTechnologies
	CivilizationBonuses
)

// BasePlusUpgrade shows a stat as the base, plus the technologies,
// then plus the unique technologies and the civilization bonuses when they change it,
// e.g. Malians attack "7+2+5" with Farimba, Khitans "7+2+2" with their double Forging and Iron Casting.
func BasePlusUpgrade(row Row, stat func(u Unit) float64) string {
	upgrades := map[UpgradeSource]float64{}
	for _, step := range row.Steps {
		upgrades[step.Source] += stat(step.After) - stat(step.Before)
	}
	result := formatNumber(stat(row.BaseUnit)) + signed(upgrades[Technologies])
	for _, source := range []UpgradeSource{UniqueTechnologies, CivilizationBonuses} {
		if math.Abs(upgrades[source]) > 1e-6 {
			result += signed(upgrades[source])
		}
	}
	return result
}

// signed is the number with its sign, e.g. "+2", "-1", "+0".
func signed(v float64) string {
	if v < 0 {
		return formatNumber(v)
	}
	return "+" + formatNumber(v)
}

// FullyUpgradedRows returns one row per civilization that has the line,
// with every available technology and the civilization bonuses applied.
func FullyUpgradedRows(tree *TechTree, line Line) []Row {
	var rows []Row
	for _, civilization := range tree.Civilizations {
		unitID := UnitID(-1)
		for _, id := range line.UnitIDs {
			if civilization.UnitIDs[id] {
				unitID = id
			}
		}
		base, isFound := tree.Units[unitID]
		if !isFound {
			continue
		}
		unit := base.clone()
		unit.Class = line.Class
		// the automatic technologies are part of the base, e.g. Scout Cavalry attack 5 from the Feudal Age
		for id := range AutomaticTechnologies {
			TechnologyEffects[id].Apply(&unit)
		}
		withoutUniqueTechnologies := unit.clone()
		row := Row{CivilizationID: civilization.ID, Civilization: civilization.Name, BaseUnit: unit.clone()}

		var technologyIDs []TechID
		for id := range TechnologyEffects {
			if civilization.TechnologyIDs[id] && !AutomaticTechnologies[id] {
				technologyIDs = append(technologyIDs, id)
			}
		}
		sort.Slice(technologyIDs, func(i, j int) bool { return technologyIDs[i] < technologyIDs[j] }) // a fixed order, e.g. Bloodlines before the Mongols HP bonus
		apply := func(effect Effect, source UpgradeSource, technology Technology) {
			before := unit.clone()
			effect.Apply(&unit)
			if source != UniqueTechnologies {
				effect.Apply(&withoutUniqueTechnologies)
			}
			name := technology.Name
			if source == CivilizationBonuses || name == "" {
				name = effect.Name
			}
			row.Steps = append(row.Steps, Step{Name: name, Technology: technology, Source: source, Before: before, After: unit.clone()})
		}
		for _, id := range technologyIDs {
			effect := TechnologyEffects[id]
			source := Technologies
			if effect.IsUniqueTechnology {
				source = UniqueTechnologies
			}
			apply(effect, source, tree.Technologies[id])
		}
		for _, effect := range CivilizationBonusEffects[civilization.ID] {
			apply(effect, CivilizationBonuses, tree.civilizationBonusTechnology(civilization.ID, effect.EffectID))
		}
		row.Unit = unit
		row.UnitWithoutUniqueTechnologies = withoutUniqueTechnologies
		rows = append(rows, row)
	}
	return rows
}

// DamagePerHit is the damage of one hit of the attacker on the target:
// the sum over the armor classes of both of attack minus armor (not below 0), at least 1.
func DamagePerHit(attacker Attacker, target Unit) float64 {
	damage := 0.0
	for class, attack := range attacker.Attacks {
		armor, isFound := target.Armors[class]
		if !isFound {
			continue
		}
		damage += math.Max(0, attack-armor)
	}
	return math.Max(1, damage)
}

// HitsToKill is how many hits of the attacker kill the target.
// It uses the HP with decimals, e.g. Mongols Hussar 117.5 HP, shown as 118 in game.
func HitsToKill(attacker Attacker, target Unit) int {
	return int(math.Ceil(cleanHP(target) / DamagePerHit(attacker, target)))
}

// ShownHP is the HP shown in game, rounded to nearest (half up),
// e.g. Mongols Hussar 117.5 HP is shown as 118.
func ShownHP(u Unit) float64 {
	return math.Floor(cleanHP(u) + 0.5)
}

// cleanHP rounds away the noise of multiplied effect values,
// e.g. the Mongols 1.08333 (for 1.08333...) makes 117.5 HP into 117.4997.
func cleanHP(u Unit) float64 {
	return math.Round(u.HP*1000) / 1000
}

// SortRows sorts by the criteria in order, the next one breaking ties of the previous,
// then by civilization name.
func SortRows(rows []Row, criteria []Criterion) {
	sort.SliceStable(rows, func(i, j int) bool {
		for _, criterion := range criteria {
			vi, vj := criterion.Value(rows[i]), criterion.Value(rows[j])
			if math.Abs(vi-vj) < 1e-6 {
				continue
			}
			if criterion.IsHigherBetter {
				return vi > vj
			}
			return vi < vj
		}
		return rows[i].Civilization < rows[j].Civilization
	})
}

// MarkdownTable prints the rows as a Markdown table:
// the rank, civilization and unit, then the columns in order,
// then the footnotes of the IsFootnote columns.
func MarkdownTable(rows []Row, columns []Column) string {
	headers := []string{"#", "Civilization", "Unit"}
	for _, column := range columns {
		headers = append(headers, column.Name)
	}
	table := [][]string{headers}
	var footnotes []string
	for i, row := range rows {
		cells := []string{strconv.Itoa(i + 1), row.Civilization, row.Unit.Name}
		for _, column := range columns {
			value := column.Value(row)
			if column.IsFootnote && value != "" {
				footnotes = append(footnotes, fmt.Sprintf("- [%v] %v %v: %v", len(footnotes)+1, row.Civilization, row.Unit.Name, value))
				value = fmt.Sprintf("[%v]", len(footnotes))
			}
			cells = append(cells, value)
		}
		table = append(table, cells)
	}
	result := alignedTable(table)
	if len(footnotes) > 0 {
		result += "\n" + strings.Join(footnotes, "\n") + "\n"
	}
	return result
}

// alignedTable pads each column to its widest cell,
// so the raw Markdown reads as a table too.
// The first line is the header.
func alignedTable(table [][]string) string {
	widths := make([]int, len(table[0]))
	for _, line := range table {
		for i, cell := range line {
			widths[i] = max(widths[i], utf8.RuneCountInString(cell))
		}
	}
	var b strings.Builder
	writeLine := func(cells []string) {
		for i, cell := range cells {
			b.WriteString("| " + cell + strings.Repeat(" ", widths[i]-utf8.RuneCountInString(cell)) + " ")
		}
		b.WriteString("|\n")
	}
	writeLine(table[0])
	for _, width := range widths {
		b.WriteString("|" + strings.Repeat("-", width+2))
	}
	b.WriteString("|\n")
	for _, line := range table[1:] {
		writeLine(line)
	}
	return b.String()
}

// formatNumber rounds to 2 decimals, without trailing zeros.
func formatNumber(v float64) string {
	return strconv.FormatFloat(math.Round(v*100)/100, 'f', -1, 64)
}

// formatCost shows the cost with the short resource names, e.g. "80f", "35f 45g".
func formatCost(c Cost) string {
	var parts []string
	for _, r := range []struct {
		name   string
		amount float64
	}{{"f", c.Food}, {"w", c.Wood}, {"g", c.Gold}, {"s", c.Stone}} {
		if math.Round(r.amount) > 0 {
			parts = append(parts, fmt.Sprintf("%v%v", math.Round(r.amount), r.name))
		}
	}
	return strings.Join(parts, " ")
}

// clone copies the unit, so applying effects does not change the base unit maps.
func (u Unit) clone() Unit {
	c := u
	c.Attacks = maps.Clone(u.Attacks)
	c.Armors = maps.Clone(u.Armors)
	return c
}
