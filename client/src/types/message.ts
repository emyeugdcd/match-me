export interface Message {
  id: number;
  sender_id: number;
  receiver_id: number;
  content: string;
  created_at: string;
}

export interface ChatPreview {
  user_id: number;
  display_name: string;
  profile_picture_url?: string;
  last_message: string;
  last_message_time: string;
  unread_count: number;
}
