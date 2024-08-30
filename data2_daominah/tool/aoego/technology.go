package aoego

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type Technology struct {
	ID               TechID
	Name             string // name without spaces, e.g. "Catapult_Tower", "Heavy_Horse_Archer", ...
	NameInGame       string // name shown in the game, e.g. "Ballista Tower", "Heavy Horse Archer", ...
	Cost             Cost
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

func (t Technology) GetLocation() UnitID {
	return t.Location
}

func (t Technology) GetEffectsName() string {
	var effectNames []string
	for _, f := range t.Effects {
		effectNames = append(effectNames, GetFunctionName(f))
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
	case StoneAge, ToolAge, BronzeAge, IronAge:
		return fmt.Sprintf("C%v", id)
	default:
		return fmt.Sprintf("R%v", id)
	}
}

// TechID enum
const (
	StoneAge  TechID = 100 // always pre-researched
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
	Jihad      TechID = 23
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
	EnableLightTransport TechID = 1
	EnableScoutShip      TechID = 3
	EnableChariotArcher  TechID = 59
	EnableChariot        TechID = 68
	EnableCavalry        TechID = 69
	EnableCamel          TechID = 124
	EnableHorseArcher    TechID = 60
	EnableElephantArcher TechID = 61
	EnableWarElephant    TechID = 70
	EnableBallista       TechID = 58
)

// the following vars are used as constants, not changed during runtime
var (
	AllStoneTechs = []TechID{
		GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt,
	}
	AllToolTechs = []TechID{
		EnableMarket, EnableArcheryRange, EnableStable,
		EnableSlinger,
		MarketBuilt, ArcheryRangeBuilt, StableBuilt,
		WatchTower, SmallWall,
		Toolworking, LeatherArmorInfantry, LeatherArmorArchers, LeatherArmorCavalry,
		Woodworking, StoneMining, GoldMining, Domestication,
	}
	AllBronzeTechs = []TechID{
		EnableGovernmentCenter, EnableTemple, EnableSiegeWorkshop, EnableAcademy,
		EnableChariot, EnableCavalry, EnableCamel, EnableChariotArcher,
		GovernmentCenterBuilt, TempleBuilt, SiegeWorkshopBuilt, AcademyBuilt,
		SentryTower, MediumWall, Axe,
		Metalworking, ScaleArmorInfantry, ScaleArmorArchers, ScaleArmorCavalry, BronzeShield,
		ShortSword, Broadsword,
		Wheel, Artisanship, Plow,
		ImprovedBow, CompositeBow,
		Nobility, Writing, Architecture, Logistics,
		Astrology, Mysticism, Polytheism,
	}
	AllIronTechs = []TechID{
		EnableHorseArcher, EnableElephantArcher, EnableWarElephant, EnableBallista,
		GuardTower, BallistaTower, FortifiedWall,
		Metallurgy, ChainMailInfantry, ChainMailArchers, ChainMailCavalry, IronShield, TowerShield,
		LongSword, Legion,
		Craftsmanship, Siegecraft, Coinage, Irrigation,
		HeavyHorseArcher,
		HeavyCalvary, Cataphract, ArmoredElephant,
		Aristocracy, Ballistics, Alchemy, Engineering,
		Medicine, Monotheism, Fanaticism, Jihad, Sacrifice,
		Catapult, MassiveCatapult, Helepolis,
		Phalanx, Centurion,
	}

	UnitEnabledByTechs = map[UnitID]TechID{
		Swordsman:      ShortSword,
		Chariot:        EnableChariot,
		Cavalry:        EnableCavalry,
		Elephant:       EnableWarElephant,
		Camel:          EnableCamel,
		ImprovedBowman: ImprovedBow,
		ChariotArcher:  EnableChariotArcher,
		HorseArcher:    EnableHorseArcher,
		ElephantArcher: EnableElephantArcher,
		Ballista:       EnableBallista,

		Market:           EnableMarket,
		ArcheryRange:     EnableArcheryRange,
		Stable:           EnableStable,
		GovernmentCenter: EnableGovernmentCenter,
		Temple:           EnableTemple,
		SiegeWorkshop:    EnableSiegeWorkshop,
		Academy:          EnableAcademy,
		Wonder:           IronAge,
	}
)

var (
	// AllTechs is initialized in func init then will be used as a constant
	AllTechs = make(map[TechID]Technology)
	// AllAutoTechs are techs not shown in the game, zero cost, zero time,
	// they will be automatically researched when a building is built
	// or another tech is researched (e.g. Iron Age researched will make
	// EnableHorseArcher, EnableBallista, ... automatically researched)
	AllAutoTechs = make(map[TechID]Technology)
)

func init() {
	for _, list := range [][]TechID{
		{StoneAge, ToolAge, BronzeAge, IronAge},
		AllStoneTechs,
		AllToolTechs,
		AllBronzeTechs,
		AllIronTechs,
	} {
		for _, id := range list {
			t, err := NewTechnology(id)
			if err != nil {
				panic(fmt.Errorf("error NewTechnology(%v): %w", id, err))
			}
			AllTechs[id] = *t
			if t.Cost.IsZero() {
				AllAutoTechs[id] = *t
			}
		}
	}
	println("len(AllTechs):", len(AllTechs))         // Output: len(AllTechs): 96
	println("len(AllAutoTechs):", len(AllAutoTechs)) // Output: len(AllAutoTechs): 28
}

// EffectFunc can modify Unit attributes, enable or disable Technology
// or modify player resources. TODO: real type EffectFunc func
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
		t.RequiredTechs = []TechID{StoneAge, GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt}
		t.MinRequiredTechs = 3
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

	case WatchTower:
		t.NameInGame, t.Name = "Watch Tower", "Watch_Tower"
		t.Cost = Cost{Food: 50}
		t.Time, t.Location = 10, Granary
		t.RequiredTechs = []TechID{ToolAge}
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
		t.RequiredTechs = []TechID{IronAge, Nobility}
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
		t.RequiredTechs = []TechID{IronAge, IronShield}

	case ImprovedBow:
		t.NameInGame, t.Name = "Improved Bow", "Improved_Bow"
		t.Cost = Cost{Wood: 80, Food: 140}
		t.Time, t.Location = 60, ArcheryRange
		t.RequiredTechs = []TechID{BronzeAge}
	case CompositeBow:
		t.NameInGame, t.Name = "Composite Bow", "Composit_bow"
		t.Cost = Cost{Wood: 100, Food: 180}
		t.Time, t.Location = 100, ArcheryRange
		t.RequiredTechs = []TechID{BronzeAge, ImprovedBow}
	case HeavyHorseArcher:
		t.NameInGame, t.Name = "Heavy Horse Archer", "Heavy_Horse_Archer"
		t.Cost = Cost{Food: 1750, Gold: 800}
		t.Time, t.Location = 150, ArcheryRange
		t.RequiredTechs = []TechID{IronAge, ChainMailArchers}

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
		t.RequiredTechs = []TechID{IronAge, Craftsmanship}

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
	case Jihad:
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
		t.NameInGame, t.Name = "GranaryBuilt", "GranaryBuilt"
	case StoragePitBuilt:
		t.NameInGame, t.Name = "StoragePitBuilt", "StoragePitBuilt"
	case BarracksBuilt:
		t.NameInGame, t.Name = "BarracksBuilt", "BarracksBuilt"
		t.Effects = []EffectFunc{BarracksBuiltEffect}
	case DockBuilt:
		t.NameInGame, t.Name = "DockBuilt", "DockBuilt"

	case EnableMarket:
		t.Name = "EnableMarket"
		t.RequiredTechs = []TechID{ToolAge, GranaryBuilt}
		t.Effects = []EffectFunc{EnableMarketEffect77}
	case EnableArcheryRange:
		t.Name = "EnableArcheryRange"
		t.RequiredTechs = []TechID{ToolAge, BarracksBuilt}
		t.Effects = []EffectFunc{EnableArcheryRangeEffect75}
	case EnableStable:
		t.Name = "EnableStable"
		t.RequiredTechs = []TechID{ToolAge, BarracksBuilt}
		t.Effects = []EffectFunc{EnableStableEffect79}

	case MarketBuilt:
		t.NameInGame, t.Name = "MarketBuilt", "MarketBuilt"
	case ArcheryRangeBuilt:
		t.NameInGame, t.Name = "ArcheryRangeBuilt", "ArcheryRangeBuilt"
	case StableBuilt:
		t.NameInGame, t.Name = "StableBuilt", "StableBuilt"

	case EnableGovernmentCenter:
		t.Name = "EnableGovernmentCenter"
		t.RequiredTechs = []TechID{BronzeAge, MarketBuilt}
		t.Effects = []EffectFunc{EnableGovernmentCenterEffect76}
	case EnableTemple:
		t.Name = "EnableTemple"
		t.RequiredTechs = []TechID{BronzeAge, MarketBuilt}
		t.Effects = []EffectFunc{EnableTempleEffect80}
	case EnableSiegeWorkshop:
		t.Name = "EnableSiegeWorkshop"
		t.RequiredTechs = []TechID{BronzeAge, ArcheryRangeBuilt}
		t.Effects = []EffectFunc{EnableSiegeWorkshopEffect78}
	case EnableAcademy:
		t.Name = "EnableAcademy"
		t.RequiredTechs = []TechID{BronzeAge, StableBuilt}
		t.Effects = []EffectFunc{EnableAcademyEffect74}

	case GovernmentCenterBuilt:
		t.Name = "GovernmentCenterBuilt"
	case TempleBuilt:
		t.Name = "TempleBuilt"
	case SiegeWorkshopBuilt:
		t.Name = "SiegeWorkshopBuilt"
	case AcademyBuilt:
		t.Name = "AcademyBuilt"

	case EnableSlinger:
		t.NameInGame, t.Name = "EnableSlinger", "EnableSlinger"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableSlingerEffect201}
	case EnableLightTransport:
		t.Name = "EnableLightTransport"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableLightTransportEffect1}
	case EnableScoutShip:
		t.Name = "EnableScoutShip"
		t.RequiredTechs = []TechID{ToolAge}
		t.Effects = []EffectFunc{EnableScoutShipEffect3}

	case EnableChariotArcher:
		t.Name = "EnableChariotArcher"
		t.RequiredTechs = []TechID{BronzeAge, Wheel}
		t.Effects = []EffectFunc{EnableChariotArcherEffect59}
	case EnableChariot:
		t.Name = "EnableChariot"
		t.RequiredTechs = []TechID{BronzeAge, Wheel}
		t.Effects = []EffectFunc{EnableChariotEffect68}
	case EnableCavalry:
		t.Name = "EnableCavalry"
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{EnableCavalryEffect69}
	case EnableCamel:
		t.Name = "EnableCamel"
		t.RequiredTechs = []TechID{BronzeAge}
		t.Effects = []EffectFunc{EnableCamelEffect209}

	case EnableHorseArcher:
		t.Name = "EnableHorseArcher"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableHorseArcherEffect60}
	case EnableElephantArcher:
		t.Name = "EnableElephantArcher"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableElephantArcherEffect61}
	case EnableWarElephant:
		t.Name = "EnableWarElephant"
		t.RequiredTechs = []TechID{IronAge}
		t.Effects = []EffectFunc{EnableWarElephantEffect70}
	case EnableBallista:
		t.Name = "EnableBallista"
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

// effects after normal tech is researched:

func ToolAgeEffect95(e *EmpireDeveloping) {
	if e.Techs[BarracksBuilt] {
		e.Techs[EnableStable] = true
		EnableStableEffect79(e)
	}
}

// effects after building is built:

func BarracksBuiltEffect(e *EmpireDeveloping) {
	if e.Techs[ToolAge] {
		e.Techs[EnableStable] = true
		EnableStableEffect79(e)
	}
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
	// e.EnabledUnits[LightTransport] = true
}

func EnableScoutShipEffect3(e *EmpireDeveloping) {
	// e.EnabledUnits[ScoutShip] = true
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

func GetFunctionName(i interface{}) string {
	// fullName example: github.com/daominah/age_of_empires_ror_hd/data2_daominah/tool/aoego.EnableStableEffect79
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 {
		return fullName
	}
	return fullName[lastDot+1:]
}
