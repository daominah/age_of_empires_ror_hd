package aoego

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Technology can be a Tool Age, Leather Armor Infantry, EnableHorseArcher, ...
// Should be initialized with func NewTechnology for default values.
type Technology struct {
	ID               TechID
	Name             string  // name without spaces, e.g. "Catapult_Tower", "Heavy_Horse_Archer", ...
	NameInGame       string  // name shown in the game, e.g. "Ballista Tower", "Heavy Horse Archer", ...
	Cost             Cost    // no civ bonus affects it technologies cost
	Time             float64 // research time in seconds
	Location         UnitID  // building that researches this technology
	RequiredTechs    []TechID
	MinRequiredTechs int // used e.g. Bronze Age needs 2 building from Tool Age
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

func (t Technology) GetFullName() string {
	return fmt.Sprintf("%v(%v)", t.NameInGame, t.ID.ActionID())
}

func (t Technology) GetLocation() UnitID {
	return t.Location
}

func (t Technology) GetCost() *Cost {
	clone := t.Cost
	return &clone
}

func (t Technology) GetEffectsName() string {
	var effectNames []string
	for _, f := range t.Effects {
		effectNames = append(effectNames, getFunctionName(f))
	}
	return strings.Join(effectNames, ", ")
}

// TechID is enum
type TechID int

func (id TechID) IntID() int { return int(id) }
func (id TechID) GetNameInGame() string {
	if t, ok := AllTechs[id]; ok {
		return t.NameInGame
	}
	return fmt.Sprintf("TechID%v", id) // should not happen
}

func (id TechID) ActionID() string {
	switch id {
	case ToolAge, BronzeAge, IronAge:
		return fmt.Sprintf("C%v", id)
	default:
		return fmt.Sprintf("R%v", id)
	}
}

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

	SmallWall     TechID = 11
	MediumWall    TechID = 13
	FortifiedWall TechID = 14

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
	Plow          TechID = 31
	Irrigation    TechID = 80

	Nobility     TechID = 34
	Writing      TechID = 114
	Architecture TechID = 112
	Logistics    TechID = 121
	Aristocracy  TechID = 113
	Ballistics   TechID = 106
	Alchemy      TechID = 37
	Engineering  TechID = 35

	FishingShip     TechID = 4
	MerchantShip    TechID = 6
	WarGalley       TechID = 5 // Bronze Medium War Ship
	Trireme         TechID = 7 // Iron Big War Ship
	CatapultTrireme TechID = 9
	Juggernaught    TechID = 25
	HeavyTransport  TechID = 8

	Axe        TechID = 63
	ShortSword TechID = 64
	Broadsword TechID = 65
	LongSword  TechID = 66
	Legion     TechID = 77

	ImprovedBow      TechID = 56
	CompositeBow     TechID = 57
	HeavyHorseArcher TechID = 38

	ScytheChariot   TechID = 126
	HeavyCalvary    TechID = 71
	Cataphract      TechID = 78
	ArmoredElephant TechID = 125

	Catapult        TechID = 54
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
	Zealotry   TechID = 23
	Sacrifice  TechID = 120

	NullTech TechID = -1

	GranaryBuilt    TechID = 10
	StoragePitBuilt TechID = 39
	BarracksBuilt   TechID = 62
	DockBuilt       TechID = 0

	MarketBuilt       TechID = 26
	ArcheryRangeBuilt TechID = 55
	StableBuilt       TechID = 67

	GovernmentCenterBuilt TechID = 33
	TempleBuilt           TechID = 17
	SiegeWorkshopBuilt    TechID = 53
	AcademyBuilt          TechID = 72

	EnableMarket           TechID = 94
	EnableArcheryRange     TechID = 95
	EnableStable           TechID = 97
	EnableGovernmentCenter TechID = 93
	EnableTemple           TechID = 98
	EnableSiegeWorkshop    TechID = 96
	EnableAcademy          TechID = 92

	EnableSlinger        TechID = 123
	EnableTransportBoat  TechID = 1
	EnableWarBoat        TechID = 3
	EnableChariotArcher  TechID = 59
	EnableChariot        TechID = 68
	EnableCavalry        TechID = 69
	EnableCamel          TechID = 124
	EnableHorseArcher    TechID = 60
	EnableElephantArcher TechID = 61
	EnableWarElephant    TechID = 70
	EnableBallista       TechID = 58
	EnableFireBoat       TechID = 118
)

// EffectFunc can modify Unit attributes, enable or disable Technology
// or modify player resources.
type EffectFunc func(empire *EmpireDeveloping)

func NewTechnology(id TechID) (*Technology, error) {
	t := &Technology{ID: id}
	switch id {
	case StoneAge:
		t.NameInGame, t.Name = "Stone Age", "Stone_Age"
	case ToolAge:
		t.NameInGame, t.Name = "Tool Age", "Tool_Age"
		t.Cost = Cost{Food: 500}
		t.Time, t.Location = 120, TownCenter
		t.RequiredTechs = []TechID{GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt}
		t.MinRequiredTechs = 2
		t.Effects = []EffectFunc{ToolAgeEffect95}
	case BronzeAge:
		t.NameInGame, t.Name = "Bronze Age", "Bronze_Age"
		t.Cost = Cost{Food: 800}
		t.Time, t.Location = 140, TownCenter
		t.RequiredTechs = []TechID{ToolAge, MarketBuilt, ArcheryRangeBuilt, StableBuilt}
		t.MinRequiredTechs = 3
	case IronAge:
		t.NameInGame, t.Name = "Iron Age", "Iron_Age"
		t.Cost = Cost{Food: 1000, Gold: 800}
		t.Time, t.Location = 160, TownCenter
		t.RequiredTechs = []TechID{BronzeAge, GovernmentCenterBuilt, TempleBuilt, SiegeWorkshopBuilt, AcademyBuilt}
		t.MinRequiredTechs = 3
		t.Effects = []EffectFunc{IronAgeEffect97}

	case WatchTower:
		t.NameInGame, t.Name = "Watch Tower", "Watch_Tower"
		t.Cost = Cost{Food: 50}
		t.Time, t.Location = 10, Granary
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableTowerEffect12}
		// if not set MinRequiredTechs, it will be set to require all at the end of this function
	case SentryTower:
		t.NameInGame, t.Name = "Sentry Tower", "Sentry_Tower"
		t.Cost = Cost{Food: 120, Stone: 50}
		t.Time, t.Location = 30, Granary
		t.RequiredTechs = []TechID{BronzeAge, WatchTower}
	case GuardTower:
		t.NameInGame, t.Name = "Guard Tower", "Guard_Tower"
		t.Cost = Cost{Food: 300, Stone: 100}
		t.Time, t.Location = 75, Granary
		t.RequiredTechs = []TechID{IronAge, SentryTower}
	case BallistaTower:
		t.NameInGame, t.Name = "Ballista Tower", "Catapult_Tower"
		t.Cost = Cost{Food: 1800, Stone: 750}
		t.Time, t.Location = 150, Granary
		t.RequiredTechs = []TechID{IronAge, GuardTower, Ballistics}
	case SmallWall:
		t.NameInGame, t.Name = "Small Wall", "Small_Wall"
		t.Cost = Cost{Food: 50}
		t.Time, t.Location = 10, Granary
		t.RequiredTechs = []TechID{ToolAge}
	case MediumWall:
		t.NameInGame, t.Name = "Medium Wall", "Medium_Wall"
		t.Cost = Cost{Food: 180, Stone: 100}
		t.Time, t.Location = 60, Granary
		t.RequiredTechs = []TechID{BronzeAge, SmallWall}
	case FortifiedWall:
		t.NameInGame, t.Name = "Fortified Wall", "Fortified_Wall"
		t.Cost = Cost{Food: 300, Stone: 175}
		t.Time, t.Location = 75, Granary
		t.RequiredTechs = []TechID{IronAge, MediumWall}

	case FishingShip:
		t.NameInGame, t.Name = "Fishing Ship", "Fishing_Ship"
		t.Cost = Cost{Wood: 100, Food: 50}
		t.Time, t.Location = 30, Dock
		t.RequiredTechs = []TechID{BronzeAge}
	case MerchantShip:
		t.NameInGame, t.Name = "Merchant Ship", "Merchant_Ship"
		t.Cost = Cost{Wood: 75, Food: 200}
		t.Time, t.Location = 60, Dock
		t.RequiredTechs = []TechID{BronzeAge}
	case WarGalley:
		t.NameInGame, t.Name = "War Galley", "Medium_War_Ship"
		t.Cost = Cost{Wood: 75, Food: 150}
		t.Time, t.Location = 75, Dock
		t.RequiredTechs = []TechID{BronzeAge}
	case Trireme:
		t.NameInGame, t.Name = "Trireme", "Trireme"
		t.Cost = Cost{Wood: 100, Food: 250}
		t.Time, t.Location = 100, Dock
		t.RequiredTechs = []TechID{IronAge, WarGalley}
	case CatapultTrireme:
		t.NameInGame, t.Name = "Catapult Trireme", "Trireme_with_Catapult"
		t.Cost = Cost{Wood: 150, Food: 300}
		t.Time, t.Location = 100, Dock
		t.RequiredTechs = []TechID{IronAge, Trireme}
		t.Effects = []EffectFunc{CatapultTriremeEffect9}
	case Juggernaught:
		t.NameInGame, t.Name = "Juggernaught", "Juggernaught"
		t.Cost = Cost{Wood: 900, Food: 2000}
		t.Time, t.Location = 180, Dock
		t.RequiredTechs = []TechID{IronAge, CatapultTrireme, Engineering}
	case HeavyTransport:
		t.NameInGame, t.Name = "Heavy Transport", "Heavy_Transport"
		t.Cost = Cost{Wood: 125, Food: 150}
		t.Time, t.Location = 75, Dock
		t.RequiredTechs = []TechID{IronAge}

	case Wheel:
		t.NameInGame, t.Name = "Wheel", "Wheel"
		t.Cost = Cost{Wood: 75, Food: 175}
		t.Time, t.Location = 75, Market
		t.RequiredTechs = []TechID{BronzeAge}
	case Woodworking:
		t.NameInGame, t.Name = "Woodworking", "Wood_Working"
		t.Cost = Cost{Wood: 75, Food: 120}
		t.Time, t.Location = 60, Market
		t.RequiredTechs = []TechID{ToolAge}
	case Artisanship:
		t.NameInGame, t.Name = "Artisanship", "Artisanship"
		t.Cost = Cost{Wood: 150, Food: 170}
		t.Time, t.Location = 80, Market
		t.RequiredTechs = []TechID{BronzeAge, Woodworking}
	case Craftsmanship:
		t.NameInGame, t.Name = "Craftsmanship", "Craftmanship"
		t.Cost = Cost{Wood: 200, Food: 240}
		t.Time, t.Location = 100, Market
		t.RequiredTechs = []TechID{IronAge, Artisanship}
	case StoneMining:
		t.NameInGame, t.Name = "Stone Mining", "Stone_Mining"
		t.Cost = Cost{Food: 100, Stone: 50}
		t.Time, t.Location = 30, Market
		t.RequiredTechs = []TechID{ToolAge}
	case Siegecraft:
		t.NameInGame, t.Name = "Siegecraft", "Siegecraft"
		t.Cost = Cost{Food: 190, Stone: 100}
		t.Time, t.Location = 60, Market
		t.RequiredTechs = []TechID{IronAge, StoneMining}
	case GoldMining:
		t.NameInGame, t.Name = "Gold Mining", "Gold_Mining"
		t.Cost = Cost{Wood: 100, Food: 120}
		t.Time, t.Location = 50, Market
		t.RequiredTechs = []TechID{ToolAge}
	case Coinage:
		t.NameInGame, t.Name = "Coinage", "Coinage"
		t.Cost = Cost{Food: 200, Gold: 100}
		t.Time, t.Location = 60, Market
		t.RequiredTechs = []TechID{IronAge, GoldMining}
	case Domestication:
		t.NameInGame, t.Name = "Domestication", "Domestication"
		t.Cost = Cost{Wood: 50, Food: 200}
		t.Time, t.Location = 40, Market
		t.RequiredTechs = []TechID{ToolAge}
	case Plow:
		t.NameInGame, t.Name = "Plow", "Plow"
		t.Cost = Cost{Wood: 75, Food: 250}
		t.Time, t.Location = 75, Market
		t.RequiredTechs = []TechID{BronzeAge, Domestication}
	case Irrigation:
		t.NameInGame, t.Name = "Irrigation", "Irrigation"
		t.Cost = Cost{Wood: 100, Food: 300}
		t.Time, t.Location = 100, Market
		t.RequiredTechs = []TechID{IronAge, Plow}

	case Nobility:
		t.NameInGame, t.Name = "Nobility", "Nobility"
		t.Cost = Cost{Food: 175, Gold: 120}
		t.Time, t.Location = 70, GovernmentCenter
		t.RequiredTechs = []TechID{BronzeAge}
	case Writing:
		t.NameInGame, t.Name = "Writing", "Writing"
		t.Cost = Cost{Food: 200, Gold: 75}
		t.Time, t.Location = 60, GovernmentCenter
		t.RequiredTechs = []TechID{BronzeAge}
	case Architecture:
		t.NameInGame, t.Name = "Architecture", "Architecture"
		t.Cost = Cost{Wood: 175, Food: 150}
		t.Time, t.Location = 50, GovernmentCenter
		t.RequiredTechs = []TechID{BronzeAge}
	case Logistics:
		t.NameInGame, t.Name = "Logistics", "Logistics"
		t.Cost = Cost{Food: 180, Gold: 100}
		t.Time, t.Location = 60, GovernmentCenter
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{LogisticsEffect200}
	case Aristocracy:
		t.NameInGame, t.Name = "Aristocracy", "Aristocracy"
		t.Cost = Cost{Food: 175, Gold: 150}
		t.Time, t.Location = 60, GovernmentCenter
		t.RequiredTechs = []TechID{IronAge}
	case Ballistics:
		t.NameInGame, t.Name = "Ballistics", "Ballistics"
		t.Cost = Cost{Food: 200, Gold: 50}
		t.Time, t.Location = 60, GovernmentCenter
		t.RequiredTechs = []TechID{IronAge}
	case Alchemy:
		t.NameInGame, t.Name = "Alchemy", "Alchemy"
		t.Cost = Cost{Food: 250, Gold: 200}
		t.Time, t.Location = 100, GovernmentCenter
		t.RequiredTechs = []TechID{IronAge}
	case Engineering:
		t.NameInGame, t.Name = "Engineering", "Engineering"
		t.Cost = Cost{Wood: 100, Food: 200}
		t.Time, t.Location = 70, GovernmentCenter
		t.RequiredTechs = []TechID{IronAge}

	case Toolworking:
		t.NameInGame, t.Name = "Toolworking", "Toolworking"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 40, StoragePit
		t.RequiredTechs = []TechID{ToolAge}
	case Metalworking:
		t.NameInGame, t.Name = "Metalworking", "Metal_Working"
		t.Cost = Cost{Food: 200, Gold: 120}
		t.Time, t.Location = 75, StoragePit
		t.RequiredTechs = []TechID{BronzeAge, Toolworking}
	case Metallurgy:
		t.NameInGame, t.Name = "Metallurgy", "Metallurgy"
		t.Cost = Cost{Food: 300, Gold: 180}
		t.Time, t.Location = 100, StoragePit
		t.RequiredTechs = []TechID{IronAge, Metalworking}

	case LeatherArmorInfantry:
		t.NameInGame, t.Name = "Leather Armor Infantry", "Leather_Armor_-_Soldiers"
		t.Cost = Cost{Food: 75}
		t.Time, t.Location = 30, StoragePit
		t.RequiredTechs = []TechID{ToolAge}
	case ScaleArmorInfantry:
		t.NameInGame, t.Name = "Scale Armor Infantry", "Scale_Armor_-_Soldiers"
		t.Cost = Cost{Food: 100, Gold: 50}
		t.Time, t.Location = 60, StoragePit
		t.RequiredTechs = []TechID{BronzeAge, LeatherArmorInfantry}
	case ChainMailInfantry:
		t.NameInGame, t.Name = "Chain Mail Infantry", "Chain_Mail_-_Soldiers"
		t.Cost = Cost{Food: 125, Gold: 100}
		t.Time, t.Location = 75, StoragePit
		t.RequiredTechs = []TechID{IronAge, ScaleArmorInfantry}

	case LeatherArmorArchers:
		t.NameInGame, t.Name = "Leather Armor Archers", "Leather_Armor_-_Archer"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 30, StoragePit
		t.RequiredTechs = []TechID{ToolAge}
	case ScaleArmorArchers:
		t.NameInGame, t.Name = "Scale Armor Archers", "Scale_Armor_-_Archers"
		t.Cost = Cost{Food: 125, Gold: 50}
		t.Time, t.Location = 60, StoragePit
		t.RequiredTechs = []TechID{BronzeAge, LeatherArmorArchers}
	case ChainMailArchers:
		t.NameInGame, t.Name = "Chain Mail Archers", "Chain_Mail_-_Archers"
		t.Cost = Cost{Food: 150, Gold: 100}
		t.Time, t.Location = 75, StoragePit
		t.RequiredTechs = []TechID{IronAge, ScaleArmorArchers}

	case LeatherArmorCavalry:
		t.NameInGame, t.Name = "Leather Armor Cavalry", "Leather_Armor_Mounted"
		t.Cost = Cost{Food: 125}
		t.Time, t.Location = 30, StoragePit
		t.RequiredTechs = []TechID{ToolAge}
	case ScaleArmorCavalry:
		t.NameInGame, t.Name = "Scale Armor Cavalry", "Scale_Armor_-_Cavalry"
		t.Cost = Cost{Food: 150, Gold: 50}
		t.Time, t.Location = 60, StoragePit
		t.RequiredTechs = []TechID{BronzeAge, LeatherArmorCavalry}
	case ChainMailCavalry:
		t.NameInGame, t.Name = "Chain Mail Cavalry", "Chain_Mail_-_Cavalry"
		t.Cost = Cost{Food: 175, Gold: 100}
		t.Time, t.Location = 75, StoragePit
		t.RequiredTechs = []TechID{IronAge, ScaleArmorCavalry}

	case BronzeShield:
		t.NameInGame, t.Name = "Bronze Shield", "Bronze_Shield"
		t.Cost = Cost{Food: 150, Gold: 180}
		t.Time, t.Location = 50, StoragePit
		t.RequiredTechs = []TechID{BronzeAge}
	case IronShield:
		t.NameInGame, t.Name = "Iron Shield", "Iron_Shield"
		t.Cost = Cost{Food: 200, Gold: 320}
		t.Time, t.Location = 75, StoragePit
		t.RequiredTechs = []TechID{IronAge, BronzeShield}
	case TowerShield:
		t.NameInGame, t.Name = "Tower Shield", "Tower_Shield"
		t.Cost = Cost{Food: 250, Gold: 400}
		t.Time, t.Location = 100, StoragePit
		t.RequiredTechs = []TechID{IronAge, IronShield}

	case Axe:
		t.NameInGame, t.Name = "Battle Axe", "Axe"
		t.Cost = Cost{Food: 100}
		t.Time, t.Location = 40, Barracks
		t.RequiredTechs = []TechID{ToolAge}
	case ShortSword:
		t.NameInGame, t.Name = "Short Sword", "Short_Sword"
		t.Cost = Cost{Food: 120, Gold: 50}
		t.Time, t.Location = 50, Barracks
		t.RequiredTechs = []TechID{BronzeAge, Axe}
		t.Effects = []EffectFunc{ShortSwordEffect64}
	case Broadsword:
		t.NameInGame, t.Name = "Broadsword", "Broad_Sword"
		t.Cost = Cost{Food: 140, Gold: 50}
		t.Time, t.Location = 80, Barracks
		t.RequiredTechs = []TechID{BronzeAge, ShortSword}
	case LongSword:
		t.NameInGame, t.Name = "Long Sword", "Long_Sword"
		t.Cost = Cost{Food: 160, Gold: 50}
		t.Time, t.Location = 90, Barracks
		t.RequiredTechs = []TechID{IronAge, Broadsword}
	case Legion:
		t.NameInGame, t.Name = "Legion", "Legion"
		t.Cost = Cost{Food: 1400, Gold: 600}
		t.Time, t.Location = 150, Barracks
		t.RequiredTechs = []TechID{IronAge, LongSword, Fanaticism}

	case ScytheChariot:
		t.NameInGame, t.Name = "Scythe Chariot", "Scythe_Chariot"
		t.Cost = Cost{Wood: 1200, Gold: 800}
		t.Time, t.Location = 150, Stable
		t.RequiredTechs = []TechID{IronAge, EnableChariot, Nobility}
	case HeavyCalvary:
		t.NameInGame, t.Name = "Heavy Calvary", "Heavy_Cavalry"
		t.Cost = Cost{Food: 350, Gold: 125}
		t.Time, t.Location = 90, Stable
		t.RequiredTechs = []TechID{IronAge}
	case Cataphract:
		t.NameInGame, t.Name = "Cataphract", "Cataphracts"
		t.Cost = Cost{Food: 2000, Gold: 850}
		t.Time, t.Location = 150, Stable
		t.RequiredTechs = []TechID{IronAge, HeavyCalvary, Metallurgy}
	case ArmoredElephant:
		t.NameInGame, t.Name = "Armored Elephant", "Armored_Elephant"
		t.Cost = Cost{Food: 1000, Gold: 1200}
		t.Time, t.Location = 150, Stable
		t.RequiredTechs = []TechID{IronAge, EnableWarElephant, IronShield}

	case ImprovedBow:
		t.NameInGame, t.Name = "Improved Bow", "Improved_Bow"
		t.Cost = Cost{Wood: 80, Food: 140}
		t.Time, t.Location = 60, ArcheryRange
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{ImprovedBowEffect56}
	case CompositeBow:
		t.NameInGame, t.Name = "Composite Bow", "Composit_bow"
		t.Cost = Cost{Wood: 100, Food: 180}
		t.Time, t.Location = 100, ArcheryRange
		t.RequiredTechs = []TechID{BronzeAge, ImprovedBow}
	case HeavyHorseArcher:
		t.NameInGame, t.Name = "Heavy Horse Archer", "Heavy_Horse_Archer"
		t.Cost = Cost{Food: 1750, Gold: 800}
		t.Time, t.Location = 150, ArcheryRange
		t.RequiredTechs = []TechID{IronAge, EnableHorseArcher, ChainMailArchers}

	case Catapult:
		t.NameInGame, t.Name = "Catapult", "Heavy_Catapult"
		t.Cost = Cost{Wood: 250, Food: 300}
		t.Time, t.Location = 100, SiegeWorkshop
		t.RequiredTechs = []TechID{IronAge}
	case MassiveCatapult:
		t.NameInGame, t.Name = "Heavy Catapult", "Massive_Catapult"
		t.Cost = Cost{Wood: 900, Food: 1800}
		t.Time, t.Location = 150, SiegeWorkshop
		t.RequiredTechs = []TechID{IronAge, Catapult, Siegecraft}
	case Helepolis:
		t.NameInGame, t.Name = "Helepolis", "Helepolis"
		t.Cost = Cost{Wood: 1000, Food: 1500}
		t.Time, t.Location = 150, SiegeWorkshop
		t.RequiredTechs = []TechID{IronAge, EnableBallista, Craftsmanship}

	case Phalanx:
		t.NameInGame, t.Name = "Phalanx", "Phalanx"
		t.Cost = Cost{Food: 300, Gold: 100}
		t.Time, t.Location = 90, Academy
		t.RequiredTechs = []TechID{IronAge}
	case Centurion:
		t.NameInGame, t.Name = "Centurion", "Centurion"
		t.Cost = Cost{Food: 1800, Gold: 700}
		t.Time, t.Location = 150, Academy
		t.RequiredTechs = []TechID{IronAge, Phalanx, Aristocracy}

	case Astrology:
		t.NameInGame, t.Name = "Astrology", "Astrology"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 50, Temple
		t.RequiredTechs = []TechID{BronzeAge}
	case Mysticism:
		t.NameInGame, t.Name = "Mysticism", "Mysticism"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 50, Temple
		t.RequiredTechs = []TechID{BronzeAge}
	case Polytheism:
		t.NameInGame, t.Name = "Polytheism", "Polytheism"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 50, Temple
		t.RequiredTechs = []TechID{BronzeAge}
	case Medicine:
		t.NameInGame, t.Name = "Medicine", "Medicine"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 50, Temple
		t.RequiredTechs = []TechID{IronAge}
	case Afterlife:
		t.NameInGame, t.Name = "Afterlife", "Afterlife"
		t.Cost = Cost{Gold: 275}
		t.Time, t.Location = 75, Temple
		t.RequiredTechs = []TechID{IronAge}
	case Monotheism:
		t.NameInGame, t.Name = "Monotheism", "Monotheism"
		t.Cost = Cost{Gold: 350}
		t.Time, t.Location = 75, Temple
		t.RequiredTechs = []TechID{IronAge}
	case Fanaticism:
		t.NameInGame, t.Name = "Fanaticism", "Fanaticism"
		t.Cost = Cost{Gold: 150}
		t.Time, t.Location = 60, Temple
		t.RequiredTechs = []TechID{IronAge}
	case Zealotry:
		t.NameInGame, t.Name = "Zealotry", "Jihad"
		t.Cost = Cost{Gold: 120}
		t.Time, t.Location = 60, Temple
		t.RequiredTechs = []TechID{IronAge}
	case Sacrifice:
		t.NameInGame, t.Name = "Sacrifice", "Martyrdom"
		t.Cost = Cost{Gold: 600}
		t.Time, t.Location = 100, Temple
		t.RequiredTechs = []TechID{IronAge}

	// the following are techs not shown in the game, zero cost, zero time,
	// they will be automatically researched when a building is built
	// or another tech is researched (e.g. Iron Age researched will make
	// EnableHorseArcher, EnableBallista, ... automatically researched)

	case GranaryBuilt:
		t.Name = "_GranaryBuilt"
	case StoragePitBuilt:
		t.Name = "_StoragePitBuilt"
	case BarracksBuilt:
		t.Name = "_BarracksBuilt"
		t.Effects = []EffectFunc{BarracksBuiltEffect62}
	case DockBuilt:
		t.Name = "_DockBuilt"
		t.Effects = []EffectFunc{DockBuiltEffect0}
	case MarketBuilt:
		t.Name = "_MarketBuilt"
		t.Effects = []EffectFunc{MarketBuiltEffect26}
	case ArcheryRangeBuilt:
		t.Name = "_ArcheryRangeBuilt"
		t.Effects = []EffectFunc{ArcheryRangeBuiltEffect55}
	case StableBuilt:
		t.Name = "_StableBuilt"
		t.Effects = []EffectFunc{StableBuiltEffect67}
	case GovernmentCenterBuilt:
		t.Name = "_GovernmentCenterBuilt"
		t.Effects = []EffectFunc{GovernmentCenterBuiltEffect33}
	case TempleBuilt:
		t.Name = "_TempleBuilt"
		t.Effects = []EffectFunc{TempleBuiltEffect17}
	case SiegeWorkshopBuilt:
		t.Name = "_SiegeWorkshopBuilt"
		t.Effects = []EffectFunc{SiegeWorkshopBuiltEffect53}
	case AcademyBuilt:
		t.Name = "_AcademyBuilt"
		t.Effects = []EffectFunc{AcademyBuiltEffect72}

	case EnableMarket:
		t.Name = "enableMarket"
		t.RequiredTechs = []TechID{ToolAge, GranaryBuilt}
		t.Effects = []EffectFunc{EnableMarketEffect77}
	case EnableArcheryRange:
		t.Name = "enableArcheryRange"
		t.RequiredTechs = []TechID{ToolAge, BarracksBuilt}
		t.Effects = []EffectFunc{EnableArcheryRangeEffect75}
	case EnableStable:
		t.Name = "enableStable"
		t.RequiredTechs = []TechID{ToolAge, BarracksBuilt}
		t.Effects = []EffectFunc{EnableStableEffect79}

	case EnableGovernmentCenter:
		t.Name = "enableGovernmentCenter"
		t.RequiredTechs = []TechID{BronzeAge, MarketBuilt}
		t.Effects = []EffectFunc{EnableGovernmentCenterEffect76}
	case EnableTemple:
		t.Name = "enableTemple"
		t.RequiredTechs = []TechID{BronzeAge, MarketBuilt}
		t.Effects = []EffectFunc{EnableTempleEffect80}
	case EnableSiegeWorkshop:
		t.Name = "enableSiegeWorkshop"
		t.RequiredTechs = []TechID{BronzeAge, ArcheryRangeBuilt}
		t.Effects = []EffectFunc{EnableSiegeWorkshopEffect78}
	case EnableAcademy:
		t.Name = "enableAcademy"
		t.RequiredTechs = []TechID{BronzeAge, StableBuilt}
		t.Effects = []EffectFunc{EnableAcademyEffect74}

	case EnableSlinger:
		t.Name = "enableSlinger"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableSlingerEffect201}
	case EnableTransportBoat:
		t.Name = "enableTransportBoat"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableLightTransportEffect1}
	case EnableWarBoat:
		t.Name = "enableWarBoat"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableWarBoatEffect3}
	case EnableFireBoat:
		t.Name = "enableFireBoat"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableFireBoatEffect202}

	case EnableChariotArcher:
		t.Name = "enableChariotArcher"
		t.RequiredTechs = []TechID{BronzeAge, Wheel}
		t.Effects = []EffectFunc{EnableChariotArcherEffect59}
	case EnableChariot:
		t.Name = "enableChariot"
		t.RequiredTechs = []TechID{BronzeAge, Wheel}
		t.Effects = []EffectFunc{EnableChariotEffect68}
	case EnableCavalry:
		t.Name = "enableCavalry"
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{EnableCavalryEffect69}
	case EnableCamel:
		t.Name = "enableCamel"
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{EnableCamelEffect209}

	case EnableHorseArcher:
		t.Name = "enableHorseArcher"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableHorseArcherEffect60}
	case EnableElephantArcher:
		t.Name = "enableElephantArcher"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableElephantArcherEffect61}
	case EnableWarElephant:
		t.Name = "enableWarElephant"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableWarElephantEffect70}
	case EnableBallista:
		t.Name = "enableBallista"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableBallistaEffect58}

	default:
		return nil, ErrTechIDNotFound
	}

	if t.MinRequiredTechs == 0 {
		t.MinRequiredTechs = len(t.RequiredTechs)
	}
	if t.NameInGame == "" {
		t.NameInGame = t.Name
	}
	return t, nil
}

// effects after building is built:

func BarracksBuiltEffect62(e *EmpireDeveloping) {
	e.EnabledUnits[Clubman] = true
}

func DockBuiltEffect0(e *EmpireDeveloping) {
	e.EnabledUnits[FishingBoat] = true
	e.EnabledUnits[TradeBoat] = true
}

func MarketBuiltEffect26(e *EmpireDeveloping) {
	e.EnabledUnits[Farm] = true
}

func ArcheryRangeBuiltEffect55(e *EmpireDeveloping) {
	e.EnabledUnits[Bowman] = true
}

func StableBuiltEffect67(e *EmpireDeveloping) {
	e.EnabledUnits[Scout] = true
}

func GovernmentCenterBuiltEffect33(e *EmpireDeveloping) {
}

func TempleBuiltEffect17(e *EmpireDeveloping) {
	e.EnabledUnits[Priest] = true
}

func SiegeWorkshopBuiltEffect53(e *EmpireDeveloping) {
	e.EnabledUnits[StoneThrower] = true
}

func AcademyBuiltEffect72(e *EmpireDeveloping) {
	e.EnabledUnits[Hoplite] = true
}

// effects enable buildings:

func EnableMarketEffect77(e *EmpireDeveloping) {
	e.EnabledUnits[Market] = true
}

func EnableArcheryRangeEffect75(e *EmpireDeveloping) {
	e.EnabledUnits[ArcheryRange] = true
}

func EnableStableEffect79(e *EmpireDeveloping) {
	e.EnabledUnits[Stable] = true
}

func EnableGovernmentCenterEffect76(e *EmpireDeveloping) {
	e.EnabledUnits[GovernmentCenter] = true
}

func EnableTempleEffect80(e *EmpireDeveloping) {
	e.EnabledUnits[Temple] = true
}

func EnableSiegeWorkshopEffect78(e *EmpireDeveloping) {
	e.EnabledUnits[SiegeWorkshop] = true
}

func EnableAcademyEffect74(e *EmpireDeveloping) {
	e.EnabledUnits[Academy] = true
}

// effects enable units:

func EnableSlingerEffect201(e *EmpireDeveloping) {
	e.EnabledUnits[Slinger] = true
}

func EnableLightTransportEffect1(e *EmpireDeveloping) {
	e.EnabledUnits[TransportBoat] = true
}

func EnableWarBoatEffect3(e *EmpireDeveloping) {
	e.EnabledUnits[WarBoat] = true
}

func EnableFireBoatEffect202(e *EmpireDeveloping) {
	e.EnabledUnits[FireBoat] = true
}

func CatapultTriremeEffect9(e *EmpireDeveloping) {
	e.EnabledUnits[CatapultBoat] = true
}

func EnableChariotArcherEffect59(e *EmpireDeveloping) {
	e.EnabledUnits[ChariotArcher] = true
}

func EnableChariotEffect68(e *EmpireDeveloping) {
	e.EnabledUnits[Chariot] = true
}

func EnableCavalryEffect69(e *EmpireDeveloping) {
	e.EnabledUnits[Cavalry] = true
}

func EnableCamelEffect209(e *EmpireDeveloping) {
	e.EnabledUnits[Camel] = true
}

func EnableWarElephantEffect70(e *EmpireDeveloping) {
	e.EnabledUnits[Elephant] = true
}

func EnableHorseArcherEffect60(e *EmpireDeveloping) {
	e.EnabledUnits[HorseArcher] = true
}

func EnableElephantArcherEffect61(e *EmpireDeveloping) {
	e.EnabledUnits[ElephantArcher] = true
}

func EnableBallistaEffect58(e *EmpireDeveloping) {
	e.EnabledUnits[Ballista] = true
}

func EnableTowerEffect12(e *EmpireDeveloping) {
	e.EnabledUnits[Tower] = true
}

func ShortSwordEffect64(e *EmpireDeveloping) {
	e.EnabledUnits[Swordsman] = true
}

func ImprovedBowEffect56(e *EmpireDeveloping) {
	e.EnabledUnits[ImprovedBowman] = true
}

// other effects:

func ToolAgeEffect95(e *EmpireDeveloping) {
}

func BronzeAgeEffect96(e *EmpireDeveloping) {
}

func IronAgeEffect97(e *EmpireDeveloping) {
	e.EnabledUnits[Wonder] = true
}

func LogisticsEffect200(e *EmpireDeveloping) {
	for _, u := range []UnitID{Clubman, Swordsman, Slinger} {
		if _, found := e.UnitStats[u]; found {
			e.UnitStats[u].Population -= 0.5
		}
	}
}

// getFunctionName returns Go function name to debug.
func getFunctionName(i interface{}) string {
	// fullName example: github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoego.EnableStableEffect79
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 {
		return fullName
	}
	return fullName[lastDot+1:]
}
