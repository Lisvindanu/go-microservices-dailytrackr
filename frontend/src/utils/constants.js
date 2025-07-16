// constants - Application constants

export const API_ENDPOINTS = {
  // Auth endpoints
  AUTH: {
    LOGIN: '/users/auth/login',
    REGISTER: '/users/auth/register',
    PROFILE: '/users/api/v1/users/profile'
  },

  // Activity endpoints
  ACTIVITIES: {
    BASE: '/activities/api/v1/activities',
    PHOTO: (id) => `/activities/api/v1/activities/${id}/photo`
  },

  // Habit endpoints
  HABITS: {
    BASE: '/habits/api/v1/habits',
    LOGS: (id) => `/habits/api/v1/habits/${id}/logs`,
    STATS: (id) => `/habits/api/v1/habits/${id}/stats`
  },

  // Stats endpoints
  STATS: {
    DASHBOARD: '/stats/api/v1/stats/dashboard',
    ACTIVITIES: '/stats/api/v1/stats/activities/summary'
  },

  // AI endpoints
  AI: {
    DAILY_SUMMARY: '/ai/api/v1/ai/daily-summary',
    INSIGHTS: '/ai/api/v1/ai/insights',
    HABIT_RECOMMENDATION: '/ai/api/v1/ai/habit-recommendation'
  }
};

export const ROUTES = {
  LOGIN: '/login',
  REGISTER: '/register',
  DASHBOARD: '/',
  ACTIVITIES: '/activities',
  HABITS: '/habits',
  STATS: '/stats',
  PROFILE: '/profile'
};

export const HABIT_STATUS = {
  DONE: 'DONE',
  SKIPPED: 'SKIPPED',
  FAILED: 'FAILED'
};

export const HABIT_STATUS_LABELS = {
  [HABIT_STATUS.DONE]: '✅ Done',
  [HABIT_STATUS.SKIPPED]: '⏭️ Skipped',
  [HABIT_STATUS.FAILED]: '❌ Failed'
};

export const TOAST_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info'
};

export const PAGINATION = {
  DEFAULT_PAGE_SIZE: 10,
  MAX_PAGE_SIZE: 50
};

export const FILE_UPLOAD = {
  MAX_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_TYPES: ['image/jpeg', 'image/png', 'image/webp'],
  ALLOWED_EXTENSIONS: ['.jpg', '.jpeg', '.png', '.webp']
};

export const VALIDATION = {
  MIN_PASSWORD_LENGTH: 6,
  MIN_USERNAME_LENGTH: 3,
  MAX_USERNAME_LENGTH: 50,
  MAX_BIO_LENGTH: 500,
  MAX_NOTE_LENGTH: 1000
};

export const LOCAL_STORAGE_KEYS = {
  TOKEN: 'token',
  USER: 'user',
  THEME: 'theme'
};

export const DATE_FORMATS = {
  DISPLAY: 'DD MMMM YYYY',
  API: 'YYYY-MM-DD',
  DATETIME_LOCAL: 'YYYY-MM-DDTHH:mm'
};

export const CURRENCIES = {
  IDR: 'IDR',
  USD: 'USD'
};