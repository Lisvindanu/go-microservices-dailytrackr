import { apiRequest } from './api';

export const aiService = {
  getDailySummary: async (date = null) => {
    const params = date ? `?date=${date}` : '';
    return await apiRequest(`/ai/api/v1/ai/daily-summary${params}`);
  },

  getHabitRecommendation: async (days = 7) => {
    return await apiRequest(`/ai/api/v1/ai/habit-recommendation?days=${days}`);
  },

  getInsights: async () => {
    return await apiRequest('/ai/api/v1/ai/insights');
  }
};