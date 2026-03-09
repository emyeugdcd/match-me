package services

import (
	"database/sql"
	"server/src/structs"
)

type ChatService struct {
	DB *sql.DB
}

func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{DB: db}
}

func (s *ChatService) getOrCreateChatID(user1ID, user2ID int) (int, error) {
	// Ensure consistent ordering for query
	u1, u2 := user1ID, user2ID
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	var chatID int
	query := `SELECT id FROM chats WHERE user_1_id = $1 AND user_2_id = $2`
	err := s.DB.QueryRow(query, u1, u2).Scan(&chatID)
	if err == sql.ErrNoRows {
		// Create new chat
		create := `INSERT INTO chats (user_1_id, user_2_id) VALUES ($1, $2) RETURNING id`
		err = s.DB.QueryRow(create, u1, u2).Scan(&chatID)
		if err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	return chatID, nil
}

func (s *ChatService) SendMessage(senderID, receiverID int, content string) (*structs.Message, error) {
	chatID, err := s.getOrCreateChatID(senderID, receiverID)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO messages (chat_id, sender_id, message_text) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, is_read`

	var msg structs.Message
	msg.ChatID = chatID
	msg.SenderID = senderID
	msg.ReceiverID = receiverID
	msg.Content = content

	err = s.DB.QueryRow(query, chatID, senderID, content).Scan(&msg.ID, &msg.CreatedAt, &msg.IsRead)
	if err != nil {
		return nil, err
	}

	// Update chat updated_at
	updateChat := `UPDATE chats SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, _ = s.DB.Exec(updateChat, chatID)

	return &msg, nil
}

func (s *ChatService) GetMessages(userID1, userID2 int, limit, offset int) ([]structs.Message, error) {
	query := `
		SELECT m.id, m.chat_id, m.sender_id, m.message_text, m.created_at, m.is_read
		FROM messages m
		JOIN chats c ON m.chat_id = c.id
		WHERE (c.user_1_id = $1 AND c.user_2_id = $2) 
		   OR (c.user_1_id = $2 AND c.user_2_id = $1)
		ORDER BY m.created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := s.DB.Query(query, userID1, userID2, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []structs.Message{}
	for rows.Next() {
		var msg structs.Message
		if err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.IsRead); err != nil {
			return nil, err
		}

		// Set computed ReceiverID
		if msg.SenderID == userID1 {
			msg.ReceiverID = userID2
		} else {
			msg.ReceiverID = userID1
		}

		messages = append(messages, msg)
	}

	// Reverse messages to show oldest first (since we fetched newest first for pagination)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	// Mark messages as read asynchronously
	go s.MarkMessagesAsRead(userID1, userID2)

	return messages, nil
}

func (s *ChatService) MarkMessagesAsRead(readerID, senderID int) error {
	query := `
		UPDATE messages m
		SET is_read = true
		FROM chats c
		WHERE m.chat_id = c.id
		AND (c.user_1_id = $1 AND c.user_2_id = $2 OR c.user_1_id = $2 AND c.user_2_id = $1)
		AND m.sender_id = $2
		AND m.is_read = false
	`
	_, err := s.DB.Exec(query, readerID, senderID)
	return err
}

func (s *ChatService) GetChats(userID int) ([]structs.ChatPreview, error) {
	query := `
		SELECT 
			c.id,
			CASE WHEN c.user_1_id = $1 THEN c.user_2_id ELSE c.user_1_id END as other_user_id,
			p.first_name || ' ' || p.last_name as display_name,
			p.profile_picture_url,
			(SELECT message_text FROM messages m WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) as last_message,
			(SELECT created_at FROM messages m WHERE m.chat_id = c.id ORDER BY m.created_at DESC LIMIT 1) as last_message_time,
			(SELECT COUNT(*) FROM messages m WHERE m.chat_id = c.id AND m.sender_id != $1 AND m.is_read = false) as unread_count
		FROM chats c
		JOIN profiles p ON p.user_id = (CASE WHEN c.user_1_id = $1 THEN c.user_2_id ELSE c.user_1_id END)
		WHERE c.user_1_id = $1 OR c.user_2_id = $1
		ORDER BY last_message_time DESC NULLS LAST
	`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []structs.ChatPreview
	for rows.Next() {
		var chat structs.ChatPreview
		var chatID int
		var lastMessage sql.NullString
		var lastMessageTime sql.NullTime

		if err := rows.Scan(&chatID, &chat.UserID, &chat.DisplayName, &chat.ProfilePictureURL, &lastMessage, &lastMessageTime, &chat.UnreadCount); err != nil {
			return nil, err
		}

		if lastMessage.Valid {
			chat.LastMessage = lastMessage.String
		}
		if lastMessageTime.Valid {
			chat.LastMessageTime = lastMessageTime.Time
		}

		chats = append(chats, chat)
	}

	return chats, nil
}
