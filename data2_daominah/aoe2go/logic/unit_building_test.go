package logic

import (
	"math"
	"testing"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
)

func TestLoadUnitWithTechnologiesBuilding(t *testing.T) {
	// GIVEN the embedded game data
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}

	for _, c := range []struct {
		name         string
		id           UnitID
		techIDs      []TechID
		wantHP       float64
		wantArmors   map[ArmorClass]float64
		wantSameWall bool
	}{
		{"Barracks without technologies", Barracks, nil, 1200,
			map[ArmorClass]float64{ArmorClassMelee: 0, ArmorClassPierce: 7, ArmorClassAllBuildings: 0, ArmorClassStandardBuildings: 0}, false},
		{"Barracks with Masonry", Barracks, []TechID{Masonry}, 1320,
			map[ArmorClass]float64{ArmorClassMelee: 1, ArmorClassPierce: 8, ArmorClassAllBuildings: 3, ArmorClassStandardBuildings: 0}, false},
		{"Barracks with Masonry and Architecture", Barracks, []TechID{Masonry, Architecture}, 1452,
			map[ArmorClass]float64{ArmorClassMelee: 2, ArmorClassPierce: 9, ArmorClassAllBuildings: 6, ArmorClassStandardBuildings: 0}, false},
		// a tower has its own unit class, with the same commands
		{"Keep with Masonry and Architecture", Keep, []TechID{Masonry, Architecture}, 2722.5,
			map[ArmorClass]float64{ArmorClassMelee: 5, ArmorClassPierce: 11, ArmorClassAllBuildings: 6, ArmorClassStoneDefense: 0}, false},
		// the Castle already has armor in the class All Buildings
		{"Castle with Masonry and Architecture", Castle, []TechID{Masonry, Architecture}, 5808,
			map[ArmorClass]float64{ArmorClassMelee: 10, ArmorClassPierce: 13, ArmorClassAllBuildings: 14, ArmorClassCastles: 0}, false},
		// the walls are not buildings of the class the technologies target
		{"Stone Wall with Masonry and Architecture", StoneWall, []TechID{Masonry, Architecture}, 1080,
			map[ArmorClass]float64{ArmorClassMelee: 8, ArmorClassPierce: 10, ArmorClassAllBuildings: 16}, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			// WHEN the building is loaded with the technologies researched
			got := tree.LoadUnitWithTechnologies(c.id, c.techIDs...)

			// THEN the HP and the armors have the technologies applied
			if math.Abs(got.HP-c.wantHP) > 1e-6 {
				t.Errorf("HP: got %v, want %v", got.HP, c.wantHP)
			}
			for class, want := range c.wantArmors {
				if got.Armors[class] != want {
					t.Errorf("armor %v: got %v, want %v", class, got.Armors[class], want)
				}
			}
		})
	}
}

func TestCalculateDamagePerHitBuilding(t *testing.T) {
	// GIVEN a Barracks at each step, and the attackers of a building
	tree, err := ParseTechTree(data.Civilizations, data.Units, data.Technologies, ParseGameStrings(data.GameStringsEnglish))
	if err != nil {
		t.Fatalf("error ParseTechTree: %v", err)
	}
	barracks := tree.LoadUnitWithTechnologies(Barracks)
	barracksMasonry := tree.LoadUnitWithTechnologies(Barracks, Masonry)
	barracksArchitecture := tree.LoadUnitWithTechnologies(Barracks, Masonry, Architecture)
	// a unit with a melee attack of 10 and the bonus attacks of the Champion and the Siege Ram
	champion := Unit{Attacks: map[ArmorClass]float64{ArmorClassMelee: 10, ArmorClassStandardBuildings: 4}}
	ram := Unit{Attacks: map[ArmorClass]float64{ArmorClassMelee: 4, ArmorClassAllBuildings: 200}}

	for _, c := range []struct {
		name     string
		attacker Unit
		target   Unit
		want     float64
	}{
		{"melee 10 and Standard Buildings 4 vs no technologies", champion, barracks, 14},
		// Masonry cuts the melee attack by 1, but not the Standard Buildings bonus
		{"melee 10 and Standard Buildings 4 vs Masonry", champion, barracksMasonry, 13},
		{"melee 10 and Standard Buildings 4 vs Architecture", champion, barracksArchitecture, 12},
		// the melee 4 of the ram is gone after 4 armor, the bonus loses 3 per technology
		{"ram vs no technologies", ram, barracks, 204},
		{"ram vs Masonry", ram, barracksMasonry, 200},
		{"ram vs Architecture", ram, barracksArchitecture, 196},
	} {
		t.Run(c.name, func(t *testing.T) {
			// WHEN the attacker hits the Barracks once
			got := CalculateDamagePerHit(c.attacker, c.target)

			// THEN the armor of the researched technologies is cut from the attack classes the Barracks lists
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
