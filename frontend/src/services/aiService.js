import { apiRequest } from './api';

export const aiService = {
  getDailySummary: async (date = null) => {
    const params = date ? `?date=${date}` : '';
    // FIXED: Remove duplicate /api - use apiRequest with /api/v1/ai
    return await apiRequest(`/api/v1/ai/daily-summary${params}`, {
      method: 'POST'
    });
  },

  getHabitRecommendation: async (days = 7) => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/ai/habit-recommendation', {
      method: 'POST'
    });
  },

  getInsights: async () => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/ai/insights');
  },

  analyzeActivities: async (days = 7) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/ai/analyze-activities?days=${days}`);
  },

  getProductivityTips: async () => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/ai/productivity-tips');
  }
};