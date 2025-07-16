// Services Index - letakkan file ini di folder: frontend/src/services/index.js

// Main API functions
export { apiRequest, checkGatewayHealth, APIError } from './api';

// Service modules
export { authService } from './authService';
export { activitiesService } from './activitiesService';
export { habitsService } from './habitsService';
export { aiService } from './aiService';
export { statsService } from './statsService';

// Legacy exports for backward compatibility
export {
    authAPI,
    userAPI,
    habitAPI,
    activityAPI,
    statsAPI
} from './api';