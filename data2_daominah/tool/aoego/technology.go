package aoego

import (
	"fmt"
)

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
	StoneAge  TechID = 100
	ToolAge   TechID = 101
	BronzeAge TechID = 102
	IronAge   TechID = 103

	WatchTower    TechID = 16
	SentryTower   TechID = 12
	GuardTower    TechID = 15
	BallistaTower TechID = 2

	LeatherArmorArchers  TechID = 41
	ScaleArmorArchers    TechID = 44
	ChainMailArchers     TechID = 49
	Toolworking          TechID = 46
	Metalworking         TechID = 51
	Metallurgy           TechID = 52
	LeatherArmorInfantry TechID = 40
	ScaleArmorInfantry   TechID = 43
	ChainMailInfantry    TechID = 48
	LeatherArmorCavalry  TechID = 42
	ScaleArmorCavalry    TechID = 45
	ChainMailCavalry     TechID = 50
	BronzeShield         TechID = 47
	IronShield           TechID = 117
	TowerShield          TechID = 122

	Wheel         TechID = 28
	Woodworking   TechID = 107
	Artisanship   TechID = 32
	Craftsmanship TechID = 110
	StoneMining   TechID = 109
	Siegecraft    TechID = 111
	GoldMining    TechID = 108
	Coinage       TechID = 30
	Domestication TechID = 81

	Nobility     TechID = 34
	Writing      TechID = 114
	Architecture TechID = 112
	Logistics    TechID = 121
	Aristocracy  TechID = 113
	Ballistics   TechID = 106
	Alchemy      TechID = 37
	Engineering  TechID = 35

	Axe        TechID = 63
	ShortSword TechID = 64
	Broadsword TechID = 65
	LongSword  TechID = 66
	Legion     TechID = 77

	ScytheChariot   TechID = 126
	HeavyCalvary    TechID = 71
	Cataphract      TechID = 78
	ArmoredElephant TechID = 125

	ImprovedBow      TechID = 56
	CompositeBow     TechID = 57
	HeavyHorseArcher TechID = 38

	Catapult        TechID = 11
	MassiveCatapult TechID = 36
	Helepolis       TechID = 27

	Phalanx   TechID = 73
	Centurion TechID = 79

	Astrology  TechID = 22
	Mysticism  TechID = 21
	Polytheism TechID = 24
	Medicine   TechID = 119
	Afterlife  TechID = 18
	Monotheism TechID = 19
	Fanaticism TechID = 20
	Jihad      TechID = 23
	Sacrifice  TechID = 120

	NullTech TechID = -1

	// the following are internal techs, not shown in the game,
	// they will be automatically researched when a building is built:

	EnableGranaryStoragePitBarracksDock TechID = 105
	GranaryTech                         TechID = 10
	StoragePitTech                      TechID = 39
	BarracksTech                        TechID = 62
	DockTech                            TechID = 0

	EnableMarket     TechID = 94
	MarketTech       TechID = 26
	ArcheryRangeTech TechID = 55
	StableTech       TechID = 67

	GovernmentCenterTech TechID = 33
	TempleTech           TechID = 17
	SiegeWorkshopTech    TechID = 53
	AcademyTech          TechID = 72
)

// EffectFunc can modify Unit attributes, enable or disable Technology
// or modify player resources. TODO: real type EffectFunc func
type EffectFunc func(empire *EmpireDeveloping)

// NewTechnology PANIC if the TechID is not found.
func NewTechnology(id TechID) *Technology {
	t := &Technology{ID: id, RequiredTechs: map[TechID]bool{}}
	switch id {
	case ToolAge:
		t.NameInGame, t.Name = "Tool Age", "Tool_Age"
		t.Cost = Cost{Food: 500}
		t.Time, t.Location = 120, TownCenter
		t.RequiredTechs, t.MinRequiredTechs = map[TechID]bool{
			GranaryTech:    true,
			StoragePitTech: true,
			BarracksTech:   true,
			DockTech:       true,
		}, 2
	case BronzeAge:
		t.NameInGame, t.Name = "Bronze Age", "Bronze_Age"
		t.Cost = Cost{Food: 800}
		t.Time, t.Location = 140, TownCenter
		t.RequiredTechs, t.MinRequiredTechs = map[TechID]bool{
			MarketTech:       true,
			ArcheryRangeTech: true,
			StableTech:       true,
		}, 2
	case IronAge:
		t.NameInGame, t.Name = "Iron Age", "Iron_Age"
		t.Cost = Cost{Food: 1000, Gold: 800}
		t.Time, t.Location = 160, TownCenter

	case WatchTower:
		t.NameInGame, t.Name = "Watch Tower", "Watch_Tower"
		t.Cost = Cost{Food: 50}
		t.Time, t.Location = 10, Granary
	case SentryTower:
		t.NameInGame, t.Name = "Sentry Tower", "Sentry_Tower"
		t.Cost = Cost{Food: 120, Stone: 50}
		t.Time, t.Location = 30, Granary
	case GuardTower:
		t.NameInGame, t.Name = "Guard Tower", "Guard_Tower"
		t.Cost = Cost{Food: 300, Stone: 100}
		t.Time, t.Location = 75, Granary
	case BallistaTower:
		t.NameInGame, t.Name = "Ballista Tower", "Catapult_Tower"
		t.Cost = Cost{Food: 1800, Stone: 750}
		t.Time, t.Location = 150, Granary

	case Wheel:
		t.NameInGame, t.Name = "Wheel", "Wheel"
		t.Cost = Cost{Wood: 75, Food: 175}
		t.Time, t.Location = 75, Market
	case Woodworking:
		t.NameInGame, t.Name = "Woodworking", "Wood_Working"
		t.Cost = Cost{Wood: 75, Food: 120}
		t.Time, t.Location = 60, Market
	case Artisanship:
		t.NameInGame, t.Name = "Artisanship", "Artisanship"
		t.Cost = Cost{Wood: 150, Food: 170}
		t.Time, t.Location = 80, Market
	case Craftsmanship:
		t.NameInGame, t.Name = "Craftsmanship", "Craftmanship"
		t.Cost = Cost{Wood: 200, Food: 240}
		t.Time, t.Location = 100, Market
	case StoneMining:
		t.NameInGame, t.Name = "Stone Mining", "Stone_Mining"
		t.Cost = Cost{Food: 100, Stone: 50}
		t.Time, t.Location = 30, Market
	case Siegecraft:
		t.NameInGame, t.Name = "Siegecraft", "Siegecraft"
		t.Cost = Cost{Food: 190, Stone: 100}
		t.Time, t.Location = 60, Market
	case GoldMining:
		t.NameInGame, t.Name = "Gold Mining", "Gold_Mining"
		t.Cost = Cost{Wood: 100, Food: 120}
		t.Time, t.Location = 50, Market
	case Coinage:
		t.NameInGame, t.Name = "Coinage", "Coinage"
		t.Cost = Cost{Food: 200, Gold: 100}
		t.Time, t.Location = 60, Market
	case Domestication:
		t.NameInGame, t.Name = "Domestication", "Domestication"
		t.Cost = Cost{Wood: 50, Food: 200}

	case Nobility:
		t.NameInGame, t.Name = "Nobility", "Nobility"
		t.Cost = Cost{Food: 175, Gold: 120}
		t.Time, t.Location = 70, GovernmentCenter
	case Writing:
		t.NameInGame, t.Name = "Writing", "Writing"
		t.Cost = Cost{Food: 200, Gold: 75}
		t.Time, t.Location = 60, GovernmentCenter
	case Architecture:
		t.NameInGame, t.Name = "Architecture", "Architecture"
		t.Cost = Cost{Wood: 175, Food: 150}
		t.Time, t.Location = 50, GovernmentCenter
	case Logistics:
		t.NameInGame, t.Name = "Logistics", "Logistics"
		t.Cost = Cost{Food: 180, Gold: 100}
		t.Time, t.Location = 60, GovernmentCenter
	case Aristocracy:
		t.NameInGame, t.Name = "Aristocracy", "Aristocracy"
		t.Cost = Cost{Food: 175, Gold: 150}
		t.Time, t.Location = 60, GovernmentCenter
	case Ballistics:
		t.NameInGame, t.Name = "Ballistics", "Ballistics"
		t.Cost = Cost{Food: 200, Gold: 50}
		t.Time, t.Location = 60, GovernmentCenter
	case Alchemy:
		t.NameInGame, t.Name = "Alchemy", "Alchemy"
		t.Cost = Cost{Food: 250, Gold: 200}
		t.Time, t.Location = 100, GovernmentCenter
	case Engineering:
		t.NameInGame, t.Name = "Engineering", "Engineering"
		t.Cost = Cost{Wood: 100, Food: 200}
		t.Time, t.Location = 70, GovernmentCenter

	case Toolworking:
		t.NameInGame, t.Name = "Toolworking", "Toolworking"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 40, StoragePit
	case Metalworking:
		t.NameInGame, t.Name = "Metalworking", "Metal_Working"
		t.Cost = Cost{Food: 200, Gold: 120}
		t.Time, t.Location = 75, StoragePit
	case Metallurgy:
		t.NameInGame, t.Name = "Metallurgy", "Metallurgy"
		t.Cost = Cost{Food: 300, Gold: 180}
		t.Time, t.Location = 100, StoragePit
	case LeatherArmorInfantry:
		t.NameInGame, t.Name = "Leather Armor Infantry", "Leather_Armor_-_Soldie"
		t.Cost = Cost{Food: 75}
		t.Time, t.Location = 30, StoragePit
	case ScaleArmorInfantry:
		t.NameInGame, t.Name = "Scale Armor Infantry", "Scale_Armor_-_Soldiers"
		t.Cost = Cost{Food: 100, Gold: 50}
		t.Time, t.Location = 60, StoragePit
	case ChainMailInfantry:
		t.NameInGame, t.Name = "Chain Mail Infantry", "Chain_Mail_-_Soldiers"
		t.Cost = Cost{Food: 125, Gold: 100}
		t.Time, t.Location = 75, StoragePit
	case LeatherArmorArchers:
		t.NameInGame, t.Name = "Leather Armor Archers", "Leather_Armor_-_Archer"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 30, StoragePit
	case ScaleArmorArchers:
		t.NameInGame, t.Name = "Scale Armor Archers", "Scale_Armor_-_Archers"
		t.Cost = Cost{Food: 125, Gold: 50}
		t.Time, t.Location = 60, StoragePit
	case ChainMailArchers:
		t.NameInGame, t.Name = "Chain Mail Archers", "Chain_Mail_-_Archers"
		t.Cost = Cost{Food: 150, Gold: 100}
		t.Time, t.Location = 75, StoragePit
	case LeatherArmorCavalry:
		t.NameInGame, t.Name = "Leather Armor Cavalry", "Leather_Armor_Mounted"
		t.Cost = Cost{Food: 125}
		t.Time, t.Location = 30, StoragePit
	case ScaleArmorCavalry:
		t.NameInGame, t.Name = "Scale Armor Cavalry", "Scale_Armor_-_Cavalry"
		t.Cost = Cost{Food: 150, Gold: 50}
		t.Time, t.Location = 60, StoragePit
	case ChainMailCavalry:
		t.NameInGame, t.Name = "Chain Mail Cavalry", "Chain_Mail_-_Cavalry"
		t.Cost = Cost{Food: 175, Gold: 100}
		t.Time, t.Location = 75, StoragePit
	case BronzeShield:
		t.NameInGame, t.Name = "Bronze Shield", "Bronze_Shield"
		t.Cost = Cost{Food: 150, Gold: 180}
		t.Time, t.Location = 50, StoragePit
	case IronShield:
		t.NameInGame, t.Name = "Iron Shield", "Iron_Shield"
		t.Cost = Cost{Food: 200, Gold: 320}
		t.Time, t.Location = 75, StoragePit
	case TowerShield:
		t.NameInGame, t.Name = "Tower Shield", "Tower_Shield"
		t.Cost = Cost{Food: 250, Gold: 400}
		t.Time, t.Location = 100, StoragePit

	case Axe:
		t.NameInGame, t.Name = "Battle Axe", "Axe"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 40, Barracks
	case ShortSword:
		t.NameInGame, t.Name = "Short Sword", "Short_Sword"
		t.Cost = Cost{Food: 120, Gold: 50}
		t.Time, t.Location = 50, Barracks
	case Broadsword:
		t.NameInGame, t.Name = "Broadsword", "Broad_Sword"
		t.Cost = Cost{Food: 140, Gold: 50}
		t.Time, t.Location = 80, Barracks
	case LongSword:
		t.NameInGame, t.Name = "Long Sword", "Long_Sword"
		t.Cost = Cost{Food: 160, Gold: 50}
		t.Time, t.Location = 90, Barracks
	case Legion:
		t.NameInGame, t.Name = "Legion", "Legion"
		t.Cost = Cost{Food: 1400, Gold: 600}
		t.Time, t.Location = 150, Barracks

	case ScytheChariot:
		t.NameInGame, t.Name = "Scythe Chariot", "Scythe_Chariot"
		t.Cost = Cost{Wood: 1200, Gold: 800}
		t.Time, t.Location = 150, Stable
	case HeavyCalvary:
		t.NameInGame, t.Name = "Heavy Calvary", "Heavy_Cavalry"
		t.Cost = Cost{Food: 350, Gold: 125}
		t.Time, t.Location = 90, Stable
	case Cataphract:
		t.NameInGame, t.Name = "Cataphract", "Cataphracts"
		t.Cost = Cost{Food: 2000, Gold: 850}
		t.Time, t.Location = 150, Stable
	case ArmoredElephant:
		t.NameInGame, t.Name = "Armored Elephant", "Armored_Elephant"
		t.Cost = Cost{Food: 1000, Gold: 1200}
		t.Time, t.Location = 150, Stable

	case ImprovedBow:
		t.NameInGame, t.Name = "Improved Bow", "Improved_Bow"
		t.Cost = Cost{Wood: 80, Food: 140}
		t.Time, t.Location = 60, ArcheryRange
	case CompositeBow:
		t.NameInGame, t.Name = "Composite Bow", "Composit_bow"
		t.Cost = Cost{Wood: 100, Food: 180}
		t.Time, t.Location = 100, ArcheryRange
	case HeavyHorseArcher:
		t.NameInGame, t.Name = "Heavy Horse Archer", "Heavy_Horse_Archer"
		t.Cost = Cost{Food: 1750, Gold: 800}
		t.Time, t.Location = 150, ArcheryRange

	case Catapult:
		t.NameInGame, t.Name = "Catapult", "Heavy_Catapult"
		t.Cost = Cost{Wood: 250, Food: 300}
		t.Time, t.Location = 100, SiegeWorkshop
	case MassiveCatapult:
		t.NameInGame, t.Name = "Heavy Catapult", "Massive_Catapult"
		t.Cost = Cost{Wood: 900, Food: 1800}
		t.Time, t.Location = 150, SiegeWorkshop
	case Helepolis:
		t.NameInGame, t.Name = "Helepolis", "Helepolis"
		t.Cost = Cost{Wood: 1000, Food: 1500}
		t.Time, t.Location = 150, SiegeWorkshop

	case Phalanx:
		t.NameInGame, t.Name = "Phalanx", "Phalanx"
		t.Cost = Cost{Food: 300, Gold: 100}
		t.Time, t.Location = 90, Academy
	case Centurion:
		t.NameInGame, t.Name = "Centurion", "Centurion"
		t.Cost = Cost{Food: 1800, Gold: 700}
		t.Time, t.Location = 150, Academy

	case Astrology:
		t.NameInGame, t.Name = "Astrology", "Astrology"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 50, Temple
	case Mysticism:
		t.NameInGame, t.Name = "Mysticism", "Mysticism"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 50, Temple
	case Polytheism:
		t.NameInGame, t.Name = "Polytheism", "Polytheism"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 50, Temple
	case Medicine:
		t.NameInGame, t.Name = "Medicine", "Medicine"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 50, Temple
	case Afterlife:
		t.NameInGame, t.Name = "Afterlife", "Afterlife"
		t.Cost = Cost{Gold: 275}
		t.Time, t.Location = 75, Temple
	case Monotheism:
		t.NameInGame, t.Name = "Monotheism", "Monotheism"
		t.Cost = Cost{Gold: 350}
		t.Time, t.Location = 75, Temple
	case Fanaticism:
		t.NameInGame, t.Name = "Fanaticism", "Fanaticism"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 60, Temple
	case Jihad:
		t.NameInGame, t.Name = "Zealotry", "Jihad"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 60, Temple
	case Sacrifice:
		t.NameInGame, t.Name = "Sacrifice", "Martyrdom"
		t.Cost = Cost{Gold: 600}
		t.Time, t.Location = 100, Temple

	// from here are internal techs, not shown in the game,
	// they will be automatically researched when a building is built:

	case EnableGranaryStoragePitBarracksDock:
		t.Effects = append(t.Effects, EnableGranaryStoragePitBarracksDockEffect)
	case EnableMarket:
		t.RequiredTechs, t.MinRequiredTechs = map[TechID]bool{
			GranaryTech: true,
			ToolAge:     true,
		}, 2
		t.Effects = append(t.Effects, EnableMarketEffect)

	default:
		panic(fmt.Errorf("NewTechnology: %v: %w", id, ErrTechIDNotFound))
	}
	return t
}

func EnableGranaryStoragePitBarracksDockEffect(e *EmpireDeveloping) {
	e.UnitStats[Granary] = NewUnit(Granary)
	e.UnitStats[StoragePit] = NewUnit(StoragePit)
	e.UnitStats[Barracks] = NewUnit(Barracks)
	e.UnitStats[Dock] = NewUnit(Dock)
}

func EnableMarketEffect(e *EmpireDeveloping) {
	e.UnitStats[Market] = NewUnit(Market)
}
