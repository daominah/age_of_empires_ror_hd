package aoego

type Civilization struct {
	ID CivilizationID
}

// CivilizationID is enum
type CivilizationID int

// CivilizationID enum
const (
	AssyrianID     CivilizationID = 81
	BabylonianID   CivilizationID = 82
	CarthaginianID CivilizationID = 205
	ChosonID       CivilizationID = 91
	EgyptianID     CivilizationID = 83
	GreekID        CivilizationID = 84
	HittiteID      CivilizationID = 85
	MacedonianID   CivilizationID = 206
	MinoanID       CivilizationID = 86
	PalmyranID     CivilizationID = 207
	PersianID      CivilizationID = 87
	PhoenicianID   CivilizationID = 88
	RomanID        CivilizationID = 208
	ShangID        CivilizationID = 89
	SumerianID     CivilizationID = 90
	YamatoID       CivilizationID = 92
)
