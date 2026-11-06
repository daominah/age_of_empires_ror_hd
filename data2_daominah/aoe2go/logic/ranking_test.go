package logic

import (
	"math"
	"testing"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
)

var (
	scoutLine      = Line{UnitIDs: []UnitID{ScoutCavalry, LightCavalry, Hussar, WingedHussar}, Class: ClassCavalry}
	skirmisherLine = Line{UnitIDs: []UnitID{Skirmisher, EliteSkirmisher, ImperialSkirmisher}, Class: ClassArcher}
)

func TestLoadRowsUpgraded(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	// WHEN each civilization gets its fully upgraded scout line unit
	rows := LoadRowsUpgraded(tree, scoutLine)

	// THEN the stats include the technologies and civilization bonuses
	byCivilization := map[CivilizationID]Row{}
	for _, row := range rows {
		byCivilization[row.CivilizationID] = row
	}
	for _, c := range []struct {
		civilization CivilizationID
		unit         string
		hp           float64 // shown in game
		pierceArmor  float64
		meleeAttack  float64
		food         float64
	}{
		// Bloodlines +20, Barding +4 pierce armor, Forging line +4 attack
		{Persians, "Hussar", 95, 6, 11, 80},
		// Bloodlines is not multiplied by the HP bonus: 75 * 1.3 + 20 = 117.5, 118 confirmed in game
		{Mongols, "Hussar", 118, 4, 11, 80},
		// Silk Armor +1 pierce armor
		{Tatars, "Hussar", 95, 7, 11, 80},
		// light cavalry cost -15%
		{Magyars, "Hussar", 95, 6, 11, 68},
		// stable units cost -15% then -20%
		{Berbers, "Hussar", 95, 6, 11, 64},
		// Farimba +5 attack, no Blast Furnace
		{Malians, "Light Cavalry", 80, 6, 14, 80},
		// no Bloodlines, no Plate Barding, +20% HP
		{Franks, "Light Cavalry", 72, 6, 11, 80},
		// no Light Cavalry, the Feudal Age gives the Scout Cavalry +2 attack, outside the Technology Tree
		{Teutons, "Scout Cavalry", 65, 6, 9, 80},
		// the Winged Hussar replaces the Hussar, no Plate Barding
		{Poles, "Winged Hussar", 100, 4, 13, 80},
	} {
		row, isFound := byCivilization[c.civilization]
		if !isFound {
			t.Errorf("%v: no row", c.civilization)
			continue
		}
		u := row.Unit
		hp := GetShownHP(u)
		if u.Name != c.unit || hp != c.hp || u.Armors[ArmorClassPierce] != c.pierceArmor ||
			u.Attacks[ArmorClassMelee] != c.meleeAttack || int(u.Cost.Food+0.5) != int(c.food) {
			t.Errorf("%v: got %v HP %v pierce armor %v attack %v food %v, want %+v",
				c.civilization, u.Name, hp, u.Armors[ArmorClassPierce], u.Attacks[ArmorClassMelee], u.Cost.Food, c)
		}
	}

	// THEN the cost saving of a unique technology is kept apart from the free civilization bonuses:
	// Gurjaras Kshatriyas -25% food, Magyars light cavalry -15%
	for _, c := range []struct {
		civilization                  CivilizationID
		food                          float64
		foodWithoutUniqueTechnologies float64
	}{
		{Gurjaras, 60, 80},
		{Magyars, 68, 68},
	} {
		row := byCivilization[c.civilization]
		if int(row.Unit.Cost.Food+0.5) != int(c.food) ||
			int(row.UnitWithoutUniqueTechnologies.Cost.Food+0.5) != int(c.foodWithoutUniqueTechnologies) {
			t.Errorf("%v: got food %v, without unique technologies %v, want %+v",
				c.civilization, row.Unit.Cost.Food, row.UnitWithoutUniqueTechnologies.Cost.Food, c)
		}
	}

	// THEN a unit outside the scout line gets its own row:
	// Magyars Elite Magyar Huszar costs 80 food and no gold with Corvinian Army, Bloodlines +20 HP
	huszarRows := LoadRowsUpgraded(tree, Line{UnitIDs: []UnitID{MagyarHuszar, EliteMagyarHuszar}, Class: ClassCavalry})
	if len(huszarRows) != 1 || huszarRows[0].CivilizationID != Magyars ||
		huszarRows[0].Unit.Cost != (Cost{Food: 80}) || GetShownHP(huszarRows[0].Unit) != 110 {
		t.Errorf("got Magyar Huszar rows %+v, want only Magyars with 80 food and 110 HP", huszarRows)
	}

	// THEN the skirmisher line gets the archer technologies and its civilization bonuses
	skirmisherRows := map[CivilizationID]Row{}
	for _, row := range LoadRowsUpgraded(tree, skirmisherLine) {
		skirmisherRows[row.CivilizationID] = row
	}
	for _, c := range []struct {
		civilization CivilizationID
		unit         string
		hp           float64
		pierceArmor  float64
		pierceAttack float64
		reloadTime   float64
	}{
		// Fletching, Bodkin Arrow, Bracer, Chemistry +4 attack, archer armors +4 pierce armor
		{Britons, "Elite Skirmisher", 35, 8, 7, 3},
		// the Imperial Skirmisher, archers +20% HP
		{Vietnamese, "Imperial Skirmisher", 42, 9, 8, 3},
		// Tower Shields +2 pierce armor
		{Lithuanians, "Elite Skirmisher", 35, 10, 7, 3},
		// +5 HP in each of Feudal, Castle and Imperial Age, no Bracer
		{Mapuche, "Elite Skirmisher", 50, 8, 6, 3},
		// skirmishers attack 25% faster
		{Dravidians, "Elite Skirmisher", 35, 8, 7, 2.4},
		// Atlatl +1 attack, no Ring Archer Armor
		{Aztecs, "Elite Skirmisher", 35, 6, 8, 3},
	} {
		u := skirmisherRows[c.civilization].Unit
		if u.Name != c.unit || GetShownHP(u) != c.hp || u.Armors[ArmorClassPierce] != c.pierceArmor ||
			u.Attacks[ArmorClassPierce] != c.pierceAttack || math.Abs(u.ReloadTime-c.reloadTime) > 1e-6 {
			t.Errorf("%v: got %v HP %v pierce armor %v attack %v reload time %v, want %+v",
				c.civilization, u.Name, GetShownHP(u), u.Armors[ArmorClassPierce], u.Attacks[ArmorClassPierce], u.ReloadTime, c)
		}
	}

	// THEN the Maya skirmisher throws an extra projectile of attack 1 with Hul'che Javelineers,
	// so against a very high pierce armor it deals 2 damage and the others 1
	ram := Unit{Armors: map[ArmorClass]float64{ArmorClassPierce: 150}, HP: 175}
	for civilization, want := range map[CivilizationID]float64{Maya: 2, Britons: 1} {
		u := skirmisherRows[civilization].Unit
		got := CalculateDamagePerHit(u, ram)
		if got != want {
			t.Errorf("%v: got %v damage to a ram, want %v", civilization, got, want)
		}
	}

	// THEN the Inca skirmisher has no minimum range after Andean Sling, the others have 1
	for civilization, want := range map[CivilizationID]float64{Inca: 0, Britons: 1} {
		if got := skirmisherRows[civilization].Unit.MinRange; got != want {
			t.Errorf("%v: got min range %v, want %v", civilization, got, want)
		}
	}

	// THEN the range includes Fletching, Bodkin Arrow, Bracer and the unique technologies
	for _, c := range []struct {
		civilization CivilizationID
		unitRange    float64
	}{
		{Persians, 7},   // no Bracer: 5 + 2
		{Vietnamese, 8}, // 5 + 3
		{Britons, 9},    // Yeomen +1, the Britons archer range bonus excludes skirmishers
		{Aztecs, 9},     // Atlatl +1
	} {
		if got := skirmisherRows[c.civilization].Unit.Range; got != c.unitRange {
			t.Errorf("%v: got range %v, want %v", c.civilization, got, c.unitRange)
		}
	}

	// THEN a bonus triggered in several ages notes each age, e.g. the Mapuche skirmisher HP
	if got, want := FormatAgeNote(skirmisherRows[Mapuche]), "Skirmisher HP 35 in Feudal Age, Elite Skirmisher HP 45 in Castle Age, 50 in Imperial Age"; got != want {
		t.Errorf("Mapuche: got note %q, want %q", got, want)
	}

	// THEN a note names the unit available in each age, e.g. no Hussar in the Castle Age:
	// the Light Cavalry has (60 base HP) * 1.2 + 20 Bloodlines = 92
	mongolsHussar := byCivilization[Mongols]
	if got, want := FormatAgeNote(mongolsHussar), "Light Cavalry HP 92 in Castle Age, Hussar HP 118 in Imperial Age"; got != want {
		t.Errorf("Mongols: got note %q, want %q", got, want)
	}

	// THEN a note counts only the technologies of the age:
	// the Teutons Scout Cavalry has no Plate Barding Armor in the Castle Age,
	// Scale and Chain Barding Armor +2 and the Teutons bonus +1 = 3
	if got, want := FormatAgeNote(byCivilization[Teutons]), "melee armor 3 in Castle Age, 5 in Imperial Age"; got != want {
		t.Errorf("Teutons: got note %q, want %q", got, want)
	}

	// THEN civilizations without the scout line have no row
	for _, civilization := range []CivilizationID{Aztecs, Inca, Maya, Mapuche, Muisca, Tupi} {
		if _, isFound := byCivilization[civilization]; isFound {
			t.Errorf("%v: has a row, but no scout line", civilization)
		}
	}
}

func TestLoadRowsUpgradedCavalryArcher(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	// WHEN each civilization gets its fully upgraded cavalry archer line unit
	line := Line{UnitIDs: []UnitID{CavalryArcher, HeavyCavalryArcher}, Class: ClassCavalryArcher}
	byCivilization := map[CivilizationID]Row{}
	for _, row := range LoadRowsUpgraded(tree, line) {
		byCivilization[row.CivilizationID] = row
	}

	// THEN the civilization bonuses and unique technologies written for the cavalry class
	// also apply to the cavalry archer class
	for _, c := range []struct {
		civilization CivilizationID
		hp           float64 // shown in game
		wood         float64
		gold         float64
		reloadTime   float64
	}{
		// Bloodlines +20 HP, Thumb Ring 10% faster attack
		{Persians, 80, 40, 60, 1.8},
		// Sipahi +20 HP
		{Turks, 100, 40, 60, 1.8},
		// wood cost -50%
		{Koreans, 60, 20, 60, 1.8},
		// attack 25% faster: 1.8 * 0.8
		{Mongols, 80, 40, 60, 1.44},
		// attack 20% faster: 1.8 * 0.8 with Thumb Ring included
		{Jurchens, 80, 40, 60, 1.6},
	} {
		u := byCivilization[c.civilization].Unit
		if GetShownHP(u) != c.hp || math.Abs(u.Cost.Wood-c.wood) > 0.5 || math.Abs(u.Cost.Gold-c.gold) > 0.5 ||
			math.Abs(u.ReloadTime-c.reloadTime) > 0.01 {
			t.Errorf("%v: got HP %v wood %v gold %v reload time %v, want %+v",
				c.civilization, GetShownHP(u), u.Cost.Wood, u.Cost.Gold, u.ReloadTime, c)
		}
	}
}

func TestCalculateHitsToKill(t *testing.T) {
	castle := Unit{Name: "Castle", Attacks: map[ArmorClass]float64{ArmorClassPierce: 15, ArmorClassSpearmen: 2}}
	for _, c := range []struct {
		name   string
		target Unit
		want   int
	}{
		// GIVEN 95 HP and 7 pierce armor: 8 damage per arrow
		{"armor reduces damage", Unit{HP: 95, Armors: map[ArmorClass]float64{ArmorClassPierce: 7}}, 12},
		// GIVEN the HP is a multiplied effect result just above an integer
		{"rounding noise is ignored", Unit{HP: 96.0000047, Armors: map[ArmorClass]float64{ArmorClassPierce: 7}}, 12},
		// GIVEN the HP has decimals: Mongols Hussar 117.5 HP and 4 pierce armor
		{"HP decimals count", Unit{HP: 117.5, Armors: map[ArmorClass]float64{ArmorClassPierce: 4}}, 11},
		// GIVEN the armor is higher than the attack
		{"at least 1 damage", Unit{HP: 10, Armors: map[ArmorClass]float64{ArmorClassPierce: 20}}, 10},
	} {
		// WHEN the Castle shoots the target
		got := CalculateHitsToKill(castle, c.target)

		// THEN it needs this many arrows, the attack class the target does not have is ignored
		if got != c.want {
			t.Errorf("%v: got %v, want %v", c.name, got, c.want)
		}
	}
}

func TestLoadUnitUpgraded_Arbalester(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	// WHEN the Arbalester gets the technologies any civilization can research
	unit := tree.LoadUnitUpgraded(Arbalester)

	// THEN it has Fletching, Bodkin Arrow, Bracer, Chemistry, Thumb Ring and the archer armors,
	// without civilization bonuses, e.g. no Britons range
	if unit.HP != 40 || unit.Attacks[ArmorClassPierce] != 10 || unit.Range != 8 ||
		unit.Armors[ArmorClassPierce] != 4 || unit.Armors[ArmorClassMelee] != 3 || math.Abs(unit.ReloadTime-1.7) > 1e-6 {
		t.Errorf("got %+v, want 40 HP, attack 10, range 8, armor 3/4, reload time 1.7", unit)
	}
}

func TestNewCriterionHitsToKillTarget(t *testing.T) {
	// GIVEN an Arbalester with 40 HP and 4 pierce armor,
	// and an Elite Skirmisher with 7 pierce attack and 4 bonus attack against archers
	arbalester := Unit{Name: "Arbalester", HP: 40, Armors: map[ArmorClass]float64{ArmorClassPierce: 4, ArmorClassArchers: 0}}
	row := Row{Unit: Unit{Attacks: map[ArmorClass]float64{ArmorClassPierce: 7, ArmorClassArchers: 4}}}

	// WHEN the skirmisher shoots the Arbalester
	criterion := NewCriterionHitsToKillTarget(arbalester)

	// THEN 7 damage per hit kills it in 6 hits, and fewer hits rank higher
	if got := criterion.Value(row); got != 6 || criterion.IsHigherBetter || criterion.Name != "Hits to kill Arbalester" {
		t.Errorf("got %v hits, criterion %+v, want 6 hits, fewer is better", got, criterion)
	}
}

func TestSortRows(t *testing.T) {
	// GIVEN civilizations tied on the first criterion
	rows := []Row{
		{Civilization: "C", Unit: Unit{HP: 10, Cost: Cost{Food: 80}}},
		{Civilization: "B", Unit: Unit{HP: 20, Cost: Cost{Food: 80}}},
		{Civilization: "A", Unit: Unit{HP: 20, Cost: Cost{Food: 60}}},
		{Civilization: "D", Unit: Unit{HP: 20, Cost: Cost{Food: 60}}},
	}
	criteria := []Criterion{
		{Name: "HP", Value: func(r Row) float64 { return r.Unit.HP }, IsHigherBetter: true},
		{Name: "Cost", Value: func(r Row) float64 { return r.Unit.Cost.Sum() }},
	}

	// WHEN the rows are sorted
	SortRows(rows, criteria)

	// THEN the next criterion breaks the ties, then the civilization name
	got := ""
	for _, row := range rows {
		got += row.Civilization
	}
	if got != "ADBC" {
		t.Errorf("got order %v, want ADBC", got)
	}
}

func TestEncodeClassAmount(t *testing.T) {
	for _, c := range []struct {
		value  float64
		class  ArmorClass
		amount float64
	}{
		{769, ArmorClassPierce, 1},
		{1029, ArmorClassMelee, 5},
		{-1025, ArmorClassMelee, -1},
	} {
		// WHEN an armor or attack command value is encoded and decoded
		encoded := encodeClassAmount(c.class, c.amount)
		class, amount := decodeClassAmount(c.value)

		// THEN it is the .dat value, and back the armor class and the amount
		if encoded != c.value || class != c.class || amount != c.amount {
			t.Errorf("%+v: got encoded %v, decoded class %v amount %v", c, encoded, class, amount)
		}
	}
}

func TestFormatTraitNote(t *testing.T) {
	// GIVEN rows of civilizations with traits the rankings do not calculate
	for _, c := range []struct {
		row  Row
		want string
	}{
		// the cost and the stats are not changed, only the note tells
		{Row{CivilizationID: Poles, Unit: Unit{ID: WingedHussar}}, "33% trample damage in 0.5 tile radius after Lechitic Legacy"},
		{Row{CivilizationID: Sicilians, Unit: Unit{ID: Hussar}}, "bonus damage taken -40%"},
		// a trait of another unit of the civilization is not noted
		{Row{CivilizationID: Khitans, Unit: Unit{ID: LightCavalry}}, "regeneration in combat 90 per minute after Ordo Cavalry"},
		{Row{CivilizationID: Persians, Unit: Unit{ID: Hussar}}, ""},
		{Row{CivilizationID: Berbers, Unit: Unit{ID: Hussar}}, ""},
	} {
		// WHEN the note is formatted
		got := FormatTraitNote(c.row)

		// THEN it lists the traits of the row unit
		if got != c.want {
			t.Errorf("%v %v: got %q, want %q", c.row.CivilizationID, c.row.Unit.ID, got, c.want)
		}
	}
}

func TestFormatBonusAttackNote(t *testing.T) {
	// GIVEN a Mangudai with the bonus attacks of the line and one special bonus
	row := Row{Unit: Unit{Attacks: map[ArmorClass]float64{
		ArmorClassPierce: 12, ArmorClassSpearmen: 3, ArmorClassSiegeWeapons: 5, ArmorClassRoyalHeirs: -3, ArmorClassRams: 0,
	}}}

	// WHEN the note skips the bonus of the whole line (Parthian Tactics vs Spearmen)
	got := FormatBonusAttackNote(row, ArmorClassSpearmen)

	// THEN only the special bonus is noted, not the pierce attack, a zero or a negative one
	if want := "bonus attack 5 vs Siege Weapons"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTechnologyCostNote(t *testing.T) {
	// GIVEN the scout line rows
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}
	line := Line{UnitIDs: []UnitID{ScoutCavalry, LightCavalry, Hussar, WingedHussar}, Class: ClassCavalry}
	notes := map[CivilizationID]string{}
	for _, row := range LoadRowsUpgraded(tree, line) {
		notes[row.CivilizationID] = FormatTechnologyCostNote(row)
	}

	// THEN the cheaper technologies are noted with the cost before and after
	for _, c := range []struct {
		civilization CivilizationID
		want         string
	}{
		{Turks, "Light Cavalry upgrade free instead of 150f 50g, Hussar upgrade free instead of 500f 600g"},
		{Burgundians, "Husbandry 75f instead of 150f, Light Cavalry upgrade 75f 25g instead of 150f 50g, Hussar upgrade 250f 300g instead of 500f 600g"},
		{Persians, ""},
	} {
		if got := notes[c.civilization]; got != c.want {
			t.Errorf("civilization %v: got %q, want %q", c.civilization, got, c.want)
		}
	}
}

func TestCalculateDamagePerTwoSeconds(t *testing.T) {
	// GIVEN an attacker of pierce attack 7 and a target of pierce armor 4
	attacker := Unit{Attacks: map[ArmorClass]float64{ArmorClassPierce: 7}, ReloadTime: 3}
	target := Unit{Armors: map[ArmorClass]float64{ArmorClassPierce: 4}}

	// WHEN the reload time is 3 seconds
	// THEN it hits for 3 every 3 seconds, so 2 in 2 seconds
	if got := CalculateDamagePerTwoSeconds(attacker, target); got != 2 {
		t.Errorf("got %v, want 2", got)
	}

	// WHEN the target armor is higher than the attack
	attacker.ReloadTime = 2
	target = Unit{Armors: map[ArmorClass]float64{ArmorClassPierce: 100}}
	// THEN each hit still deals the minimum 1 damage, so 1 per second with a reload time of 2 seconds
	if got := CalculateDamagePerTwoSeconds(attacker, target); got != 1 {
		t.Errorf("got %v, want 1", got)
	}

	// WHEN the attacker also fires 2 extra arrows of pierce attack 5 with each attack
	attacker.ExtraProjectiles = 2
	attacker.ExtraAttacks = map[ArmorClass]float64{ArmorClassPierce: 5}
	target = Unit{Armors: map[ArmorClass]float64{ArmorClassPierce: 4}}
	// THEN each attack deals 3 + 1 + 1, so 5 in 2 seconds
	if got := CalculateDamagePerTwoSeconds(attacker, target); got != 5 {
		t.Errorf("got %v, want 5", got)
	}

	// WHEN the attacker also has a charged attack of pierce 34 every 29 seconds, without extra arrows
	attacker.ChargeAttack = &ChargeAttack{Shot: Shot{Attacks: map[ArmorClass]float64{ArmorClassPierce: 34}}, ChargeTime: 29}
	// THEN every 30 seconds, the charge time and half a reload, a charged shot of 30 replaces a normal shot of 5,
	// adding 25 per 30 seconds, so 5 + 25*2/30 in 2 seconds
	if got, want := CalculateDamagePerTwoSeconds(attacker, target), 5+25*2.0/30; math.Abs(got-want) > 1e-9 {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestLoadUnitUpgraded_ExtraArrowsAndChargedAttack(t *testing.T) {
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}
	hussar := tree.LoadUnitUpgraded(Hussar)

	// WHEN loading the Elite Kipchak, with every technology of its civilization
	kipchak := tree.LoadUnitUpgraded(EliteKipchak)
	// THEN each attack fires 3 extra arrows of 3 pierce, not improved by the Blacksmith
	if kipchak.ExtraProjectiles != 3 || kipchak.ExtraAttacks[ArmorClassPierce] != 3 || kipchak.ChargeAttack != nil {
		t.Errorf("Elite Kipchak: got %v extra arrows of %v, charged attack %v, want 3 of 3 pierce and none",
			kipchak.ExtraProjectiles, kipchak.ExtraAttacks, kipchak.ChargeAttack)
	}

	// WHEN loading the Xianbei Raider
	xianbei := tree.LoadUnitUpgraded(XianbeiRaider)
	// THEN a normal attack fires 1 arrow,
	// and the charged volley fires 1 arrow of 5 pierce and 5 of 1 pierce,
	// each improved by the Blacksmith +4 pierce, and by Parthian Tactics +2 vs Spearmen
	charge := xianbei.ChargeAttack
	if xianbei.ExtraProjectiles != 0 || charge == nil || charge.ChargeTime != 30 || charge.ExtraProjectiles != 5 ||
		charge.Attacks[ArmorClassPierce] != 9 || charge.ExtraAttacks[ArmorClassPierce] != 5 ||
		charge.Attacks[ArmorClassSpearmen] != 5 || charge.ExtraAttacks[ArmorClassSpearmen] != 6 {
		t.Errorf("Xianbei Raider: got %v extra arrows, charged attack %+v", xianbei.ExtraProjectiles, charge)
	}

	// WHEN loading the Elite Bolas Rider of the Mapuche, who have no Bracer
	bolas := tree.LoadUnitUpgraded(EliteBolasRider)
	// THEN its charged bolas of 9 pierce gets Fletching, Bodkin Arrow and Chemistry +3 pierce, without extra arrows
	if bolas.ChargeAttack == nil || bolas.ChargeAttack.Attacks[ArmorClassPierce] != 12 || bolas.ChargeAttack.ExtraProjectiles != 0 {
		t.Errorf("Elite Bolas Rider: got charged attack %+v, want 12 pierce", bolas.ChargeAttack)
	}

	// THEN against a fully upgraded Hussar, the charged bolas deals 12 - 6 = 6 plus 3 vs cavalry,
	// and the Xianbei Raider charged volley 9 - 6 = 3 plus 5 arrows of 5 - 6, at least 1 each
	if got := CalculateShotDamage(bolas.ChargeAttack.Shot, hussar); got != 9 {
		t.Errorf("Elite Bolas Rider charged bolas: got %v damage to a Hussar, want 9", got)
	}
	if got := CalculateShotDamage(charge.Shot, hussar); got != 8 {
		t.Errorf("Xianbei Raider charged volley: got %v damage to a Hussar, want 8", got)
	}
}

func TestNewCriterionWeightedSum(t *testing.T) {
	// GIVEN a baseline unit of range 7, reload time 2 and melee armor 0
	baseline := NewRowFromUnit(Unit{Range: 7, ReloadTime: 2, Armors: map[ArmorClass]float64{ArmorClassMelee: 0}})
	meleeArmor := Criterion{Name: "Melee armor", Value: func(r Row) float64 { return r.Unit.Armors[ArmorClassMelee] }, IsHigherBetter: true}

	// WHEN the score weighs the range by 7, then the reload time and the melee armor by the default weights
	score := NewCriterionWeightedSum("Score", baseline, []WeightedCriterion{
		{Criterion: CriterionRange, Weight: 7},
		{Criterion: CriterionReloadTime},
		{Criterion: meleeArmor},
	})

	// THEN the baseline scores 0
	if got := score.Value(baseline); got != 0 {
		t.Errorf("baseline: got %v, want 0", got)
	}
	// THEN a range of 8 is 1/7 better than 7, worth 1
	if got := score.Value(NewRowFromUnit(Unit{Range: 8, ReloadTime: 2})); math.Abs(got-1) > 1e-9 {
		t.Errorf("range 8: got %v, want 1", got)
	}
	// THEN a reload time of 1.5 is 25% better as lower is better, worth 0.25 * 4 by default above the last criterion
	if got := score.Value(NewRowFromUnit(Unit{Range: 7, ReloadTime: 1.5})); math.Abs(got-1) > 1e-9 {
		t.Errorf("reload time 1.5: got %v, want 1", got)
	}
	// THEN a melee armor of 2 vs 0 counts the difference, as no ratio exists, worth 2 * 1 as the last criterion
	if got := score.Value(NewRowFromUnit(Unit{Range: 7, ReloadTime: 2, Armors: map[ArmorClass]float64{ArmorClassMelee: 2}})); math.Abs(got-2) > 1e-9 {
		t.Errorf("melee armor 2: got %v, want 2", got)
	}
}

func TestNewCriterionWeightedSum_Limit(t *testing.T) {
	// GIVEN a baseline unit of range 7, and a score of the range capped at +-50%
	baseline := NewRowFromUnit(Unit{Range: 7})
	score := NewCriterionWeightedSum("Score", baseline, []WeightedCriterion{{Criterion: CriterionRange, Weight: 2, Limit: 0.5}})

	// WHEN a unit has range 21, 200% better than the baseline
	got := score.Value(NewRowFromUnit(Unit{Range: 21}))

	// THEN it counts as 50% better, worth 0.5 * 2
	if math.Abs(got-1) > 1e-9 {
		t.Errorf("range 21: got %v, want 1", got)
	}
}

func TestNewCriterionWeightedSum_WeightIfWorse(t *testing.T) {
	// GIVEN a baseline unit of range 7, and a score where a range below it counts twice as much
	baseline := NewRowFromUnit(Unit{Range: 7})
	score := NewCriterionWeightedSum("Score", baseline, []WeightedCriterion{{Criterion: CriterionRange, Weight: 7, WeightIfWorse: 14}})

	// WHEN units have a tile more and a tile less than the baseline
	better := score.Value(NewRowFromUnit(Unit{Range: 8}))
	worse := score.Value(NewRowFromUnit(Unit{Range: 6}))

	// THEN the tile more adds 1, the tile less costs 2
	if math.Abs(better-1) > 1e-9 {
		t.Errorf("range 8: got %v, want 1", better)
	}
	if math.Abs(worse+2) > 1e-9 {
		t.Errorf("range 6: got %v, want -2", worse)
	}
}
