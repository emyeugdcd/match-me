import React, { useState } from 'react';
import type { UserBasicInfo, UserProfileInfo, UserBioInfo } from '../types';
import profileService from '../services/profileService';
import OnlineIndicator from './OnlineIndicator';

interface ProfileViewProps {
  basicInfo: UserBasicInfo;
  profileInfo?: UserProfileInfo | null;
  bioInfo?: UserBioInfo | null;
  onEdit?: () => void;
  isOwnProfile?: boolean;
}

const ProfileView: React.FC<ProfileViewProps> = ({
  basicInfo,
  profileInfo,
  bioInfo,
  onEdit,
  isOwnProfile = false
}) => {
  const [showFullAboutMe, setShowFullAboutMe] = useState(false);

  const profilePictureUrl = profileService.getProfilePictureUrl(basicInfo.profile_picture_url);

  const aboutMePreview = profileInfo?.about_me && profileInfo.about_me.length > 200
    ? profileInfo.about_me.substring(0, 200) + '...'
    : profileInfo?.about_me;

  const renderArrayItems = (items: string[] | undefined, label: string) => {
    if (!items || items.length === 0) {
      return <p className="empty-text">No {label.toLowerCase()} added yet</p>;
    }

    return (
      <div className="tags-container">
        {items.map((item, index) => (
          <span key={index} className="tag">
            {item}
          </span>
        ))}
      </div>
    );
  };

  return (
    <div className="profile-view">
      <div className="profile-header">
        <div className="profile-picture-container">
          <img
            src={profilePictureUrl || 'https://via.placeholder.com/150'}
            alt={`${basicInfo.first_name} ${basicInfo.last_name}`}
            className="profile-picture"
          />
        </div>

        <div className="profile-info">
          <div className="name-section" style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <h1 className="profile-name">{basicInfo.first_name} {basicInfo.last_name}</h1>
            {!isOwnProfile && <OnlineIndicator userId={basicInfo.id} showLabel={true} />}
          </div>

          {basicInfo.username && (
            <p className="profile-username">@{basicInfo.username}</p>
          )}

          {isOwnProfile && basicInfo.email && (
            <p className="profile-email">📧 {basicInfo.email}</p>
          )}

          {profileInfo?.location && (
            <div className="profile-location">
              <span>📍</span> {profileInfo.location}
            </div>
          )}

          {isOwnProfile && onEdit && (
            <button onClick={onEdit} className="edit-profile-btn">
              Edit Profile
            </button>
          )}
        </div>
      </div>

      <div className="profile-content">
        <div className="profile-main-column">
          {/* About Me Section */}
          <section className="profile-section">
            <h2>About Me</h2>
            {profileInfo?.about_me ? (
              <div className="about-me-section">
                <p className="about-me-text">
                  {showFullAboutMe ? profileInfo.about_me : aboutMePreview}
                </p>
                {profileInfo.about_me.length > 200 && (
                  <button
                    onClick={() => setShowFullAboutMe(!showFullAboutMe)}
                    className="read-more-btn"
                  >
                    {showFullAboutMe ? 'Show Less' : 'Read More'}
                  </button>
                )}
              </div>
            ) : (
              <p className="empty-text">No about me information yet</p>
            )}
          </section>

          {/* Looking For Section */}
          {bioInfo && (
            <section className="profile-section">
              <h2>Looking For</h2>
              {bioInfo.looking_for ? (
                <p className="about-me-text">{bioInfo.looking_for}</p>
              ) : (
                <p className="empty-text">Not specified</p>
              )}
            </section>
          )}
        </div>

        <div className="profile-sidebar-column">
          {bioInfo && (
            <>
              <section className="profile-section">
                <h2>Interests</h2>
                {renderArrayItems(bioInfo.interests, 'interests')}
              </section>

              <section className="profile-section">
                <h2>Programming Languages</h2>
                {renderArrayItems(bioInfo.hobbies, 'programming languages')}
              </section>

              <section className="profile-section">
                <h2>Music Preferences</h2>
                {renderArrayItems(bioInfo.music_preferences, 'music preferences')}
              </section>

              <section className="profile-section">
                <h2>Food Preferences</h2>
                {renderArrayItems(bioInfo.food_preferences, 'food preferences')}
              </section>

              {bioInfo.age_range_min && bioInfo.age_range_max && (
                <section className="profile-section">
                  <h2>Age Preference</h2>
                  <p>{bioInfo.age_range_min} - {bioInfo.age_range_max} years old</p>
                </section>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileView;