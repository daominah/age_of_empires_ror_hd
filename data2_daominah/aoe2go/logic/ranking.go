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

// Column is a column of the ranking table, apart from the sorting criteria.
type Column struct {
	Name  string
	Value func(row Row) string
	// IsFootnote shows a long value as a numbered note below the table, e.g. "[1]" in the cell.
	IsFootnote bool
}

// Columns of the unit stats.
var (
	ColumnHP          = Column{Name: "HP", Value: func(row Row) string { return formatNumber(GetShownHP(row.Unit)) }}
	ColumnMeleeArmor  = Column{Name: "Melee armor", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getMeleeArmor) }}
	ColumnPierceArmor = Column{Name: "Pierce armor", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getPierceArmor) }}
	ColumnAttack      = Column{Name: "Attack", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getAttack) }}
	ColumnRange       = Column{Name: "Range", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getRange) }}
	ColumnReloadTime  = Column{Name: "Reload time", Value: func(row Row) string { return formatNumber(row.Unit.ReloadTime) }}
	ColumnSpeed       = Column{Name: "Speed", Value: func(row Row) string { return formatNumber(row.Unit.Speed) }}
	ColumnCost        = Column{Name: "Cost", Value: func(row Row) string { return formatCost(row.Unit.Cost) }}
	ColumnNote        = Column{Name: "Note", Value: FormatNote, IsFootnote: true}
)

// NewColumnBonusAttack shows the bonus attacks against armor classes, joined by ", ",
// e.g. the Elite Skirmisher "Bonus vs Archer, CA" is "4, 2".
// Each value is the total after upgrades, not base+upgrade.
func NewColumnBonusAttack(name string, classes ...ArmorClass) Column {
	return Column{
		Name: name,
		Value: func(row Row) string {
			var values []string
			for _, class := range classes {
				values = append(values, formatNumber(row.Unit.Attacks[class]))
			}
			return strings.Join(values, ", ")
		},
	}
}

// Note is FormatCostNote then FormatAgeNote.
func FormatNote(row Row) string {
	var notes []string
	for _, note := range []string{FormatCostNote(row), FormatAgeNote(row)} {
		if note != "" {
			notes = append(notes, note)
		}
	}
	return strings.Join(notes, "; ")
}

// FormatAgeNote shows a stat that civilization bonuses improve in several ages,
// with the value after the bonuses of each age, e.g. Berbers "cost 68f in Castle Age, 64f in Imperial Age".
// The table shows the last one.
func FormatAgeNote(row Row) string {
	var notes []string
	for _, stat := range []struct {
		name  string
		value func(u Unit) string
	}{
		{"cost", func(u Unit) string { return formatCost(u.Cost) }},
		{"HP", func(u Unit) string { return formatNumber(GetShownHP(u)) }},
		{"attack", func(u Unit) string { return formatNumber(getAttack(u)) }},
		{"melee armor", func(u Unit) string { return formatNumber(getMeleeArmor(u)) }},
		{"pierce armor", func(u Unit) string { return formatNumber(getPierceArmor(u)) }},
		{"range", func(u Unit) string { return formatNumber(u.Range) }},
		{"reload time", func(u Unit) string { return formatNumber(u.ReloadTime) }},
		{"speed", func(u Unit) string { return formatNumber(u.Speed) }},
	} {
		var ages []TechID
		valueByAge := map[TechID]string{}
		for _, step := range row.Steps {
			age := step.Technology.FindAge()
			if step.Source != CivilizationBonuses || age == 0 || stat.value(step.UnitBefore) == stat.value(step.UnitAfter) {
				continue
			}
			if _, isFound := valueByAge[age]; !isFound {
				ages = append(ages, age)
			}
			valueByAge[age] = stat.value(step.UnitAfter)
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

// FormatCostNote explains a cost changed by a technology to research, empty if none,
// e.g. "base cost 35f 45g, 80f after Corvinian Army (200f 300g)".
// Civilization bonuses always apply, so they need no note, e.g. Magyars Hussar 68f.
func FormatCostNote(row Row) string {
	var changes []string
	for _, step := range row.Steps {
		if step.Source == CivilizationBonuses || step.UnitBefore.Cost == step.UnitAfter.Cost {
			continue
		}
		changes = append(changes, fmt.Sprintf("%v after %v (%v)", formatCost(step.UnitAfter.Cost), step.Name, formatCost(step.Technology.Cost)))
	}
	if len(changes) == 0 {
		return ""
	}
	return "base cost " + formatCost(row.BaseUnit.Cost) + ", " + strings.Join(changes, ", ")
}

func getMeleeArmor(u Unit) float64  { return u.Armors[ArmorClassMelee] }
func getPierceArmor(u Unit) float64 { return u.Armors[ArmorClassPierce] }
func getRange(u Unit) float64       { return u.Range }

// attack is the main attack of the unit: melee, or pierce for a unit without a melee attack, e.g. a skirmisher.
func getAttack(u Unit) float64 {
	if melee, isFound := u.Attacks[ArmorClassMelee]; isFound {
		return melee
	}
	return u.Attacks[ArmorClassPierce]
}

// Step is an effect applied to a row unit, with the unit before and after it.
type Step struct {
	Name       string     // the technology name, or the effect name of a civilization bonus
	Technology Technology // the researched technology, empty for a civilization bonus
	Source     UpgradeSource
	UnitBefore Unit
	UnitAfter  Unit
}

// UpgradeSource groups the effects shown as one term of a base+upgrade stat.
type UpgradeSource int

// UpgradeSource values, in the order of their term.
const (
	Technologies UpgradeSource = iota
	UniqueTechnologies
	CivilizationBonuses
)

// FormatBasePlusUpgrade shows a stat as the base, plus the technologies,
// then plus the unique technologies and the civilization bonuses when they change it,
// e.g. Malians attack "7+2+5" with Farimba, Khitans "7+2+2" with their double Forging and Iron Casting.
func FormatBasePlusUpgrade(row Row, stat func(u Unit) float64) string {
	upgrades := map[UpgradeSource]float64{}
	for _, step := range row.Steps {
		upgrades[step.Source] += stat(step.UnitAfter) - stat(step.UnitBefore)
	}
	result := formatNumber(stat(row.BaseUnit)) + formatSigned(upgrades[Technologies])
	for _, source := range []UpgradeSource{UniqueTechnologies, CivilizationBonuses} {
		if math.Abs(upgrades[source]) > 1e-6 {
			result += formatSigned(upgrades[source])
		}
	}
	return result
}

// formatSigned is the number with its sign, e.g. "+2", "-1", "+0".
func formatSigned(v float64) string {
	if v < 0 {
		return formatNumber(v)
	}
	return "+" + formatNumber(v)
}

// LoadRowsUpgraded returns one row per civilization that has the line,
// with every available technology and the civilization bonuses applied.
func LoadRowsUpgraded(tree *TechTree, line Line) []Row {
	var rows []Row
	for _, civilization := range tree.Civilizations {
		unitID := UnitID(-1)
		for _, id := range line.UnitIDs {
			if civilization.UnitIDs[id] {
				unitID = id
			}
		}
		if _, isFound := tree.Units[unitID]; !isFound {
			continue
		}
		// the Scout Cavalry is class 47 in the data, the effects copied in effect.go use class 12 for the whole line
		rows = append(rows, tree.loadRowUpgraded(civilization, unitID, line.Class))
	}
	return rows
}

// loadRowUpgraded is the unit of the civilization
// with every technology it has in its Technology Tree and its civilization bonuses,
// step by step. The class replaces the unit class in the data.
func (tree *TechTree) loadRowUpgraded(civilization Civilization, unitID UnitID, class UnitClass) Row {
	// the automatic technologies are part of the base, e.g. Scout Cavalry attack 5 from the Feudal Age
	unit := tree.LoadUnitBase(unitID)
	unit.Class = class
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
		row.Steps = append(row.Steps, Step{Name: name, Technology: technology, Source: source, UnitBefore: before, UnitAfter: unit.clone()})
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
		technology := tree.findCivilizationBonusTechnology(civilization.ID, effect.EffectID)
		if effect.TechnologyID != 0 {
			technology = tree.Technologies[effect.TechnologyID]
		}
		apply(effect, CivilizationBonuses, technology)
	}
	row.Unit = unit
	row.UnitWithoutUniqueTechnologies = withoutUniqueTechnologies
	return row
}

// CalculateDamagePerHit is the damage of one hit of the attacker on the target:
// the sum over the armor classes of both of attack minus armor (not below 0), at least 1.
func CalculateDamagePerHit(attacker Attacker, target Unit) float64 {
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

// CalculateHitsToKill is how many hits of the attacker kill the target.
// It uses the HP with decimals, e.g. Mongols Hussar 117.5 HP, shown as 118 in game.
func CalculateHitsToKill(attacker Attacker, target Unit) int {
	return int(math.Ceil(cleanHP(target) / CalculateDamagePerHit(attacker, target)))
}

// GetShownHP is the HP shown in game, rounded to nearest (half up),
// e.g. Mongols Hussar 117.5 HP is shown as 118.
func GetShownHP(u Unit) float64 {
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

// RenderMarkdownTable prints the rows as a Markdown table:
// the rank, civilization and unit, then the columns in order,
// then the footnotes of the IsFootnote columns.
func RenderMarkdownTable(rows []Row, columns []Column) string {
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
	result := alignTable(table)
	if len(footnotes) > 0 {
		result += "\n" + strings.Join(footnotes, "\n") + "\n"
	}
	return result
}

// alignTable pads each column to its widest cell,
// so the raw Markdown reads as a table too.
// The first line is the header.
func alignTable(table [][]string) string {
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
