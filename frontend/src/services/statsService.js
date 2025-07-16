import { apiRequest } from './api';

export const statsService = {
  getDashboard: async () => {
    // FIXED: Remove duplicate /api - use apiRequest with /api/v1/stats
    return await apiRequest('/api/v1/stats/dashboard');
  },

  getActivityChart: async (type = 'daily', period = 7) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/stats/activities/chart?type=${type}&period=${period}`);
  },

  getActivitySummary: async (startDate = null, endDate = null) => {
    let url = '/api/v1/stats/activities/summary';
    if (startDate && endDate) {
      const params = new URLSearchParams({ start_date: startDate, end_date: endDate });
      url += `?${params}`;
    }
    // FIXED: Remove duplicate /api
    return await apiRequest(url);
  },

  getHabitProgress: async () => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/stats/habits/progress');
  },

  getExpenseReport: async () => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/stats/expenses/report');
  }
};