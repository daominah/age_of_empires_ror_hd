package aoego

import (
	"reflect"
	"runtime"
	"strings"
)

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

// effects upgrade units:

func BroadSwordEffect65(e *EmpireDeveloping) {
	e.UnitStats[Swordsman].NameInGame = "Broad Swordsman"
}

func LongSwordEffect66(e *EmpireDeveloping) {
	e.UnitStats[Swordsman].NameInGame = "Long Swordsman"
}

func LegionEffect123(e *EmpireDeveloping) {
	e.UnitStats[Swordsman].NameInGame = "Legion"
}

func CompositeBowEffect57(e *EmpireDeveloping) {
	e.UnitStats[ImprovedBowman].NameInGame = "Composite Bowman"
}

func HeavyHorseArcherEffect124(e *EmpireDeveloping) {
	e.UnitStats[HorseArcher].NameInGame = "Heavy Horse Archer"
}

func ScytheChariotEffect204(e *EmpireDeveloping) {
	e.UnitStats[Chariot].NameInGame = "Scythe Chariot"
}

func HeavyCalvaryEffect71(e *EmpireDeveloping) {
	e.UnitStats[Cavalry].NameInGame = "Heavy Cavalry"
}

func CataphractEffect126(e *EmpireDeveloping) {
	e.UnitStats[Cavalry].NameInGame = "Cataphract"
}

func ArmoredElephantEffect203(e *EmpireDeveloping) {
	e.UnitStats[Elephant].NameInGame = "Armored Elephant"
}

func CatapultEffect54(e *EmpireDeveloping) {
	e.UnitStats[StoneThrower].NameInGame = "Catapult"
}

func MassiveCatapultEffect122(e *EmpireDeveloping) {
	e.UnitStats[StoneThrower].NameInGame = "Massive Catapult"
}

func HelepolisEffect125(e *EmpireDeveloping) {
	e.UnitStats[Ballista].NameInGame = "Helepolis"
}

func PhalanxEffect73(e *EmpireDeveloping) {
	e.UnitStats[Hoplite].NameInGame = "Phalanx"
}

func CenturionEffect25(e *EmpireDeveloping) {
	e.UnitStats[Hoplite].NameInGame = "Centurion"
}

func ZealotryEffect23(e *EmpireDeveloping) {
	e.UnitStats[Villager].NameInGame = "Jihad"
}

// other effects:

func ToolAgeEffect95(_ *EmpireDeveloping) {
	// Scout +2 sight
}

func BronzeAgeEffect96(_ *EmpireDeveloping) {
	// Scout +2 sight
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
