// profileService.ts
// This service handles all profile-related API calls
// It uses the authService for authenticated requests

import authService from './authService';
import type {
  Profile,
  BioData,
  UserBasicInfo,
  UserProfileInfo,
  UserBioInfo,
  ProfileUpdateRequest,
  BioDataUpdateRequest,
} from '../types';

const API_BASE_URL = 'http://localhost:8080';

class ProfileService {

  /**
   * Get current user's basic info (name + picture)
   * GET /me
   */
  async getMyBasicInfo(): Promise<UserBasicInfo> {
    const response = await authService.authenticatedFetch(`/me?t=${new Date().getTime()}`);

    if (!response.ok) {
      throw new Error('Failed to fetch user basic info');
    }

    return response.json();
  }

  /**
   * Get current user's profile info (about me section)
   * GET /me/profile
   */
  async getMyProfile(): Promise<UserProfileInfo> {
    const response = await authService.authenticatedFetch(`/me/profile?t=${new Date().getTime()}`);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Profile not found - please complete your profile');
      }
      throw new Error('Failed to fetch profile');
    }

    return response.json();
  }

  /**
   * Get current user's bio data (matching information)
   * GET /me/bio
   */
  async getMyBioData(): Promise<UserBioInfo> {
    const response = await authService.authenticatedFetch(`/me/bio?t=${new Date().getTime()}`);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Bio data not found - please complete your bio');
      }
      throw new Error('Failed to fetch bio data');
    }

    return response.json();
  }

  /**
   * Get current user's complete profile (basic + profile + bio)
   * This combines all three endpoints for convenience
   */
  async getMyCompleteProfile(): Promise<{
    basic: UserBasicInfo;
    profile: UserProfileInfo | null;
    bio: UserBioInfo | null;
  }> {
    try {
      const basic = await this.getMyBasicInfo();
      console.log('Basic info response:', basic); // Debug log

      let profile: UserProfileInfo | null = null;
      let bio: UserBioInfo | null = null;

      // Try to get profile (might not exist yet)
      try {
        profile = await this.getMyProfile();
        console.log('Profile info response:', profile); // Debug log
      } catch (error) {
        console.log('Profile not found, needs to be created');
      }

      // Try to get bio (might not exist yet)
      try {
        bio = await this.getMyBioData();
        console.log('Bio info response:', bio); // Debug log
      } catch (error) {
        console.log('Bio not found, needs to be created');
      }

      return { basic, profile, bio };
    } catch (error) {
      console.error('Error fetching complete profile:', error);
      throw error;
    }
  }

  /**
   * Create or update current user's profile
   * POST /me/profile
   */
  async updateMyProfile(profileData: ProfileUpdateRequest): Promise<Profile> {
    const response = await authService.authenticatedFetch('/me/profile', {
      method: 'POST',
      body: JSON.stringify(profileData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update profile');
    }

    return response.json();
  }

  /**
   * Create or update current user's bio data
   * POST /me/bio
   */
  async updateMyBioData(bioData: BioDataUpdateRequest): Promise<BioData> {
    const response = await authService.authenticatedFetch('/me/bio', {
      method: 'POST',
      body: JSON.stringify(bioData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update bio data');
    }

    return response.json();
  }

  /**
   * Upload profile picture
   * POST /me/profile-picture
   */
  async uploadProfilePicture(file: File): Promise<{ profile_picture_url: string }> {
    // Create FormData for file upload
    const formData = new FormData();
    formData.append('profile_picture', file);

    const token = authService.getToken();
    if (!token) {
      throw new Error('Not authenticated');
    }

    const response = await fetch(`${API_BASE_URL}/me/profile-picture`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to upload profile picture');
    }

    return response.json();
  }

  /**
   * Delete profile picture
   * DELETE /me/profile-picture
   */
  async deleteProfilePicture(): Promise<void> {
    const response = await authService.authenticatedFetch('/me/profile-picture', {
      method: 'DELETE',
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete profile picture');
    }
  }


  /**
   * Get another user's basic info (name + picture)
   * GET /users/{id}
   */
  async getUserBasicInfo(userId: number): Promise<UserBasicInfo> {
    const response = await authService.authenticatedFetch(`/users/${userId}`);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('User not found or you do not have permission to view this profile');
      }
      throw new Error('Failed to fetch user info');
    }

    return response.json();
  }

  /**
   * Get another user's profile info (about me section)
   * GET /users/{id}/profile
   */
  async getUserProfile(userId: number): Promise<UserProfileInfo> {
    const response = await authService.authenticatedFetch(`/users/${userId}/profile`);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Profile not found or you do not have permission to view it');
      }
      throw new Error('Failed to fetch profile');
    }

    return response.json();
  }

  /**
   * Get another user's bio data (matching information)
   * GET /users/{id}/bio
   */
  async getUserBioData(userId: number): Promise<UserBioInfo> {
    const response = await authService.authenticatedFetch(`/users/${userId}/bio`);

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Bio data not found or you do not have permission to view it');
      }
      throw new Error('Failed to fetch bio data');
    }

    return response.json();
  }

  /**
   * Get another user's complete profile (basic + profile + bio)
   * This combines all three endpoints for convenience
   */
  async getUserCompleteProfile(userId: number): Promise<{
    basic: UserBasicInfo;
    profile: UserProfileInfo;
    bio: UserBioInfo;
  }> {
    try {
      // Fetch all three in parallel for better performance
      const [basic, profile, bio] = await Promise.all([
        this.getUserBasicInfo(userId),
        this.getUserProfile(userId),
        this.getUserBioData(userId),
      ]);

      return { basic, profile, bio };
    } catch (error) {
      console.error('Error fetching complete profile:', error);
      throw error;
    }
  }

 
  /**
   * Check if profile data is complete enough for matching
   * This is a client-side helper to validate before submission
   */
  isProfileDataComplete(profile: ProfileUpdateRequest, bio: BioDataUpdateRequest): boolean {
    // Profile must have name and location
    const hasName = Boolean((profile.first_name && profile.last_name) || profile.username);
    const hasLocation = Boolean(profile.location);

    // Bio must have at least 5 data points
    let bioDataPoints = 0;

    if (bio.interests && bio.interests.length > 0 && bio.interests[0] !== '') {
      bioDataPoints++;
    }
    if (bio.hobbies && bio.hobbies.length > 0 && bio.hobbies[0] !== '') {
      bioDataPoints++;
    }
    if (bio.music_preferences && bio.music_preferences.length > 0 && bio.music_preferences[0] !== '') {
      bioDataPoints++;
    }
    if (bio.food_preferences && bio.food_preferences.length > 0 && bio.food_preferences[0] !== '') {
      bioDataPoints++;
    }
    if (bio.looking_for && bio.looking_for.trim() !== '') {
      bioDataPoints++;
    }
    if (bio.age_range_min !== undefined && bio.age_range_max !== undefined) {
      bioDataPoints++;
    }

    return hasName && hasLocation && bioDataPoints >= 5;
  }

  /**
   * Validate image file before upload
   * Returns error message if invalid, null if valid
   */
  validateImageFile(file: File): string | null {
    // Check file size (max 5MB)
    const maxSize = 5 * 1024 * 1024; // 5MB in bytes
    if (file.size > maxSize) {
      return 'File size must be less than 5MB';
    }

    // Check file type
    const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
      return 'File must be an image (JPEG, PNG, GIF, or WebP)';
    }

    return null;
  }

  /**
   * Get profile picture URL (handles empty/placeholder cases)
   */
  getProfilePictureUrl(pictureUrl: string | null | undefined): string {
    if (!pictureUrl || pictureUrl === '' || pictureUrl === '👤') {
      return ''; // Return empty string, component can show placeholder
    }

    // If it's a relative path, prepend API base URL
    if (pictureUrl.startsWith('/uploads/')) {
      return `${API_BASE_URL}${pictureUrl}`;
    }

    return pictureUrl;
  }
}

// Create and export singleton instance
const profileService = new ProfileService();
export default profileService;
