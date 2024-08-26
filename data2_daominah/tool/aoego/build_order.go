// Package aoego is used to define the data structure in the game
// Age of Empires: The Rise of Rome.
// This can be used to validate and calculate the cost of
// build order (".ai" file), which defines computer player's strategy.
package aoego

// BuildOrder is equivalent to a ".ai" file, each step is a line in the file.
type BuildOrder []Step

type Step struct {
	Action       Action
	UnitOrTech   UnitOrTech
	Location     UnitID // parent building that trains unit or researches technology
	LimitRebuild int    // only used when Action is BuildLimit or TrainLimit
}

// Action is enum: build or train or research
type Action string

// Action enum
const (
	Build            Action = "B" // will rebuild if destroyed
	BuildLimit       Action = "A" // will rebuild up to the limit times specified if destroyed
	Train            Action = "U" // will always be replaced if killed
	TrainLimit       Action = "T" // will be trained up to the limit times specified
	Research         Action = "R" // item will be researched if possible, can be skipped if failed too many attempts
	ResearchCritical Action = "C" // the build stuck until this research is done
)

type UnitOrTech int // UnitOrTech is a UnitID or TechID
