package main

import (
	"flag"
	"log"
	"os"

	"github.com/daominah/age_of_empires_ror_hd/data2_daominah/aoego"
)

// inputFilePath is the path to the "*.ai" file that defines the strategy
var inputFilePath string

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Lmicroseconds)

	flag.StringVar(&inputFilePath,
		"i",
		`D:\game\age_of_empires_ror_hd\data2_daominah\Assyria_Archer.ai`,
		`the path to the "*.ai" file that defines the strategy`)
	flag.Parse()
	if inputFilePath == "" {
		log.Printf("empty input file path")
		return
	}

	// daominah optimized strategies:

	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Assyria_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Babylon_Tower_Priest.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Carthage_Helepolis.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Choson_Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Egypt_Chariot_Priest.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Greek_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Hittite_Horse_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Hittite_Catapult.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Macedon_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Minoa_Composite_Bowmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Palmyra_Camel.ai`
	inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Palmyra_Stable.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Persia_War_Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Phoenicia_Elephant_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Rome_Legion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Rome_Siege.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Shang_Stable.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Sumeria_Catapult.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Yamato_Cavalry.ai`

	// bull shit original strategies, e.g. Rome Legion without researching ShortSword

	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Assyria.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Babylon Scouts.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Babylon Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Carthage Phalanx.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Carthage War Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Choson Axemen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Choson Priests.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Choson Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Egypt.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Greek.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Hittite Bowmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Hittite Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Hittite Horse Archers.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Macedon Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Macedon Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Macedon Phalanx.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Minoa.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Palmyra Composite Bow.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Palmyra Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Palmyra Horse Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Persia Elephant Archers.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Persia Priests.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Persia War Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Phoenicia Elephants.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Rome Axemen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Rome Legion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Rome Siege.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Shang Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Shang Clubmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Shang Heavy Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Sumeria.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_no_edit_20240801\Immortal Yamato.ai`

	log.Printf("input file path: %v", inputFilePath)
	strategyBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	strategy, errs := aoego.NewStrategy(string(strategyBytes))
	if len(errs) > 0 {
		for _, err := range errs {
			log.Printf("error parsing strategy: %v", err)
		}
		log.Printf("____________________________________________________")
		log.Printf("ERROR PARSING STRATEGY len(errs): %v", len(errs))
		log.Printf("____________________________________________________")
	}
	civilizationID := aoego.GuessCivilization(inputFilePath)
	empire, err := aoego.NewEmpireDeveloping(aoego.WithCivilization(civilizationID))
	if err != nil {
		log.Fatalf("error NewEmpireDeveloping: %v", err)
	}

	for _, step := range strategy {
		if step.Action == aoego.PrintSummary {
			log.Printf(empire.Summary())
			continue
		}

		//if step.UnitOrTechID == aoego.Wonder {
		//	log.Printf("----------------")
		//	log.Printf("empire: %v", empire.Summary())
		//	break // origin AI files add nonsense number of units after Wonder
		//}

		err := empire.Do(step)
		if err != nil {
			log.Printf("error line %v empire.Do(%v): %v", step.OriginLineNo, step, err)
		}
	}

	// log.Printf("\n\n----------------")
	// log.Printf("empire end: %v", empire.Summary())
}
