package aoego

import (
	"fmt"
)

// Unit can be a Villager, Swordsman, Cavalry, ...
// or a building like TownCenter, Granary, ...
// Should be initialized with func NewUnit for default values.
type Unit struct {
	ID           UnitID
	Name         string  // name without spaces, e.g. "Man", "Soldier-Chariot2", ...
	NameInGame   string  // name shown in the game, e.g. "Villager", "Chariot Archer", ...
	Cost         Cost    // unit's cost as pointer so easier to apply civilization bonus
	Time         float64 // train time in seconds
	Population   float64 // almost all units needs 1 population, except Barracks units after Logistics researched
	Location     UnitID  // building that trains this unit
	IsBuilding   bool
	InitiateTech TechID // when the building is created, this tech is automatically researched
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

func (u Unit) GetFullName() string {
	return fmt.Sprintf("%v(%v)", u.NameInGame, u.ID.ActionID())
}

func (u Unit) GetLocation() UnitID {
	return u.Location
}

func (u Unit) GetCost() *Cost {
	clone := u.Cost
	return &clone
}

// Cost holds a certain amount of collectible resources
type Cost struct {
	Wood  float64
	Food  float64
	Gold  float64
	Stone float64
}

// Add adds the argument to the receiver, then return the receiver for chaining
func (c *Cost) Add(d Cost) *Cost {
	c.Wood += d.Wood
	c.Food += d.Food
	c.Gold += d.Gold
	c.Stone += d.Stone
	return c
}

// Multiply multiplies the receiver with the argument, then return the receiver for chaining
func (c *Cost) Multiply(m float64) *Cost {
	c.Wood *= m
	c.Food *= m
	c.Gold *= m
	c.Stone *= m
	return c
}

func (c *Cost) CheckEqual(d Cost) bool {
	return c.Wood == d.Wood && c.Food == d.Food && c.Gold == d.Gold && c.Stone == d.Stone
}

func (c *Cost) IsZero() bool {
	return c.Wood == 0 && c.Food == 0 && c.Gold == 0 && c.Stone == 0
}

func (c *Cost) String() string {
	return fmt.Sprintf("%.0f wood, %.0f food, %.0f gold, %.0f stone", c.Wood, c.Food, c.Gold, c.Stone)
}

// UnitID is enum
type UnitID int

func (id UnitID) IntID() int { return int(id) }
func (id UnitID) GetNameInGame() string {
	if u, found := AllUnits[id]; found {
		return u.NameInGame
	}
	return fmt.Sprintf("UnitID%v", id) // should not happen
}

func (id UnitID) ActionID() string {
	if u, found := AllUnits[id]; found {
		if u.IsBuilding {
			return fmt.Sprintf("B%v", id)
		}
		return fmt.Sprintf("U%v", id)
	}
	return fmt.Sprintf("UnitID%v", id) // should not happen
}

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

	FishingBoat   UnitID = 13
	TradeBoat     UnitID = 15
	TransportBoat UnitID = 17
	WarBoat       UnitID = 19
	CatapultBoat  UnitID = 250
	FireBoat      UnitID = 360

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

// NewUnit returns a Unit based on the given UnitID,
// with default attributes values (not considering civilization bonus)
func NewUnit(id UnitID) (*Unit, error) {
	u := &Unit{
		ID:           id,
		Cost:         Cost{},
		Population:   1,
		Location:     NullUnit, // will be corrected later in the switch
		IsBuilding:   AllBuildings[id],
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

	case FishingBoat:
		u.NameInGame, u.Name = "Fishing Boat", "Boat-Fishing1"
		u.Cost = Cost{Wood: 50}
		u.Time, u.Location = 40, Dock
	case TradeBoat:
		u.NameInGame, u.Name = "Trade Boat", "Boat-Trade1"
		u.Cost = Cost{Wood: 100}
		u.Time, u.Location = 50, Dock
	case TransportBoat:
		u.NameInGame, u.Name = "Light Transport", "Boat-Transport1"
		u.Cost = Cost{Wood: 150}
		u.Time, u.Location = 75, Dock
	case WarBoat:
		u.NameInGame, u.Name = "Scout Ship", "Boat-War1"
		u.Cost = Cost{Wood: 135}
		u.Time, u.Location = 60, Dock
	case CatapultBoat:
		u.NameInGame, u.Name = "Catapult Trireme", "Boat-War4"
		u.Cost = Cost{Wood: 135, Gold: 75}
		u.Time, u.Location = 90, Dock
	case FireBoat:
		u.NameInGame, u.Name = "Fire Galley", "Boat-War6"
		u.Cost = Cost{Wood: 115, Gold: 40}
		u.Time, u.Location = 45, Dock

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
		u.NameInGame, u.Name = "Camel", "Soldier-Camel"
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
	case House:
		u.NameInGame, u.Name = "House", "House"
		u.Cost = Cost{Wood: 30}
		u.Time = 20

	case Granary:
		u.NameInGame, u.Name = "Granary", "Granary"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
		u.InitiateTech = GranaryBuilt
	case StoragePit:
		u.NameInGame, u.Name = "Storage Pit", "Storage_Pit1"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
		u.InitiateTech = StoragePitBuilt
	case Barracks:
		u.NameInGame, u.Name = "Barracks", "Barracks1"
		u.Cost = Cost{Wood: 125}
		u.Time = 30
		u.InitiateTech = BarracksBuilt
	case Dock:
		u.NameInGame, u.Name = "Dock", "Dock_1"
		u.Cost = Cost{Wood: 100}
		u.Time = 50
		u.InitiateTech = DockBuilt

	case ArcheryRange:
		u.NameInGame, u.Name = "Range", "Range1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = ArcheryRangeBuilt
	case Stable:
		u.NameInGame, u.Name = "Stable", "Stable1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = StableBuilt
	case Market:
		u.NameInGame, u.Name = "Market", "Market1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
		u.InitiateTech = MarketBuilt
	case Farm:
		u.NameInGame, u.Name = "Farm", "Farm"
		u.Cost = Cost{Wood: 75}
		u.Time = 30
	case Tower:
		u.NameInGame, u.Name = "Tower", "Watch_Tower"
		u.Cost = Cost{Stone: 150}
		u.Time = 80
	case Wall:
		u.NameInGame, u.Name = "Wall", "Wall_Small"
		u.Cost = Cost{Stone: 5}
		u.Time = 7

	case GovernmentCenter:
		u.NameInGame, u.Name = "Government Center", "Government_Center"
		u.Cost = Cost{Wood: 175}
		u.Time = 60
		u.InitiateTech = GovernmentCenterBuilt
	case Temple:
		u.NameInGame, u.Name = "Temple", "Temple1"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = TempleBuilt
	case SiegeWorkshop:
		u.NameInGame, u.Name = "Siege Workshop", "Siege_Workshop"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = SiegeWorkshopBuilt
	case Academy:
		u.NameInGame, u.Name = "Academy", "Academy"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
		u.InitiateTech = AcademyBuilt

	case Wonder:
		u.NameInGame, u.Name = "Wonder", "Wonder"
		u.Cost = Cost{Wood: 1000, Gold: 1000, Stone: 1000}
		u.Time = 8000

	default:
		return nil, ErrUnitIDNotFound
	}
	return u, nil
}
