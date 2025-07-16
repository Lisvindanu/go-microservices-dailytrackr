const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000';
const GATEWAY_URL = process.env.REACT_APP_GATEWAY_URL || 'http://localhost:3000';

class APIError extends Error {
  constructor(message, status, data) {
    super(message);
    this.status = status;
    this.data = data;
  }
}

export const apiRequest = async (endpoint, options = {}) => {
  const token = localStorage.getItem('token');

  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
    ...options,
  };

  try {
    // SIMPLIFIED: Direct URL construction without complex mapping
    let fullUrl;

    // Handle auth endpoints specifically
    if (endpoint.startsWith('/auth/')) {
      // Auth routes go to /api/users/auth/*
      fullUrl = `${API_BASE_URL}/api/users${endpoint}`;
    }
    // Handle API v1 endpoints
    else if (endpoint.startsWith('/api/v1/users')) {
      // User API routes go to /api/users/api/v1/users/*
      fullUrl = `${API_BASE_URL}/api/users${endpoint}`;
    }
    else if (endpoint.startsWith('/api/v1/habits')) {
      // Habit API routes go to /api/habits/api/v1/habits/*
      fullUrl = `${API_BASE_URL}/api/habits${endpoint}`;
    }
    else if (endpoint.startsWith('/api/v1/activities')) {
      // Activity API routes go to /api/activities/api/v1/activities/*
      fullUrl = `${API_BASE_URL}/api/activities${endpoint}`;
    }
    else if (endpoint.startsWith('/api/v1/stats')) {
      // Stats API routes go to /api/stats/api/v1/stats/*
      fullUrl = `${API_BASE_URL}/api/stats${endpoint}`;
    }
    else if (endpoint.startsWith('/api/v1/ai')) {
      // AI API routes go to /api/ai/api/v1/ai/*
      fullUrl = `${API_BASE_URL}/api/ai${endpoint}`;
    }
    else {
      // Default: use endpoint as is
      fullUrl = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
    }

    console.log('API Request:', fullUrl); // Debug log

    const response = await fetch(fullUrl, config);

    // Handle non-JSON responses (like 404 HTML pages)
    let data;
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      data = await response.json();
    } else {
      // If not JSON, create error response
      data = {
        success: false,
        message: `Service not available (${response.status})`,
        error: await response.text()
      };
    }

    if (!response.ok) {
      throw new APIError(data.message || 'Request failed', response.status, data);
    }

    return data;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }

    // Handle fetch errors (network issues)
    console.log('Network error:', error.message);
    throw new APIError('Network error - Make sure Gateway is running on port 3000', 0, null);
  }
};

export const checkGatewayHealth = async () => {
  try {
    const response = await fetch(`${GATEWAY_URL}/`);
    return await response.json();
  } catch (error) {
    throw new APIError('Gateway not available on port 3000', 0, null);
  }
};

// Specific API functions
export const authAPI = {
  login: (credentials) => apiRequest('/auth/login', {
    method: 'POST',
    body: JSON.stringify(credentials)
  }),
  register: (userData) => apiRequest('/auth/register', {
    method: 'POST',
    body: JSON.stringify(userData)
  })
};

export const userAPI = {
  getProfile: () => apiRequest('/api/v1/users/profile'),
  updateProfile: (data) => apiRequest('/api/v1/users/profile', {
    method: 'PUT',
    body: JSON.stringify(data)
  })
};

export const habitAPI = {
  getAll: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiRequest(`/api/v1/habits${query ? '?' + query : ''}`);
  },
  getById: (id) => apiRequest(`/api/v1/habits/${id}`),
  create: (data) => apiRequest('/api/v1/habits', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  update: (id, data) => apiRequest(`/api/v1/habits/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  delete: (id) => apiRequest(`/api/v1/habits/${id}`, {
    method: 'DELETE'
  }),
  logHabit: (id, data) => apiRequest(`/api/v1/habits/${id}/logs`, {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  getStats: (id) => apiRequest(`/api/v1/habits/${id}/stats`)
};

export const activityAPI = {
  getAll: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiRequest(`/api/v1/activities${query ? '?' + query : ''}`);
  },
  getById: (id) => apiRequest(`/api/v1/activities/${id}`),
  create: (data) => apiRequest('/api/v1/activities', {
    method: 'POST',
    body: JSON.stringify(data)
  }),
  update: (id, data) => apiRequest(`/api/v1/activities/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data)
  }),
  delete: (id) => apiRequest(`/api/v1/activities/${id}`, {
    method: 'DELETE'
  })
};

export const statsAPI = {
  getDashboard: () => apiRequest('/api/v1/stats/dashboard'),
  getActivityChart: (params = {}) => {
    const query = new URLSearchParams(params).toString();
    return apiRequest(`/api/v1/stats/activities/chart${query ? '?' + query : ''}`);
  },
  getHabitProgress: () => apiRequest('/api/v1/stats/habits/progress')
};

export { APIError };