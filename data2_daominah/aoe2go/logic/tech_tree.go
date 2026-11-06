package logic

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// TechTree is the parsed data/dat_original_to_json:
// base stats of units, and what is available to each civilization.
type TechTree struct {
	Units         map[UnitID]Unit
	Technologies  map[TechID]Technology
	Civilizations []Civilization // the AoE2 ones (not Chronicles), sorted by name
}

// Technology is what a ranking shows of a technology.
type Technology struct {
	ID             TechID
	Name           string // English name, e.g. "Corvinian Army"
	Cost           Cost
	EffectID       int
	CivilizationID CivilizationID // -1 is any civilization
	// RequiredTechnologyIDs include the age, e.g. CastleAge for a Castle Age civilization bonus.
	RequiredTechnologyIDs []TechID
}

// FindAge is the age technology the technology requires, 0 if none,
// e.g. CastleAge for the Berbers "Stable units -15%" bonus.
func (t Technology) FindAge() TechID {
	age := TechID(0)
	for _, id := range t.RequiredTechnologyIDs {
		if _, isAge := AgeNames[id]; isAge {
			age = max(age, id)
		}
	}
	return age
}

// AgeNames are the English names of the age technologies.
// The .dat internal names are shifted, e.g. "Feudal Age" is CastleAge.
var AgeNames = map[TechID]string{FeudalAge: "Feudal Age", CastleAge: "Castle Age", ImperialAge: "Imperial Age"}

// AnyCivilization is the CivilizationID of a technology that any civilization can research,
// if its Technology Tree allows it.
const AnyCivilization CivilizationID = -1

// findCivilizationBonusTechnology is the technology of a civilization that triggers the effect,
// e.g. the Berbers Imperial Age technology for "Stable units -20%".
// Some effects are also triggered by a technology for any civilization,
// e.g. Khitans "Double Forging" triggers the Forging effect again.
func (tree *TechTree) findCivilizationBonusTechnology(civilizationID CivilizationID, effectID int) Technology {
	result := Technology{}
	for _, t := range tree.Technologies {
		if t.EffectID == effectID && t.CivilizationID == civilizationID && (result.Name == "" || t.ID < result.ID) {
			result = t
		}
	}
	return result
}

// Civilization lists what is available to a civilization in its Technology Tree.
type Civilization struct {
	ID            CivilizationID
	Name          string // English name, e.g. "Hindustanis"
	UnitIDs       map[UnitID]bool
	TechnologyIDs map[TechID]bool
}

// Unit is the stats of a unit, base or after applying effects.
type Unit struct {
	ID    UnitID
	Name  string    // English name, e.g. "Hussar"
	Class UnitClass // unit class in the .dat, e.g. ClassCavalry
	HP    float64
	// Attacks and Armors are amounts by armor class, e.g. ArmorClassPierce.
	// Damage only counts the classes in both the attacker Attacks and the target Armors.
	Attacks    map[ArmorClass]float64
	Armors     map[ArmorClass]float64
	Range      float64 // max range in tiles, e.g. 0 for a Knight, 1 for a Steppe Lancer, 5 for a Throwing Axeman
	ReloadTime float64 // seconds between attacks
	Speed      float64
	TrainTime  float64 // seconds
	Cost       Cost
	// IsIgnoringArmor is set by Wootz Steel: the unit attack ignores the target armor.
	IsIgnoringArmor bool
}

// Cost is the resources to train a unit.
type Cost struct {
	Food  float64
	Wood  float64
	Gold  float64
	Stone float64
}

// Sum is the total of all resources.
func (c Cost) Sum() float64 {
	return c.Food + c.Wood + c.Gold + c.Stone
}

// ParseTechTree parses civilizations.json, units.json and technologies.json
// from data/dat_original_to_json, names are from ParseGameStrings.
func ParseTechTree(civilizationsJSON, unitsJSON, technologiesJSON []byte, names map[int]string) (*TechTree, error) {
	var civilizations []civilizationJSON
	if err := json.Unmarshal(civilizationsJSON, &civilizations); err != nil {
		return nil, fmt.Errorf("error json.Unmarshal civilizations.json: %w", err)
	}
	var units []unitJSON
	if err := json.Unmarshal(unitsJSON, &units); err != nil {
		return nil, fmt.Errorf("error json.Unmarshal units.json: %w", err)
	}
	var technologies []technologyJSON
	if err := json.Unmarshal(technologiesJSON, &technologies); err != nil {
		return nil, fmt.Errorf("error json.Unmarshal technologies.json: %w", err)
	}

	tree := &TechTree{Units: map[UnitID]Unit{}, Technologies: map[TechID]Technology{}}
	for _, t := range technologies {
		technology := Technology{
			ID:             TechID(t.ID),
			Name:           names[t.NameStringID],
			Cost:           parseCost(t.Cost),
			EffectID:       t.EffectID,
			CivilizationID: CivilizationID(t.CivilizationID),
		}
		for _, id := range t.RequiredTechnologyIDs {
			technology.RequiredTechnologyIDs = append(technology.RequiredTechnologyIDs, TechID(id))
		}
		if technology.Name == "" {
			technology.Name = t.InternalName
		}
		tree.Technologies[technology.ID] = technology
	}
	for _, u := range units {
		unit := Unit{
			ID:         UnitID(u.ID),
			Name:       names[u.NameStringID],
			HP:         u.HP,
			Attacks:    map[ArmorClass]float64{},
			Armors:     map[ArmorClass]float64{},
			Class:      UnitClass(u.Class),
			Range:      u.Range,
			ReloadTime: u.ReloadTime,
			Speed:      u.Speed,
			TrainTime:  u.TrainTime,
			Cost:       parseCost(u.Cost),
		}
		if unit.Name == "" {
			unit.Name = u.InternalName
		}
		for class, amount := range u.Attacks {
			unit.Attacks[ArmorClass(class)] = amount
		}
		for class, amount := range u.Armors {
			unit.Armors[ArmorClass(class)] = amount
		}
		tree.Units[unit.ID] = unit
	}
	for _, c := range civilizations {
		if c.Era != baseEra || c.ID == gaiaID {
			continue
		}
		civilization := Civilization{
			ID:            CivilizationID(c.ID),
			Name:          names[c.NameStringID],
			UnitIDs:       map[UnitID]bool{},
			TechnologyIDs: map[TechID]bool{},
		}
		for _, id := range c.UnitIDs {
			civilization.UnitIDs[UnitID(id)] = true
		}
		for _, id := range c.TechnologyIDs {
			civilization.TechnologyIDs[TechID(id)] = true
		}
		tree.Civilizations = append(tree.Civilizations, civilization)
	}
	sort.Slice(tree.Civilizations, func(i, j int) bool {
		return tree.Civilizations[i].Name < tree.Civilizations[j].Name
	})
	return tree, nil
}

// parseCost reads a cost of the JSON files, e.g. {"food": 200, "gold": 300}.
func parseCost(cost map[string]float64) Cost {
	return Cost{Food: cost["food"], Wood: cost["wood"], Gold: cost["gold"], Stone: cost["stone"]}
}

// The AoE2 civilizations: era "base", without Gaia (Chronicles are "antiquity").
const (
	baseEra = "base"
	gaiaID  = 0
)

// ParseGameStrings parses the game strings file, lines like `5661 "Hussar"`,
// into text by string ID. Comment lines and other lines are skipped.
func ParseGameStrings(content []byte) map[int]string {
	names := map[int]string{}
	for _, match := range gameStringLine.FindAllSubmatch(content, -1) {
		id, err := strconv.Atoi(string(match[1]))
		if err != nil {
			continue
		}
		names[id] = string(match[2])
	}
	return names
}

// gameStringLine matches a line of the game strings file, e.g. `5661 "Hussar"`.
var gameStringLine = regexp.MustCompile(`(?m)^(\d+)\s+"(.*)"\r?$`)

// civilizationJSON is the part of civilizations.json used here.
type civilizationJSON struct {
	ID            int    `json:"id"`
	NameStringID  int    `json:"name_string_id"`
	Era           string `json:"era"`
	UnitIDs       []int  `json:"unit_ids"`
	TechnologyIDs []int  `json:"technology_ids"`
}

// technologyJSON is the part of technologies.json used here.
type technologyJSON struct {
	ID                    int                `json:"id"`
	InternalName          string             `json:"internal_name"`
	NameStringID          int                `json:"name_string_id"`
	Cost                  map[string]float64 `json:"cost"`
	EffectID              int                `json:"effect_id"`
	CivilizationID        int                `json:"civilization_id"`
	RequiredTechnologyIDs []int              `json:"required_technology_ids"`
}

// unitJSON is the part of units.json used here.
type unitJSON struct {
	ID           int                `json:"id"`
	InternalName string             `json:"internal_name"`
	NameStringID int                `json:"name_string_id"`
	Class        int                `json:"class"`
	HP           float64            `json:"hp"`
	Speed        float64            `json:"speed"`
	Range        float64            `json:"range"`
	ReloadTime   float64            `json:"reload_time"`
	TrainTime    float64            `json:"train_time"`
	Attacks      map[int]float64    `json:"attacks"`
	Armors       map[int]float64    `json:"armors"`
	Cost         map[string]float64 `json:"cost"`
}
