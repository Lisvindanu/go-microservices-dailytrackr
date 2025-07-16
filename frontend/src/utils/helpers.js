// helpers - Utility helper functions

export const formatDate = (date, options = {}) => {
    if (!date) return '';

    const dateObj = new Date(date);
    const defaultOptions = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        ...options
    };

    return dateObj.toLocaleDateString('id-ID', defaultOptions);
};

export const formatCurrency = (amount, currency = 'IDR') => {
    if (amount === null || amount === undefined) return 'Rp 0';

    return new Intl.NumberFormat('id-ID', {
        style: 'currency',
        currency: currency,
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
    }).format(amount);
};

export const formatTime = (timeString) => {
    if (!timeString) return '';

    const date = new Date(timeString);
    return date.toLocaleTimeString('id-ID', {
        hour: '2-digit',
        minute: '2-digit'
    });
};

export const formatDuration = (minutes) => {
    if (!minutes) return '0 min';

    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;

    if (hours > 0) {
        return `${hours}h ${mins}m`;
    }
    return `${mins}m`;
};

export const getRelativeTime = (date) => {
    if (!date) return '';

    const now = new Date();
    const targetDate = new Date(date);
    const diffInSeconds = Math.floor((now - targetDate) / 1000);

    if (diffInSeconds < 60) return 'Just now';
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)} min ago`;
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)} hours ago`;
    if (diffInSeconds < 2592000) return `${Math.floor(diffInSeconds / 86400)} days ago`;

    return formatDate(date);
};

export const calculateDaysBetween = (startDate, endDate) => {
    const start = new Date(startDate);
    const end = new Date(endDate);
    const diffTime = Math.abs(end - start);
    return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
};

export const isToday = (date) => {
    const today = new Date();
    const targetDate = new Date(date);

    return today.toDateString() === targetDate.toDateString();
};

export const truncateText = (text, maxLength = 100) => {
    if (!text || text.length <= maxLength) return text;

    return text.substring(0, maxLength).trim() + '...';
};

export const validateEmail = (email) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
};

export const debounce = (func, wait) => {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
};

export const generateId = () => {
    return Math.random().toString(36).substr(2, 9);
};

export const classNames = (...classes) => {
    return classes.filter(Boolean).join(' ');
};