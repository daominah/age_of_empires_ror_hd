package logic

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// CostResource is the resource of a technology cost that a TechnologyCostChange applies to,
// the same numbers as the technology cost commands in the .dat.
type CostResource int

// CostResource values.
const (
	AllResources CostResource = -1
	FoodCost     CostResource = 0
	WoodCost     CostResource = 1
	StoneCost    CostResource = 2
	GoldCost     CostResource = 3
)

// TechnologyCostChange is a civilization bonus that makes technologies cheaper to research,
// copied by hand from the technology cost commands (type 101) of the effects in the .dat,
// only the technologies that matter to the ranked lines.
// The ranking shows a technology only when it changes the row unit, or upgrades to it.
type TechnologyCostChange struct {
	CivilizationID CivilizationID
	TechnologyIDs  []TechID
	Resource       CostResource
	Multiplier     float64 // 0 is free
}

// TechnologyCostChanges are the hand copied cost changes.
var TechnologyCostChanges = []TechnologyCostChange{
	{Magyars, []TechID{Forging, IronCasting, BlastFurnace}, AllResources, 0},
	{Turks, []TechID{Chemistry, TechLightCavalry, TechHussar}, AllResources, 0},
	{Italians, []TechID{Chemistry}, AllResources, 0.75},
	// Spanish: the Blacksmith technologies cost no gold
	{Spanish, []TechID{Fletching, BodkinArrow, Bracer, IronCasting, BlastFurnace, ChainMailArmor, PlateMailArmor,
		ChainBardingArmor, PlateBardingArmor, LeatherArcherArmor, RingArcherArmor}, GoldCost, 0},
	{Koreans, []TechID{PaddedArcherArmor, LeatherArcherArmor, RingArcherArmor}, AllResources, 0},
	{Bulgarians, []TechID{Forging, IronCasting, BlastFurnace, Fletching, BodkinArrow, Bracer,
		ScaleMailArmor, ChainMailArmor, PlateMailArmor, ScaleBardingArmor, ChainBardingArmor, PlateBardingArmor,
		PaddedArcherArmor, LeatherArcherArmor, RingArcherArmor}, FoodCost, 0.5},
	{Tatars, []TechID{ParthianTactics, ThumbRing}, AllResources, 0},
	{Burgundians, []TechID{Husbandry, Bloodlines, TechLightCavalry, TechHussar}, AllResources, 0.5},
	{Poles, []TechID{Bloodlines, TechLightCavalry, TechWingedHussar}, FoodCost, 0.5},
	{Khitans, []TechID{TechHeavyCavalryArcher}, AllResources, 0.5},
	{Shu, []TechID{Fletching, BodkinArrow, Bracer, PaddedArcherArmor, LeatherArcherArmor, RingArcherArmor,
		ThumbRing, ParthianTactics, TechHeavyCavalryArcher}, AllResources, 0.75},
	{Tupi, []TechID{TechEliteSkirmisher, TechImperialSkirmisher, ParthianTactics, ThumbRing,
		TechHeavyCavalryArcher}, FoodCost, 0.5},
}

// Apply is the cost after the change.
func (change TechnologyCostChange) Apply(cost Cost) Cost {
	scale := func(amount float64, resource CostResource) float64 {
		if change.Resource == AllResources || change.Resource == resource {
			return amount * change.Multiplier
		}
		return amount
	}
	return Cost{
		Food:  scale(cost.Food, FoodCost),
		Wood:  scale(cost.Wood, WoodCost),
		Gold:  scale(cost.Gold, GoldCost),
		Stone: scale(cost.Stone, StoneCost),
	}
}

// FormatTechnologyCostNote explains the technologies that the civilization researches cheaper,
// empty if none, e.g. "Hussar upgrade free instead of 500f 600g".
// It covers the technologies that change the row unit and the upgrades to the unit.
func FormatTechnologyCostNote(row Row) string {
	var notes []string
	for _, item := range listCheaperTechnologies(row) {
		text := formatCost(item.cheaper)
		if text == "" {
			text = "free"
		}
		notes = append(notes, fmt.Sprintf("%v %v instead of %v", item.name, text, formatCost(item.cost)))
	}
	return strings.Join(notes, ", ")
}

// CalculateTechnologySaving is the resources the civilization saves on the technologies of FormatTechnologyCostNote,
// e.g. 1300 for the Turks Hussar (Light Cavalry and Hussar upgrades free).
func CalculateTechnologySaving(row Row) float64 {
	saving := 0.0
	for _, item := range listCheaperTechnologies(row) {
		saving += item.cost.Sum() - item.cheaper.Sum()
	}
	return saving
}

// cheaperTechnology is a technology the civilization researches cheaper, one per cost change.
type cheaperTechnology struct {
	name    string
	cost    Cost
	cheaper Cost
}

// listCheaperTechnologies lists the technologies that change the row unit and the upgrades to the unit,
// that the civilization researches cheaper.
func listCheaperTechnologies(row Row) []cheaperTechnology {
	type researched struct {
		technology Technology
		name       string
	}
	var list []researched
	for _, step := range row.Steps {
		if step.Source == Technologies && step.Technology.ID != 0 && !reflect.DeepEqual(step.UnitBefore, step.UnitAfter) {
			list = append(list, researched{step.Technology, step.Name})
		}
	}
	for _, earlier := range append(append([]Row{}, row.EarlierRows...), row) {
		for _, technology := range earlier.UpgradeTechnologies {
			list = append(list, researched{technology, technology.Name + " upgrade"})
		}
	}

	var cheaperList []cheaperTechnology
	for _, item := range list {
		cost := item.technology.Cost
		for _, change := range TechnologyCostChanges {
			if change.CivilizationID != row.CivilizationID || !containsTechID(change.TechnologyIDs, item.technology.ID) {
				continue
			}
			cheaper := change.Apply(cost)
			if cheaper == cost {
				continue
			}
			cheaperList = append(cheaperList, cheaperTechnology{name: item.name, cost: cost, cheaper: cheaper})
			cost = cheaper
		}
	}
	return cheaperList
}

func containsTechID(ids []TechID, id TechID) bool {
	return slices.Contains(ids, id)
}
