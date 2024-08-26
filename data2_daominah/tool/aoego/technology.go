package aoego

type Technology struct {
	ID               TechID
	Name             string // name shown in the game, e.g. Ballista Tower, Heavy Horse Archer, ...
	NameInternal     string // internal code name, e.g. Catapult_Tower, Heavy_Horse_Archer, ...
	Cost             Cost
	Time             float64 // research time in seconds
	Location         UnitID  // building that researches this technology
	RequiredTechs    map[TechID]bool
	MinRequiredTechs int // e.g. Bronze Age needs 2 building from Tool Age
}

// TechID enum
type TechID int

// TechID enum
const (
	StoneAgeID  TechID = 100
	ToolAgeID   TechID = 101
	BronzeAgeID TechID = 102
	IronAgeID   TechID = 103

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
