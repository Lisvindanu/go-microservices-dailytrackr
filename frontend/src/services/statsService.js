import { apiRequest } from './api';

export const statsService = {
  getDashboard: async () => {
    return await apiRequest('/api/stats/api/v1/stats/dashboard');
  },

  getActivitySummary: async (startDate, endDate) => {
    const params = new URLSearchParams({ start_date: startDate, end_date: endDate });
    return await apiRequest(`/api/stats/api/v1/stats/activities/summary?${params}`);
  }
};