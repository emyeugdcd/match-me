package services

import (
	"database/sql"
	"errors"
	"math"

	structs "server/src/structs"
)

// LocationService handles location-based candidate filtering
// (same-city and proximity-based radius).
type LocationService struct {
	db *sql.DB
}

func NewLocationService(db *sql.DB) *LocationService {
	return &LocationService{db: db}
}

// Light struct used by recommendation logic
type LocationCandidate struct {
	UserID    int
	Location  string
	Latitude  *float64
	Longitude *float64
}

// internal helper: load profile + bio for a user
func (s *LocationService) getUserLocationData(userID int) (*structs.Profile, *structs.BioData, error) {
	var p structs.Profile
	err := s.db.QueryRow(`
	  SELECT id, user_id, first_name, last_name, username, about_me,
	         profile_picture_url, location, latitude, longitude,
	         is_complete, created_at, updated_at
	  FROM profiles
	  WHERE user_id = $1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.FirstName, &p.LastName, &p.Username,
		&p.AboutMe, &p.ProfilePictureURL, &p.Location,
		&p.Latitude, &p.Longitude,
		&p.IsComplete, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	var b structs.BioData
	err = s.db.QueryRow(`
	  SELECT id, user_id, interests, hobbies, music_preferences,
	         food_preferences, looking_for, age_range_min, age_range_max,
	         max_distance_km, additional_data, created_at, updated_at
	  FROM bio_data
	  WHERE user_id = $1
	`, userID).Scan(
		&b.ID, &b.UserID, &b.Interests, &b.Hobbies,
		&b.MusicPreferences, &b.FoodPreferences,
		&b.LookingFor, &b.AgeRangeMin, &b.AgeRangeMax,
		&b.MaxDistanceKm, &b.AdditionalData,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	return &p, &b, nil
}

// GetCityCandidates returns all complete profiles in the same city.
func (s *LocationService) GetCityCandidates(userID int, city string) ([]LocationCandidate, error) {
	rows, err := s.db.Query(`
	  SELECT p.user_id, p.location, p.latitude, p.longitude
	  FROM profiles p
	  JOIN bio_data b ON b.user_id = p.user_id
	  WHERE p.user_id != $1
	    AND p.is_complete = TRUE
	    AND b.id IS NOT NULL
	    AND p.location = $2
	`, userID, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []LocationCandidate
	for rows.Next() {
		var c LocationCandidate
		if err := rows.Scan(&c.UserID, &c.Location, &c.Latitude, &c.Longitude); err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

// GetRadiusCandidates filters all complete profiles by distance (km) using Haversine.
func (s *LocationService) GetRadiusCandidates(userID int, lat, lon float64, maxKm int) ([]LocationCandidate, error) {
	rows, err := s.db.Query(`
	  SELECT p.user_id, p.location, p.latitude, p.longitude
	  FROM profiles p
	  JOIN bio_data b ON b.user_id = p.user_id
	  WHERE p.user_id != $1
	    AND p.is_complete = TRUE
	    AND b.id IS NOT NULL
	    AND p.latitude IS NOT NULL
	    AND p.longitude IS NOT NULL
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []LocationCandidate
	for rows.Next() {
		var c LocationCandidate
		if err := rows.Scan(&c.UserID, &c.Location, &c.Latitude, &c.Longitude); err != nil {
			return nil, err
		}
		if c.Latitude == nil || c.Longitude == nil {
			continue
		}
		d := haversineKm(lat, lon, *c.Latitude, *c.Longitude)
		if d <= float64(maxKm) {
			res = append(res, c)
		}
	}
	return res, nil
}

// GetLocationCandidatesForUser decides which logic to use:
//
// - If user has latitude+longitude and max_distance_km > 0 -> radius / proximity-based
// - Otherwise -> simple same-city matching
func (s *LocationService) GetLocationCandidatesForUser(userID int) ([]LocationCandidate, error) {
	profile, bio, err := s.getUserLocationData(userID)
	if err != nil {
		return nil, err
	}
	if profile.Location == "" {
		return nil, errors.New("user has no location set")
	}

	if profile.Latitude != nil && profile.Longitude != nil &&
		bio.MaxDistanceKm != nil && *bio.MaxDistanceKm > 0 {
		return s.GetRadiusCandidates(userID, *profile.Latitude, *profile.Longitude, *bio.MaxDistanceKm)
	}

	return s.GetCityCandidates(userID, profile.Location)
}

// Haversine distance in km (good enough accuracy for this app)
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0 // Earth radius in kilometers

	toRad := func(d float64) float64 { return d * math.Pi / 180 }

	φ1 := toRad(lat1)
	φ2 := toRad(lat2)
	dφ := toRad(lat2 - lat1)
	dλ := toRad(lon2 - lon1)

	a := math.Sin(dφ/2)*math.Sin(dφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(dλ/2)*math.Sin(dλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
