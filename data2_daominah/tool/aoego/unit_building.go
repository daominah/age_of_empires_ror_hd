package aoego

// Unit can be a Villager, Swordsman, Cavalry, ...
// or a building like TownCenter, Granary, ...
type Unit struct {
	ID           UnitID
	Name         string // name shown in the game, e.g. Villager, Chariot Archer, ...
	NameInternal string // internal code name, e.g. Man, Soldier-Chariot2, ...
	Cost         Storage
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

func (u Unit) GetNameInternal() string {
	return u.NameInternal
}

func (u Unit) GetLocation() UnitID {
	return u.Location
}

func (u Unit) GetCost() Storage {
	return u.Cost
}

// Storage holds a certain amount of collectible resources: wood, food, gold, and stone;
// Storage can be used to represent cost of a Unit or a Technology
type Storage struct {
	Wood  float64
	Food  float64
	Gold  float64
	Stone float64
}

// UnitID is enum
type UnitID int

// ID is just a convenient method so IDE can suggest code completion
func (id UnitID) ID() UnitOrTechID { return UnitOrTechID(id) }

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
		u.Name, u.NameInternal = "Villager", "Man"
		u.Cost = Storage{Food: 50}
		u.Time, u.Location = 20, TownCenterID

	case SwordsmanID:
		u.Name, u.NameInternal = "Short Swordsman", "Soldier-Inf3"
		u.Cost = Storage{Food: 35, Gold: 15}
		u.Time, u.Location = 26, BarracksID

	case BowmanID:
		u.Name, u.NameInternal = "Bowman", "Soldier-Archer1"
		u.Cost = Storage{Wood: 20, Food: 40}
		u.Time, u.Location = 30, ArcheryRangeID
	case ImprovedBowmanID:
		u.Name, u.NameInternal = "Improved Bowman", "Soldier-Archer2"
		u.Cost = Storage{Food: 40, Gold: 20}
		u.Time, u.Location = 30, ArcheryRangeID
	case ChariotArcherID:
		u.Name, u.NameInternal = "Chariot Archer", "Soldier-Chariot2"
		u.Cost = Storage{Wood: 70, Food: 40}
		u.Time, u.Location = 40, ArcheryRangeID
	case HorseArcherID:
		u.Name, u.NameInternal = "Horse Archer", "Soldier-Cavalry3_Arc"
		u.Cost = Storage{Food: 50, Gold: 70}
		u.Time, u.Location = 40, ArcheryRangeID
	case ElephantArcherID:
		u.Name, u.NameInternal = "Elephant Archer", "Soldier-Elephant1"
		u.Cost = Storage{Food: 180, Gold: 60}
		u.Time, u.Location = 50, ArcheryRangeID

	case ScoutID:
		u.Name, u.NameInternal = "Scout", "Soldier-Scout"
		u.Cost = Storage{Food: 100}
		u.Time, u.Location = 30, StableID
	case ChariotID:
		u.Name, u.NameInternal = "Chariot", "Soldier-Chariot1"
		u.Cost = Storage{Wood: 60, Food: 40}
		u.Time, u.Location = 40, StableID
	case CavalryID:
		u.Name, u.NameInternal = "Cavalry", "Soldier-Cavalry1"
		u.Cost = Storage{Food: 70, Gold: 80}
		u.Time, u.Location = 40, StableID
	case ElephantID:
		u.Name, u.NameInternal = "War Elephant", "Soldier-Elephant"
		u.Cost = Storage{Food: 170, Gold: 40}
		u.Time, u.Location = 50, StableID
	case CamelID:
		u.Name, u.NameInternal = "Camel Rider", "Soldier-Camel"
		u.Cost = Storage{Food: 70, Gold: 60}
		u.Time, u.Location = 30, StableID

	case PriestID:
		u.Name, u.NameInternal = "Priest", "Priest"
		u.Cost = Storage{Gold: 125}
		u.Time, u.Location = 50, TempleID

	case StoneThrowerID:
		u.Name, u.NameInternal = "Stone Thrower", "Soldier-Catapult1"
		u.Cost = Storage{Wood: 180, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshopID
	case BallistaID:
		u.Name, u.NameInternal = "Ballista", "Soldier-Ballista"
		u.Cost = Storage{Wood: 100, Gold: 80}
		u.Time, u.Location = 60, SiegeWorkshopID

	case HopliteID:
		u.Name, u.NameInternal = "Hoplite", "Soldier-Phal1"
		u.Cost = Storage{Food: 60, Gold: 40}
		u.Time, u.Location = 36, AcademyID

	case TownCenterID:
		u.Name, u.NameInternal = "Town Center", "Town_Center1"
		u.Cost = Storage{Wood: 200}
		u.Time = 60
	case HouseID:
		u.Name, u.NameInternal = "House", "House"
		u.Cost = Storage{Wood: 30}
		u.Time = 20
	case GranaryID:
		u.Name, u.NameInternal = "Granary", "Granary"
		u.Cost = Storage{Wood: 120}
		u.Time = 30
	case StoragePitID:
		u.Name, u.NameInternal = "Storage Pit", "Storage_Pit1"
		u.Cost = Storage{Wood: 120}
		u.Time = 30
	case BarracksID:
		u.Name, u.NameInternal = "Barracks", "Barracks1"
		u.Cost = Storage{Wood: 125}
		u.Time = 30
	case DockID:
		u.Name, u.NameInternal = "Dock", "Dock_1"
		u.Cost = Storage{Wood: 100}
		u.Time = 50
	case ArcheryRangeID:
		u.Name, u.NameInternal = "Archery Range", "Range1"
		u.Cost = Storage{Wood: 150}
		u.Time = 40
	case StableID:
		u.Name, u.NameInternal = "Stable", "Stable1"
		u.Cost = Storage{Wood: 150}
		u.Time = 40
	case MarketID:
		u.Name, u.NameInternal = "Market", "Market1"
		u.Cost = Storage{Wood: 150}
		u.Time = 40
	case FarmID:
		u.Name, u.NameInternal = "Farm", "Farm"
		u.Cost = Storage{Wood: 75}
		u.Time = 30
	case TowerID:
		u.Name, u.NameInternal = "Watch Tower", "Watch_Tower"
		u.Cost = Storage{Stone: 150}
		u.Time = 80
	case WallID:
		u.Name, u.NameInternal = "Small Wall", "Wall_Small"
		u.Cost = Storage{Stone: 5}
		u.Time = 7
	case GovernmentCenterID:
		u.Name, u.NameInternal = "Government Center", "Government_Center"
		u.Cost = Storage{Wood: 175}
		u.Time = 60
	case TempleID:
		u.Name, u.NameInternal = "Temple", "Temple1"
		u.Cost = Storage{Wood: 200}
		u.Time = 60
	case SiegeWorkshopID:
		u.Name, u.NameInternal = "Siege Workshop", "Siege_Workshop"
		u.Cost = Storage{Wood: 200}
		u.Time = 60
	case AcademyID:
		u.Name, u.NameInternal = "Academy", "Academy"
		u.Cost = Storage{Wood: 200}
		u.Time = 60
	case WonderID:
		u.Name, u.NameInternal = "Wonder", "Wonder"
		u.Cost = Storage{Wood: 1000, Gold: 1000, Stone: 1000}
		u.Time = 8000

	default:
		return nil
	}
	return u
}
