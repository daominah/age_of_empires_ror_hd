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

	// halberdier is a fully upgraded Halberdier, speed 1.1 with Squires,
	// the counter that chases the cavalry archer, for its kite time.
	halberdier = tree.LoadUnitUpgraded(logic.Halberdier)
	// knight is a fully upgraded Knight, speed 1.485 with Husbandry, the most common chaser of cavalry archers,
	// a bit slower than a standard Heavy Cavalry Archer, so a slower unit cannot run away from it.
	knight = tree.LoadUnitUpgraded(logic.Knight)

	// trebuchet and siegeRam are the siege that cavalry archers raid,
	// pierce armor 150 and 195, so every arrow does the minimum 1 damage plus the bonus attack.
	// The Siege Ram has armor 2 vs Rams, which cancels the Elite Bolas Rider +2 vs Rams.
	trebuchet = tree.LoadUnitUpgraded(logic.Trebuchet)
	siegeRam  = tree.LoadUnitUpgraded(logic.SiegeRam)

	// baseline is a standard fully upgraded Heavy Cavalry Archer, the 100% of each score weight:
	// attack 7+4, reload 1.8, range 7, speed 1.54, attack delay 0.9, 60 gold.
	baseline = logic.NewRowFromUnit(tree.LoadUnitUpgraded(logic.HeavyCavalryArcher))

	// scoreWeights are the stats of the score, each weight is the score of that stat 100% better than the baseline.
	// A weight of 0 is 1 for the last stat and 4 times the next one for each stat above it.
	//   - Damage vs the Hussar, and per gold: the main job, killing the raiders that hunt cavalry archers,
	//     so a cheap unit with a weaker shot can still pay off, e.g. the Xianbei Raider (35 gold).
	//   - Range: a tile less costs more than a tile more gains, as a shorter range has to walk into the others to fight,
	//     e.g. the Magyars range 8 gains, the Elite Kipchak range 6 loses more, despite its cheap gold and its speed.
	//   - Kite time vs the Halberdier: the seconds per shot to stop for the attack delay, fire,
	//     and win back the gap on its counter, attack delay × speed / (speed - 1.1), 3.14 seconds for the baseline.
	//     It combines the attack delay and the speed, e.g. the Elite Mangudai (delay 0.5) 1.74 seconds,
	//     the Elite Ratha (delay 0.62, speed 1.43) 2.69 seconds, its short delay makes up for its speed.
	//     A speed near the Halberdier hurts a lot, e.g. the Elite War Wagon (delay 1, speed 1.32) 5.98 seconds.
	//     The cap keeps the slowest units from losing without end at the bottom,
	//     e.g. the Castle Age Mounted Crossbowman (speed 1.25) 7.58 seconds.
	//   - Speed deficit vs the Knight: how much slower than a Knight, which then forces the fight,
	//     a straight line so a bit slower loses a little, much slower a lot,
	//     e.g. the Elite Ratha (1.43) 0.055 slower, a melee tank that rarely needs to run,
	//     the Elite War Wagon (1.32) 0.165 slower, three times the Ratha.
	//     A unit faster than the Knight loses nothing, e.g. the Elite Bolas Rider (1.5).
	//   - Damage vs the baseline, small on purpose: the cavalry archer mirror only breaks ties,
	//     e.g. the Japanese +2 vs Archers, the Elite Camel Archer +6 vs Cavalry Archers.
	//   - Tank vs the baseline: HP and pierce armor together, vs other archers,
	//     e.g. the Tatars +1 pierce armor, the Turks +20 HP, the Khitans without Bloodlines (60 HP).
	//     It is weighted so a tanky Heavy Cavalry Archer beats a cheap fragile unit,
	//     e.g. the Turks (20 hits to kill) above the Elite Bolas Rider (13 hits).
	//     The cap keeps a huge tank from deciding alone, as nobody fights it with cavalry archers,
	//     e.g. the Elite War Wagon (200 HP).
	//   - Tank vs the Hussar, small: HP and melee armor vs the raiders that hunt cavalry archers,
	//     the baseline needs 14 hits to die, e.g. the Xianbei Raider (melee armor 8) 20 hits, the Elite War Wagon 25 hits.
	//   - Wood cost, small and capped: the gold is already in the damage per gold,
	//     and cheap wood helps a little only, e.g. the Koreans Heavy Cavalry Archer (20w).
	//     The cap keeps a small base of wood from swinging the score, e.g. the Elite War Wagon (100w).
	//   - Technology saving, small: free or cheaper upgrades break the ties of the combat,
	//     e.g. the Tatars (Parthian Tactics and Thumb Ring free) vs the Turks (Chemistry free),
	//     as the two tie head to head: 20 hits each, Sipahi HP vs Silk Armor.
	//   - Damage vs the Trebuchet and vs the Siege Ram, the Trebuchet more as the more common raid target:
	//     the minimum 1 damage per arrow counts the arrows per second and the bonus attack,
	//     e.g. the Elite Mangudai +5 vs Siege Weapons, the Elite Kipchak several arrows each shot,
	//     the Xianbei Raider charged volley of 6 arrows every 30 seconds, the Heavy Mounted Crossbowman reload 2.25.
	//     The cap is high so a big bonus vs Siege Weapons stands out from several arrows.
	scoreWeights = []logic.WeightedCriterion{
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(hussar), Weight: 3.97},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(baseline.Unit), Weight: 0.5},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsPer100GoldVs(hussar), Weight: 2.65},
		{Criterion: logic.CriterionRange, Weight: 7, WeightIfWorse: 15},
		{Criterion: logic.NewCriterionKiteTimeVs(halberdier), Weight: 1.76, Limit: 1},
		{Criterion: logic.NewCriterionSpeedDeficitVs(knight), Weight: 6},
		{Criterion: logic.NewCriterionTankVs(baseline.Unit), Weight: 2.5, Limit: 0.5},
		{Criterion: logic.NewCriterionTankVs(hussar), Weight: 0.3, Limit: 0.5},
		{Criterion: logic.CriterionWoodCost, Weight: 0.1, Limit: 0.5},
		{Criterion: logic.CriterionTechnologySaving, Weight: 0.1},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(trebuchet), Weight: 0.07, Limit: 6},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(siegeRam), Weight: 0.035, Limit: 6},
	}

	// extraScores are the strengths the score leaves out:
	//   - Bengalis: the Ratha switches to a melee mode, to fight skirmishers, siege and buildings up close,
	//     which no other cavalry archer can. It also makes up a bit for the tank cap,
	//     as the Ratha needs 34 hits of the baseline to die, more than twice the baseline.
	//     The Bengalis rank only the Ratha here, so the score of the civilization is the score of the Ratha.
	extraScores = map[logic.CivilizationID]float64{logic.Bengalis: 0.3}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// The score decides, the others only break exact ties with what the score leaves out.
	// The cost is in the score through the damage per gold.
	criteria = []logic.Criterion{
		logic.NewCriterionWithExtraScores(logic.NewCriterionWeightedSum("Score", baseline, scoreWeights), extraScores),
		// the strengths the score leaves out, roughly ordered:
		// the bonus attacks (vs Siege Weapons, Rams, Infantry),
		// then the handwritten traits (regeneration, bonus damage taken)
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
		// seconds from the start of the attack to the shot, lower is better
		logic.ColumnAttackDelay,
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
