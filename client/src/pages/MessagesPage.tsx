import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { chatService } from '../services/chatService';
import type { ChatPreview } from '../types/message';
import profileService from '../services/profileService';
import OnlineIndicator from '../components/OnlineIndicator';

const MessagesPage = () => {
    const [chats, setChats] = useState<ChatPreview[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchChats = async () => {
            try {
                const data = await chatService.getChats();
                setChats(data || []);
            } catch (err: any) {
                setError(err.message || 'Failed to load messages');
            } finally {
                setLoading(false);
            }
        };

        fetchChats();
    }, []);

    if (loading) return (
        <div className="page-container">
            <div className="loading-container">
                <div className="spinner"></div>
                <p>Loading conversations...</p>
            </div>
        </div>
    );

    if (error) return (
        <div className="page-container">
            <div className="error-container">
                <p>{error}</p>
                <button onClick={() => window.location.reload()} className="btn-primary">Retry</button>
            </div>
        </div>
    );

    return (
        <div className="page-container">
            <div className="page-header">
                <div className="header-left">
                    <h1>Messages</h1>
                    <p className="page-subtitle">Your active conversations</p>
                </div>
            </div>

            {chats.length === 0 ? (
                <div className="empty-state">
                    <p>No conversations yet.</p>
                    <button onClick={() => navigate('/connections')} className="btn-primary">
                        Start a chat with connections
                    </button>
                </div>
            ) : (
                <div className="messages-list-container" style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {chats.map(chat => (
                        <div
                            key={chat.user_id}
                            className="message-preview-card"
                            onClick={() => navigate(`/chat/${chat.user_id}`)}
                            style={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: '1rem',
                                padding: '1rem',
                                background: 'rgba(255, 255, 255, 0.03)',
                                borderRadius: '12px',
                                border: '1px solid rgba(255, 255, 255, 0.05)',
                                cursor: 'pointer',
                                transition: 'all 0.2s'
                            }}
                            onMouseEnter={(e) => {
                                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.05)';
                                e.currentTarget.style.transform = 'translateY(-2px)';
                            }}
                            onMouseLeave={(e) => {
                                e.currentTarget.style.background = 'rgba(255, 255, 255, 0.03)';
                                e.currentTarget.style.transform = 'none';
                            }}
                        >
                            <div className="preview-avatar">
                                <img
                                    src={profileService.getProfilePictureUrl(chat.profile_picture_url)}
                                    alt={chat.display_name}
                                    style={{
                                        width: '60px',
                                        height: '60px',
                                        borderRadius: '50%',
                                        objectFit: 'cover',
                                        border: '2px solid var(--background-color)'
                                    }}
                                />
                            </div>

                            <div className="preview-content" style={{ flex: 1 }}>
                                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.25rem', alignItems: 'center' }}>
                                    <div style={{display: 'flex', alignItems: 'center', gap: '8px'}}>
                                        <h3 style={{ margin: 0, fontSize: '1.1rem', color: '#fff' }}>{chat.display_name}</h3>
                                        <OnlineIndicator userId={chat.user_id} />
                                    </div>
                                    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                                        <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                                            {new Date(chat.last_message_time).toLocaleDateString()}
                                        </span>
                                        {chat.unread_count > 0 && (
                                            <span style={{
                                                backgroundColor: '#007bff',
                                                color: 'white',
                                                borderRadius: '50%',
                                                padding: '2px 6px',
                                                fontSize: '0.75rem',
                                                fontWeight: 'bold',
                                                minWidth: '20px',
                                                textAlign: 'center'
                                            }}>
                                                {chat.unread_count}
                                            </span>
                                        )}
                                    </div>
                                </div>
                                <p style={{
                                    margin: 0,
                                    color: 'var(--text-secondary)',
                                    whiteSpace: 'nowrap',
                                    overflow: 'hidden',
                                    textOverflow: 'ellipsis',
                                    maxWidth: '600px'
                                }}>
                                    {chat.last_message || 'No messages yet'}
                                </p>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default MessagesPage;
