import { apiRequest } from './api';

export const activitiesService = {
  getAll: async (params = {}) => {
    const queryString = new URLSearchParams(params).toString();
    // FIXED: Remove duplicate /api - use apiRequest with /api/v1/activities
    return await apiRequest(`/api/v1/activities${queryString ? `?${queryString}` : ''}`);
  },

  getById: async (id) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/activities/${id}`);
  },

  create: async (activityData) => {
    // FIXED: Remove duplicate /api
    return await apiRequest('/api/v1/activities', {
      method: 'POST',
      body: JSON.stringify(activityData),
    });
  },

  update: async (id, activityData) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/activities/${id}`, {
      method: 'PUT',
      body: JSON.stringify(activityData),
    });
  },

  delete: async (id) => {
    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/activities/${id}`, {
      method: 'DELETE',
    });
  },

  uploadPhoto: async (id, photoFile) => {
    const formData = new FormData();
    formData.append('photo', photoFile);

    // FIXED: Remove duplicate /api
    return await apiRequest(`/api/v1/activities/${id}/photo`, {
      method: 'POST',
      headers: {}, // Remove Content-Type for FormData
      body: formData,
    });
  }
};