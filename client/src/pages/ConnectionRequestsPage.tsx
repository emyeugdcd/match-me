import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import ConnectionRequestsList from '../components/ConnectionRequestsList';

const ConnectionRequestsPage: React.FC = () => {
  const navigate = useNavigate();
  const [requestCount, setRequestCount] = useState(0);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleRequestsUpdate = (count: number) => {
    setRequestCount(count);
  };

  const handleRefresh = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="page-container connection-requests-page">
      <div className="page-header">
        <div className="header-left">
          <button onClick={() => navigate('/connections')} className="back-btn">
            ← Back
          </button>
          <h1>Connection Requests</h1>
        </div>
        <div className="header-right">
          <button onClick={handleRefresh} className="refresh-btn">
            🔄 Refresh
          </button>
        </div>
      </div>

      {requestCount > 0 && (
        <div className="page-info">
          <p>You have {requestCount} pending connection request{requestCount !== 1 ? 's' : ''}</p>
        </div>
      )}

      <ConnectionRequestsList
        refreshTrigger={refreshTrigger}
        onRequestsUpdate={handleRequestsUpdate}
      />
    </div>
  );
};

export default ConnectionRequestsPage;