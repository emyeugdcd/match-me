// types/index.ts
// Centralized type definitions for the entire application
// This keeps types consistent across all components

export interface User {
  id: number;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
  message: string;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignupCredentials {
  email: string;
  password: string;
}

export interface Profile {
  id: number;
  user_id: number;
  first_name: string;
  last_name: string;
  username: string;
  about_me: string;
  profile_picture_url: string;
  location: string;
  latitude: number | null;
  longitude: number | null;
  is_complete: boolean;
  created_at: string;
  updated_at: string;
}

export interface BioData {
  id: number;
  user_id: number;
  interests: string[];
  hobbies: string[];
  music_preferences: string[];
  food_preferences: string[];
  looking_for: string;
  age_range_min?: number;
  age_range_max?: number;
  max_distance_km?: number;
  additional_data?: Record<string, unknown>;
  created_at: string;
  updated_at: string;
}

// API Response types 
export interface UserBasicInfo {
  id: number;
  name: string;  // Display name
  first_name: string;
  last_name: string;
  username?: string;
  email?: string;
  profile_picture_url: string;
}

export interface UserProfileInfo {
  id: number;
  about_me: string;
  location: string;
  latitude?: number | null;
  longitude?: number | null;
}

export interface UserBioInfo {
  id: number;
  interests: string[];
  hobbies: string[];
  music_preferences: string[];
  food_preferences: string[];
  looking_for: string;
  age_range_min?: number;
  age_range_max?: number;
  max_distance_km?: number;
}

// Complete user profile (combination of all three)
export interface CompleteUserProfile {
  basic: UserBasicInfo;
  profile: UserProfileInfo;
  bio: UserBioInfo;
}

// Request types for creating/updating
export interface ProfileUpdateRequest {
  first_name?: string;
  last_name?: string;
  username?: string;
  about_me?: string;
  location?: string;
  profile_picture_url?: string;
  latitude?: number | null;
  longitude?: number | null;
}

export interface BioDataUpdateRequest {
  interests: string[];
  hobbies: string[];
  music_preferences: string[];
  food_preferences: string[];
  looking_for: string;
  age_range_min?: number;
  age_range_max?: number;
  max_distance_km?: number;
  additional_data?: Record<string, unknown>;
}

export interface ConnectionRequest {
  id: number;
  from_user_id: number;
  display_name: string;
  profile_picture_url: string;
  created_at: string;
}

export interface SentConnectionRequest {
  id: number;
  to_user_id: number;
  display_name: string;
  profile_picture_url: string;
  created_at: string;
}

export interface Connection {
  id: number;
  user_1_id: number;
  user_2_id: number;
  created_at: string;
}

export interface DetailedConnection {
  user_id: number;
  display_name: string;
  profile_picture_url: string;
  location: string;
  connected_at: string;
}

export interface ConnectionsResponse {
  connections: number[];
}

export interface ConnectionStatus {
  status: 'none' | 'connected' | 'request_sent' | 'request_received';
  connected: boolean;
  request_id?: number;
}

// Action responses
export interface SendRequestResponse {
  message: string;
  request_id?: number;
  status: string;
  created_at?: string;
  connected_at?: string;
}

export interface AcceptRequestResponse {
  message: string;
  connection_id: number;
  connected_user_id: number;
  connected_at: string;
}

export interface RecommendationResponse {
  recommendations: number[];
}

export interface Recommendation {
  user_id: number;
  profile: Profile;
  bio_data: BioData;
  match_score: number;
  is_viewed: boolean;
  is_dismissed: boolean;
}

export interface Chat {
  id: number;
  user_1_id: number;
  user_2_id: number;
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: number;
  chat_id: number;
  sender_id: number;
  message_text: string;
  is_read: boolean;
  created_at: string;
}

export interface ChatWithDetails {
  id: number;
  user_1_id: number;
  user_2_id: number;
  other_user: UserBasicInfo;
  last_message?: Message;
  unread_count: number;
  created_at: string;
  updated_at: string;
}

export interface LoadingState {
  isLoading: boolean;
  error: string | null;
}

export interface PaginationState {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}

export interface FormError {
  field: string;
  message: string;
}

export interface ProfileFormState {
  first_name: string;
  last_name: string;
  username: string;
  about_me: string;
  location: string;
  profile_picture_url: string;
  latitude: number | null;
  longitude: number | null;
}

export interface BioFormState {
  interests: string[];
  hobbies: string[];
  music_preferences: string[];
  food_preferences: string[];
  looking_for: string;
  age_range_min: number | undefined;
  age_range_max: number | undefined;
  max_distance_km: number | undefined;
}

export interface ApiError {
  error: string;
  message?: string;
  status?: number;
}

export interface UserCardProps {
  userId: number;
  name: string;
  profilePictureUrl: string;
  location?: string;
  onClick?: () => void;
}

export interface ConnectionRequestCardProps {
  request: ConnectionRequest;
  onAccept: (requestId: number) => void;
  onDismiss: (requestId: number) => void;
}

export interface ConnectedUserCardProps {
  connection: DetailedConnection;
  onDisconnect: (userId: number) => void;
  onViewProfile: (userId: number) => void;
  onStartChat: (userId: number) => void;
}

export type AsyncStatus = 'idle' | 'loading' | 'success' | 'error';

export interface AsyncState<T> {
  status: AsyncStatus;
  data: T | null;
  error: string | null;
}

// Helper to create initial async state
export const createAsyncState = <T>(): AsyncState<T> => ({
  status: 'idle',
  data: null,
  error: null,
});

export const CONNECTION_STATUS = {
  NONE: 'none' as const,
  CONNECTED: 'connected' as const,
  REQUEST_SENT: 'request_sent' as const,
  REQUEST_RECEIVED: 'request_received' as const,
};

export const ASYNC_STATUS = {
  IDLE: 'idle' as const,
  LOADING: 'loading' as const,
  SUCCESS: 'success' as const,
  ERROR: 'error' as const,
};

// Location options
export const LOCATIONS = [
  'Helsinki',
  'Espoo',
  'Vantaa',
  'Tampere',
  'Turku',
  'Oulu',
  'Jyväskylä',
  'Kuopio',
  'Lahti',
  'Pori',
  'Lappeenranta',
  'Vaasa',
  'Joensuu',
  'Hämeenlinna',
  'Rovaniemi',
  'Seinäjoki',
  'Mikkeli',
  'Kotka',
  'Salo',
  'Kouvola',
  'Porvoo',
  'Lohja',
  'Kokkola',
  'Hyvinkää',
  'Järvenpää',
  'Nurmijärvi',
  'Kirkkonummi',
  'Kerava',
  'Rauma',
  'Kajaani',
  'Savonlinna',
  'Kemi',
  'Imatra',
  'Riihimäki',
  'Kangasala',
  'Ylöjärvi',
  'Sipoo',
  'Tuusula',
  'Lempäälä',
  'Nokia',
  'Karhula',
  'Jakobstad (Pietarsaari)',
  'Sastamala',
  'Raisio',
  'Kuusamo',
  'Iisalmi',
  'Valkeakoski',
  'Hanko',
] as const;

export type Location = typeof LOCATIONS[number];

export const isApiError = (error: unknown): error is ApiError => {
  return typeof error === 'object' && error !== null && 'error' in error && typeof (error as Record<string, unknown>).error === 'string';
};

export const isConnectionRequest = (obj: unknown): obj is ConnectionRequest => {
  return typeof obj === 'object' && obj !== null && 'id' in obj && typeof (obj as Record<string, unknown>).id === 'number';
};

export const isProfile = (obj: unknown): obj is Profile => {
  return typeof obj === 'object' && obj !== null && 'user_id' in obj && typeof (obj as Record<string, unknown>).user_id === 'number';
};
