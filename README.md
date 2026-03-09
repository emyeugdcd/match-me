# Web Social App

This project contains a Go server and a React client. This is a full-stack application, which matches user profiles based on your own recommendation algorithm. It allows users to register, view potential matches, connect with other coders and chat. It is built on the concept of a "man cave", and is intended to be viewed as Tinder for Coders.

## Registration
It enables users to sign up with a unique email address, and a password. The password will be securely protected. Users can then log in with their email address and password.Sessions are managed with a JWT, and user can log out from any page in the application.

## User Profile
Users must complete their profile before they are able to see recommendations or connect with other users. Users will be matched based on their profiles, and the profile must be completed before any recommendations are made.

The contents of the bio are up to you. You may choose to have a first and last name if it is for professional connections, or a username if it is for coders or gamers. You may also chose to have some kind of "About me", again this is up to you.

You'll need to capture biographical information about yourself before being able to match with other coders. What are your interests, hobbies or favorite coding languages? What food or music do you like? Do they like long walks? Do they like interplanetary travel?

You will be able to add/remove/change your profile picture. If no picture is uploaded, a placeholder image will be shown in its place "👤". But worry not, you are able to modify your profile information at any time.

## Recommendations
Users will be able to see a list of recommended connections. This list will be prioritized to show the strongest matches first (and also closest possible in terms of locations). The user can get a maximum of 10 recommendations at a time.

## Connections
When a user sees a recommendation that they find interesting, they can request to connect with them. Users then will be able to see a list of connection requests, where they can accept or dismiss requests. It will.be possible to disconnect with a user, if they are no longer interesting.

## Chat
Once two users are connected, they are able to chat with one another.
From the connected user's profile, it must be possible to start or resume a chat. Only one chat history can exist between two users.
The chat history will be visible to both users. Users will be able to see a list of all of their chats, with the most recent chats appearing first.

## Running the Project

You can use the `Makefile` to run the project in various modes.

Make sure you have Docker installed and running on your machine.

### Run Everything (Docker + Local)
Starts Docker containers (DB, etc.) and runs both server and client locally.
```bash
make run
```

### Run Components Individually
**Server Only:**
```bash
make run-server
```

**Client Only:**
```bash
make run-client
```

### Run in Parallel
Starts Docker containers and instructs you to run server and client in separate terminals.
```bash
make run-parallel
```

## Seeding Data

To populate the database with fictitious users for testing, you can use the seeding endpoint.

**Endpoint:** `POST /dev/seed-users`

**Example using curl:**
```bash
# Seed 100 users (default)
curl -X POST http://localhost:8080/dev/seed-users

# Seed specific number of users (e.g., 50)
curl -X POST http://localhost:8080/dev/seed-users \
  -H "Content-Type: application/json" \
  -d '{"count": 50}'
```

The server must be running for this to work.

## Matching Logic
### recommendationService.go

### Exclusions
Existing connections are excluded.
Pending connection requests (sent or received) are excluded.
Dismissed recommendations are excluded.
### Distance Filtering
If both users have location coordinates (latitude/longitude), the distance between them is calculated using the Haversine formula.
If the distance exceeds the user's max_distance_km (defaulting to 100km if not set), the user is excluded.
### Scoring
A score is calculated based on overlapping items in profile lists:
Interests: 3 points per match.
Hobbies: 2 points per match.
Music Preferences: 1 point per match.
Food Preferences: 1 point per match.
Matching is case-insensitive.
### Normalization
The raw score is normalized to a percentage (0-100%).
The maximum possible score is assumed to be 26 points (based on a theoretical max of 4 interests, 4 hobbies, 3 music, 3 food items).
percentageScore = (rawScore / 26.0) * 100.0
Capped at 100%.
### Threshold
Matches with a percentage score below 5% are excluded.
### Sorting & Limiting
Recommendations are sorted by MatchScore in descending order.
The list is limited to the top 10 recommendations.
### Viewed Tracking
Returned recommendations are marked as "viewed" in the database.

Notes:
- Online presence: Users' online status is shown in profile and chat views via WebSocket realtime updates. Ensure your browser is connected to the server WebSocket for presence to work.
