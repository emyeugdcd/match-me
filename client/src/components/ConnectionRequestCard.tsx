import React, { useState } from 'react';
import { type ConnectionRequest } from '../types';
import connectionService from '../services/connectionService';
import profileService from '../services/profileService';

interface ConnectionRequestCardProps {
  request: ConnectionRequest;
  onAccept: (requestId: number) => void;
  onDismiss: (requestId: number) => void;
}

const ConnectionRequestCard: React.FC<ConnectionRequestCardProps> = ({
  request,
  onAccept,
  onDismiss
}) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const profilePictureUrl = profileService.getProfilePictureUrl(request.profile_picture_url);
  const formattedDate = connectionService.formatConnectionDate(request.created_at);

  const handleAccept = async () => {
    setIsProcessing(true);
    setActionError(null);

    try {
      await onAccept(request.id);
    } catch (err: any) {
      setActionError(err.message || 'Failed to accept request');
    } finally {
      setIsProcessing(false);
    }
  };

  const handleDismiss = async () => {
    setIsProcessing(true);
    setActionError(null);

    try {
      await onDismiss(request.id);
    } catch (err: any) {
      setActionError(err.message || 'Failed to dismiss request');
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="user-card">
      <div className="card-header">
        {profilePictureUrl ? (
          <img
            src={profilePictureUrl}
            alt={request.display_name}
            className="card-profile-pic"
          />
        ) : (
          <div className="card-profile-pic" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', background: 'var(--primary-gradient)', fontSize: '2rem', fontWeight: 'bold', color: '#000' }}>
            {request.display_name.charAt(0).toUpperCase()}
          </div>
        )}
        <h3 className="card-name">{request.display_name}</h3>
        <p className="request-date" style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>Sent {formattedDate}</p>
      </div>

      <div className="card-body">
        {actionError && (
          <div className="action-error" style={{ color: 'var(--error-color)', textAlign: 'center', fontSize: '0.9rem' }}>
            {actionError}
          </div>
        )}
      </div>

      <div className="card-actions">
        <button
          onClick={handleAccept}
          disabled={isProcessing}
          className="btn-primary-card"
        >
          {isProcessing ? 'Accepting...' : 'Accept'}
        </button>

        <button
          onClick={handleDismiss}
          disabled={isProcessing}
          className="btn-secondary-card"
        >
          {isProcessing ? 'Dismissing...' : 'Dismiss'}
        </button>
      </div>
    </div>
  );
};

export default ConnectionRequestCard;