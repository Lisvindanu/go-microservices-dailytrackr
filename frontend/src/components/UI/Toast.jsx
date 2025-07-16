// Toast component - Notification toast (success, error, info, warning)
// Props: type, message, onClose
import React from 'react';

const Toast = ({ type, message, onClose }) => {
  const getIcon = () => {
    switch (type) {
      case 'success': return '✅';
      case 'error': return '❌';
      case 'warning': return '⚠️';
      case 'info': return 'ℹ️';
      default: return 'ℹ️';
    }
  };

  const getBackgroundColor = () => {
    switch (type) {
      case 'success': return 'bg-green-500';
      case 'error': return 'bg-red-500';
      case 'warning': return 'bg-yellow-500';
      case 'info': return 'bg-blue-500';
      default: return 'bg-blue-500';
    }
  };

  return (
      <div className={`flex items-center p-4 rounded-lg shadow-lg text-white ${getBackgroundColor()}`}>
        <span className="text-lg mr-3">{getIcon()}</span>
        <span className="flex-1 text-sm font-medium">{message}</span>
        {onClose && (
            <button
                onClick={onClose}
                className="ml-4 text-white hover:text-gray-200 focus:outline-none"
            >
              ✕
            </button>
        )}
      </div>
  );
};

export default Toast;