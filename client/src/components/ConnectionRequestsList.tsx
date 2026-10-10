import React, { useState, useEffect, useCallback } from 'react';
import connectionService, { type ConnectionRequest } from '../services/connectionService';
import ConnectionRequestCard from './ConnectionRequestCard.tsx';

interface ConnectionRequestsListProps {
  refreshTrigger?: number;
  onRequestsUpdate?: (count: number) => void;
}

const ConnectionRequestsList: React.FC<ConnectionRequestsListProps> = ({
  refreshTrigger = 0,
  onRequestsUpdate
}) => {
  const [requests, setRequests] = useState<ConnectionRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRequests = useCallback(async () => {
    try {
      setLoading(true);
      const receivedRequests = await connectionService.getReceivedRequests();
      setRequests(receivedRequests);
      onRequestsUpdate?.(receivedRequests.length);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load connection requests');
    } finally {
      setLoading(false);
    }
  }, [onRequestsUpdate]);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests, refreshTrigger]);

  const handleAccept = async (requestId: number) => {
    try {
      await connectionService.acceptConnectionRequest(requestId);
      // Remove the accepted request from the list
      setRequests(prev => prev.filter(req => req.id !== requestId));
      onRequestsUpdate?.(requests.length - 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to accept request');
    }
  };

  const handleDismiss = async (requestId: number) => {
    try {
      await connectionService.dismissConnectionRequest(requestId);
      // Remove the dismissed request from the list
      setRequests(prev => prev.filter(req => req.id !== requestId));
      onRequestsUpdate?.(requests.length - 1);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to dismiss request');
    }
  };

  const handleAcceptAll = async () => {
    try {
      const acceptedCount = await connectionService.acceptAllRequests();
      if (acceptedCount > 0) {
        setRequests([]);
        onRequestsUpdate?.(0);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to accept all requests');
    }
  };

  const handleDismissAll = async () => {
    try {
      const dismissedCount = await connectionService.dismissAllRequests();
      if (dismissedCount > 0) {
        setRequests([]);
        onRequestsUpdate?.(0);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to dismiss all requests');
    }
  };

  if (loading) {
    return (
      <div className="connection-requests-list">
        <div className="loading">Loading connection requests...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="connection-requests-list">
        <div className="error-message">
          {error}
          <button onClick={fetchRequests} className="retry-btn">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="connection-requests-list">
      <div className="requests-header">
        <h2>Connection Requests ({requests.length})</h2>

        {requests.length > 0 && (
          <div className="bulk-actions">
            <button onClick={handleAcceptAll} className="accept-all-btn">
              Accept All
            </button>
            <button onClick={handleDismissAll} className="dismiss-all-btn">
              Dismiss All
            </button>
          </div>
        )}
      </div>

      {requests.length === 0 ? (
        <div className="no-requests">
          <p>No pending connection requests</p>
        </div>
      ) : (
        <div className="users-grid">
          {requests.map(request => (
            <ConnectionRequestCard
              key={request.id}
              request={request}
              onAccept={handleAccept}
              onDismiss={handleDismiss}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export default ConnectionRequestsList;