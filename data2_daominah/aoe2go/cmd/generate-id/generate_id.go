// Command generate-id writes the Go constants of the AoE2 DE IDs to package logic,
// so the code says Hussar instead of 441:
// UnitID for units and buildings (a building is a unit in the .dat), TechID, CivilizationID, ArmorClass.
// Run it after refreshing the data files: go generate ./logic
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoe2go/data"
)

// Naming of duplicate names (edit after a game patch when the generator warns):
var (
	// variants maps a normal unit or building to its variants sharing the in-game name:
	// the normal one keeps the bare name, a variant gets the suffix (its ID if empty),
	// e.g. Huskarl is trained at the Castle, HuskarlBarracks needs a unique technology.
	// The suffix is mostly the building that trains the variant, read from the .dat.
	// It is also generated as logic.UnitVariants.
	variants = map[int][]variant{
		41:   {{759, "Barracks"}},     // Huskarl, Goths
		555:  {{761, "Barracks"}},     // Elite Huskarl
		755:  {{886, "Stable"}},       // Tarkan, Huns
		757:  {{887, "Stable"}},       // Elite Tarkan
		1225: {{1254, "Krepost"}},     // Konnik, Bulgarians, the normal one is trained at the Castle and Krepost
		1227: {{1255, "Krepost"}},     // Elite Konnik
		93:   {{1786, "Donjon"}},      // Spearman, Sicilians
		358:  {{1787, "Donjon"}},      // Pikeman
		359:  {{1788, "Donjon"}},      // Halberdier
		1658: {{1660, "Donjon"}},      // Serjeant, the normal one is trained at the Castle and Donjon
		1659: {{1661, "Donjon"}},      // Elite Serjeant
		1233: {{1260, "Mercenaries"}}, // Elite Kipchak, Cumans team bonus Cuman Mercenaries
		// Town Center: a villager places 621, it becomes 109 when finished,
		// the age technologies replace 109 with 71, 141 and 142, each with +1 melee and pierce armor
		109: {{621, "Unfinished"}, {71, "FeudalAge"}, {141, "CastleAge"}, {142, "ImperialAge"}},
		// the age technologies also replace these buildings, each age with more armor and HP
		12:  {{498, "FeudalAge"}, {132, "CastleAge"}, {20, "ImperialAge"}}, // Barracks
		104: {{30, "FeudalAge"}, {31, "CastleAge"}, {32, "ImperialAge"}},   // Monastery
	}
)

// techPrefix is added to a technology named like a unit or building,
// e.g. the upgrade Hussar is TechHussar, the unit is Hussar.
const techPrefix = "Tech"

// armorClassPrefix is added to all armor classes,
// e.g. ArmorClassCavalry, as many are named like units, e.g. Monks.
const armorClassPrefix = "ArmorClass"

// armorClassStringID is the game string ID of armor class 0, "Unused Armor 0".
const armorClassStringID = 12400

// Which civilizations get constants: the AoE2 ones, not Gaia or Chronicles.
const (
	baseEra = "base"
	gaiaID  = 0
)

// IDs not in any civilization Technology Tree that still get a constant,
// the same as the Technology Tree website adds, see data/dat_original_to_json/README.md.
var (
	// e.g. Packed Trebuchet, dismounted Konnik, Ratha, War Chariot modes
	additionalUnitIDs     = []int{42, 1252, 1253, 1570, 1738, 1740, 1962, 1980}
	additionalBuildingIDs = []int{71, 141, 142, 498, 132, 20, 30, 31, 32} // the age stages of the Town Center, Barracks and Monastery, the age technologies replace the building with them
	additionalTechIDs     = []int{19, 20, 90}                             // Cartography, Auto upgrade Scout Feudal Age, Tracking
)

func main() {
	var civilizations []civilizationJSON
	if err := json.Unmarshal(data.Civilizations, &civilizations); err != nil {
		log.Fatalf("error json.Unmarshal civilizations: %v", err)
	}
	var unitList, techList []itemJSON
	if err := json.Unmarshal(data.Units, &unitList); err != nil {
		log.Fatalf("error json.Unmarshal units: %v", err)
	}
	if err := json.Unmarshal(data.Technologies, &techList); err != nil {
		log.Fatalf("error json.Unmarshal technologies: %v", err)
	}
	names := parseGameStrings(data.GameStringsEnglish)

	// constants for what the AoE2 civilizations can use, as in their Technology Tree
	unitIDs, buildingIDs, techIDs, upgradeTechIDs := map[int]bool{}, map[int]bool{}, map[int]bool{}, map[int]bool{}
	var baseCivilizations []civilizationJSON
	for _, c := range civilizations {
		if c.Era != baseEra || c.ID == gaiaID {
			continue
		}
		baseCivilizations = append(baseCivilizations, c)
		addAll(unitIDs, c.UnitIDs)
		addAll(buildingIDs, c.BuildingIDs)
		addAll(techIDs, c.TechnologyIDs)
		for _, techID := range c.UpgradeTechnologyIDs {
			upgradeTechIDs[techID] = true
		}
	}
	// like the Technology Tree website, an upgrade without a cost is not a real technology,
	// e.g. 1156 is the trigger of a node no civilization has
	for _, t := range techList {
		if upgradeTechIDs[t.ID] && len(t.Cost) > 0 {
			techIDs[t.ID] = true
		}
	}
	addAll(unitIDs, additionalUnitIDs)
	addAll(buildingIDs, additionalBuildingIDs)
	addAll(techIDs, additionalTechIDs)

	units := collectNamedIDs(unitList, unitIDs, names)
	buildings := collectNamedIDs(unitList, buildingIDs, names)
	techs := collectNamedIDs(techList, techIDs, names)

	unitConstants := assignConstantNames(units)
	buildingConstants := assignConstantNames(buildings)
	taken := map[string]bool{}
	for _, c := range unitConstants {
		taken[c.Constant] = true
	}
	for _, c := range buildingConstants {
		if taken[c.Constant] {
			log.Fatalf("building %v has the name of a unit, add a suffix in the generator", c.Constant)
		}
		taken[c.Constant] = true
	}
	techConstants := assignConstantNames(techs)
	for i, c := range techConstants {
		if taken[c.Constant] {
			techConstants[i].Constant = techPrefix + c.Constant
		}
	}
	for _, c := range techConstants {
		if taken[c.Constant] {
			log.Fatalf("technology %v is still a duplicate name", c.Constant)
		}
		taken[c.Constant] = true
	}
	var civilizationConstants []namedID
	for _, c := range baseCivilizations {
		constant := makeConstantName(names[c.NameStringID])
		if taken[constant] {
			log.Fatalf("civilization %v has the name of a unit or technology", constant)
		}
		civilizationConstants = append(civilizationConstants,
			namedID{ID: c.ID, Name: names[c.NameStringID], InternalName: c.DatName, Constant: constant})
	}

	logicDir := filepath.Join(findSourceDir(), "..", "..", "logic")
	writeGoFile(filepath.Join(logicDir, "unit_id.go"), renderUnitIDFile(unitConstants, buildingConstants))
	writeGoFile(filepath.Join(logicDir, "building_id.go"), renderConstantsFile(
		"Buildings, as UnitID: a building is a unit in the .dat.", "UnitID", buildingConstants))
	writeGoFile(filepath.Join(logicDir, "tech_id.go"), renderTechIDFile(techConstants))
	writeGoFile(filepath.Join(logicDir, "civilization_id.go"), renderCivilizationFile(civilizationConstants))
	writeGoFile(filepath.Join(logicDir, "armor_class.go"), renderArmorClassFile(names))
}

func addAll(set map[int]bool, ids []int) {
	for _, id := range ids {
		set[id] = true
	}
}

// collectNamedIDs reads the ID and English name of the items with the wanted IDs.
func collectNamedIDs(items []itemJSON, wanted map[int]bool, names map[int]string) map[int]namedID {
	result := map[int]namedID{}
	for _, item := range items {
		if !wanted[item.ID] {
			continue
		}
		name := names[item.NameStringID]
		if name == "" {
			name = item.InternalName
		}
		result[item.ID] = namedID{ID: item.ID, Name: name, InternalName: item.InternalName}
	}
	for id := range wanted {
		if _, isFound := result[id]; !isFound {
			log.Fatalf("ID %v is not in the extracted data", id)
		}
	}
	return result
}

// assignConstantNames turns the names into Go constants, sorted by ID.
// A duplicate name listed in variants gets its suffix,
// an unlisted duplicate gets its ID with a warning, except the lowest ID.
func assignConstantNames(items map[int]namedID) []namedID {
	suffixes := map[int]string{}
	isNormal := map[int]bool{}
	for normalID, vs := range variants {
		isNormal[normalID] = true
		for _, v := range vs {
			suffixes[v.ID] = v.Suffix
			if v.Suffix == "" {
				suffixes[v.ID] = strconv.Itoa(v.ID)
			}
		}
	}
	byName := map[string][]int{}
	for id, item := range items {
		c := makeConstantName(item.Name)
		byName[c] = append(byName[c], id)
	}
	var result []namedID
	for c, ids := range byName {
		sort.Ints(ids)
		for i, id := range ids {
			item := items[id]
			item.Constant = c
			switch {
			case len(ids) == 1 || isNormal[id]:
			case suffixes[id] != "":
				item.Constant = c + suffixes[id]
			case i == 0:
				log.Printf("warning: %v %v has the bare name, the same name as %v, add them to variants", c, id, ids[1:])
			default:
				item.Constant = c + strconv.Itoa(id)
				log.Printf("warning: %v %v is named %v, add it to variants", c, id, item.Constant)
			}
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

// makeConstantName is the name in PascalCase without punctuation,
// e.g. "Two-Handed Swordsman" is TwoHandedSwordsman, "Villager (Male)" is VillagerMale.
func makeConstantName(name string) string {
	var b strings.Builder
	isWordStart := true
	for _, r := range name {
		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			isWordStart = true
			continue
		}
		if isWordStart {
			r = unicode.ToUpper(r)
		}
		b.WriteRune(r)
		isWordStart = false
	}
	c := b.String()
	if c == "" || unicode.IsDigit(rune(c[0])) {
		c = "ID" + c
	}
	return c
}

func renderUnitIDFile(units []namedID, buildings []namedID) string {
	constants := map[int]string{}
	for _, c := range append(units, buildings...) {
		constants[c.ID] = c.Constant
	}
	var normalIDs []int
	for id := range variants {
		normalIDs = append(normalIDs, id)
	}
	sort.Ints(normalIDs)
	var b strings.Builder
	b.WriteString(`// UnitID is a unit or building ID in the game .dat.
type UnitID int

// UnitVariants maps a normal unit or building to its variants sharing the in-game name,
// e.g. Huskarl to HuskarlBarracks: Goths train it at the Barracks after a unique technology.
var UnitVariants = map[UnitID][]UnitID{
`)
	for _, id := range normalIDs {
		var names []string
		for _, v := range variants[id] {
			names = append(names, constants[v.ID])
		}
		fmt.Fprintf(&b, "\t%v: {%v},\n", constants[id], strings.Join(names, ", "))
	}
	b.WriteString("}\n\n")
	return b.String() + renderConstantsFile("Units.", "UnitID", units)
}

func renderTechIDFile(techs []namedID) string {
	header := `// TechID is a technology ID in the game .dat.
// A technology named like a unit or building has the prefix Tech, e.g. TechHussar.
type TechID int

`
	return header + renderConstantsFile("Technologies, including the unit upgrades.", "TechID", techs)
}

func renderConstantsFile(title string, typeName string, items []namedID) string {
	var b strings.Builder
	fmt.Fprintf(&b, "// %v\n// The comment is the internal name in the .dat.\nconst (\n", title)
	for _, item := range items {
		fmt.Fprintf(&b, "\t%v %v = %v // %v\n", item.Constant, typeName, item.ID, item.InternalName)
	}
	b.WriteString(")\n")
	return b.String()
}

func renderCivilizationFile(civilizations []namedID) string {
	var b strings.Builder
	b.WriteString(`// CivilizationID is a civilization ID in the .dat.
type CivilizationID int

// Civilizations of AoE2 (not Chronicles), named by the game strings.
// The comment is the name in the .dat, some are old, e.g. "French" for Franks.
const (
`)
	for _, c := range civilizations {
		fmt.Fprintf(&b, "\t%v CivilizationID = %v // %v\n", c.Constant, c.ID, c.InternalName)
	}
	b.WriteString(")\n")
	return b.String()
}

// renderArmorClassFile names the armor classes from the game strings,
// the name of armor class N is string armorClassStringID + N, e.g. 12403 "Pierce".
func renderArmorClassFile(names map[int]string) string {
	var b strings.Builder
	b.WriteString(`import "fmt"

// ArmorClass is an armor class in the .dat: attacks and armors are amounts by class,
// a hit only counts the classes in both the attacker attacks and the target armors.
type ArmorClass int

// Armor classes, named by the game strings.
const (
`)
	var classNames []string
	for class := 0; ; class++ {
		name, isFound := names[armorClassStringID+class]
		if !isFound {
			break
		}
		fmt.Fprintf(&b, "\t%v%v ArmorClass = %v // %v\n", armorClassPrefix, makeConstantName(name), class, name)
		classNames = append(classNames, name)
	}
	b.WriteString(")\n")
	b.WriteString(`
// String is the game name of the armor class, e.g. "Pierce",
// so a printed unit reads map[Pierce:4 Melee:3] instead of map[3:4 4:3].
func (c ArmorClass) String() string {
	if name, isFound := armorClassNames[c]; isFound {
		return name
	}
	return fmt.Sprintf("ArmorClass(%d)", int(c))
}

// armorClassNames are the game names of the armor classes.
var armorClassNames = map[ArmorClass]string{
`)
	for class, name := range classNames {
		fmt.Fprintf(&b, "\t%v: %q,\n", class, name)
	}
	b.WriteString("}\n")
	return b.String()
}

// writeGoFile writes the generated code formatted, with the generated header.
func writeGoFile(path string, body string) {
	source := "// Code generated by cmd/generate-id from data/dat_original_to_json and the game strings; DO NOT EDIT.\n\npackage logic\n\n" + body
	formatted, err := format.Source([]byte(source))
	if err != nil {
		log.Fatalf("error format.Source %v: %v", path, err)
	}
	if err := os.WriteFile(path, formatted, 0o644); err != nil {
		log.Fatalf("error os.WriteFile: %v", err)
	}
	log.Printf("wrote %v", path)
}

// findSourceDir is the directory of this source file,
// so the output lands in logic whatever the working directory of "go run".
func findSourceDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

// parseGameStrings is a copy of logic.ParseGameStrings:
// importing logic would stop the generator from building when a generated file is broken.
func parseGameStrings(content []byte) map[int]string {
	names := map[int]string{}
	for _, match := range gameStringLine.FindAllSubmatch(bytes.TrimPrefix(content, []byte("\xef\xbb\xbf")), -1) {
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

// namedID is an ID with its English name and the generated constant name.
type namedID struct {
	ID           int
	Name         string
	InternalName string
	Constant     string
}

// variant is a unit or building sharing the in-game name of a normal one.
type variant struct {
	ID     int
	Suffix string // added to the bare name, the ID if empty
}

// civilizationJSON is the part of civilizations.json used here.
type civilizationJSON struct {
	ID                   int            `json:"id"`
	DatName              string         `json:"dat_name"`
	NameStringID         int            `json:"name_string_id"`
	Era                  string         `json:"era"`
	UnitIDs              []int          `json:"unit_ids"`
	BuildingIDs          []int          `json:"building_ids"`
	TechnologyIDs        []int          `json:"technology_ids"`
	UpgradeTechnologyIDs map[string]int `json:"upgrade_technology_ids"`
}

// itemJSON is the part of units.json and technologies.json used here.
type itemJSON struct {
	ID           int                `json:"id"`
	InternalName string             `json:"internal_name"`
	NameStringID int                `json:"name_string_id"`
	Cost         map[string]float64 `json:"cost"`
}
