package services

import (
	"database/sql"
	"math"
	"server/src/structs"
	"sort"
	"strings"
)

type RecommendationService struct {
	DB *sql.DB
}

func NewRecommendationService(db *sql.DB) *RecommendationService {
	return &RecommendationService{DB: db}
}

func (s *RecommendationService) getExcludedUserIDs(userID int) (map[int]bool, error) {
	excluded := make(map[int]bool)

	// 1. Exclude existing connections
	rows, err := s.DB.Query(`
		SELECT user_1_id, user_2_id FROM connections 
		WHERE user_1_id = $1 OR user_2_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u1, u2 int
		if err := rows.Scan(&u1, &u2); err != nil {
			continue
		}
		if u1 == userID {
			excluded[u2] = true
		} else {
			excluded[u1] = true
		}
	}

	// 2. Exclude pending connection requests (both sent and received)
	rowsReq, err := s.DB.Query(`
		SELECT from_user_id, to_user_id FROM connection_requests 
		WHERE (from_user_id = $1 OR to_user_id = $1) AND status = 'pending'
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rowsReq.Close()

	for rowsReq.Next() {
		var from, to int
		if err := rowsReq.Scan(&from, &to); err != nil {
			continue
		}
		if from == userID {
			excluded[to] = true
		} else {
			excluded[from] = true
		}
	}

	// 3. Exclude dismissed recommendations
	rowsDis, err := s.DB.Query(`
		SELECT dismissed_user_id FROM dismissed_recommendations 
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rowsDis.Close()

	for rowsDis.Next() {
		var dismissedID int
		if err := rowsDis.Scan(&dismissedID); err != nil {
			continue
		}
		excluded[dismissedID] = true
	}

	return excluded, nil
}

func (s *RecommendationService) GetRecommendations(userID int) ([]structs.Recommendation, error) {
	// 1. Fetch current user's profile and bio
	// Note: We need to handle potential nulls if bio doesn't exist, but for now assuming it does or we handle it
	// For simplicity, we'll fetch basic data.
	// A better approach is to use the existing ProfileService, but we need raw access for efficiency.

	// Let's use a simplified struct for internal calculation
	type UserData struct {
		ID             int
		Lat, Lon       *float64
		Interests      []string
		Hobbies        []string
		Music          []string
		Food           []string
		AgeMin, AgeMax *int
		MaxDist        *int
	}

	rows, err := s.DB.Query(`
		SELECT u.id, p.latitude, p.longitude, 
		       COALESCE(b.interests, '{}'), COALESCE(b.hobbies, '{}'), 
		       COALESCE(b.music_preferences, '{}'), COALESCE(b.food_preferences, '{}'),
		       b.age_range_min, b.age_range_max, b.max_distance_km
		FROM users u
		JOIN profiles p ON u.id = p.user_id
		LEFT JOIN bio_data b ON u.id = b.user_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allUsers []UserData
	for rows.Next() {
		var u UserData

		// Let's use the structs.BioData fields which are pq.StringArray
		var bio structs.BioData
		var lat, lon *float64
		var id int

		if err := rows.Scan(&id, &lat, &lon, &bio.Interests, &bio.Hobbies, &bio.MusicPreferences, &bio.FoodPreferences, &bio.AgeRangeMin, &bio.AgeRangeMax, &bio.MaxDistanceKm); err != nil {
			continue // Skip bad rows
		}

		u = UserData{
			ID:        id,
			Lat:       lat,
			Lon:       lon,
			Interests: bio.Interests,
			Hobbies:   bio.Hobbies,
			Music:     bio.MusicPreferences,
			Food:      bio.FoodPreferences,
			AgeMin:    bio.AgeRangeMin,
			AgeMax:    bio.AgeRangeMax,
			MaxDist:   bio.MaxDistanceKm,
		}
		allUsers = append(allUsers, u)
	}

	// Find me
	var myData *UserData
	for i := range allUsers {
		if allUsers[i].ID == userID {
			myData = &allUsers[i]
			break
		}
	}

	if myData == nil {
		return nil, nil // User not found
	}

	// Filter and Score
	var recommendations []structs.Recommendation

	excluded, err := s.getExcludedUserIDs(userID)
	if err != nil {
		return nil, err
	}

	for _, other := range allUsers {
		if other.ID == userID {
			continue
		}

		// Check if excluded
		if excluded[other.ID] {
			continue
		}

		// 1. Distance Filter
		// Only check distance if both users have coordinates
		if myData.Lat != nil && myData.Lon != nil && other.Lat != nil && other.Lon != nil {
			dist := calculateDistance(myData.Lat, myData.Lon, other.Lat, other.Lon)
			maxDist := 100 // Default
			if myData.MaxDist != nil {
				maxDist = *myData.MaxDist
			}
			if dist > float64(maxDist) {
				continue
			}
		}

		// 2. Score Calculation
		score := 0.0
		score += calculateOverlap(myData.Interests, other.Interests) * 3.0
		score += calculateOverlap(myData.Hobbies, other.Hobbies) * 2.0
		score += calculateOverlap(myData.Music, other.Music) * 1.0
		score += calculateOverlap(myData.Food, other.Food) * 1.0

		// Normalize score or keep raw? Raw is fine for sorting.

		// Fetch full profile for response
		// For now, let's construct a basic Recommendation object.
		// We need to fetch the full profile/bio to return it.

		// Normalize score to percentage (0-100)
		// Max possible score:
		// Interests (4 items * 3 pts) = 12
		// Hobbies (4 items * 2 pts) = 8
		// Music (3 items * 1 pt) = 3
		// Food (3 items * 1 pt) = 3
		// Total max = 26
		// Let's scale it roughly. If score is 13, that's 50%.

		// For now, let's just say if raw score < 5, skip it.
		maxPossibleScore := 26.0
		percentageScore := (score / maxPossibleScore) * 100.0
		if percentageScore > 100 {
			percentageScore = 100
		}

		// Filter: Only allow matches above 5%
		if percentageScore < 5 {
			continue
		}

		rec := structs.Recommendation{
			UserID:     other.ID,
			MatchScore: percentageScore,
		}
		recommendations = append(recommendations, rec)
	}

	// Sort by score descending
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].MatchScore > recommendations[j].MatchScore
	})

	// Limit to top 10
	if len(recommendations) > 10 {
		recommendations = recommendations[:10]
	}

	// Now populate the full details for these top users
	var recommendedUserIDs []int
	for i := range recommendations {
		// a helper query here to be self-contained
		s.fillRecommendationDetails(&recommendations[i])
		recommendedUserIDs = append(recommendedUserIDs, recommendations[i].UserID)
	}

	// Mark as viewed
	if len(recommendedUserIDs) > 0 {
		go s.markRecommendationsAsViewed(userID, recommendedUserIDs)
	}

	return recommendations, nil
}

func (s *RecommendationService) markRecommendationsAsViewed(userID int, viewedIDs []int) {
	if len(viewedIDs) == 0 {
		return
	}

	// Simple loop for MVP. Bulk insert would be better.
	for _, viewedID := range viewedIDs {
		_, err := s.DB.Exec(`
			INSERT INTO viewed_recommendations (user_id, viewed_user_id) 
			VALUES ($1, $2) 
			ON CONFLICT (user_id, viewed_user_id) DO NOTHING
		`, userID, viewedID)
		if err != nil {
			// Log error but don't stop
			// log.Printf("Error marking recommendation as viewed: %v", err)
		}
	}
}

func (s *RecommendationService) DismissRecommendation(userID, dismissedUserID int) error {
	_, err := s.DB.Exec(`
		INSERT INTO dismissed_recommendations (user_id, dismissed_user_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, dismissed_user_id) DO NOTHING
	`, userID, dismissedUserID)
	return err
}

func (s *RecommendationService) fillRecommendationDetails(rec *structs.Recommendation) {
	// Fetch Profile
	query := `SELECT id, user_id, first_name, last_name, username, about_me, profile_picture_url, location, latitude, longitude, is_complete 
	          FROM profiles WHERE user_id = $1`
	s.DB.QueryRow(query, rec.UserID).Scan(
		&rec.Profile.ID, &rec.Profile.UserID, &rec.Profile.FirstName, &rec.Profile.LastName,
		&rec.Profile.Username, &rec.Profile.AboutMe, &rec.Profile.ProfilePictureURL,
		&rec.Profile.Location, &rec.Profile.Latitude, &rec.Profile.Longitude, &rec.Profile.IsComplete,
	)

	// Fetch Bio
	queryBio := `SELECT id, user_id, interests, hobbies, music_preferences, food_preferences, looking_for, age_range_min, age_range_max, max_distance_km 
	             FROM bio_data WHERE user_id = $1`
	s.DB.QueryRow(queryBio, rec.UserID).Scan(
		&rec.BioData.ID, &rec.BioData.UserID, &rec.BioData.Interests, &rec.BioData.Hobbies,
		&rec.BioData.MusicPreferences, &rec.BioData.FoodPreferences, &rec.BioData.LookingFor,
		&rec.BioData.AgeRangeMin, &rec.BioData.AgeRangeMax, &rec.BioData.MaxDistanceKm,
	)
}

func calculateDistance(lat1, lon1, lat2, lon2 *float64) float64 {
	if lat1 == nil || lon1 == nil || lat2 == nil || lon2 == nil {
		return 999999 // Infinite distance if coords missing
	}

	// Haversine formula
	const R = 6371 // Earth radius in km
	dLat := (*lat2 - *lat1) * (math.Pi / 180)
	dLon := (*lon2 - *lon1) * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(*lat1*(math.Pi/180))*math.Cos(*lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func calculateOverlap(list1, list2 []string) float64 {
	count := 0
	// Simple O(N*M) is fine for small lists
	for _, item1 := range list1 {
		for _, item2 := range list2 {
			if strings.EqualFold(strings.TrimSpace(item1), strings.TrimSpace(item2)) {
				count++
				break
			}
		}
	}
	return float64(count)
}
