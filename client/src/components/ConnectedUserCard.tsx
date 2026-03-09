import React, { useState } from 'react';
import { type DetailedConnection } from '../types';
import connectionService from '../services/connectionService';
import profileService from '../services/profileService';
import OnlineIndicator from './OnlineIndicator';

interface ConnectedUserCardProps {
  connection: DetailedConnection;
  onDisconnect: (userId: number) => void;
  onViewProfile: (userId: number) => void;
  onStartChat: (userId: number) => void;
}

const ConnectedUserCard: React.FC<ConnectedUserCardProps> = ({
  connection,
  onDisconnect,
  onViewProfile,
  onStartChat
}) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const profilePictureUrl = profileService.getProfilePictureUrl(connection.profile_picture_url);
  const connectedDate = connectionService.formatConnectionDate(connection.connected_at);

  const handleDisconnect = async () => {
    if (!window.confirm(`Are you sure you want to disconnect from ${connection.display_name}?`)) {
      return;
    }

    setIsProcessing(true);
    setActionError(null);

    try {
      await onDisconnect(connection.user_id);
    } catch (err: any) {
      setActionError(err.message || 'Failed to disconnect');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleViewProfile = () => {
    onViewProfile(connection.user_id);
  };

  const handleStartChat = () => {
    onStartChat(connection.user_id);
  };

  return (
    <div className="user-card fade-in-up">
      <div className="card-header">
        {profilePictureUrl ? (
          <img
            src={profilePictureUrl}
            alt={connection.display_name}
            className="card-profile-pic"
          />
        ) : (
          <div className="card-profile-pic" style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'var(--primary-gradient)',
            fontSize: '2.5rem',
            fontWeight: 'bold',
            color: '#000',
            textShadow: '0 2px 4px rgba(0,0,0,0.1)'
          }}>
            {connection.display_name.charAt(0).toUpperCase()}
          </div>
        )}
        <div style={{display: 'flex', gap: '8px', alignItems: 'center'}}>
          <h3 className="card-name" style={{margin: 0}}>{connection.display_name}</h3>
          <OnlineIndicator userId={connection.user_id} />
        </div>
        <p className="card-location">
          <span style={{ marginRight: '4px' }}>📍</span>
          {connection.location || 'Unknown Location'}
        </p>
      </div>

      <div className="card-body">
        <div style={{
          background: 'rgba(255,255,255,0.05)',
          padding: '0.75rem',
          borderRadius: '12px',
          textAlign: 'center',
          marginBottom: '1rem'
        }}>
          <p className="connection-date" style={{
            color: 'rgba(255,255,255,0.7)',
            fontSize: '0.85rem',
            margin: 0
          }}>
            Connected {connectedDate}
          </p>
        </div>

        {actionError && (
          <div className="action-error" style={{
            color: 'var(--error-color)',
            textAlign: 'center',
            fontSize: '0.9rem',
            background: 'rgba(231, 76, 60, 0.1)',
            padding: '0.5rem',
            borderRadius: '8px'
          }}>
            {actionError}
          </div>
        )}
      </div>

      <div className="card-actions" style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.75rem' }}>
          <button
            onClick={handleStartChat}
            className="btn-primary-card btn-glow"
            style={{ flex: 1 }}
          >
            💬 Chat
          </button>

          <button
            onClick={handleViewProfile}
            className="btn-secondary-card"
            style={{ flex: 1 }}
          >
            👤 Profile
          </button>
        </div>

        <button
          onClick={handleDisconnect}
          disabled={isProcessing}
          className="btn-text"
          style={{
            background: 'transparent',
            border: 'none',
            color: 'rgba(255, 255, 255, 0.4)',
            fontSize: '0.8rem',
            padding: '0.5rem',
            cursor: 'pointer',
            transition: 'color 0.2s',
            alignSelf: 'center'
          }}
          onMouseEnter={(e) => e.currentTarget.style.color = 'var(--error-color)'}
          onMouseLeave={(e) => e.currentTarget.style.color = 'rgba(255, 255, 255, 0.4)'}
        >
          {isProcessing ? 'Disconnecting...' : 'Disconnect'}
        </button>
      </div>
    </div>
  );
};

export default ConnectedUserCard;