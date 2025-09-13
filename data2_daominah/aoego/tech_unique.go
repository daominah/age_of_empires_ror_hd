package aoego

const (
	Choson_cheap_ATK2      TechID = 144
	Choson_cheap_ATK4      TechID = 145
	Choson_cheap_ATK7      TechID = 146
	Choson_Cataphracts     TechID = 147
	Choson_cheap_Armor2Inf TechID = 148
	Choson_cheap_Armor2Arc TechID = 149
	Choson_cheap_Armor2Cav TechID = 150
	Choson_cheap_Armor4Inf TechID = 151
	Choson_cheap_Armor4Arc TechID = 152
	Choson_cheap_Armor4Cav TechID = 153
	Choson_cheap_Shield1   TechID = 154

	Assyrian_cheap_MassiveCatapult TechID = 155
	Assyrian_cheap_Catapult        TechID = 156
	Assyrian_cheap_Helepolis       TechID = 157

	Babylonian_cheap_Wheel         TechID = 158
	Babylonian_cheap_Coinage       TechID = 159 // gold2
	Babylonian_cheap_Plow          TechID = 160 // farm2
	Babylonian_cheap_Artisanship   TechID = 161 // wood2
	Babylonian_cheap_Irrigation    TechID = 162 // farm3
	Babylonian_cheap_Domestication TechID = 163 // farm1
	Babylonian_cheap_WoodWorking   TechID = 164 // wood1
	Babylonian_cheap_GoldMining    TechID = 165 // gold1
	Babylonian_cheap_StoneMining   TechID = 166 // stone1
	Babylonian_cheap_Craftsmanship TechID = 167 // wood3
	Babylonian_cheap_SiegeCraft    TechID = 168 // stone2
	Babylonian_unlock_ChariotArc   TechID = 169 // unchanged, but Wheel -> ChariotArcher
	Babylonian_unlock_Chariot      TechID = 170 // unchanged, but Wheel -> Chariot
	Babylonian_MassiveCatapult     TechID = 171 // unchanged, but SiegeCraft -> MassiveCatapult
	Babylonian_ScytheChariot       TechID = 172 // unchanged, but Wheel -> Chariot -> ScytheChariot

	Yamato_cheap_HHorseArc   TechID = 173
	Yamato_cheap_ImprovedBow TechID = 174
	Yamato_cheap_CompositBow TechID = 175
	Yamato_cheap_Cavalry2    TechID = 176
	Yamato_cheap_Cataphracts TechID = 177

	Carthaginian_cheap_Nobility TechID = 178

	Greek_cheap_Astrology  TechID = 179
	Greek_cheap_Polytheism TechID = 180

	Hittite_cheap_Wheel       TechID = 181
	Hittite_unlock_ChariotArc TechID = 182 // unchanged, but Wheel -> ChariotArcher
	Hittite_unlock_Chariot    TechID = 183 // unchanged, but Wheel -> Chariot
	Hittite_ScytheChariot     TechID = 184 // unchanged, but Wheel -> Chariot -> ScytheChariot
)

// CheckUniqueTechID returns the uniqueTechID for a given baseTechID and civID,
// or NullTech if not found civilization-specific unique tech for the baseTechID
func CheckUniqueTechID(baseTechID TechID, civID CivilizationID) TechID {
	switch civID {
	// Choson has 10 Storage Pit techs -40% cost,
	// Metallurgy is requirement for Cataphract,
	// so total 11 unique techs
	case Choson:
		switch baseTechID {
		case Toolworking:
			return Choson_cheap_ATK2
		case Metalworking:
			return Choson_cheap_ATK4
		case Metallurgy:
			return Choson_cheap_ATK7
		case Cataphract:
			return Choson_Cataphracts
		case LeatherArmorInfantry:
			return Choson_cheap_Armor2Inf
		case LeatherArmorArchers:
			return Choson_cheap_Armor2Arc
		case LeatherArmorCavalry:
			return Choson_cheap_Armor2Cav
		case ScaleArmorInfantry:
			return Choson_cheap_Armor4Inf
		case ScaleArmorArchers:
			return Choson_cheap_Armor4Arc
		case ScaleArmorCavalry:
			return Choson_cheap_Armor4Cav
		case BronzeShield:
			return Choson_cheap_Shield1
		}

	// Assyrian has 3 Siege Workshop techs -50% cost,
	// they are not requirements for other techs except themselves,
	// so total 3 unique techs
	case Assyrian:
		switch baseTechID {
		case MassiveCatapult:
			return Assyrian_cheap_MassiveCatapult
		case Catapult:
			return Assyrian_cheap_Catapult
		case Helepolis:
			return Assyrian_cheap_Helepolis
		}

	// Babylonian has 11 Market techs -30% cost,
	// Wheel is requirement for ChariotArcher, Chariot, ScytheChariot,
	// SiegeCraft is a requirement for MassiveCatapult,
	// so total 15 unique techs
	case Babylonian:
		switch baseTechID {
		case Wheel:
			return Babylonian_cheap_Wheel
		case Coinage:
			return Babylonian_cheap_Coinage
		case Plow:
			return Babylonian_cheap_Plow
		case Artisanship:
			return Babylonian_cheap_Artisanship
		case Irrigation:
			return Babylonian_cheap_Irrigation
		case Domestication:
			return Babylonian_cheap_Domestication
		case Woodworking:
			return Babylonian_cheap_WoodWorking
		case GoldMining:
			return Babylonian_cheap_GoldMining
		case StoneMining:
			return Babylonian_cheap_StoneMining
		case Craftsmanship:
			return Babylonian_cheap_Craftsmanship
		case Siegecraft:
			return Babylonian_cheap_SiegeCraft
		case EnableChariotArcher:
			return Babylonian_unlock_ChariotArc
		case EnableChariot:
			return Babylonian_unlock_Chariot
		case MassiveCatapult:
			return Babylonian_MassiveCatapult
		case ScytheChariot:
			return Babylonian_ScytheChariot
		}

	// Yamato has 3 Archery Range techs and 2 Stable techs -30% cost,
	// they are not requirements for other techs except themselves,
	// so total 5 unique techs
	case Yamato:
		switch baseTechID {
		case HeavyHorseArcher:
			return Yamato_cheap_HHorseArc
		case ImprovedBow:
			return Yamato_cheap_ImprovedBow
		case CompositeBow:
			return Yamato_cheap_CompositBow
		case HeavyCalvary:
			return Yamato_cheap_Cavalry2
		case Cataphract:
			return Yamato_cheap_Cataphracts
		}

	// Carthaginian has 1 Government Center tech -100% cost,
	// Nobility is requirement for ScytheChariot but is not available to Carthage anyway,
	// so total 1 unique tech
	case Carthaginian:
		switch baseTechID {
		case Nobility:
			return Carthaginian_cheap_Nobility
		}
	// Greek has 2 Temple techs -100% cost,
	// they are not requirements for other techs except themselves,
	// so total 2 unique techs
	case Greek:
		switch baseTechID {
		case Astrology:
			return Greek_cheap_Astrology
		case Polytheism:
			return Greek_cheap_Polytheism
		}

	// Hittite has 1 Market tech -50% cost and research -50% time,
	// Wheel is requirement for ChariotArcher, Chariot, ScytheChariot,
	// so total 4 unique techs
	case Hittite:
		switch baseTechID {
		case Wheel:
			return Hittite_cheap_Wheel
		case EnableChariotArcher:
			return Hittite_unlock_ChariotArc
		case EnableChariot:
			return Hittite_unlock_Chariot
		case ScytheChariot:
			return Hittite_ScytheChariot
		}
	}
	return NullTech
}

// GetReplacementUniqueTechIfNeeded returns unique techID if needed,
// otherwise returns NullTech
func GetReplacementUniqueTechIfNeeded(civ CivilizationID, step Step) TechID {
	if step.Action == PrintSummary {
		return NullTech
	}
	target, err := step.determineUnitOrTech(step.UnitOrTechID.IntID())
	if err != nil {
		return NullTech
	}
	switch baseTech := target.(type) {
	case *Technology:
		return CheckUniqueTechID(baseTech.ID, civ)
	default:
		return NullTech
	}
}
