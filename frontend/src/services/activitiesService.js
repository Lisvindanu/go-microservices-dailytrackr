import { apiRequest } from './api';

export const activitiesService = {
  getAll: async (params = {}) => {
    const queryString = new URLSearchParams(params).toString();
    return await apiRequest(`/api/activities/api/v1/activities${queryString ? `?${queryString}` : ''}`);
  },

  getById: async (id) => {
    return await apiRequest(`/api/activities/api/v1/activities/${id}`);
  },

  create: async (activityData) => {
    return await apiRequest('/api/activities/api/v1/activities', {
      method: 'POST',
      body: JSON.stringify(activityData),
    });
  },

  update: async (id, activityData) => {
    return await apiRequest(`/api/activities/api/v1/activities/${id}`, {
      method: 'PUT',
      body: JSON.stringify(activityData),
    });
  },

  delete: async (id) => {
    return await apiRequest(`/api/activities/api/v1/activities/${id}`, {
      method: 'DELETE',
    });
  },

  uploadPhoto: async (id, photoFile) => {
    const formData = new FormData();
    formData.append('photo', photoFile);

    return await apiRequest(`/api/activities/api/v1/activities/${id}/photo`, {
      method: 'POST',
      headers: {}, // Remove Content-Type for FormData
      body: formData,
    });
  }
};