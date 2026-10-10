import React, { useState, useEffect } from 'react';
import socketService from '../services/socketService';
import authService from '../services/authService';

interface OnlineIndicatorProps {
    userId: number;
    showLabel?: boolean;
}

const OnlineIndicator: React.FC<OnlineIndicatorProps> = ({ userId, showLabel = false }) => {
    const [isOnline, setIsOnline] = useState(false);

    useEffect(() => {
        console.log('OnlineIndicator: register handlers for user', userId);
        // Listen for status updates
        const handleStatusUpdate = (payload: unknown) => {
            const data = payload as { user_id?: number; status?: string };
            console.log('OnlineIndicator: user_status received', payload);
            if (data?.user_id === userId) {
                setIsOnline(data?.status === 'online');
            }
        };

        socketService.on('user_status', handleStatusUpdate);
        
        // Listen for presence list when first connecting
        const handlePresenceList = (payload: unknown) => {
            const data = payload as { user_ids?: number[] };
            console.log('OnlineIndicator: presence_list received', payload);
            const ids: number[] = data?.user_ids || [];
            setIsOnline(ids.includes(userId));
        };

        socketService.on('presence_list', handlePresenceList);

        // Also fetch current status from the server to avoid race conditions
        const fetchStatus = async () => {
            try {
                const res = await authService.authenticatedFetch(`/users/${userId}/status`);
                if (!res.ok) return;
                const data = await res.json();
                if (typeof data.online === 'boolean') {
                    setIsOnline(data.online);
                }
            } catch {
                // ignore
            }
        };

        // If socketService has cached status, use it immediately
        try {
            const cached = socketService.getCachedUserStatus(userId);
            if (typeof cached === 'boolean') {
                console.log('OnlineIndicator: using cached socket status', cached);
                setIsOnline(cached);
            }
        } catch {
            // swallow
        }

        fetchStatus();



        return () => {
            console.log('OnlineIndicator: unregister handlers for user', userId);
            socketService.off('user_status', handleStatusUpdate);
            socketService.off('presence_list', handlePresenceList);
        };
    }, [userId]);

    return (
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <div
                style={{
                    width: '10px',
                    height: '10px',
                    borderRadius: '999px',
                    backgroundColor: isOnline ? 'var(--success-color)' : 'rgba(128,128,128,0.6)',
                    border: '2px solid #fff',
                }}
                title={isOnline ? 'Online' : 'Offline'}
            />
            {showLabel && (
                <span style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
                    {isOnline ? 'Online' : 'Offline'}
                </span>
            )}
        </div>
    );
};

export default OnlineIndicator;
