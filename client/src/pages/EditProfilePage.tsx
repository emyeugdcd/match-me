import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import ProfileForm from '../components/ProfileForm';
import ProfilePictureUpload from '../components/ProfilePictureUpload';
import profileService from '../services/profileService';
import type { ProfileUpdateRequest, BioDataUpdateRequest } from '../types';

const EditProfilePage: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [initialProfile, setInitialProfile] = useState<ProfileUpdateRequest | undefined>();
  const [initialBio, setInitialBio] = useState<BioDataUpdateRequest | undefined>();
  const [currentPictureUrl, setCurrentPictureUrl] = useState<string>('');
  const [uploadError, setUploadError] = useState<string | null>(null);

  const loadCurrentProfile = async () => {
    try {
      setLoading(true);
      const { basic, profile, bio } = await profileService.getMyCompleteProfile();

      // Set current picture URL
      setCurrentPictureUrl(basic.profile_picture_url);

      // Set initial profile data if exists
      if (profile) {
        setInitialProfile({
          first_name: basic.first_name || '',
          last_name: basic.last_name || '',
          username: basic.username || basic.name,
          about_me: profile.about_me || '',
          location: profile.location || '',
          profile_picture_url: basic.profile_picture_url || '',
          latitude: profile.latitude ?? null,
          longitude: profile.longitude ?? null,
        });
      } else {
        // Set empty profile if none exists
        setInitialProfile({
          first_name: basic.first_name || '',
          last_name: basic.last_name || '',
          username: basic.username || basic.name,
          about_me: '',
          location: '',
          profile_picture_url: basic.profile_picture_url || '',
          latitude: null,
          longitude: null,
        });
      }

      // Set initial bio data if exists
      if (bio) {
        setInitialBio({
          interests: bio.interests && bio.interests.length > 0 ? bio.interests : [''],
          hobbies: bio.hobbies && bio.hobbies.length > 0 ? bio.hobbies : [''],
          music_preferences:
            bio.music_preferences && bio.music_preferences.length > 0
              ? bio.music_preferences
              : [''],
          food_preferences:
            bio.food_preferences && bio.food_preferences.length > 0
              ? bio.food_preferences
              : [''],
          looking_for: bio.looking_for || '',
          age_range_min: bio.age_range_min,
          age_range_max: bio.age_range_max,
          max_distance_km: bio.max_distance_km,
        });
      } else {
        // Set empty bio if none exists
        setInitialBio({
          interests: [''],
          hobbies: [''],
          music_preferences: [''],
          food_preferences: [''],
          looking_for: '',
          age_range_min: undefined,
          age_range_max: undefined,
          max_distance_km: undefined,
        });
      }
    } catch (error) {
      console.error('Error loading profile:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadCurrentProfile();
  }, []);

  const handleProfileSave = () => {
    // Navigate to profile view after successful save
    navigate('/profile/me');
  };

  const handleCancel = () => {
    navigate('/profile/me');
  };

  const handlePictureUploadSuccess = (newUrl: string) => {
    setCurrentPictureUrl(newUrl);
    setUploadError(null);
  };

  const handlePictureUploadError = (error: string) => {
    setUploadError(error);
  };

  const handleRemovePicture = () => {
    setCurrentPictureUrl('');
    setUploadError(null);
  };

  if (loading) {
    return (
      <div className="page-container edit-profile-page">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>Loading your profile...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container edit-profile-page">
      <div className="page-header">
        <h1>Edit Your Profile</h1>
        <p className="page-subtitle">
          Complete your profile to start connecting with others
        </p>
      </div>

      <div className="profile-edit-content">
        {/* Profile Picture Section */}
        <section className="picture-section">
          <h2>Profile Picture</h2>
          <ProfilePictureUpload
            currentPictureUrl={currentPictureUrl}
            onUploadSuccess={handlePictureUploadSuccess}
            onUploadError={handlePictureUploadError}
            onRemovePicture={handleRemovePicture}
          />
          {uploadError && <div className="error-message">{uploadError}</div>}
        </section>

        {/* Profile Form Section */}
        <div className="profile-form-container">
          <ProfileForm
            initialProfile={initialProfile}
            initialBio={initialBio}
            profilePictureUrl={currentPictureUrl}
            onSave={handleProfileSave}
            onCancel={handleCancel}
          />
        </div>
      </div>
    </div>
  );
};

export default EditProfilePage;
