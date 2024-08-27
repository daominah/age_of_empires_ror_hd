package aoego

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
	TownCenterID UnitID = 109
	HouseID      UnitID = 70
	GranaryID    UnitID = 68
	StoragePitID UnitID = 103
	BarracksID   UnitID = 12
	DockID       UnitID = 45

	ArcheryRangeID UnitID = 87
	StableID       UnitID = 101
	MarketID       UnitID = 84
	FarmID         UnitID = 50
	TowerID        UnitID = 79
	WallID         UnitID = 72

	GovernmentCenterID UnitID = 82
	TempleID           UnitID = 104
	SiegeWorkshopID    UnitID = 49
	AcademyID          UnitID = 0

	WonderID UnitID = 276

	VillagerID UnitID = 83

	SwordsmanID UnitID = 75

	BowmanID         UnitID = 4
	ImprovedBowmanID UnitID = 5
	ChariotArcherID  UnitID = 41
	HorseArcherID    UnitID = 39
	ElephantArcherID UnitID = 25

	ScoutID    UnitID = 299
	ChariotID  UnitID = 40
	CavalryID  UnitID = 37
	ElephantID UnitID = 46
	CamelID    UnitID = 338

	PriestID UnitID = 125

	StoneThrowerID UnitID = 35
	BallistaID     UnitID = 11

	HopliteID UnitID = 93

	NullUnitID UnitID = -1
)

func CheckIsBuilding(id UnitID) bool {
	switch id {
	case TownCenterID, HouseID, GranaryID, StoragePitID, BarracksID, DockID:
		return true
	case ArcheryRangeID, StableID, MarketID, FarmID, TowerID, WallID:
		return true
	case GovernmentCenterID, TempleID, SiegeWorkshopID, AcademyID:
		return true
	case WonderID:
		return true
	default:
		return false
	}
}

// NewUnit returns a Unit based on the given UnitID,
// with default attributes values (not considering civilization bonus),
// if the UnitID is not found, returns nil
func NewUnit(id UnitID) *Unit {
	u := &Unit{
		ID:           id,
		Location:     NullUnitID, // will be corrected later in the switch
		IsBuilding:   CheckIsBuilding(id),
		InitiateTech: NullTechID,
	}
	switch id {
	case VillagerID:
		u.NameInGame, u.Name = "Villager", "Man"
		u.Cost = Cost{Food: 50}
		u.Time, u.Location = 20, TownCenterID

	case SwordsmanID:
		u.NameInGame, u.Name = "Short Swordsman", "Soldier-Inf3"
		u.Cost = Cost{Food: 35, Gold: 15}
		u.Time, u.Location = 26, BarracksID

	case BowmanID:
		u.NameInGame, u.Name = "Bowman", "Soldier-Archer1"
		u.Cost = Cost{Wood: 20, Food: 40}
		u.Time, u.Location = 30, ArcheryRangeID
	case ImprovedBowmanID:
		u.NameInGame, u.Name = "Improved Bowman", "Soldier-Archer2"
		u.Cost = Cost{Food: 40, Gold: 20}
		u.Time, u.Location = 30, ArcheryRangeID
	case ChariotArcherID:
		u.NameInGame, u.Name = "Chariot Archer", "Soldier-Chariot2"
		u.Cost = Cost{Wood: 70, Food: 40}
		u.Time, u.Location = 40, ArcheryRangeID
	case HorseArcherID:
		u.NameInGame, u.Name = "Horse Archer", "Soldier-Cavalry3_Arc"
		u.Cost = Cost{Food: 50, Gold: 70}
		u.Time, u.Location = 40, ArcheryRangeID
	case ElephantArcherID:
		u.NameInGame, u.Name = "Elephant Archer", "Soldier-Elephant1"
		u.Cost = Cost{Food: 180, Gold: 60}
		u.Time, u.Location = 50, ArcheryRangeID

	case ScoutID:
		u.NameInGame, u.Name = "Scout", "Soldier-Scout"
		u.Cost = Cost{Food: 100}
		u.Time, u.Location = 30, StableID
	case ChariotID:
		u.NameInGame, u.Name = "Chariot", "Soldier-Chariot1"
		u.Cost = Cost{Wood: 60, Food: 40}
		u.Time, u.Location = 40, StableID
	case CavalryID:
		u.NameInGame, u.Name = "Cavalry", "Soldier-Cavalry1"
		u.Cost = Cost{Food: 70, Gold: 80}
		u.Time, u.Location = 40, StableID
	case ElephantID:
		u.NameInGame, u.Name = "War Elephant", "Soldier-Elephant"
		u.Cost = Cost{Food: 170, Gold: 40}
		u.Time, u.Location = 50, StableID
	case CamelID:
		u.NameInGame, u.Name = "Camel Rider", "Soldier-Camel"
		u.Cost = Cost{Food: 70, Gold: 60}
		u.Time, u.Location = 30, StableID

	case PriestID:
		u.NameInGame, u.Name = "Priest", "Priest"
		u.Cost = Cost{Gold: 125}
		u.Time, u.Location = 50, TempleID

	case StoneThrowerID:
		u.NameInGame, u.Name = "Stone Thrower", "Soldier-Catapult1"
		u.Cost = Cost{Wood: 180, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshopID
	case BallistaID:
		u.NameInGame, u.Name = "Ballista", "Soldier-Ballista"
		u.Cost = Cost{Wood: 100, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshopID

	case HopliteID:
		u.NameInGame, u.Name = "Hoplite", "Soldier-Phal1"
		u.Cost = Cost{Food: 60, Gold: 40}
		u.Time, u.Location = 36, AcademyID

	case TownCenterID:
		u.NameInGame, u.Name = "Town Center", "Town_Center1"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
	case HouseID:
		u.NameInGame, u.Name = "House", "House"
		u.Cost = Cost{Wood: 30}
		u.Time = 20
	case GranaryID:
		u.NameInGame, u.Name = "Granary", "Granary"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
	case StoragePitID:
		u.NameInGame, u.Name = "Storage Pit", "Storage_Pit1"
		u.Cost = Cost{Wood: 120}
		u.Time = 30
	case BarracksID:
		u.NameInGame, u.Name = "Barracks", "Barracks1"
		u.Cost = Cost{Wood: 125}
		u.Time = 30
	case DockID:
		u.NameInGame, u.Name = "Dock", "Dock_1"
		u.Cost = Cost{Wood: 100}
		u.Time = 50
	case ArcheryRangeID:
		u.NameInGame, u.Name = "Archery Range", "Range1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
	case StableID:
		u.NameInGame, u.Name = "Stable", "Stable1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
	case MarketID:
		u.NameInGame, u.Name = "Market", "Market1"
		u.Cost = Cost{Wood: 150}
		u.Time = 40
	case FarmID:
		u.NameInGame, u.Name = "Farm", "Farm"
		u.Cost = Cost{Wood: 75}
		u.Time = 30
	case TowerID:
		u.NameInGame, u.Name = "Watch Tower", "Watch_Tower"
		u.Cost = Cost{Stone: 150}
		u.Time = 80
	case WallID:
		u.NameInGame, u.Name = "Small Wall", "Wall_Small"
		u.Cost = Cost{Stone: 5}
		u.Time = 7
	case GovernmentCenterID:
		u.NameInGame, u.Name = "Government Center", "Government_Center"
		u.Cost = Cost{Wood: 175}
		u.Time = 60
	case TempleID:
		u.NameInGame, u.Name = "Temple", "Temple1"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
	case SiegeWorkshopID:
		u.NameInGame, u.Name = "Siege Workshop", "Siege_Workshop"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
	case AcademyID:
		u.NameInGame, u.Name = "Academy", "Academy"
		u.Cost = Cost{Wood: 200}
		u.Time = 60
	case WonderID:
		u.NameInGame, u.Name = "Wonder", "Wonder"
		u.Cost = Cost{Wood: 1000, Gold: 1000, Stone: 1000}
		u.Time = 8000

	default:
		return nil
	}
	return u
}
