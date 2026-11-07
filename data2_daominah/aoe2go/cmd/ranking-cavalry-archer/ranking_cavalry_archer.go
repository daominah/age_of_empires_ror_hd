// Command ranking-cavalry-archer writes the AoE2 DE civilizations ranked by their fully upgraded cavalry archer line unit,
// as a Markdown table next to its source file, one row per civilization, plus extra rows for unique cavalry archers.
package main

import (
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

	// hussar is a fully upgraded Hussar, both the target of the cavalry archer and the attacker for its tank vs melee:
	// 95 HP, pierce armor 6, attack 7+4.
	hussar = tree.LoadUnitUpgraded(logic.Hussar)

	// baseline is a standard fully upgraded Heavy Cavalry Archer, the 100% of each score weight:
	// attack 7+4, reload 1.8, range 7, speed 1.54, 60 gold.
	baseline = logic.NewRowFromUnit(tree.LoadUnitUpgraded(logic.HeavyCavalryArcher))

	// scoreWeights are the stats of the score, each weight is the score of that stat 100% better than the baseline.
	// A weight of 0 is 1 for the last stat and 4 times the next one for each stat above it.
	//   - Damage 3.97 and per gold 2.65: the Xianbei Raider (35 gold) ranks a bit below the baseline,
	//     above a Heavy Cavalry Archer missing an attack upgrade.
	//   - Range 7, so 1 per tile: the Magyars Heavy Cavalry Archer (range 8) ranks 2nd.
	//     A tile less costs 15/7, about 2: a cavalry archer of range 6 has to move to fight the others,
	//     e.g. the Elite Kipchak drops below the Heavy Cavalry Archers without Parthian Tactics,
	//     despite its cheap gold and its speed, as the range counts more than the tank.
	//   - Speed 6.16, so 0.4 per 0.1 speed: the Danes and Saxons Heavy Mounted Crossbowman
	//     (range 8, speed 1.38) rank just after the baseline.
	//   - Damage vs the baseline 0.5, small on purpose: the cavalry archer mirror only breaks ties,
	//     e.g. the Japanese +2 vs Archers makes them the best of the Heavy Cavalry Archers of attack 7+4 and reload 1.8,
	//     while the Elite Camel Archer +6 vs Cavalry Archers keeps it 3rd, below the Magyars Heavy Cavalry Archer.
	//   - Tank vs the baseline 2, capped at +-50%: HP and pierce armor together,
	//     e.g. the Tatars +1 pierce armor and the Turks +20 HP rank above the baseline,
	//     the Khitans without Bloodlines (60 HP) below it.
	//     The cap keeps the Elite War Wagon (+317%) from deciding alone, as nobody fights it with cavalry archers.
	//   - Wood cost 0.1, capped at +-50%: the gold is already in the damage per gold,
	//     and cheap wood helps a little only, e.g. the Koreans Heavy Cavalry Archer (20w) ranks above the standard ones.
	//     The cap keeps a small base of wood from swinging the score, e.g. the Elite War Wagon (100w).
	//   - Technology saving 0.1 per 1000 resources, small: the free or cheaper upgrades break the ties of the combat,
	//     e.g. the Tatars (Parthian Tactics and Thumb Ring free) above the Turks (Chemistry free),
	//     as the two tie head to head: 20 hits each, Sipahi HP vs Silk Armor.
	//     A larger weight lets the Koreans free archer armors lift the Elite War Wagon to the top.
	scoreWeights = []logic.WeightedCriterion{
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(hussar), Weight: 3.97},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(baseline.Unit), Weight: 0.5},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsPer100GoldVs(hussar), Weight: 2.65},
		{Criterion: logic.CriterionRange, Weight: 7, WeightIfWorse: 15},
		{Criterion: logic.CriterionSpeed, Weight: 6.16},
		{Criterion: logic.NewCriterionTankVs(baseline.Unit), Weight: 2, Limit: 0.5},
		{Criterion: logic.CriterionWoodCost, Weight: 0.1, Limit: 0.5},
		{Criterion: logic.CriterionTechnologySaving, Weight: 0.1},
	}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// The score decides, the others only break exact ties with what the score leaves out.
	// The cost is in the score through the damage per gold.
	criteria = []logic.Criterion{
		logic.NewCriterionWeightedSum("Score", baseline, scoreWeights),
		// the strengths the score leaves out, roughly ordered:
		// melee armor vs a Hussar (Teutons), the bonus attacks (vs Siege Weapons, Rams, Infantry),
		// then the handwritten traits (regeneration, bonus damage taken)
		logic.NewCriterionTankVs(hussar),
		logic.CriterionBonusAttackSum,
		logic.CriterionTraitWeightSum,
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		// the score, 0 for the baseline, a standard fully upgraded unit
		logic.NewColumnFromCriterion(criteria[0]),
		// the basic stats, add the stats that matter for this line later
		// a range unit shows its offense first
		logic.ColumnRange,
		logic.ColumnAttack,
		// the bonus attacks that matter for a cavalry archer: vs the Hussar, vs the other cavalry archers,
		// a Heavy Cavalry Archer being both an archer and a cavalry archer,
		// and vs siege weapons and rams, which cavalry archers raid
		logic.NewColumnBonusAttackShort(baseline,
			logic.ArmorClassLabel{Class: logic.ArmorClassCavalry, Label: "Cav"},
			logic.ArmorClassLabel{Class: logic.ArmorClassArchers, Label: "Archers"},
			logic.ArmorClassLabel{Class: logic.ArmorClassCavalryArchers, Label: "CavArcher"},
			logic.ArmorClassLabel{Class: logic.ArmorClassSiegeWeapons, Label: "Siege"},
			logic.ArmorClassLabel{Class: logic.ArmorClassRams, Label: "Ram"},
		),
		logic.ColumnReloadTime,
		logic.NewColumnFromCriterionNamed("DP2Svs Scout", logic.NewCriterionDamagePerTwoSecondsVs(hussar)),
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnCost,
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
// e.g. tree.LoadUnitUpgraded(logic.Hussar).
var tree = parseTechTreeOrExit()

// parseTechTreeOrExit parses the embedded game data, it exits on error.
func parseTechTreeOrExit() *logic.TechTree {
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	return tree
}
