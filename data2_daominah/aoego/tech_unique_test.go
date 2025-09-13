package aoego

import "testing"

func TestCheckUniqueTechID(t *testing.T) {
	tests := []struct {
		baseTechID     TechID
		civID          CivilizationID
		expectedUnique TechID
	}{
		{Toolworking, Choson, Choson_cheap_ATK2},
		{Toolworking, Babylonian, NullTech},
		{Cataphract, Choson, Choson_Cataphracts},
		{Catapult, Assyrian, Assyrian_cheap_Catapult},
		{MassiveCatapult, Greek, NullTech},
		{Wheel, Babylonian, Babylonian_cheap_Wheel},
		{EnableChariot, Babylonian, Babylonian_unlock_Chariot},
		{ScytheChariot, Babylonian, Babylonian_ScytheChariot},
		{ScytheChariot, Egyptian, NullTech},
	}

	for _, tt := range tests {
		got := CheckUniqueTechID(tt.baseTechID, tt.civID)
		if got != tt.expectedUnique {
			t.Errorf("CheckUniqueTechID(%v, %v): got %v, want %v", tt.baseTechID, tt.civID, got, tt.expectedUnique)
		}
	}
}
