package aoego

import (
	"fmt"
)

// Unit can be a Villager, Swordsman, Cavalry, ...
// or a building like TownCenter, Granary, ...
type Unit struct {
	ID           UnitID
	Name         string // name without spaces, e.g. "Man", "Soldier-Chariot2", ...
	NameInGame   string // name shown in the game, e.g. "Villager", "Chariot Archer", ...
	Cost         Cost
	Time         float64 // train time in seconds
	Location     UnitID  // building that trains this unit
	IsBuilding   bool
	InitiateTech TechID // when the building is created, this tech is researched
}

func (u Unit) IsUnit() bool {
	return true
}

func (u Unit) GetID() UnitOrTechID {
	return UnitOrTechID(u.ID)
}

func (u Unit) GetName() string {
	return u.Name
}

func (u Unit) GetLocation() UnitID {
	return u.Location
}

func (u Unit) GetCost() Cost {
	return u.Cost
}

// Cost holds a certain amount of collectible resources
type Cost struct {
	Wood  float64
	Food  float64
	Gold  float64
	Stone float64
}

func (c Cost) Add(d Cost) Cost {
	return Cost{
		Wood:  c.Wood + d.Wood,
		Food:  c.Food + d.Food,
		Gold:  c.Gold + d.Gold,
		Stone: c.Stone + d.Stone,
	}
}

func (c Cost) Multiply(m float64) Cost {
	return Cost{
		Wood:  c.Wood * m,
		Food:  c.Food * m,
		Gold:  c.Gold * m,
		Stone: c.Stone * m,
	}
}

// UnitID is enum
type UnitID int

func (id UnitID) IntID() int { return int(id) }

// UnitID enum
const (
	TownCenter UnitID = 109
	House      UnitID = 70
	Granary    UnitID = 68
	StoragePit UnitID = 103
	Barracks   UnitID = 12
	Dock       UnitID = 45

	ArcheryRange UnitID = 87
	Stable       UnitID = 101
	Market       UnitID = 84
	Farm         UnitID = 50
	Tower        UnitID = 79
	Wall         UnitID = 72

	GovernmentCenter UnitID = 82
	Temple           UnitID = 104
	SiegeWorkshop    UnitID = 49
	Academy          UnitID = 0

	Wonder UnitID = 276

	Villager UnitID = 83

	Clubman   UnitID = 73
	Swordsman UnitID = 75
	Slinger   UnitID = 347

	Bowman         UnitID = 4
	ImprovedBowman UnitID = 5
	ChariotArcher  UnitID = 41
	HorseArcher    UnitID = 39
	ElephantArcher UnitID = 25

	Scout    UnitID = 299
	Chariot  UnitID = 40
	Cavalry  UnitID = 37
	Elephant UnitID = 46
	Camel    UnitID = 338

	Priest UnitID = 125

	StoneThrower UnitID = 35
	Ballista     UnitID = 11

	Hoplite UnitID = 93

	NullUnit UnitID = -1
)

// Buildings is a list of all building,
// used as a constant (do not change this var in runtime)
var Buildings = map[UnitID]bool{
	TownCenter: true, House: true, Granary: true, StoragePit: true, Barracks: true, Dock: true,
	ArcheryRange: true, Stable: true, Market: true, Farm: true, Tower: true, Wall: true,
	GovernmentCenter: true, Temple: true, SiegeWorkshop: true, Academy: true,
	Wonder: true,
}

// Combatants is a list of all units that is not a building,
// used as a constant (do not change this var in runtime)
var Combatants = map[UnitID]bool{
	Villager: true,
	Clubman:  true, Swordsman: true, Slinger: true,
	Bowman: true, ImprovedBowman: true, ChariotArcher: true, HorseArcher: true, ElephantArcher: true,
	Scout: true, Chariot: true, Cavalry: true, Elephant: true, Camel: true,
	Priest:       true,
	StoneThrower: true, Ballista: true,
	Hoplite: true,
}

// NewUnit returns a Unit based on the given UnitID,
// with default attributes values (not considering civilization bonus),
// this func PANIC if the UnitID is not found.
func NewUnit(id UnitID) *Unit {
	u := &Unit{
		ID:           id,
		Location:     NullUnit, // will be corrected later in the switch
		IsBuilding:   Buildings[id],
		InitiateTech: NullTech,
	}
	switch id {
	case Villager:
		u.NameInGame, u.Name = "Villager", "Man"
		u.Cost = Cost{Food: 50}
		u.Time, u.Location = 20, TownCenter

	case Clubman:
		u.NameInGame, u.Name = "Clubman", "Soldier-Inf1"
		u.Cost = Cost{Food: 50}
		u.Time, u.Location = 26, Barracks
	case Swordsman:
		u.NameInGame, u.Name = "Short Swordsman", "Soldier-Inf3"
		u.Cost = Cost{Food: 35, Gold: 15}
		u.Time, u.Location = 26, Barracks
	case Slinger:
		u.NameInGame, u.Name = "Slinger", "Soldier-Slinger"
		u.Cost = Cost{Food: 40, Stone: 10}
		u.Time, u.Location = 24, Barracks

	case Bowman:
		u.NameInGame, u.Name = "Bowman", "Soldier-Archer1"
		u.Cost = Cost{Wood: 20, Food: 40}
		u.Time, u.Location = 30, ArcheryRange
	case ImprovedBowman:
		u.NameInGame, u.Name = "Improved Bowman", "Soldier-Archer2"
		u.Cost = Cost{Food: 40, Gold: 20}
		u.Time, u.Location = 30, ArcheryRange
	case ChariotArcher:
		u.NameInGame, u.Name = "Chariot Archer", "Soldier-Chariot2"
		u.Cost = Cost{Wood: 70, Food: 40}
		u.Time, u.Location = 40, ArcheryRange
	case HorseArcher:
		u.NameInGame, u.Name = "Horse Archer", "Soldier-Cavalry3_Arc"
		u.Cost = Cost{Food: 50, Gold: 70}
		u.Time, u.Location = 40, ArcheryRange
	case ElephantArcher:
		u.NameInGame, u.Name = "Elephant Archer", "Soldier-Elephant1"
		u.Cost = Cost{Food: 180, Gold: 60}
		u.Time, u.Location = 50, ArcheryRange

	case Scout:
		u.NameInGame, u.Name = "Scout", "Soldier-Scout"
		u.Cost = Cost{Food: 100}
		u.Time, u.Location = 30, Stable
	case Chariot:
		u.NameInGame, u.Name = "Chariot", "Soldier-Chariot1"
		u.Cost = Cost{Wood: 60, Food: 40}
		u.Time, u.Location = 40, Stable
	case Cavalry:
		u.NameInGame, u.Name = "Cavalry", "Soldier-Cavalry1"
		u.Cost = Cost{Food: 70, Gold: 80}
		u.Time, u.Location = 40, Stable
	case Elephant:
		u.NameInGame, u.Name = "War Elephant", "Soldier-Elephant"
		u.Cost = Cost{Food: 170, Gold: 40}
		u.Time, u.Location = 50, Stable
	case Camel:
		u.NameInGame, u.Name = "Camel Rider", "Soldier-Camel"
		u.Cost = Cost{Food: 70, Gold: 60}
		u.Time, u.Location = 30, Stable

	case Priest:
		u.NameInGame, u.Name = "Priest", "Priest"
		u.Cost = Cost{Gold: 125}
		u.Time, u.Location = 50, Temple

	case StoneThrower:
		u.NameInGame, u.Name = "Stone Thrower", "Soldier-Catapult1"
		u.Cost = Cost{Wood: 180, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshop
	case Ballista:
		u.NameInGame, u.Name = "Ballista", "Soldier-Ballista"
		u.Cost = Cost{Wood: 100, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshop

	case Hoplite:
		u.NameInGame, u.Name = "Hoplite", "Soldier-Phal1"
		u.Cost = Cost{Food: 60, Gold: 40}
		u.Time, u.Location = 36, Academy

	case TownCenter:
		u.NameInGame, u.Name = "Town Center", "Town_Center1"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = EnableGranaryStoragePitBarracksDock
	case House:
		u.NameInGame, u.Name = "House", "House"
		u.Cost = Cost{Wood: 30}
		u.Time = 20

	case Granary:
		u.NameInGame, u.Name = "Granary", "Granary"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
		u.InitiateTech = GranaryTech
	case StoragePit:
		u.NameInGame, u.Name = "Storage Pit", "Storage_Pit1"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
		u.InitiateTech = StoragePitTech
	case Barracks:
		u.NameInGame, u.Name = "Barracks", "Barracks1"
		u.Cost = Cost{Wood: 125}
		u.Time = 30
		u.InitiateTech = BarracksTech
	case Dock:
		u.NameInGame, u.Name = "Dock", "Dock_1"
		u.Cost = Cost{Wood: 100}
		u.Time = 50
		u.InitiateTech = DockTech

	case ArcheryRange:
		u.NameInGame, u.Name = "Archery Range", "Range1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = ArcheryRangeTech
	case Stable:
		u.NameInGame, u.Name = "Stable", "Stable1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = StableTech
	case Market:
		u.NameInGame, u.Name = "Market", "Market1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = MarketTech
	case Farm:
		u.NameInGame, u.Name = "Farm", "Farm"
		u.Cost = Cost{Wood: 75}
		u.Time = 30
	case Tower:
		u.NameInGame, u.Name = "Watch Tower", "Watch_Tower"
		u.Cost = Cost{Stone: 150}
		u.Time = 80
	case Wall:
		u.NameInGame, u.Name = "Small Wall", "Wall_Small"
		u.Cost = Cost{Stone: 5}
		u.Time = 7

	case GovernmentCenter:
		u.NameInGame, u.Name = "Government Center", "Government_Center"
		u.Cost = Cost{Wood: 175}
		u.Time = 60
		u.InitiateTech = GovernmentCenterTech
	case Temple:
		u.NameInGame, u.Name = "Temple", "Temple1"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = TempleTech
	case SiegeWorkshop:
		u.NameInGame, u.Name = "Siege Workshop", "Siege_Workshop"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = SiegeWorkshopTech
	case Academy:
		u.NameInGame, u.Name = "Academy", "Academy"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = AcademyTech

	case Wonder:
		u.NameInGame, u.Name = "Wonder", "Wonder"
		u.Cost = Cost{Wood: 1000, Gold: 1000, Stone: 1000}
		u.Time = 8000

	default:
		panic(fmt.Errorf("NewUnit: %v: %w", id, ErrUnitIDNotFound))
	}
	return u
}
