// Command ranking-scout writes the AoE2 DE civilizations ranked by their fully upgraded scout line unit,
// as a Markdown table next to its source file, one row per civilization.
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

	// castle15 is a fully upgraded Castle arrow: 11 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	castle15 = logic.Unit{Name: "Castle 15", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 15}}
	// townCenter9 is a fully upgraded Town Center arrow: 5 + Fletching, Bodkin Arrow, Bracer, Chemistry.
	townCenter9 = logic.Unit{Name: "Town Center 9", Attacks: map[logic.ArmorClass]float64{logic.ArmorClassPierce: 9}}
	// arbalester is a fully upgraded Arbalester: attack 6+4.
	arbalester = tree.LoadUnitUpgraded(logic.Arbalester)
	// heavyCavalryArcher is a fully upgraded Heavy Cavalry Archer: attack 7+4.
	heavyCavalryArcher = tree.LoadUnitUpgraded(logic.HeavyCavalryArcher)
	// hussar is a fully upgraded Hussar, for the tank in the Hussar mirror: attack 7+4.
	hussar = tree.LoadUnitUpgraded(logic.Hussar)
	// villager is a Villager with Loom, the raid target of a scout: 25 HP + 15, armor 0+1/0+2.
	// It is written by hand as tree.LoadUnitUpgraded applies only the military technologies.
	villager = logic.Unit{
		Name:   "Villager",
		HP:     40,
		Armors: map[logic.ArmorClass]float64{logic.ArmorClassMelee: 1, logic.ArmorClassPierce: 2},
	}
	// halberdier is a fully upgraded Halberdier, the counter a scout fights in a raid.
	halberdier = tree.LoadUnitUpgraded(logic.Halberdier)
	// cavalier is a fully upgraded Cavalier, the Knight line a scout meets more often than the Paladin.
	cavalier = tree.LoadUnitUpgraded(logic.Cavalier)
	// heavyCamel is a fully upgraded Heavy Camel Rider, the counter to the scout line at the Stable.
	heavyCamel = tree.LoadUnitUpgraded(logic.HeavyCamelRider)
	// commonBonusClasses are the classes the whole line has a bonus attack against,
	// so they are neither scored nor noted.
	commonBonusClasses = []logic.ArmorClass{logic.ArmorClassMonks}
	// baseline is a standard fully upgraded Hussar, the 100% of each score weight.
	baseline = logic.NewRowFromUnit(hussar)
	// scoreWeights weight the tank and the damage of a scout.
	// The tank is in whole hits to kill, vs the arrows that hit a scout most:
	// a Castle, an Arbalester, a Heavy Cavalry Archer and a Town Center,
	// then vs the melee units: the Hussar in the mirror, the Halberdier, the Cavalier and the Heavy Camel Rider,
	// all fully upgraded.
	// The damage is in 2 seconds, so it includes the reload time, e.g. Bulgarians Stirrups,
	// vs the units a scout fights: the Hussar in the mirror and the Halberdier weigh most,
	// then the Villager in a raid and the Arbalester (the bonus vs Archers, e.g. Burmese).
	// The tank weighs much more than the damage,
	// so +1 pierce armor (Tatars) ranks above a faster reload (Bulgarians).
	// The bonus attacks out of the line count a little per point, as they help vs a few targets only,
	// e.g. the Lithuanians Winged Hussar vs Gunpowder Units and the Elite Magyar Huszar vs Rams.
	// The free or cheaper upgrades count 0.8 per 1000 resources saved, as they feel strong in practice,
	// e.g. the Turks (Light Cavalry and Hussar upgrades free) rank above the Tatars (Silk Armor),
	// though the Tatars Hussar wins head to head.
	// The cost counts a lot, as cheap scouts are very good, e.g. the Berbers (-20%) rank above a standard Hussar.
	// The speed counts a little, as the standard Hussar is already among the fastest units, e.g. the Cumans (+5%).
	scoreWeights = []logic.WeightedCriterion{
		{Criterion: logic.NewCriterionTankVs(castle15), Weight: 8},
		{Criterion: logic.NewCriterionTankVs(arbalester), Weight: 5},
		{Criterion: logic.NewCriterionTankVs(heavyCavalryArcher), Weight: 5},
		{Criterion: logic.NewCriterionTankVs(townCenter9), Weight: 2},
		{Criterion: logic.NewCriterionTankVs(hussar), Weight: 5},
		{Criterion: logic.NewCriterionTankVs(halberdier), Weight: 3},
		{Criterion: logic.NewCriterionTankVs(heavyCamel), Weight: 3},
		{Criterion: logic.NewCriterionTankVs(cavalier), Weight: 2},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(hussar), Weight: 4},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(halberdier), Weight: 3},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(villager), Weight: 3},
		{Criterion: logic.NewCriterionDamagePerTwoSecondsVs(arbalester), Weight: 2},
		{Criterion: logic.NewCriterionBonusAttackSumExcept(commonBonusClasses...), Weight: 0.25},
		{Criterion: logic.CriterionCost, Weight: 10},
		{Criterion: logic.CriterionTechnologySaving, Weight: 0.8},
		{Criterion: logic.CriterionSpeed, Weight: 5},
	}

	// extraScores are the civilization bonuses the stats leave out, guessed by hand,
	// to mass the scouts (creation speed, long-term food) or in combat (trample, regeneration).
	// On the scale of the score: the Berbers -20% cost scores about 2,
	// a short-term bonus scores 0, e.g. the Mongols hunters or the Lithuanians starting food.
	// The Hussar is a late game unit, so a permanent bonus scores much more than a one-time saving,
	// e.g. the Slavs farmers working faster above the Vikings Wheelbarrow and Hand Cart free.
	extraScores = map[logic.CivilizationID]float64{
		logic.Poles:       4,   // 33% trample damage in 0.5 tile (Lechitic Legacy), Folwark collects instantly 10% farm food
		logic.Khitans:     2,   // Pasture food needs no drop-off (like the Khmer farmers), trained 15% faster, regeneration in combat (Ordo Cavalry)
		logic.Khmer:       2,   // farmers need no drop-off
		logic.Slavs:       2,   // farmers work faster
		logic.Franks:      0.5, // Feudal Scout stronger without upgrades, farm upgrades free, trained 40% faster (Chivalry)
		logic.Georgians:   0.5, // regeneration and population -20% (Aznauri Cavalry), less damage from lower elevation
		logic.Vikings:     0.5, // Wheelbarrow and Hand Cart free
		logic.Gurjaras:    0.2, // livestock garrisoned in the Mills generates food
		logic.Huns:        0.2, // trained 20% faster, no houses
		logic.Mongols:     0.2, // +2 line of sight
		logic.Tatars:      0.2, // +50% damage from higher elevation instead of +25%
		logic.Cumans:      0.1, // trained 100% faster (Steppe Husbandry)
		logic.Hindustanis: 0.1, // villagers cheaper
		logic.Danes:       0.1, // gets 25% of the cost of the buildings it destroys
		logic.Vietnamese:  0.1, // Conscription free
		logic.Wei:         0.1, // 1.8 food per minute per unit (Tuntian)
		logic.Sicilians:   0.1, // farm upgrades +100% food, saves a little wood on reseeding over a long time
	}

	// criteria sort the rows in order, the next one breaking ties of the previous.
	// Reorder or comment out lines to rank differently,
	// e.g. put logic.CriterionTraitWeightSum first.
	criteria = []logic.Criterion{
		// a stat locked behind a unique technology counts only partly, e.g. the Gurjaras cost (Kshatriyas)
		logic.NewCriterionWithExtraScores(
			logic.NewCriterionDiscountUniqueTechnologies(logic.NewCriterionWeightedSum("Score", baseline, scoreWeights), 0.75),
			extraScores),
		// the handwritten traits only break the exact ties of the score,
		// as the score already counts most of them through the extra scores (regeneration, trample, economy)
		// and the bonus damage taken through the damage
		logic.CriterionTraitWeightSum,
	}

	// columns of the output table, after the rank, civilization and unit,
	// apart from the criteria: a criterion sorts without a column of its own.
	columns = []logic.Column{
		// the score, 0 for the baseline, a standard fully upgraded unit
		logic.NewColumnFromCriterion(criteria[0]),
		// the basic stats, add the stats that matter for this line later
		logic.ColumnHP,
		logic.ColumnMeleeArmor,
		logic.ColumnPierceArmor,
		logic.ColumnAttack,
		// the bonus attacks the scout line uses, shown only when they differ from the baseline
		logic.NewColumnBonusAttackShort(baseline,
			logic.ArmorClassLabel{Class: logic.ArmorClassArchers, Label: "Archers"},
			logic.ArmorClassLabel{Class: logic.ArmorClassSkirmishers, Label: "Skirm"},
			logic.ArmorClassLabel{Class: logic.ArmorClassGunpowderUnits, Label: "Gun"},
			logic.ArmorClassLabel{Class: logic.ArmorClassSiegeWeapons, Label: "Siege"},
			logic.ArmorClassLabel{Class: logic.ArmorClassRams, Label: "Ram"},
			logic.ArmorClassLabel{Class: logic.ArmorClassStandardBuildings, Label: "Building"},
		),
		logic.ColumnReloadTime,
		logic.ColumnCost,
		logic.ColumnSpeed,
		// the damage vs a Halberdier stands for the damage in general:
		// it has the same melee armor as the Hussar and the Arbalester, so the header drops the target
		logic.NewColumnFromCriterionNamed("DP2S", logic.NewCriterionDamagePerTwoSecondsVs(halberdier)),
		logic.NewColumnNoteWithBonusAttack(commonBonusClasses...),
	}
)

func main() {
	rows := logic.LoadRowsUpgraded(tree, scoutLine)
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
// e.g. tree.LoadUnitUpgraded(logic.Arbalester).
var tree = parseTechTreeOrExit()

// parseTechTreeOrExit parses the embedded game data, it exits on error.
func parseTechTreeOrExit() *logic.TechTree {
	tree, err := logic.ParseTechTree(data.Civilizations, data.Units, data.Technologies, logic.ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		log.Fatalf("error logic.ParseTechTree: %v", err)
	}
	return tree
}
