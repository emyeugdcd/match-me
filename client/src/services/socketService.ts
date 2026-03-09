import authService from './authService';

type EventHandler = (payload: any) => void;

class SocketService {
    private socket: WebSocket | null = null;
    private listeners: Map<string, EventHandler[]> = new Map();
    private lastPayloads: Map<string, any> = new Map();
    // For user_status we want to cache per-user last status
    private lastUserStatuses: Map<number, any> = new Map();
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

    connect() {
        if (this.socket?.readyState === WebSocket.OPEN) return;

        const token = authService.getToken();
        if (!token) return;

        // Derive WS URL from VITE_API_URL if provided, otherwise default to localhost
        const apiBase = (import.meta.env.VITE_API_URL as string) || 'http://localhost:8080';
        const wsUrl = apiBase.replace(/^http/, 'ws') + `/ws?token=${token}`;

        console.log('SocketService: connecting to', wsUrl);
        this.socket = new WebSocket(wsUrl);

        this.socket.onopen = () => {
            console.log('WebSocket connected');
            console.log('SocketService: connected, will send/receive presence and status events');
            this.reconnectAttempts = 0;
            this.emit('connected', {});
        };

        this.socket.onmessage = (event) => {
            const raw = event.data as string;
            try {
                const data = JSON.parse(raw);
                console.log('SocketService: received event', data.type, data);
                this.emit(data.type, data);
            } catch (error) {
                // Attempt to recover from multiple JSON objects concatenated together
                console.error('Error parsing WS message:', error, 'raw:', raw);
                try {
                    // Extract balanced JSON objects
                    const chunks: string[] = [];
                    let depth = 0;
                    let start = -1;
                    for (let i = 0; i < raw.length; i++) {
                        const ch = raw[i];
                        if (ch === '{') {
                            if (depth === 0) start = i;
                            depth++;
                        } else if (ch === '}') {
                            depth--;
                            if (depth === 0 && start !== -1) {
                                chunks.push(raw.slice(start, i + 1));
                                start = -1;
                            }
                        }
                    }
                    for (const chunk of chunks) {
                        try {
                            const data = JSON.parse(chunk);
                            console.log('SocketService: recovered event', data.type, data);
                            this.emit(data.type, data);
                        } catch (err2) {
                            console.error('SocketService: failed to parse recovered chunk', err2, 'chunk:', chunk);
                        }
                    }
                } catch (recErr) {
                    console.error('SocketService: failed to recover from malformed message', recErr);
                }
            }
        };

        this.socket.onclose = () => {
            console.log('WebSocket disconnected');
            this.attemptReconnect();
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
        }
    }

    send(type: string, payload: any) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ type, payload }));
        } else {
            console.warn('WebSocket not connected, cannot send message');
        }
    }

    on(event: string, handler: EventHandler) {
        if (!this.listeners.has(event)) {
            this.listeners.set(event, []);
        }
        this.listeners.get(event)?.push(handler);
        // If we have cached payloads, call the handler immediately
        if (event === 'user_status') {
            // Call handler for each cached user status
            try {
                console.log('SocketService: invoking handler with cached user_status entries', this.lastUserStatuses.size);
                this.lastUserStatuses.forEach((val, uid) => {
                    handler({ user_id: uid, status: val.status });
                });
            } catch (err) {
                console.error('Error running cached user_status handlers:', err);
            }
        } else if (this.lastPayloads.has(event)) {
            try {
                console.log('SocketService: invoking handler with cached event', event);
                handler(this.lastPayloads.get(event));
            } catch (err) {
                console.error('Error running cached event handler:', err);
            }
        }
    }

    // Returns cached last status for a specific user, if available
    getCachedUserStatus(userId: number): boolean | undefined {
        const entry = this.lastUserStatuses.get(userId);
        if (entry && typeof entry.status !== 'undefined') {
            return entry.status === 'online';
        }
        return undefined;
    }

    off(event: string, handler: EventHandler) {
        const handlers = this.listeners.get(event);
        if (handlers) {
            this.listeners.set(
                event,
                handlers.filter((h) => h !== handler)
            );
        }
    }

    private emit(event: string, payload: any) {
        // Cache the last payload for this event
        try {
            this.lastPayloads.set(event, payload);
        } catch (e) {}
        // Special-case: cache per-user status for user_status events
        if (event === 'user_status' && payload && payload.user_id !== undefined) {
            try {
                this.lastUserStatuses.set(payload.user_id, payload);
            } catch (e) {}
        }
        const handlers = this.listeners.get(event);
        if (handlers) {
            handlers.forEach((handler) => handler(payload));
        }
    }

    private attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 10000);
            console.log(`Attempting reconnect in ${delay}ms...`);
            this.reconnectTimeout = setTimeout(() => this.connect(), delay);
        }
    }
}

const socketService = new SocketService();
export default socketService;
