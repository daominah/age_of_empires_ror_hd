package logic

import (
	"testing"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
)

var scoutLine = Line{UnitIDs: []UnitID{ScoutCavalry, LightCavalry, Hussar, WingedHussar}, Class: ClassCavalry}

func TestFullyUpgradedRows(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	// WHEN each civilization gets its fully upgraded scout line unit
	rows := FullyUpgradedRows(tree, scoutLine)

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
		hp := ShownHP(u)
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
	huszarRows := FullyUpgradedRows(tree, Line{UnitIDs: []UnitID{MagyarHuszar, EliteMagyarHuszar}, Class: ClassCavalry})
	if len(huszarRows) != 1 || huszarRows[0].CivilizationID != Magyars ||
		huszarRows[0].Unit.Cost != (Cost{Food: 80}) || ShownHP(huszarRows[0].Unit) != 110 {
		t.Errorf("got Magyar Huszar rows %+v, want only Magyars with 80 food and 110 HP", huszarRows)
	}

	// THEN civilizations without the scout line have no row
	for _, civilization := range []CivilizationID{Aztecs, Inca, Maya, Mapuche, Muisca, Tupi} {
		if _, isFound := byCivilization[civilization]; isFound {
			t.Errorf("%v: has a row, but no scout line", civilization)
		}
	}
}

func TestHitsToKill(t *testing.T) {
	castle := Attacker{Name: "Castle", Attacks: map[ArmorClass]float64{ArmorClassPierce: 15, ArmorClassSpearmen: 2}}
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
		got := HitsToKill(castle, c.target)

		// THEN it needs this many arrows, the attack class the target does not have is ignored
		if got != c.want {
			t.Errorf("%v: got %v, want %v", c.name, got, c.want)
		}
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
		{Name: "Cost", Value: func(r Row) float64 { return r.Unit.Cost.Total() }},
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

func TestClassAmount(t *testing.T) {
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
		encoded := classAmount(c.class, c.amount)
		class, amount := decodeClassAmount(c.value)

		// THEN it is the .dat value, and back the armor class and the amount
		if encoded != c.value || class != c.class || amount != c.amount {
			t.Errorf("%+v: got encoded %v, decoded class %v amount %v", c, encoded, class, amount)
		}
	}
}
