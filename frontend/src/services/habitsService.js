import { apiRequest } from './api';

export const habitsService = {
  getAll: async (active = false) => {
    return await apiRequest(`/api/habits/api/v1/habits${active ? '?active=true' : ''}`);
  },

  create: async (habitData) => {
    return await apiRequest('/api/habits/api/v1/habits', {
      method: 'POST',
      body: JSON.stringify(habitData),
    });
  },

  createLog: async (habitId, logData) => {
    return await apiRequest(`/api/habits/api/v1/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(logData),
    });
  },

  getStats: async (habitId) => {
    return await apiRequest(`/api/habits/api/v1/habits/${habitId}/stats`);
  }
};