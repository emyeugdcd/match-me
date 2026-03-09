// connectionService.ts
// This service handles all connection-related API calls
// Manages connection requests, acceptances, dismissals, and disconnections

import authService from './authService';


export interface ConnectionRequest {
  id: number;
  from_user_id: number;
  display_name: string;
  profile_picture_url: string;
  created_at: string;
}

export interface SentConnectionRequest {
  id: number;
  to_user_id: number;
  display_name: string;
  profile_picture_url: string;
  created_at: string;
}

export interface ConnectionsResponse {
  connections: number[];  // Array of user IDs (per spec)
}

export interface DetailedConnection {
  user_id: number;
  display_name: string;
  profile_picture_url: string;
  location: string;
  connected_at: string;
}

export interface ConnectionStatus {
  status: 'none' | 'connected' | 'request_sent' | 'request_received';
  connected: boolean;
  request_id?: number;
}

// Response types for actions
export interface SendRequestResponse {
  message: string;
  request_id?: number;
  status: string;
  created_at?: string;
  connected_at?: string;  // For auto-accepted case
}

export interface AcceptRequestResponse {
  message: string;
  connection_id: number;
  connected_user_id: number;
  connected_at: string;
}


class ConnectionService {

  /**
   * Send a connection request to another user
   * POST /connections/request/{userId}
   * 
   * Note: This can auto-accept if there's a mutual request!
   */
  async sendConnectionRequest(userId: number): Promise<SendRequestResponse> {
    const response = await authService.authenticatedFetch(
      `/connections/request/${userId}`,
      { method: 'POST' }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to send connection request');
    }

    const data = await response.json();

    // Check if it was auto-accepted (mutual interest)
    if (data.status === 'connected') {
      console.log('🎉 Connection automatically established (mutual interest)!');
    }

    return data;
  }


  /**
   * Get all pending connection requests received by current user
   * GET /connections/requests
   */
  async getReceivedRequests(): Promise<ConnectionRequest[]> {
    const response = await authService.authenticatedFetch('/connections/requests');

    if (!response.ok) {
      throw new Error('Failed to fetch connection requests');
    }

    const data = await response.json();
    return data.requests || [];
  }

  /**
   * Get all pending connection requests sent by current user
   * GET /connections/requests/sent
   */
  async getSentRequests(): Promise<SentConnectionRequest[]> {
    const response = await authService.authenticatedFetch('/connections/requests/sent');

    if (!response.ok) {
      throw new Error('Failed to fetch sent requests');
    }

    const data = await response.json();
    return data.requests || [];
  }

  
  /**
   * Accept a connection request
   * POST /connections/requests/{requestId}/accept
   */
  async acceptConnectionRequest(requestId: number): Promise<AcceptRequestResponse> {
    const response = await authService.authenticatedFetch(
      `/connections/requests/${requestId}/accept`,
      { method: 'POST' }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to accept connection request');
    }

    const data = await response.json();
    console.log('✅ Connection request accepted!');
    return data;
  }

  /**
   * Dismiss (reject) a connection request
   * POST /connections/requests/{requestId}/dismiss
   */
  async dismissConnectionRequest(requestId: number): Promise<void> {
    const response = await authService.authenticatedFetch(
      `/connections/requests/${requestId}/dismiss`,
      { method: 'POST' }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to dismiss connection request');
    }

    console.log('Connection request dismissed');
  }

  /**
   * Get all connections for current user (just IDs, per spec)
   * GET /connections
   */
  async getConnections(): Promise<number[]> {
    const response = await authService.authenticatedFetch('/connections');

    if (!response.ok) {
      throw new Error('Failed to fetch connections');
    }

    const data: ConnectionsResponse = await response.json();
    return data.connections || [];
  }

  /**
   * Get detailed connection information (with profiles)
   * GET /connections/detailed
   */
  async getDetailedConnections(): Promise<DetailedConnection[]> {
    const response = await authService.authenticatedFetch('/connections/detailed');

    if (!response.ok) {
      throw new Error('Failed to fetch detailed connections');
    }

    const data = await response.json();
    return data.connections || [];
  }

  /**
   * Disconnect from a user (remove connection)
   * DELETE /connections/{userId}
   */
  async disconnectUser(userId: number): Promise<void> {
    const response = await authService.authenticatedFetch(
      `/connections/${userId}`,
      { method: 'DELETE' }
    );

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to disconnect user');
    }

    console.log('Connection removed');
  }

  /**
   * Check connection status with another user
   * GET /connections/status/{userId}
   * 
   * Returns:
   * - 'none': No relationship
   * - 'connected': Already connected
   * - 'request_sent': You sent them a request
   * - 'request_received': They sent you a request
   */
  async checkConnectionStatus(userId: number): Promise<ConnectionStatus> {
    const response = await authService.authenticatedFetch(`/connections/status/${userId}`);

    if (!response.ok) {
      throw new Error('Failed to check connection status');
    }

    return response.json();
  }

  /**
   * Check if connected with a specific user (quick boolean check)
   */
  async isConnectedWith(userId: number): Promise<boolean> {
    try {
      const status = await this.checkConnectionStatus(userId);
      return status.connected;
    } catch (error) {
      console.error('Error checking connection:', error);
      return false;
    }
  }


  /**
   * Get the count of pending received requests
   * Useful for notification badges
   */
  async getReceivedRequestsCount(): Promise<number> {
    try {
      const requests = await this.getReceivedRequests();
      return requests.length;
    } catch (error) {
      console.error('Error getting request count:', error);
      return 0;
    }
  }

  /**
   * Get the count of sent pending requests
   */
  async getSentRequestsCount(): Promise<number> {
    try {
      const requests = await this.getSentRequests();
      return requests.length;
    } catch (error) {
      console.error('Error getting sent request count:', error);
      return 0;
    }
  }

  /**
   * Get total connections count
   */
  async getConnectionsCount(): Promise<number> {
    try {
      const connections = await this.getConnections();
      return connections.length;
    } catch (error) {
      console.error('Error getting connections count:', error);
      return 0;
    }
  }

  /**
   * Bulk action: Accept all pending requests
   * Returns number of requests accepted
   */
  async acceptAllRequests(): Promise<number> {
    try {
      const requests = await this.getReceivedRequests();
      let acceptedCount = 0;

      for (const request of requests) {
        try {
          await this.acceptConnectionRequest(request.id);
          acceptedCount++;
        } catch (error) {
          console.error(`Failed to accept request ${request.id}:`, error);
        }
      }

      return acceptedCount;
    } catch (error) {
      console.error('Error accepting all requests:', error);
      throw error;
    }
  }

  /**
   * Bulk action: Dismiss all pending requests
   * Returns number of requests dismissed
   */
  async dismissAllRequests(): Promise<number> {
    try {
      const requests = await this.getReceivedRequests();
      let dismissedCount = 0;

      for (const request of requests) {
        try {
          await this.dismissConnectionRequest(request.id);
          dismissedCount++;
        } catch (error) {
          console.error(`Failed to dismiss request ${request.id}:`, error);
        }
      }

      return dismissedCount;
    } catch (error) {
      console.error('Error dismissing all requests:', error);
      throw error;
    }
  }

  /**
   * Get connection request notification message
   * Helper for displaying user-friendly messages
   */
  getRequestNotificationMessage(count: number): string {
    if (count === 0) return 'No pending requests';
    if (count === 1) return '1 new connection request';
    return `${count} new connection requests`;
  }

  /**
   * Format date for display
   * Helper to format the created_at/connected_at timestamps
   */
  formatConnectionDate(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins} minute${diffMins > 1 ? 's' : ''} ago`;
    if (diffHours < 24) return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
    if (diffDays < 7) return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;

    return date.toLocaleDateString();
  }

  /**
   * Get action button state for a user
   * Returns what action button should be shown (if any)
   */
  async getActionButtonState(userId: number): Promise<{
    showButton: boolean;
    buttonText: string;
    buttonAction: 'send_request' | 'accept_request' | 'view_connection' | null;
    requestId?: number;
  }> {
    try {
      const status = await this.checkConnectionStatus(userId);

      switch (status.status) {
        case 'none':
          return {
            showButton: true,
            buttonText: 'Connect',
            buttonAction: 'send_request',
          };

        case 'request_sent':
          return {
            showButton: true,
            buttonText: 'Request Sent',
            buttonAction: null,
          };

        case 'request_received':
          return {
            showButton: true,
            buttonText: 'Accept Request',
            buttonAction: 'accept_request',
            requestId: status.request_id,
          };

        case 'connected':
          return {
            showButton: true,
            buttonText: 'View Connection',
            buttonAction: 'view_connection',
          };

        default:
          return {
            showButton: false,
            buttonText: '',
            buttonAction: null,
          };
      }
    } catch (error) {
      console.error('Error getting button state:', error);
      return {
        showButton: false,
        buttonText: '',
        buttonAction: null,
      };
    }
  }
}

// Create and export singleton instance
const connectionService = new ConnectionService();
export default connectionService;