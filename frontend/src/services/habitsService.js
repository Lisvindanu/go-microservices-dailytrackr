import { apiRequest } from './api';

export const habitsService = {
  getAll: async (active = false) => {
    // FIXED: Remove duplicate /api - use apiRequest with /api/v1/habits
    return await apiRequest(`/api/v1/habits${active ? '?active=true' : ''}`);
  },

  getById: async (id) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${id}`);
  },

  create: async (habitData) => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/habits', {
      method: 'POST',
      body: JSON.stringify(habitData),
    });
  },

  update: async (id, habitData) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${id}`, {
      method: 'PUT',
      body: JSON.stringify(habitData),
    });
  },

  delete: async (id) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${id}`, {
      method: 'DELETE',
    });
  },

  createLog: async (habitId, logData) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${habitId}/logs`, {
      method: 'POST',
      body: JSON.stringify(logData),
    });
  },

  getLogs: async (habitId) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${habitId}/logs`);
  },

  getStats: async (habitId) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${habitId}/stats`);
  },

  getComplete: async (habitId) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/habits/${habitId}/complete`);
  }
};