import authService from './authService';
import type { Message, ChatPreview } from '../types/message';

const API_URL = 'http://localhost:8080';

const getAuthHeaders = () => {
    const headers = authService.getAuthHeaders();
    return {
        'Content-Type': 'application/json',
        ...headers,
    };
};

export const chatService = {
    sendMessage: async (receiverId: number, content: string): Promise<void> => {
        const response = await fetch(`${API_URL}/chats/messages`, {
            method: 'POST',
            headers: getAuthHeaders(),
            body: JSON.stringify({ receiver_id: receiverId, content }),
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.error || 'Failed to send message');
        }
    },

    getMessages: async (userId: number, limit = 50, offset = 0): Promise<Message[]> => {
        const response = await fetch(`${API_URL}/chats/${userId}/messages?limit=${limit}&offset=${offset}`, {
            headers: getAuthHeaders(),
        });

        if (!response.ok) {
            throw new Error('Failed to fetch messages');
        }

        return response.json();
    },

    getChats: async (): Promise<ChatPreview[]> => {
        const response = await fetch(`${API_URL}/chats`, {
            headers: getAuthHeaders(),
        });

        if (!response.ok) {
            throw new Error('Failed to fetch chats');
        }

        return response.json();
    },
};
