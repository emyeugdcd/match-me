import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import ProfileView from '../components/ProfileView';
import profileService from '../services/profileService';
import connectionService from '../services/connectionService';
import type { UserBasicInfo, UserProfileInfo, UserBioInfo, ConnectionStatus } from '../types';

const ProfilePage: React.FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();

  const [basicInfo, setBasicInfo] = useState<UserBasicInfo | null>(null);
  const [profileInfo, setProfileInfo] = useState<UserProfileInfo | null>(null);
  const [bioInfo, setBioInfo] = useState<UserBioInfo | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isOwnProfile, setIsOwnProfile] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    loadProfile();
  }, [userId]);

  const loadProfile = async () => {
    try {
      setLoading(true);
      setError(null);

      // Determine if this is the current user's profile

      // Check if viewing own profile
      const isOwn = !userId || userId === 'me';
      setIsOwnProfile(isOwn);

      // Load profile data
      if (isOwn) {
        // Load current user's profile
        const { basic, profile, bio } = await profileService.getMyCompleteProfile();
        setBasicInfo(basic);
        setProfileInfo(profile);
        setBioInfo(bio);
      } else {
        // Load another user's profile
        const targetUserId = parseInt(userId!);

        if (isNaN(targetUserId)) {
          throw new Error('Invalid user ID');
        }

        const { basic, profile, bio } = await profileService.getUserCompleteProfile(targetUserId);
        setBasicInfo(basic);
        setProfileInfo(profile);
        setBioInfo(bio);

        // Load connection status for other users
        const status = await connectionService.checkConnectionStatus(targetUserId);
        setConnectionStatus(status);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load profile');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = () => {
    navigate('/profile/edit');
  };

  const handleSendConnectionRequest = async () => {
    if (!basicInfo) return;

    try {
      setActionLoading(true);
      await connectionService.sendConnectionRequest(basicInfo.id);

      // Reload connection status
      const status = await connectionService.checkConnectionStatus(basicInfo.id);
      setConnectionStatus(status);
    } catch (err: any) {
      alert(err.message || 'Failed to send connection request');
    } finally {
      setActionLoading(false);
    }
  };

  const handleAcceptRequest = async () => {
    if (!connectionStatus?.request_id) return;

    try {
      setActionLoading(true);
      await connectionService.acceptConnectionRequest(connectionStatus.request_id);

      // Reload connection status
      if (basicInfo) {
        const status = await connectionService.checkConnectionStatus(basicInfo.id);
        setConnectionStatus(status);
      }
    } catch (err: any) {
      alert(err.message || 'Failed to accept connection request');
    } finally {
      setActionLoading(false);
    }
  };

  const handleStartChat = () => {
    if (basicInfo) {
      navigate(`/chat/${basicInfo.id}`);
    }
  };

  const renderConnectionActions = () => {
    if (isOwnProfile || !connectionStatus) return null;

    switch (connectionStatus.status) {
      case 'none':
        return (
          <button
            onClick={handleSendConnectionRequest}
            disabled={actionLoading}
            className="action-btn connect-btn"
          >
            {actionLoading ? 'Sending...' : 'Send Connection Request'}
          </button>
        );

      case 'request_sent':
        return (
          <button disabled className="action-btn request-sent-btn">
            Connection Request Sent
          </button>
        );

      case 'request_received':
        return (
          <button
            onClick={handleAcceptRequest}
            disabled={actionLoading}
            className="action-btn accept-btn"
          >
            {actionLoading ? 'Accepting...' : 'Accept Connection Request'}
          </button>
        );

      case 'connected':
        return (
          <div className="connected-actions">
            <span className="connected-badge">✓ Connected</span>
            <button
              onClick={handleStartChat}
              className="action-btn chat-btn"
            >
              Start Chat
            </button>
          </div>
        );

      default:
        return null;
    }
  };

  if (loading) {
    return (
      <div className="page-container profile-page">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>Loading profile...</p>
        </div>
      </div>
    );
  }

  if (error || !basicInfo) {
    return (
      <div className="page-container profile-page">
        <div className="error-container">
          <h2>Profile Not Found</h2>
          <p>{error || 'This profile does not exist or you do not have permission to view it.'}</p>

          {isOwnProfile ? (
            <div style={{ marginTop: '2rem' }}>
              <p style={{ marginBottom: '1rem', color: 'rgba(255,255,255,0.7)' }}>
                It looks like you haven't created your profile yet.
              </p>
              <button
                onClick={() => navigate('/profile/edit')}
                className="action-btn"
                style={{ marginRight: '1rem' }}
              >
                Create Your Profile
              </button>
              <button onClick={() => navigate('/')} className="back-btn">
                Back to Home
              </button>
            </div>
          ) : (
            <button onClick={() => navigate('/connections')} className="back-btn">
              Back to Connections
            </button>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className="page-container profile-page">
      {!isOwnProfile && (
        <div className="profile-actions-bar">
          {renderConnectionActions()}
        </div>
      )}

      <ProfileView
        basicInfo={basicInfo}
        profileInfo={profileInfo}
        bioInfo={bioInfo}
        onEdit={isOwnProfile ? handleEdit : undefined}
        isOwnProfile={isOwnProfile}
      />
    </div>
  );
};

export default ProfilePage;
