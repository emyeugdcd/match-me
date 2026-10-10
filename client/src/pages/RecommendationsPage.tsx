import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { recommendationService } from '../services/recommendationService';
import type { Recommendation } from '../services/recommendationService';
import connectionService from '../services/connectionService';
import profileService from '../services/profileService';

const RecommendationsPage: React.FC = () => {
    const [recommendations, setRecommendations] = useState<Recommendation[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        loadRecommendations();
    }, []);

    const loadRecommendations = async () => {
        try {
            setLoading(true);
            const data = await recommendationService.getRecommendations();
            setRecommendations(data || []);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load recommendations');
        } finally {
            setLoading(false);
        }
    };

    const handleConnect = async (userId: number) => {
        try {
            await connectionService.sendConnectionRequest(userId);
            // Remove from list or show success state
            setRecommendations(prev => prev.filter(rec => rec.user_id !== userId));
            alert('Connection request sent!');
        } catch (err) {
            alert(err instanceof Error ? err.message : 'Failed to send request');
        }
    };

    const handleDismiss = async (userId: number) => {
        try {
            await recommendationService.dismissRecommendation(userId);
            setRecommendations(prev => prev.filter(rec => rec.user_id !== userId));
        } catch (err) {
            console.error('Failed to dismiss recommendation:', err);
            // Optional: show error to user
        }
    };

    if (loading) {
        return (
            <div className="page-container">
                <div className="loading-container">
                    <div className="spinner"></div>
                    <p>Finding your best matches...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="page-container">
            <div className="page-header">
                <div className="header-left">
                    <h1>Suggested Matches</h1>
                    <p className="page-subtitle">People we think you'll click with based on your interests.</p>
                </div>
            </div>

            {error && <div className="error-message">{error}</div>}

            {recommendations.length === 0 ? (
                <div className="empty-state">
                    <p>No new recommendations found at the moment.</p>
                    <p>Try updating your profile with more interests!</p>
                </div>
            ) : (
                <div className="recommendations-grid">
                    {recommendations.map((rec) => (
                        <div key={rec.user_id} className="user-card">
                            <div className="match-score-badge">
                                {Math.round(rec.match_score)}% Match
                            </div>
                            <div className="card-header">
                                <img
                                    src={profileService.getProfilePictureUrl(rec.profile?.profile_picture_url) || 'https://via.placeholder.com/150'}
                                    alt={rec.profile?.first_name || 'User'}
                                    className="card-profile-pic"
                                />
                                <h3 className="card-name">{rec.profile?.first_name} {rec.profile?.last_name}</h3>
                                <span className="card-username">@{rec.profile?.username}</span>
                            </div>

                            <div className="card-body">
                                <p className="card-location">📍 {rec.profile?.location || 'Unknown Location'}</p>

                                {rec.bio_data?.interests && rec.bio_data.interests.length > 0 && (
                                    <div className="card-interests">
                                        {rec.bio_data.interests.slice(0, 3).map((interest, i) => (
                                            <span key={i} className="card-tag">{interest}</span>
                                        ))}
                                    </div>
                                )}
                            </div>

                            <div className="card-actions">
                                <button
                                    onClick={() => handleDismiss(rec.user_id)}
                                    className="btn-secondary-card"
                                    style={{ marginRight: '8px' }}
                                >
                                    Dismiss
                                </button>
                                <button
                                    onClick={() => navigate(`/profile/${rec.user_id}`)}
                                    className="btn-secondary-card"
                                >
                                    View Profile
                                </button>
                                <button
                                    onClick={() => handleConnect(rec.user_id)}
                                    className="btn-primary-card"
                                >
                                    Connect
                                </button>
                            </div>
                        </div>
                    ))}
                </div>
            )}


        </div>
    );
};

export default RecommendationsPage;
