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
    // Build full URL with proper path
    const fullUrl = `${API_BASE_URL}${endpoint.startsWith('/') ? endpoint : '/' + endpoint}`;
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
    const response = await fetch(`${GATEWAY_URL}/health`);
    return await response.json();
  } catch (error) {
    throw new APIError('Gateway not available on port 3000', 0, null);
  }
};

export { APIError };