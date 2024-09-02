package main

import (
	"flag"
	"log"
	"os"
	"strings"

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

	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Assyria_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Babylon_Tower_Priest.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Carthage.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Choson_Swordsmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Egypt_Chariot_Priest.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Greek_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Hittite_Horse_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Macedon_Centurion.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Minoa_Composite_Bowmen.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Palmyra.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Persia_War_Elephant.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Phoenicia_Elephant_Archer.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Rome.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Shang.ai`
	// inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Sumeria_Catapult.ai`
	inputFilePath = `D:\game\age_of_empires_ror_hd\data2_daominah\Yamato_Cavalry.ai`

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
		log.Fatalf("error parsing strategy len(errs): %v", len(errs))
	}
	civilizationID := aoego.GuessCivilization(inputFilePath)
	empire, err := aoego.NewEmpireDeveloping(aoego.WithCivilization(civilizationID))
	if err != nil {
		log.Fatalf("error NewEmpireDeveloping: %v", err)
	}

	for i, step := range strategy {
		if step.Action == aoego.PrintSummary {
			if !strings.Contains(step.OriginStr, "spent army") {
				// continue
			}
			log.Printf("----------------")
			log.Printf("empire: %v", empire.Summary())
			continue
		}
		err := empire.Do(step)
		if err != nil {
			log.Printf("error i %v empire.Do(%v): %v", i, step, err)
		}
	}

	// log.Printf("\n\n----------------")
	// log.Printf("empire end: %v", empire.Summary())
}
