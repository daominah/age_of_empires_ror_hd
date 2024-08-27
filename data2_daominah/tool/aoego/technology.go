package aoego

type Technology struct {
	ID               TechID
	Name             string // name shown in the game, e.g. Ballista Tower, Heavy Horse Archer, ...
	NameInternal     string // name without space, e.g. Catapult_Tower, Heavy_Horse_Archer, ...
	Cost             Storage
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

func (t Technology) GetNameInternal() string {
	return t.NameInternal
}

func (t Technology) GetLocation() UnitID {
	return t.Location
}

func (t Technology) GetCost() Storage {
	return t.Cost
}

// TechID is enum
type TechID int

// ID is just a convenient method so IDE can suggest code completion
func (id TechID) ID() UnitOrTechID { return UnitOrTechID(id) }

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
// or modify player resources
type EffectFunc func()

func NewTechnology(id TechID) *Technology {
	t := &Technology{ID: id, RequiredTechs: map[TechID]bool{}}
	switch id {
	case ToolAgeID:
		t.Name, t.NameInternal = "Tool Age", "Tool_Age"
		t.Cost = Storage{Food: 500}
		t.Time, t.Location = 120, TownCenterID
	case BronzeAgeID:
		t.Name, t.NameInternal = "Bronze Age", "Bronze_Age"
		t.Cost = Storage{Food: 500}
		t.Time, t.Location = 120, TownCenterID
	case IronAgeID:
		t.Name, t.NameInternal = "Iron Age", "Iron_Age"
		t.Cost = Storage{Food: 500}
		t.Time, t.Location = 120, TownCenterID

	case WatchTowerID:
		t.Name, t.NameInternal = "Watch Tower", "Watch_Tower"
		t.Cost = Storage{Food: 50}
		t.Time, t.Location = 10, GranaryID
	case SentryTowerID:
		t.Name, t.NameInternal = "Sentry Tower", "Sentry_Tower"
		t.Cost = Storage{Food: 120, Stone: 50}
		t.Time, t.Location = 30, GranaryID
	case GuardTowerID:
		t.Name, t.NameInternal = "Guard Tower", "Guard_Tower"
		t.Cost = Storage{Food: 300, Stone: 100}
		t.Time, t.Location = 75, GranaryID
	case BallistaTower:
		t.Name, t.NameInternal = "Ballista Tower", "Catapult_Tower"
		t.Cost = Storage{Food: 1800, Stone: 750}
		t.Time, t.Location = 150, GranaryID

	default:
		return nil
	}
	return t
}
