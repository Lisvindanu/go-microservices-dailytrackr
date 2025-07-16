// ActivityCard component - Display single activity item
// Props: activity, onEdit, onDelete
import React from 'react';
import Button from '../UI/Button';

const ActivityCard = ({ activity, onEdit, onDelete }) => {
  const formatTime = (timeString) => {
    if (!timeString) return '-';
    const date = new Date(timeString);
    return date.toLocaleTimeString('id-ID', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatCurrency = (amount) => {
    if (!amount) return null;
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0
    }).format(amount);
  };

  return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {activity.title}
            </h3>

            <div className="flex items-center space-x-4 text-sm text-gray-500 mb-3">
            <span className="flex items-center space-x-1">
              <span>🕐</span>
              <span>{formatTime(activity.start_time)}</span>
            </span>

              <span className="flex items-center space-x-1">
              <span>⏱️</span>
              <span>{activity.duration_mins} mins</span>
            </span>

              {activity.cost > 0 && (
                  <span className="flex items-center space-x-1 text-green-600">
                <span>💰</span>
                <span>{formatCurrency(activity.cost)}</span>
              </span>
              )}
            </div>

            {activity.note && (
                <p className="text-gray-600 text-sm mb-3 leading-relaxed">
                  {activity.note}
                </p>
            )}

            {activity.photo_url && (
                <div className="mb-3">
                  <img
                      src={activity.photo_url}
                      alt="Activity"
                      className="w-24 h-24 object-cover rounded-lg border shadow-sm"
                  />
                </div>
            )}
          </div>

          <div className="flex flex-col space-y-2 ml-4">
            {onEdit && (
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => onEdit(activity)}
                >
                  Edit
                </Button>
            )}

            {onDelete && (
                <Button
                    variant="danger"
                    size="sm"
                    onClick={() => onDelete(activity.id)}
                >
                  Delete
                </Button>
            )}
          </div>
        </div>
      </div>
  );
};

export default ActivityCard;