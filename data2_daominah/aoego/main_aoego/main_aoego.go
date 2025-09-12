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

	//inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Assyria_Archer.ai`
	//inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Babylon_Tower_Priest.ai`
	inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Babylon_Chariot.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Carthage_Helepolis.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Choson_Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Choson_Tower.ai`
	//inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Egypt_Chariot_Priest.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Greek_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Hittite_Horse_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Hittite_Catapult.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Macedon_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Minoa_Bowmen_Helepolis.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Minoa_Bowmen_Catapult.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Palmyra_Camel.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Palmyra_Horse_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Persia_War_Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Phoenicia_Elephant_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Rome_Legion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Rome_Siege.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Shang_Stable.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Sumeria_Catapult.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Yamato_Cavalry.ai`

	// bull shit original strategies, e.g. Rome Legion without researching ShortSword

	//inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Assyria.ai`
	//inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Babylon Scouts.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Babylon Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Carthage Phalanx.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Carthage War Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Choson Axemen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Choson Priests.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Choson Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Egypt.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Greek.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Hittite Bowmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Hittite Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Hittite Horse Archers.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Macedon Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Macedon Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Macedon Phalanx.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Minoa.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Palmyra Composite Bow.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Palmyra Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Palmyra Horse Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Persia Elephant Archers.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Persia Priests.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Persia War Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Phoenicia Elephants.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Rome Axemen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Rome Legion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Rome Siege.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Shang Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Shang Clubmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Shang Heavy Cavalry.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Sumeria.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\backup_data2\Immortal Yamato.ai`

	log.Printf("input file path: %v", inputFilePath)
	strategyBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	strategy, errs := aoego.NewStrategy(string(strategyBytes))
	// errs will be printed at the end for better readability

	civilizationID := aoego.GuessCivilization(inputFilePath)
	empire, err := aoego.NewEmpireDeveloping(aoego.WithCivilization(civilizationID))
	if err != nil {
		log.Fatalf("error NewEmpireDeveloping: %v", err)
	}

	// only print empire summary the first time PrintSummary is called
	onlyOncePrintSummary := false
	for _, step := range strategy {
		if step.Action == aoego.PrintSummary && !onlyOncePrintSummary {
			log.Printf(empire.Summary())
			onlyOncePrintSummary = true
			continue
		}

		//if step.UnitOrTechID == aoego.Wonder {
		//	log.Printf("----------------")
		//	log.Printf("empire: %v", empire.Summary())
		//	break // origin AI files add nonsense number of units after Wonder
		//}

		err := empire.Do(step)
		if err != nil {
			log.Printf("error line %-3v empire.Do(%v): %v", step.OriginLineNo, step, err)
		}
	}

	if len(errs) > 0 {
		log.Printf("____________________________________________________")
		for _, err := range errs {
			log.Printf("error parsing strategy: %v", err)
		}
		log.Printf("ERROR PARSING STRATEGY: %v", len(errs))
		log.Printf("____________________________________________________")
	}
}
