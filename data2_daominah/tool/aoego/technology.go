package aoego

type Technology struct {
	ID               TechID
	Name             string // name without spaces, e.g. "Catapult_Tower", "Heavy_Horse_Archer", ...
	NameInGame       string // name shown in the game, e.g. "Ballista Tower", "Heavy Horse Archer", ...
	Cost             Cost
	Time             float64 // research time in seconds
	Location         UnitID  // building that researches this technology
	RequiredTechs    map[TechID]bool
	MinRequiredTechs int // e.g. Bronze Age needs 2 building from Tool Age
	Effects          []EffectFunc
}

func (t Technology) IsUnit() bool {
	return false
}

func (t Technology) GetID() UnitOrTechID {
	return UnitOrTechID(t.ID)
}

func (t Technology) GetName() string {
	return t.Name
}

func (t Technology) GetLocation() UnitID {
	return t.Location
}

func (t Technology) GetCost() Cost {
	return t.Cost
}

// TechID is enum
type TechID int

func (id TechID) IntID() int { return int(id) }

// TechID enum
const (
	StoneAgeID  TechID = 100
	ToolAgeID   TechID = 101
	BronzeAgeID TechID = 102
	IronAgeID   TechID = 103

	WatchTowerID  TechID = 16
	SentryTowerID TechID = 12
	GuardTowerID  TechID = 15
	BallistaTower TechID = 2

	WheelID         TechID = 28
	WoodworkingID   TechID = 107
	ArtisanshipID   TechID = 32
	CraftsmanshipID TechID = 110
	StoneMiningID   TechID = 109
	SiegecraftID    TechID = 111
	GoldMiningID    TechID = 108
	CoinageID       TechID = 30
	DomesticationID TechID = 81

	// TODO: add more technologies ID

	NullTechID TechID = -1
)

// EffectFunc can modify Unit attributes, enable or disable Technology
// or modify player resources. TODO: real type EffectFunc func
type EffectFunc func()

func NewTechnology(id TechID) *Technology {
	t := &Technology{ID: id, RequiredTechs: map[TechID]bool{}}
	switch id {
	case ToolAgeID:
		t.NameInGame, t.Name = "Tool Age", "Tool_Age"
		t.Cost = Cost{Food: 500}
		t.Time, t.Location = 120, TownCenterID
	case BronzeAgeID:
		t.NameInGame, t.Name = "Bronze Age", "Bronze_Age"
		t.Cost = Cost{Food: 500}
		t.Time, t.Location = 120, TownCenterID
	case IronAgeID:
		t.NameInGame, t.Name = "Iron Age", "Iron_Age"
		t.Cost = Cost{Food: 500}
		t.Time, t.Location = 120, TownCenterID

	case WatchTowerID:
		t.NameInGame, t.Name = "Watch Tower", "Watch_Tower"
		t.Cost = Cost{Food: 50}
		t.Time, t.Location = 10, GranaryID
	case SentryTowerID:
		t.NameInGame, t.Name = "Sentry Tower", "Sentry_Tower"
		t.Cost = Cost{Food: 120, Stone: 50}
		t.Time, t.Location = 30, GranaryID
	case GuardTowerID:
		t.NameInGame, t.Name = "Guard Tower", "Guard_Tower"
		t.Cost = Cost{Food: 300, Stone: 100}
		t.Time, t.Location = 75, GranaryID
	case BallistaTower:
		t.NameInGame, t.Name = "Ballista Tower", "Catapult_Tower"
		t.Cost = Cost{Food: 1800, Stone: 750}
		t.Time, t.Location = 150, GranaryID

	case WheelID:
		t.NameInGame, t.Name = "Wheel", "Wheel"
		t.Cost = Cost{Wood: 75, Food: 175}
		t.Time, t.Location = 75, MarketID

	default:
		return nil
	}
	return t
}
