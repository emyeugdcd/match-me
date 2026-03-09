package services

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	structs "server/src/structs"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type SeedUsers struct {
	db             *sql.DB
	profileService *ProfileService
}

func NewSeedUsers(db *sql.DB, profileService *ProfileService) *SeedUsers {
	return &SeedUsers{
		db:             db,
		profileService: profileService,
	}
}

// SeedUsers ensures there are at least `target` users in the system.
// It creates new fictitious users (with profile + bio) until the target is reached.
func (s *SeedUsers) SeedUsers(target int) (int, error) {
	if target <= 0 {
		target = 100
	}

	var existing int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&existing); err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	toCreate := target - existing
	if toCreate <= 0 {
		log.Printf("SeedUsers: already have %d users, nothing to create\n", existing)
		return 0, nil
	}

	log.Printf("SeedUsers: seeding %d fictitious users (target=%d, existing=%d)\n", toCreate, target, existing)

	rand.Seed(time.Now().UnixNano())
	created := 0

	for i := 0; i < toCreate; i++ {
		// --- 1. Generate base data ---
		firstName, lastName := randomName()
		username := randomUsername(firstName, lastName, i)
		email := fmt.Sprintf(
			"%s.%s+seed%d@example.com",
			strings.ToLower(firstName),
			strings.ToLower(lastName),
			time.Now().UnixNano()+int64(i),
		)

		// --- 2. Insert into users table ---
		passwordHash, err := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("SeedUsers: failed to hash password: %v\n", err)
			continue
		}

		var userID int
		err = s.db.QueryRow(
			"INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id",
			email, string(passwordHash),
		).Scan(&userID)
		if err != nil {
			// Skip duplicates or other issues and continue
			if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
				log.Printf("SeedUsers: duplicate email %s, skipping\n", email)
				continue
			}
			log.Printf("SeedUsers: failed to insert user: %v\n", err)
			continue
		}

		// --- 3. Create profile (so profile table has a row) ---
		city, lat, lon := randomLocation()
		profile := structs.Profile{
			UserID:    userID,
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			AboutMe:   "Fictitious user account for testing and development.",
			Location:  city,
			Latitude:  &lat,
			Longitude: &lon,
		}

		if _, err := s.profileService.CreateOrUpdateProfile(userID, profile); err != nil {
			log.Printf("SeedUsers: failed to create profile for user %d: %v\n", userID, err)
			continue
		}

		// --- 4. Create BioData (used by recommendation logic) ---
		bio := randomBioData(userID)
		if _, err := s.profileService.CreateBioData(userID, bio); err != nil {
			log.Printf("SeedUsers: failed to create bio data for user %d: %v\n", userID, err)
			continue
		}

		created++
	}

	log.Printf("SeedUsers: successfully created %d new fictitious users\n", created)
	return created, nil
}

// ---- helper data ----

var firstNames = []string{
	"Alex", "Sam", "Mira", "Leo", "Niko", "Sofia",
	"Liam", "Emma", "Anna", "Matti", "Laura", "Oona",
}

var lastNames = []string{
	"Korhonen", "Virtanen", "Niemi", "Laine", "Heikkinen", "Kallio",
	"Järvinen", "Mäkinen", "Hämäläinen", "Koski",
}

var locations = []string{
	"Helsinki", "Espoo", "Vantaa", "Tampere", "Turku",
	"Oulu", "Jyväskylä", "Lahti", "Kuopio",
}

var interestsOptions = []string{
	"coding", "hiking", "music", "board games", "fitness", "cooking",
	"reading", "travelling", "movies", "photography",
}

var hobbiesOptions = []string{
	"Python", "Go", "JavaScript", "TypeScript", "Rust", "Java",
	"C++", "C#", "Swift", "Kotlin", "PHP", "Ruby", "Scala",
	"Haskell", "Lua", "Dart", "Elixir", "Clojure",
}

var musicOptions = []string{
	"rock", "pop", "jazz", "electronic", "hip hop", "classical",
}

var foodOptions = []string{
	"pizza", "sushi", "burger", "vegan", "pasta", "tacos", "salad",
}

var lookingForOptions = []string{
	"friendship", "networking", "study buddy", "sports buddy", "coffee chat",
}

// ---- helper functions ----

func randomName() (string, string) {
	return firstNames[rand.Intn(len(firstNames))],
		lastNames[rand.Intn(len(lastNames))]
}

func randomUsername(first, last string, i int) string {
	return fmt.Sprintf(
		"%s_%s_%d",
		strings.ToLower(first),
		strings.ToLower(last),
		i+rand.Intn(1000),
	)
}

var cityCoords = map[string]struct{ Lat, Lon float64 }{
	"Helsinki":  {60.1699, 24.9384},
	"Espoo":     {60.2055, 24.6559},
	"Vantaa":    {60.2934, 25.0378},
	"Tampere":   {61.4978, 23.7610},
	"Turku":     {60.4518, 22.2666},
	"Oulu":      {65.0121, 25.4651},
	"Jyväskylä": {62.2426, 25.7473},
	"Lahti":     {60.9827, 25.6612},
	"Kuopio":    {62.8924, 27.6770},
}

func randomLocation() (string, float64, float64) {
	idx := rand.Intn(len(locations))
	city := locations[idx]
	coords := cityCoords[city]
	return city, coords.Lat, coords.Lon
}

func pickSome(options []string, min, max int) pq.StringArray {
	if min < 1 {
		min = 1
	}
	if max < min {
		max = min
	}
	n := min + rand.Intn(max-min+1)
	if n > len(options) {
		n = len(options)
	}

	shuffled := append([]string{}, options...)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return pq.StringArray(shuffled[:n])
}

func randomBioData(userID int) structs.BioData {
	minAge := 20 + rand.Intn(15)        // 20–34
	maxAge := minAge + 3 + rand.Intn(7) // +3..+9

	return structs.BioData{
		UserID:           userID,
		Interests:        pickSome(interestsOptions, 2, 4),
		Hobbies:          pickSome(hobbiesOptions, 2, 4),
		MusicPreferences: pickSome(musicOptions, 1, 3),
		FoodPreferences:  pickSome(foodOptions, 1, 3),
		LookingFor:       lookingForOptions[rand.Intn(len(lookingForOptions))],
		AgeRangeMin:      &minAge,
		AgeRangeMax:      &maxAge,
	}
}
