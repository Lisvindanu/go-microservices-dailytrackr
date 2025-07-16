const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:3000/api';
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
    const response = await fetch(`${API_BASE_URL}${endpoint}`, config);
    const data = await response.json();

    if (!response.ok) {
      throw new APIError(data.message || 'Request failed', response.status, data);
    }

    return data;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
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

export { APIError };