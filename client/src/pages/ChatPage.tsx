import { useState, useEffect, useRef, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { chatService } from '../services/chatService';
import socketService from '../services/socketService';
import type { Message } from '../types/message';
import authService from '../services/authService';
import OnlineIndicator from '../components/OnlineIndicator';

const ChatPage = () => {
    const { userId } = useParams<{ userId: string }>();
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const [loading, setLoading] = useState(true);
    const [isTyping, setIsTyping] = useState(false);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const typingTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const currentUser = authService.getUser();

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    const fetchMessages = useCallback(async () => {
        if (!userId) return;
        try {
            const data = await chatService.getMessages(parseInt(userId));
            setMessages(data || []);
        } catch (error) {
            console.error('Failed to fetch messages', error);
            alert('Failed to fetch messages: ' + error);
        } finally {
            setLoading(false);
        }
    }, [userId]);

    useEffect(() => {
        fetchMessages();

        // Listen for real-time messages
        const handleNewMessage = (payload: unknown) => {
            const data = payload as { payload: import("../types/message").Message };
            const msg = data.payload;
            // Check if message belongs to this chat
            // Handled above

            // If the message is from the current chat partner OR from us (echoed)
            if (msg.sender_id === parseInt(userId || '0') || msg.sender_id === currentUser?.id) {
                setMessages(prev => {
                    // Avoid duplicates
                    if (prev.some(m => m.id === msg.id)) return prev;
                    return [...prev, msg];
                });
            }
        };

        const handleTyping = (payload: unknown) => {
            const data = payload as { from_user_id?: number };
            if (data?.from_user_id === parseInt(userId || '0')) {
                setIsTyping(true);
                // Clear existing timeout
                if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
                // Auto-hide after 3 seconds
                typingTimeoutRef.current = setTimeout(() => setIsTyping(false), 3000);
            }
        };

        socketService.on('message', handleNewMessage); 
        socketService.on('typing', handleTyping);

        return () => {
            socketService.off('message', handleNewMessage);
            socketService.off('typing', handleTyping);
            if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
        };
    }, [userId, currentUser?.id, fetchMessages]);

    useEffect(() => {
        scrollToBottom();
    }, [messages, isTyping]);

    const handleSend = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!newMessage.trim() || !userId) return;

        try {
            await chatService.sendMessage(parseInt(userId), newMessage);
            setNewMessage('');
            fetchMessages(); // Fallback to ensure UI updates
        } catch (error) {
            console.error('Failed to send message', error);
            alert('Failed to send message: ' + error);
        }
    };

    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNewMessage(e.target.value);

        // Emit typing event
        if (userId) {
            socketService.send('typing', { to_user_id: parseInt(userId) });
        }
    };

    if (loading) return <div>Loading chat...</div>;

    return (
        <div className="chat-container" style={{ maxWidth: '600px', margin: '0 auto', padding: '20px' }}>
            <div className="chat-header" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '10px' }}>
                <h2>Chat</h2>
                {userId && <OnlineIndicator userId={parseInt(userId)} showLabel />}
            </div>

            <div className="messages-list" style={{
                height: '400px',
                overflowY: 'auto',
                border: '1px solid #ccc',
                padding: '10px',
                marginBottom: '20px',
                borderRadius: '8px',
                display: 'flex',
                flexDirection: 'column',
                gap: '10px'
            }}>
                {messages.map((msg) => {
                    const isOwn = msg.sender_id === currentUser?.id;
                    return (
                        <div
                            key={msg.id}
                            style={{
                                alignSelf: isOwn ? 'flex-end' : 'flex-start',
                                backgroundColor: isOwn ? '#007bff' : '#f1f0f0',
                                color: isOwn ? 'white' : 'black',
                                padding: '8px 12px',
                                borderRadius: '12px',
                                maxWidth: '70%',
                                display: 'flex',
                                flexDirection: 'column',
                            }}
                        >
                            <span>{msg.content}</span>
                            <span style={{
                                fontSize: '0.7em',
                                alignSelf: 'flex-end',
                                marginTop: '4px',
                                opacity: 0.7
                            }}>
                                {new Date(msg.created_at).toLocaleString([], { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })}
                            </span>
                        </div>
                    );
                })}
                {isTyping && (
                    <div style={{ alignSelf: 'flex-start', color: '#888', fontStyle: 'italic', fontSize: '0.8em' }}>
                        User is typing...
                    </div>
                )}
                <div ref={messagesEndRef} />
            </div>

            <form onSubmit={handleSend} style={{ display: 'flex', gap: '10px' }}>
                <input
                    type="text"
                    value={newMessage}
                    onChange={handleInput}
                    placeholder="Type a message..."
                    style={{ flex: 1, padding: '10px', borderRadius: '4px', border: '1px solid #ccc' }}
                />
                <button type="submit" style={{ padding: '10px 20px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}>
                    Send
                </button>
            </form>
        </div>
    );
};

export default ChatPage;
