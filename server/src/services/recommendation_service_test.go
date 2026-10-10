package services

import (
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	// Test nil coordinates
	if dist := calculateDistance(nil, nil, nil, nil); dist != 999999 {
		t.Errorf("expected 999999 for nil coordinates, got %f", dist)
	}

	lat1 := 60.1699
	lon1 := 24.9384
	// Same coordinate -> distance should be 0
	if dist := calculateDistance(&lat1, &lon1, &lat1, &lon1); math.Abs(dist) > 0.001 {
		t.Errorf("expected 0 for identical coordinates, got %f", dist)
	}

	// Helsinki (60.1699, 24.9384) to Espoo (60.2055, 24.6559) ~ 16-17 km
	lat2 := 60.2055
	lon2 := 24.6559
	dist := calculateDistance(&lat1, &lon1, &lat2, &lon2)
	if dist < 14 || dist > 18 {
		t.Errorf("expected distance between Helsinki and Espoo to be ~16km, got %f", dist)
	}
}

func TestCalculateOverlap(t *testing.T) {
	tests := []struct {
		name     string
		list1    []string
		list2    []string
		expected float64
	}{
		{
			name:     "Empty lists",
			list1:    []string{},
			list2:    []string{},
			expected: 0.0,
		},
		{
			name:     "No overlap",
			list1:    []string{"Go", "Kubernetes"},
			list2:    []string{"Python", "Java"},
			expected: 0.0,
		},
		{
			name:     "Exact matches",
			list1:    []string{"Go", "Docker", "DevOps"},
			list2:    []string{"Docker", "Go", "Azure"},
			expected: 2.0,
		},
		{
			name:     "Case-insensitive with whitespace",
			list1:    []string{"  Golang ", "KUBERNETES"},
			list2:    []string{"golang", "kubernetes  "},
			expected: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateOverlap(tt.list1, tt.list2)
			if result != tt.expected {
				t.Errorf("calculateOverlap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
