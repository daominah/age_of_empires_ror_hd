package logic

import (
	"fmt"
	"math"
	"slices"
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

// NewRowFromUnit is a row of the unit alone, without a civilization,
// e.g. the baseline of a NewCriterionWeightedSum from LoadUnitUpgraded.
func NewRowFromUnit(unit Unit) Row {
	return Row{Unit: unit, BaseUnit: unit, UnitWithoutUniqueTechnologies: unit}
}

// Row is a civilization's fully upgraded unit of the line.
type Row struct {
	CivilizationID CivilizationID
	Civilization   string // English name
	// AvailableAge is the age that the unit becomes available, 0 if unknown,
	// so the ages before it are not noted, e.g. no Castle Age Hussar.
	AvailableAge TechID
	// EarlierRows are the rows of the units of the line before this one that the civilization also has,
	// e.g. the Light Cavalry for the Hussar, to note the stats of each age.
	EarlierRows []Row
	// UpgradeTechnologies is the technology that upgrades to the row unit, empty for the first unit of the line.
	UpgradeTechnologies []Technology
	Unit                Unit
	BaseUnit            Unit   // the unit before any technology or civilization bonus
	Steps               []Step // the applied effects, in order
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
	ColumnMeleeArmor  = Column{Name: "MArmor", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getMeleeArmor) }}
	ColumnPierceArmor = Column{Name: "PArmor", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getPierceArmor) }}
	ColumnAttack      = Column{Name: "Attack", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getAttack) + formatExtraArrows(row.Unit) }}
	ColumnRange       = Column{Name: "Range", Value: func(row Row) string { return FormatBasePlusUpgrade(row, getRange) }}
	ColumnMinRange    = Column{Name: "Min range", Value: func(row Row) string { return formatNumber(row.Unit.MinRange) }}
	ColumnReloadTime  = Column{Name: "Reload", Value: func(row Row) string { return formatNumber(row.Unit.ReloadTime) }}
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

// ArmorClassLabel is an armor class with a short label for a narrow column, e.g. "CavArcher" for the Cavalry Archers.
type ArmorClassLabel struct {
	Class ArmorClass
	Label string
}

// NewColumnBonusAttackShort shows the bonus attacks against the important classes
// that differ from the baseline unit, empty if none, so the bonuses common to the line are not repeated,
// e.g. the Elite Bolas Rider "+3 vs Cav" on top of its attack in the cavalry archer ranking,
// the Vietnamese Imperial Skirmisher "+5 vs Archers" instead of the usual +4 in the skirmisher ranking.
// It is a short version for the table, the note keeps the special bonus attacks in detail.
func NewColumnBonusAttackShort(baseline Row, classes ...ArmorClassLabel) Column {
	return Column{
		Name: "BonusATK",
		Value: func(row Row) string {
			var values []string
			for _, class := range classes {
				amount := row.Unit.Attacks[class.Class]
				if amount != baseline.Unit.Attacks[class.Class] {
					values = append(values, formatSigned(amount)+" vs "+class.Label)
				}
			}
			return strings.Join(values, ", ")
		},
	}
}

// NewColumnNoteWithBonusAttack is the note column that also notes the special bonus attacks,
// for a line without a column of bonus attacks.
// The classes every unit of the line has a bonus against are not noted, e.g. Monks for the scout line.
func NewColumnNoteWithBonusAttack(commonClasses ...ArmorClass) Column {
	return Column{
		Name:       "Note",
		IsFootnote: true,
		Value: func(row Row) string {
			return joinNotes(FormatCostNote(row), FormatTechnologyCostNote(row), FormatAgeNote(row), FormatBonusAttackNote(row, commonClasses...), FormatTraitNote(row))
		},
	}
}

// FormatNote is FormatCostNote, FormatTechnologyCostNote, FormatAgeNote then FormatTraitNote.
func FormatNote(row Row) string {
	return joinNotes(FormatCostNote(row), FormatTechnologyCostNote(row), FormatAgeNote(row), FormatTraitNote(row))
}

// joinNotes joins the notes that are not empty.
func joinNotes(notes ...string) string {
	var kept []string
	for _, note := range notes {
		if note != "" {
			kept = append(kept, note)
		}
	}
	return strings.Join(kept, "; ")
}

// FormatBonusAttackNote shows the bonus attacks of the unit against armor classes
// that the whole line has, e.g. Mangudai "bonus attack 5 vs Siege Weapons".
// Melee and pierce are the attack itself, and the classes in commonClasses are skipped.
func FormatBonusAttackNote(row Row, commonClasses ...ArmorClass) string {
	var classes []ArmorClass
	for class, amount := range row.Unit.Attacks {
		if amount > 0 && class != ArmorClassMelee && class != ArmorClassPierce && !slices.Contains(commonClasses, class) {
			classes = append(classes, class)
		}
	}
	if len(classes) == 0 {
		return ""
	}
	slices.Sort(classes)
	var parts []string
	for _, class := range classes {
		parts = append(parts, formatNumber(row.Unit.Attacks[class])+" vs "+class.String())
	}
	return "bonus attack " + strings.Join(parts, ", ")
}

// FormatAgeNote shows a stat that civilization bonuses improve in several ages,
// with the value of the unit available in each age, so the note never has a unit before its age,
// e.g. Mongols "Light Cavalry HP 110 in Castle Age, Hussar HP 118 in Imperial Age".
// A line with a single unit in the notes drops its name, e.g. Berbers "cost 68f in Castle Age, 64f in Imperial Age".
// The table shows the last one.
func FormatAgeNote(row Row) string {
	chain := append(slices.Clone(row.EarlierRows), row)
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
		// the ages where a civilization bonus changes the stat
		var ages []TechID
		for _, step := range row.Steps {
			age := step.Age
			if step.Source == CivilizationBonuses && age != 0 && !slices.Contains(ages, age) &&
				stat.value(step.UnitBefore) != stat.value(step.UnitAfter) {
				ages = append(ages, age)
			}
		}
		if len(ages) < 2 {
			continue
		}
		slices.Sort(ages)

		type agedValue struct {
			unitName string
			text     string
		}
		var values []agedValue
		var unitNames []string
		for _, age := range ages {
			// the unit of the line that is available in the age is the last one that is already available
			var available *Row
			for i := range chain {
				if chain[i].AvailableAge <= age {
					available = &chain[i]
				}
			}
			if available == nil {
				continue
			}
			values = append(values, agedValue{available.Unit.Name, stat.value(available.FindUnitAtAge(age)) + " in " + AgeNames[age]})
			if !slices.Contains(unitNames, available.Unit.Name) {
				unitNames = append(unitNames, available.Unit.Name)
			}
		}
		if len(values) < 2 {
			continue
		}
		if len(unitNames) == 1 {
			var parts []string
			for _, v := range values {
				parts = append(parts, v.text)
			}
			notes = append(notes, stat.name+" "+strings.Join(parts, ", "))
			continue
		}
		// one group for each unit: "Light Cavalry HP 110 in Castle Age", the next values of the same unit join the group
		var groups []string
		for i, v := range values {
			if i > 0 && values[i-1].unitName == v.unitName {
				groups[len(groups)-1] += ", " + v.text
			} else {
				groups = append(groups, v.unitName+" "+stat.name+" "+v.text)
			}
		}
		notes = append(notes, strings.Join(groups, ", "))
	}
	return strings.Join(notes, "; ")
}

// FindUnitAtAge is the unit after the technologies and the civilization bonuses
// that are available in the age or before, e.g. the Mongols Light Cavalry in the Castle Age.
// It applies the effects again in the order of the steps, skipping the later ages,
// as a technology of a later age is applied before the civilization bonuses,
// e.g. the Teutons Scout Cavalry has no Plate Barding Armor in the Castle Age.
func (row Row) FindUnitAtAge(age TechID) Unit {
	unit := row.BaseUnit.clone()
	for _, step := range row.Steps {
		if step.Age > age {
			continue
		}
		step.Effect.Apply(&unit)
	}
	return unit
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

// getAttack is the melee or pierce attack shown in game, the main arrow only,
// e.g. the Elite Kipchak 5+3, its extra arrows are shown apart by formatExtraArrows.
// It takes the larger of the two, as some ranged units also list a melee attack of 0,
// e.g. the Chu Ko Nu, so that melee 0 still counts against a ram negative melee armor.
func getAttack(u Unit) float64 {
	melee := u.Attacks[ArmorClassMelee]
	pierce, isRanged := u.Attacks[ArmorClassPierce]
	if !isRanged || melee > pierce {
		return melee
	}
	return pierce
}

// formatExtraArrows is the arrows besides the main one, empty if none,
// as their count x their attack, e.g. the Elite Kipchak " (+3x3 arrows)", the Maya Elite Skirmisher " (+1x1 arrow)".
// A charged attack is left out.
func formatExtraArrows(u Unit) string {
	if u.ExtraProjectiles <= 0 {
		return ""
	}
	attacks := u.ExtraAttacks
	if attacks == nil {
		attacks = u.Attacks
	}
	noun := "arrows"
	if u.ExtraProjectiles == 1 {
		noun = "arrow"
	}
	return fmt.Sprintf(" (+%vx%v %v)", u.ExtraProjectiles, formatNumber(attacks[ArmorClassPierce]), noun)
}

// Step is an effect applied to a row unit, with the unit before and after it.
type Step struct {
	Name       string     // the technology name, or the effect name of a civilization bonus
	Technology Technology // the researched technology, empty for a civilization bonus
	Source     UpgradeSource
	// Age is the age needed to get the step, 0 if none or unknown, e.g. ImperialAge for Plate Barding Armor.
	Age        TechID
	Effect     Effect
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
	var result strings.Builder
	result.WriteString(formatNumber(stat(row.BaseUnit)) + formatSigned(upgrades[Technologies]))
	for _, source := range []UpgradeSource{UniqueTechnologies, CivilizationBonuses} {
		if math.Abs(upgrades[source]) > 1e-6 {
			result.WriteString(formatSigned(upgrades[source]))
		}
	}
	return result.String()
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
		row := tree.loadRowUpgraded(civilization, unitID, line.Class)
		for _, id := range line.UnitIDs {
			if id == unitID {
				break
			}
			if _, isFound := tree.Units[id]; isFound && civilization.UnitIDs[id] {
				row.EarlierRows = append(row.EarlierRows, tree.loadRowUpgraded(civilization, id, line.Class))
			}
		}
		rows = append(rows, row)
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
	row := Row{CivilizationID: civilization.ID, Civilization: civilization.Name, BaseUnit: unit.clone(), AvailableAge: tree.FindUnitAge(civilization, unitID)}
	if id, isFound := civilization.UpgradeTechnologyIDs[unitID]; isFound {
		row.UpgradeTechnologies = append(row.UpgradeTechnologies, tree.Technologies[id])
	}

	var technologyIDs []TechID
	for id := range TechnologyEffects {
		if civilization.TechnologyIDs[id] && !AutomaticTechnologies[id] {
			technologyIDs = append(technologyIDs, id)
		}
	}
	slices.Sort(technologyIDs) // a fixed order, e.g. Bloodlines before the Mongols HP bonus
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
		row.Steps = append(row.Steps, Step{
			Name: name, Technology: technology, Source: source,
			Age: tree.findTechnologyAge(technology.ID, map[TechID]bool{}), Effect: effect,
			UnitBefore: before, UnitAfter: unit.clone(),
		})
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
		headers = append(headers, formatHeader(column.Name))
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

// formatHeader is the column name in CamelCase, to save width of the table,
// e.g. "Tank vs Arbalester 10" is "TankVsArbalester10".
func formatHeader(name string) string {
	var b strings.Builder
	for word := range strings.FieldsSeq(name) {
		b.WriteString(strings.ToUpper(word[:1]) + word[1:])
	}
	return b.String()
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
	rounded := math.Round(v*100) / 100
	if rounded == 0 {
		// a small negative value rounds to -0, e.g. a score of -0.003 would read "-0"
		rounded = 0
	}
	return strconv.FormatFloat(rounded, 'f', -1, 64)
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
