import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ConnectedUsersList from '../components/ConnectedUsersList';
import connectionService from '../services/connectionService';

const ConnectionsPage: React.FC = () => {
  const navigate = useNavigate();
  const [connectionCount, setConnectionCount] = useState(0);
  const [requestCount, setRequestCount] = useState(0);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  // Load request count on mount
  React.useEffect(() => {
    loadRequestCount();
  }, []);

  const loadRequestCount = async () => {
    try {
      const count = await connectionService.getReceivedRequestsCount();
      setRequestCount(count);
    } catch (error) {
      console.error('Error loading request count:', error);
    }
  };

  const handleConnectionsUpdate = (count: number) => {
    setConnectionCount(count);
  };

  const handleViewProfile = (userId: number) => {
    navigate(`/profile/${userId}`);
  };

  const handleStartChat = (userId: number) => {
    navigate(`/chat/${userId}`);
  };

  const handleRefresh = () => {
    setRefreshTrigger(prev => prev + 1);
    loadRequestCount();
  };

  const handleViewRequests = () => {
    navigate('/connections/requests');
  };

  return (
    <div className="page-container connections-page">
      <div className="page-header">
        <h1>Your Network</h1>
        <div className="header-actions">
          <button onClick={handleRefresh} className="refresh-btn">
            🔄 Refresh
          </button>
        </div>
      </div>

      {/* Connection Stats */}
      <div className="connection-stats fade-in-up">
        <div className="stat-card glass-card">
          <div className="stat-number">{connectionCount}</div>
          <div className="stat-label">Connections</div>
        </div>
        <div className="stat-card glass-card clickable" onClick={handleViewRequests}>
          <div className="stat-number">
            {requestCount}
            {requestCount > 0 && <span className="notification-badge"></span>}
          </div>
          <div className="stat-label">Pending Requests</div>
        </div>
      </div>

      {/* Pending Requests Banner */}
      {requestCount > 0 && (
        <div className="requests-banner">
          <p>
            You have {requestCount} pending connection request{requestCount !== 1 ? 's' : ''}
          </p>
          <button onClick={handleViewRequests} className="view-requests-btn">
            View Requests
          </button>
        </div>
      )}

      {/* Connections List */}
      <ConnectedUsersList
        refreshTrigger={refreshTrigger}
        onConnectionsUpdate={handleConnectionsUpdate}
        onViewProfile={handleViewProfile}
        onStartChat={handleStartChat}
      />
    </div>
  );
};

export default ConnectionsPage;
