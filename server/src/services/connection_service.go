package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	structs "server/src/structs"
)

// ConnectionService handles business logic for user connections
type ConnectionService struct {
	db *sql.DB
}

// NewConnectionService creates a new connection service
func NewConnectionService(db *sql.DB) *ConnectionService {
	return &ConnectionService{db: db}
}

// SendConnectionRequest sends a connection request from one user to another
func (s *ConnectionService) SendConnectionRequest(fromUserID, toUserID int) (*structs.ConnectionRequest, error) {
	// Validate users exist and are different
	if fromUserID == toUserID {
		return nil, errors.New("cannot send connection request to yourself")
	}

	if !s.usersExist(fromUserID, toUserID) {
		return nil, errors.New("one or both users do not exist")
	}

	// Check if both profiles are complete
	if !s.isProfileComplete(fromUserID) {
		return nil, errors.New("your profile must be complete to send connection requests")
	}

	if !s.isProfileComplete(toUserID) {
		return nil, errors.New("cannot send request to user with incomplete profile")
	}

	// Check if users are already connected
	if connected, _ := s.AreUsersConnected(fromUserID, toUserID); connected {
		return nil, errors.New("users are already connected")
	}

	// Try to insert - let database enforce uniqueness
	query := `
        INSERT INTO connection_requests (from_user_id, to_user_id, status)
        VALUES ($1, $2, $3)
        RETURNING id, from_user_id, to_user_id, status, created_at, updated_at
    `

	var request structs.ConnectionRequest
	err := s.db.QueryRow(
		query,
		fromUserID,
		toUserID,
		structs.ConnectionStatusPending,
	).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		// Check if it's a unique constraint violation
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			// Check if there's a request in the opposite direction
			existingRequest, _ := s.GetPendingRequestBetweenUsers(fromUserID, toUserID)
			if existingRequest != nil {
				if existingRequest.FromUserID == toUserID {
					// They sent us a request! Auto-accept and connect
					log.Printf("Auto-accepting mutual connection request between %d and %d",
						fromUserID, toUserID)

					// Start transaction to accept their request
					tx, err := s.db.Begin()
					if err != nil {
						return nil, err
					}
					defer tx.Rollback()

					// Delete their request
					_, err = tx.Exec("DELETE FROM connection_requests WHERE id = $1",
						existingRequest.ID)
					if err != nil {
						return nil, err
					}

					// Create connection
					user1, user2 := fromUserID, toUserID
					if user1 > user2 {
						user1, user2 = user2, user1
					}

					var conn structs.Connection
					err = tx.QueryRow(
						"INSERT INTO connections (user_1_id, user_2_id) VALUES ($1, $2) RETURNING id, user_1_id, user_2_id, created_at",
						user1, user2,
					).Scan(&conn.ID, &conn.User1ID, &conn.User2ID, &conn.CreatedAt)

					if err != nil {
						return nil, err
					}

					tx.Commit()

					// Return a special response indicating auto-connection
					return &structs.ConnectionRequest{
						ID:         0, // Special marker
						FromUserID: fromUserID,
						ToUserID:   toUserID,
						Status:     "auto_accepted",
						CreatedAt:  conn.CreatedAt,
					}, nil
				}
				return nil, errors.New("you already sent a connection request to this user")
			}
		}
		return nil, fmt.Errorf("error sending connection request: %w", err)
	}

	log.Printf("Connection request sent from user %d to user %d", fromUserID, toUserID)
	return &request, nil
}

// AcceptConnectionRequest accepts a pending connection request
func (s *ConnectionService) AcceptConnectionRequest(requestID, acceptingUserID int) (*structs.Connection, error) {
	// Get the request first
	request, err := s.GetConnectionRequest(requestID)
	if err != nil {
		return nil, err
	}

	// Verify the accepting user is the recipient
	if request.ToUserID != acceptingUserID {
		return nil, errors.New("only the recipient can accept connection requests")
	}

	// Verify request is still pending
	if request.Status != structs.ConnectionStatusPending {
		return nil, fmt.Errorf("connection request is not pending (current status: %s)", request.Status)
	}

	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// DELETE the request (not update)
	_, err = tx.Exec("DELETE FROM connection_requests WHERE id = $1", requestID)
	if err != nil {
		return nil, fmt.Errorf("error deleting connection request: %w", err)
	}

	// Create connection (ensuring consistent ordering)
	user1ID, user2ID := request.FromUserID, request.ToUserID
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	var connection structs.Connection
	err = tx.QueryRow(
		"INSERT INTO connections (user_1_id, user_2_id) VALUES ($1, $2) RETURNING id, user_1_id, user_2_id, created_at",
		user1ID,
		user2ID,
	).Scan(
		&connection.ID,
		&connection.User1ID,
		&connection.User2ID,
		&connection.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating connection: %w", err)
	}

	// Create chat for the connection (Person A's domain, but needed for connections)
	result, err := tx.Exec(
		"INSERT INTO chats (user_1_id, user_2_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		user1ID,
		user2ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating chat: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("Chat created between user %d and user %d", user1ID, user2ID)
	} else {
		log.Printf("Chat already exists between user %d and user %d", user1ID, user2ID)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	log.Printf("Connection request %d accepted, connection created between user %d and user %d",
		requestID, request.FromUserID, request.ToUserID)
	return &connection, nil
}

// RejectConnectionRequest rejects a pending connection request
func (s *ConnectionService) RejectConnectionRequest(requestID, rejectingUserID int) error {
	// Get the request first
	request, err := s.GetConnectionRequest(requestID)
	if err != nil {
		return err
	}

	// Verify the rejecting user is the recipient
	if request.ToUserID != rejectingUserID {
		return errors.New("only the recipient can reject connection requests")
	}

	// Verify request is still pending
	if request.Status != structs.ConnectionStatusPending {
		return fmt.Errorf("connection request is not pending (current status: %s)", request.Status)
	}

	// DELETE the request (not update)
	result, err := s.db.Exec(
		"DELETE FROM connection_requests WHERE id = $1 AND status = $2",
		requestID,
		structs.ConnectionStatusPending,
	)
	if err != nil {
		return fmt.Errorf("error rejecting connection request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("connection request not found or already processed")
	}

	log.Printf("Connection request %d rejected by user %d", requestID, rejectingUserID)
	return nil
}

// GetConnectionRequest retrieves a specific connection request
func (s *ConnectionService) GetConnectionRequest(requestID int) (*structs.ConnectionRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM connection_requests 
		WHERE id = $1
	`

	var request structs.ConnectionRequest
	err := s.db.QueryRow(query, requestID).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("connection request not found")
		}
		return nil, fmt.Errorf("error getting connection request: %w", err)
	}

	return &request, nil
}

// GetPendingRequestBetweenUsers checks for pending requests between two users
func (s *ConnectionService) GetPendingRequestBetweenUsers(user1ID, user2ID int) (*structs.ConnectionRequest, error) {
	query := `
		SELECT id, from_user_id, to_user_id, status, created_at, updated_at
		FROM connection_requests 
		WHERE status = $1 
		AND ((from_user_id = $2 AND to_user_id = $3) OR (from_user_id = $3 AND to_user_id = $2))
		LIMIT 1
	`

	var request structs.ConnectionRequest
	err := s.db.QueryRow(query, structs.ConnectionStatusPending, user1ID, user2ID).Scan(
		&request.ID,
		&request.FromUserID,
		&request.ToUserID,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error checking pending request: %w", err)
	}

	return &request, nil
}

// GetReceivedRequests gets pending requests received by a user
func (s *ConnectionService) GetReceivedRequests(userID int) ([]structs.ConnectionRequestWithUsers, error) {
	query := `
        SELECT 
            cr.id, cr.from_user_id, cr.to_user_id, cr.status, cr.created_at, cr.updated_at,
            p.id, p.user_id, p.first_name, p.last_name, p.username, p.about_me,
            p.profile_picture_url, p.location, p.is_complete, p.created_at, p.updated_at
        FROM connection_requests cr
        JOIN profiles p ON cr.from_user_id = p.user_id
        WHERE cr.status = $1 AND cr.to_user_id = $2
        ORDER BY cr.created_at DESC
    `

	rows, err := s.db.Query(query, structs.ConnectionStatusPending, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting received requests: %w", err)
	}
	defer rows.Close()

	var requests []structs.ConnectionRequestWithUsers
	for rows.Next() {
		var req structs.ConnectionRequestWithUsers
		req.FromProfile = &structs.Profile{} // Initialize pointer

		err := rows.Scan(
			// Connection request fields
			&req.ID,
			&req.FromUserID,
			&req.ToUserID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
			// From profile fields
			&req.FromProfile.ID,
			&req.FromProfile.UserID,
			&req.FromProfile.FirstName,
			&req.FromProfile.LastName,
			&req.FromProfile.Username,
			&req.FromProfile.AboutMe,
			&req.FromProfile.ProfilePictureURL,
			&req.FromProfile.Location,
			&req.FromProfile.IsComplete,
			&req.FromProfile.CreatedAt,
			&req.FromProfile.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning request: %w", err)
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// GetSentRequests gets pending requests sent by a user
func (s *ConnectionService) GetSentRequests(userID int) ([]structs.ConnectionRequestWithUsers, error) {
	query := `
        SELECT 
            cr.id, cr.from_user_id, cr.to_user_id, cr.status, cr.created_at, cr.updated_at,
            p.id, p.user_id, p.first_name, p.last_name, p.username, p.about_me,
            p.profile_picture_url, p.location, p.is_complete, p.created_at, p.updated_at
        FROM connection_requests cr
        JOIN profiles p ON cr.to_user_id = p.user_id
        WHERE cr.status = $1 AND cr.from_user_id = $2
        ORDER BY cr.created_at DESC
    `

	rows, err := s.db.Query(query, structs.ConnectionStatusPending, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting sent requests: %w", err)
	}
	defer rows.Close()

	var requests []structs.ConnectionRequestWithUsers
	for rows.Next() {
		var req structs.ConnectionRequestWithUsers
		req.ToProfile = &structs.Profile{} // Initialize pointer

		err := rows.Scan(
			// Connection request fields
			&req.ID,
			&req.FromUserID,
			&req.ToUserID,
			&req.Status,
			&req.CreatedAt,
			&req.UpdatedAt,
			// To profile fields
			&req.ToProfile.ID,
			&req.ToProfile.UserID,
			&req.ToProfile.FirstName,
			&req.ToProfile.LastName,
			&req.ToProfile.Username,
			&req.ToProfile.AboutMe,
			&req.ToProfile.ProfilePictureURL,
			&req.ToProfile.Location,
			&req.ToProfile.IsComplete,
			&req.ToProfile.CreatedAt,
			&req.ToProfile.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning request: %w", err)
		}

		requests = append(requests, req)
	}

	return requests, nil
}

// GetUserConnections gets all connections for a user
func (s *ConnectionService) GetUserConnections(userID int) ([]structs.ConnectionWithUsers, error) {
	query := `
        SELECT 
            c.id, c.user_1_id, c.user_2_id, c.created_at,
            p.id, p.user_id, p.first_name, p.last_name, p.username, p.about_me,
            p.profile_picture_url, p.location, p.is_complete, p.created_at, p.updated_at,
            COALESCE(ch.updated_at, c.created_at) as last_activity,
            (
                SELECT COUNT(*) 
                FROM messages m 
                WHERE m.chat_id = ch.id 
                AND m.sender_id != $1 
                AND m.is_read = FALSE
            ) as unread_count
        FROM connections c
        JOIN profiles p ON (
            CASE 
                WHEN c.user_1_id = $1 THEN c.user_2_id 
                ELSE c.user_1_id 
            END = p.user_id
        )
        LEFT JOIN chats ch ON (
            (ch.user_1_id = c.user_1_id AND ch.user_2_id = c.user_2_id) OR
            (ch.user_1_id = c.user_2_id AND ch.user_2_id = c.user_1_id)
        )
        WHERE c.user_1_id = $1 OR c.user_2_id = $1
        ORDER BY last_activity DESC
    `

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting connections: %w", err)
	}
	defer rows.Close()

	var connections []structs.ConnectionWithUsers
	for rows.Next() {
		var conn structs.ConnectionWithUsers
		conn.User2Profile = &structs.Profile{} // The "other" user
		var lastActivity interface{}           // Placeholder for scan

		err := rows.Scan(
			// Connection fields
			&conn.ID,
			&conn.User1ID,
			&conn.User2ID,
			&conn.CreatedAt,
			// Other user's profile
			&conn.User2Profile.ID,
			&conn.User2Profile.UserID,
			&conn.User2Profile.FirstName,
			&conn.User2Profile.LastName,
			&conn.User2Profile.Username,
			&conn.User2Profile.AboutMe,
			&conn.User2Profile.ProfilePictureURL,
			&conn.User2Profile.Location,
			&conn.User2Profile.IsComplete,
			&conn.User2Profile.CreatedAt,
			&conn.User2Profile.UpdatedAt,
			// Extra fields
			&lastActivity,
			&conn.UnreadCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning connection: %w", err)
		}

		connections = append(connections, conn)
	}

	return connections, nil
}

// AreUsersConnected checks if two users are connected
func (s *ConnectionService) AreUsersConnected(user1ID, user2ID int) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM connections 
			WHERE (user_1_id = $1 AND user_2_id = $2) 
			OR (user_1_id = $2 AND user_2_id = $1)
		)
	`

	var connected bool
	err := s.db.QueryRow(query, user1ID, user2ID).Scan(&connected)
	if err != nil {
		return false, fmt.Errorf("error checking connection: %w", err)
	}

	return connected, nil
}

// RemoveConnection removes a connection between two users
func (s *ConnectionService) RemoveConnection(user1ID, user2ID int) error {
	// Ensure consistent ordering
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	// Delete connection
	result, err := s.db.Exec(
		"DELETE FROM connections WHERE user_1_id = $1 AND user_2_id = $2",
		user1ID,
		user2ID,
	)
	if err != nil {
		return fmt.Errorf("error deleting connection: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("connection not found")
	}

	log.Printf("Connection removed between user %d and user %d", user1ID, user2ID)
	return nil
}

// Helper function to check if users exist
func (s *ConnectionService) usersExist(userIDs ...int) bool {
	for _, userID := range userIDs {
		var exists bool
		err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
		if err != nil || !exists {
			return false
		}
	}
	return true
}

// Helper function to check profile completeness
func (s *ConnectionService) isProfileComplete(userID int) bool {
	var isComplete bool
	err := s.db.QueryRow(
		"SELECT is_complete FROM profiles WHERE user_id = $1",
		userID,
	).Scan(&isComplete)

	return err == nil && isComplete
}

// GetConnectionCount returns the number of connections for a user
func (s *ConnectionService) GetConnectionCount(userID int) (int, error) {
	var count int
	err := s.db.QueryRow(`
        SELECT COUNT(*) FROM connections 
        WHERE user_1_id = $1 OR user_2_id = $1
    `, userID).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error counting connections: %w", err)
	}

	return count, nil
}
