import React, { useState, useEffect } from 'react';
import connectionService, { type DetailedConnection } from '../services/connectionService';
import ConnectedUserCard from './ConnectedUserCard.tsx';

interface ConnectedUsersListProps {
  refreshTrigger?: number;
  onConnectionsUpdate?: (count: number) => void;
  onViewProfile?: (userId: number) => void;
  onStartChat?: (userId: number) => void;
}

const ConnectedUsersList: React.FC<ConnectedUsersListProps> = ({
  refreshTrigger = 0,
  onConnectionsUpdate,
  onViewProfile,
  onStartChat
}) => {
  const [connections, setConnections] = useState<DetailedConnection[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchConnections = async () => {
    try {
      setLoading(true);
      const detailedConnections = await connectionService.getDetailedConnections();
      setConnections(detailedConnections);
      onConnectionsUpdate?.(detailedConnections.length);
    } catch (err: any) {
      setError(err.message || 'Failed to load connections');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConnections();
  }, [refreshTrigger]);

  const handleDisconnect = async (userId: number) => {
    if (!window.confirm('Are you sure you want to disconnect from this user?')) {
      return;
    }

    try {
      await connectionService.disconnectUser(userId);
      // Remove the disconnected user from the list
      setConnections(prev => prev.filter(conn => conn.user_id !== userId));
      onConnectionsUpdate?.(connections.length - 1);
    } catch (err: any) {
      setError(err.message || 'Failed to disconnect');
    }
  };

  const handleViewProfile = (userId: number) => {
    onViewProfile?.(userId);
  };

  const handleStartChat = (userId: number) => {
    onStartChat?.(userId);
  };

  if (loading) {
    return (
      <div className="connected-users-list">
        <div className="loading">Loading connections...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="connected-users-list">
        <div className="error-message">
          {error}
          <button onClick={fetchConnections} className="retry-btn">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="connected-users-list">
      <div className="connections-header">
        <h2>Your Connections ({connections.length})</h2>
      </div>

      {connections.length === 0 ? (
        <div className="no-connections">
          <p>You haven't connected with anyone yet.</p>
          <p>Start sending connection requests to build your network!</p>
        </div>
      ) : (
        <div className="users-grid">
          {connections.map(connection => (
            <ConnectedUserCard
              key={connection.user_id}
              connection={connection}
              onDisconnect={handleDisconnect}
              onViewProfile={handleViewProfile}
              onStartChat={handleStartChat}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default ConnectedUsersList;