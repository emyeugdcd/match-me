import authService from './authService';
import type { Profile, BioData } from '../types';

// const API_URL = 'http://localhost:8080';

export interface Recommendation {
    user_id: number;
    profile: Profile;
    bio_data: BioData;
    match_score: number;
    is_viewed: boolean;
    is_dismissed: boolean;
}

export const recommendationService = {
    getRecommendations: async (): Promise<Recommendation[]> => {
        const response = await authService.authenticatedFetch('/recommendations');
        if (!response.ok) {
            throw new Error('Failed to fetch recommendations');
        }
        const data = await response.json();
        const items: { user_id: number; match_score: number }[] = data.recommendations;

        if (!items || items.length === 0) {
            return [];
        }

        // Fetch details for each item
        const promises = items.map(async (item) => {
            const id = item.user_id;
            try {
                const [profileRes, bioRes, basicRes] = await Promise.all([
                    authService.authenticatedFetch(`/users/${id}/profile`),
                    authService.authenticatedFetch(`/users/${id}/bio`),
                    authService.authenticatedFetch(`/users/${id}`)
                ]);

                if (!profileRes.ok || !bioRes.ok || !basicRes.ok) return null;

                const profileData = await profileRes.json();
                const bioData = await bioRes.json();
                const basicData = await basicRes.json();

                // Merge basic info into profile to match the Profile interface
                const fullProfile: Profile = {
                    ...profileData,
                    first_name: basicData.first_name,
                    last_name: basicData.last_name,
                    username: basicData.username,
                    profile_picture_url: basicData.profile_picture_url,
                    user_id: id,
                    // Ensure other required fields are present if missing from partial responses
                    id: basicData.id || id,
                    is_complete: true, // Assumed
                    created_at: '', // Not critical for display
                    updated_at: ''
                };

                // Construct a recommendation object. 
                // Note: We lose the server-calculated match_score here unless we fetch it separately or recalculate it.
                // For this requirement, we might need to accept that match_score is missing or calculate it client-side.
                // Or, the IDs endpoint could return {id, score}. 
                // But the requirement says "only a list of ids". 
                // So we'll assume score is not critical or we mock it.

                return {
                    user_id: id,
                    profile: fullProfile,
                    bio_data: bioData,
                    match_score: item.match_score,
                    is_viewed: true,
                    is_dismissed: false
                } as Recommendation;
            } catch {
                return null;
            }
        });

        const results = await Promise.all(promises);
        return results.filter((r): r is Recommendation => r !== null);
    },

    dismissRecommendation: async (userId: number): Promise<void> => {
        const response = await authService.authenticatedFetch(`/recommendations/${userId}/dismiss`, {
            method: 'POST'
        });
        if (!response.ok) {
            throw new Error('Failed to dismiss recommendation');
        }
    }
};
