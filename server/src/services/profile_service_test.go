package services

import (
	"server/src/structs"
	"testing"
)

func TestBioDataIsComplete(t *testing.T) {
	tests := []struct {
		name     string
		bio      structs.BioData
		expected bool
	}{
		{
			name: "Complete bio with 5 fields",
			bio: structs.BioData{
				Interests:        []string{"coding", "gaming"},
				Hobbies:          []string{"hiking"},
				MusicPreferences: []string{"rock"},
				FoodPreferences:  []string{"pizza"},
				LookingFor:       "friendship",
			},
			expected: true,
		},
		{
			name: "Incomplete bio with 4 fields",
			bio: structs.BioData{
				Interests:        []string{"coding"},
				Hobbies:          []string{"hiking"},
				MusicPreferences: []string{"rock"},
				FoodPreferences:  []string{"pizza"},
			},
			expected: false,
		},
		{
			name: "Empty arrays should not count",
			bio: structs.BioData{
				Interests:        []string{},
				Hobbies:          []string{},
				MusicPreferences: []string{},
				FoodPreferences:  []string{},
				LookingFor:       "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bio.IsComplete()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
