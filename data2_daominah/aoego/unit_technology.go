package aoego

import "fmt"

// CheckIsBuiltTech returns true if the tech is researched when a building is built.
func CheckIsBuiltTech(techID TechID) bool {
	switch techID {
	case GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt,
		MarketBuilt, ArcheryRangeBuilt, StableBuilt,
		GovernmentCenterBuilt, TempleBuilt, SiegeWorkshopBuilt, AcademyBuilt:
		return true
	default:
		return false
	}
}

// CheckIsAutoTech returns true if the tech is automatically researched when
// their required techs are researched. Auto techs have zero cost, zero time.
// Example: EnableHorseArcher will be automatically researched when Iron Age is
// researched, EnableAcademy will be automatically researched when Bronze Age
// and StableBuilt are researched.
func CheckIsAutoTech(techID TechID) bool {
	switch techID {
	case EnableMarket, EnableArcheryRange, EnableStable,
		EnableGovernmentCenter, EnableTemple, EnableSiegeWorkshop, EnableAcademy,
		EnableSlinger, EnableTransportBoat, EnableWarBoat,
		EnableChariotArcher, EnableChariot, EnableCavalry, EnableCamel,
		EnableHorseArcher, EnableElephantArcher, EnableWarElephant, EnableBallista,
		EnableFireBoat:
		return true
	default:
		return false
	}
}

// UnitEnabledByTechs is used to know how to enable a unit
var UnitEnabledByTechs = map[UnitID]TechID{
	Slinger:   EnableSlinger,
	Swordsman: ShortSword, // needs manually research

	ImprovedBowman: ImprovedBow, // needs  manually researched
	ChariotArcher:  EnableChariotArcher,
	HorseArcher:    EnableHorseArcher,
	ElephantArcher: EnableElephantArcher,

	Chariot:  EnableChariot,
	Cavalry:  EnableCavalry,
	Elephant: EnableWarElephant,
	Camel:    EnableCamel,

	Ballista: EnableBallista,

	TransportBoat: EnableTransportBoat,
	WarBoat:       EnableWarBoat,
	CatapultBoat:  CatapultTrireme, // needs manually research
	FireBoat:      EnableFireBoat,

	Market:           EnableMarket,
	Farm:             MarketBuilt, // needs manually build
	Tower:            WatchTower,  // needs manually research
	Wall:             SmallWall,   // needs manually research
	ArcheryRange:     EnableArcheryRange,
	Stable:           EnableStable,
	GovernmentCenter: EnableGovernmentCenter,
	Temple:           EnableTemple,
	SiegeWorkshop:    EnableSiegeWorkshop,
	Academy:          EnableAcademy,
	Wonder:           IronAge, // needs manually research
}

// GetAge returns StoneAge, ToolAge, BronzeAge or IronAge.
func (id TechID) GetAge() TechID {
	switch id {
	case StoneAge,
		GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt:
		return StoneAge
	case ToolAge,
		EnableMarket, EnableArcheryRange, EnableStable,
		EnableSlinger, EnableTransportBoat, EnableWarBoat,
		MarketBuilt, ArcheryRangeBuilt, StableBuilt,
		WatchTower, SmallWall,
		Toolworking, LeatherArmorInfantry, LeatherArmorArchers, LeatherArmorCavalry,
		Woodworking, StoneMining, GoldMining, Domestication:
		return ToolAge
	case BronzeAge,
		EnableGovernmentCenter, EnableTemple, EnableSiegeWorkshop, EnableAcademy,
		EnableChariotArcher, EnableChariot, EnableCavalry, EnableCamel,
		GovernmentCenterBuilt, TempleBuilt, SiegeWorkshopBuilt, AcademyBuilt,
		SentryTower, MediumWall, Axe,
		Metalworking, ScaleArmorInfantry, ScaleArmorArchers, ScaleArmorCavalry, BronzeShield,
		ShortSword, Broadsword,
		FishingShip, MerchantShip, WarGalley,
		Wheel, Artisanship, Plow,
		ImprovedBow, CompositeBow,
		Nobility, Writing, Architecture, Logistics,
		Astrology, Mysticism, Polytheism:
		return BronzeAge
	case IronAge,
		EnableHorseArcher, EnableElephantArcher, EnableWarElephant, EnableBallista,
		GuardTower, BallistaTower, FortifiedWall,
		Metallurgy, ChainMailInfantry, ChainMailArchers, ChainMailCavalry, IronShield, TowerShield,
		LongSword, Legion,
		HeavyTransport, Trireme, CatapultTrireme, Juggernaught,
		Craftsmanship, Siegecraft, Coinage, Irrigation,
		HeavyHorseArcher,
		HeavyCalvary, Cataphract, ArmoredElephant,
		Aristocracy, Ballistics, Alchemy, Engineering,
		Medicine, Monotheism, Fanaticism, Zealotry, Sacrifice,
		Catapult, MassiveCatapult, Helepolis,
		Phalanx, Centurion:
		return IronAge
	default: // should not happen
		return StoneAge
	}
}

// GetAge returns StoneAge, ToolAge, BronzeAge or IronAge.
func (id UnitID) GetAge() TechID {
	switch id {
	case
		TownCenter, Villager, House,
		Granary,
		StoragePit,
		Barracks, Clubman,
		Dock, FishingBoat, TradeBoat:
		return StoneAge
	case
		Tower, Wall,
		Slinger,
		ArcheryRange, Bowman,
		Stable, Scout,
		Market, Farm,
		TransportBoat, WarBoat:
		return ToolAge
	case
		Swordsman,
		ImprovedBowman, ChariotArcher,
		Chariot, Cavalry, Camel,
		GovernmentCenter,
		Temple, Priest,
		SiegeWorkshop, StoneThrower,
		Academy, Hoplite:
		return BronzeAge
	case
		HorseArcher, ElephantArcher,
		Elephant,
		Ballista,
		CatapultBoat, FireBoat,
		Wonder:
		return IronAge
	default: // should not happen
		return StoneAge
	}
}

// AllBuildings is used as a constant (do not change this var in runtime)
var AllBuildings = map[UnitID]bool{
	TownCenter: true, House: true,
	Granary: true, StoragePit: true, Barracks: true, Dock: true,
	Market: true, ArcheryRange: true, Stable: true,
	Farm: true, Tower: true, Wall: true,
	GovernmentCenter: true, Temple: true, SiegeWorkshop: true, Academy: true,
	Wonder: true,
}

// AllCombatants is a list of all units that is not a building,
// used as a constant (do not change this var in runtime)
var AllCombatants = map[UnitID]bool{
	Villager: true,
	Clubman:  true, Swordsman: true, Slinger: true,
	FishingBoat: true, TradeBoat: true, TransportBoat: true, WarBoat: true, CatapultBoat: true, FireBoat: true,
	Bowman: true, ImprovedBowman: true, ChariotArcher: true, HorseArcher: true, ElephantArcher: true,
	Scout: true, Chariot: true, Cavalry: true, Elephant: true, Camel: true,
	Priest:       true,
	Hoplite:      true,
	StoneThrower: true, Ballista: true,
}

var (
	// AllUnits includes all units in the game (including buildings),
	// initialized in func init then will be used as a constant
	AllUnits = make(map[UnitID]Unit)

	// AllTechs is a convenient way to read Tech info instead of func NewTechnology.
	// This map is initialized in func init then will be used as a constant.
	AllTechs = make(map[TechID]Technology)

	// AllAutoTechs are techs not shown in the game, zero cost, zero time,
	// they will be automatically researched when their required techs are
	// researched (e.g. Iron Age researched will automatically make
	// EnableHorseArcher, EnableBallista, ... automatically researched).
	// This map is initialized in func init then will be used as a constant.
	AllAutoTechs = make(map[TechID]Technology)
)

func init() {
	for _, a := range []map[UnitID]bool{AllBuildings, AllCombatants} {
		for unitID := range a {
			tmp, err := NewUnit(unitID)
			if err != nil {
				panic(fmt.Sprintf("error init AllUnits: %v", err))
			}
			AllUnits[unitID] = *tmp
		}
	}

	for _, id := range []TechID{
		StoneAge,
		GranaryBuilt, StoragePitBuilt, BarracksBuilt, DockBuilt,
		ToolAge,
		EnableMarket, EnableArcheryRange, EnableStable,
		EnableSlinger, EnableTransportBoat, EnableWarBoat,
		MarketBuilt, ArcheryRangeBuilt, StableBuilt,
		WatchTower, SmallWall,
		Toolworking, LeatherArmorInfantry, LeatherArmorArchers, LeatherArmorCavalry,
		Woodworking, StoneMining, GoldMining, Domestication,
		BronzeAge,
		EnableGovernmentCenter, EnableTemple, EnableSiegeWorkshop, EnableAcademy,
		EnableChariotArcher, EnableChariot, EnableCavalry, EnableCamel,
		GovernmentCenterBuilt, TempleBuilt, SiegeWorkshopBuilt, AcademyBuilt,
		SentryTower, MediumWall, Axe,
		Metalworking, ScaleArmorInfantry, ScaleArmorArchers, ScaleArmorCavalry, BronzeShield,
		ShortSword, Broadsword,
		FishingShip, MerchantShip, WarGalley,
		Wheel, Artisanship, Plow,
		ImprovedBow, CompositeBow,
		Nobility, Writing, Architecture, Logistics,
		Astrology, Mysticism, Polytheism,
		IronAge,
		EnableHorseArcher, EnableElephantArcher, EnableWarElephant, EnableBallista,
		GuardTower, BallistaTower, FortifiedWall,
		Metallurgy, ChainMailInfantry, ChainMailArchers, ChainMailCavalry, IronShield, TowerShield,
		LongSword, Legion,
		HeavyTransport, Trireme, CatapultTrireme, Juggernaught,
		Craftsmanship, Siegecraft, Coinage, Irrigation,
		HeavyHorseArcher,
		HeavyCalvary, Cataphract, ArmoredElephant,
		Aristocracy, Ballistics, Alchemy, Engineering,
		Medicine, Monotheism, Fanaticism, Zealotry, Sacrifice,
		Catapult, MassiveCatapult, Helepolis,
		Phalanx, Centurion,
	} {
		t, err := NewTechnology(id)
		if err != nil {
			panic(fmt.Errorf("error NewTechnology(%v): %w", id, err))
		}
		AllTechs[id] = *t
		if CheckIsAutoTech(id) {
			AllAutoTechs[id] = *t
		}
	}
	// println("len(AllTechs):", len(AllTechs))         // Output: len(AllTechs): 105
	// println("len(AllAutoTechs):", len(AllAutoTechs)) // Output: len(AllAutoTechs): 18
}
