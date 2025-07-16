// HabitCard component - Display single habit dengan progress
// Props: habit, onLog, onEdit, onDelete
import React from 'react';
import Button from '../UI/Button';

const HabitCard = ({ habit, onLog, onEdit, onDelete }) => {
  const formatDate = (dateString) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('id-ID');
  };

  const calculateProgress = () => {
    if (!habit.total_days || habit.total_days === 0) return 0;
    return Math.round((habit.completed_days / habit.total_days) * 100);
  };

  const getDaysRemaining = () => {
    const endDate = new Date(habit.end_date);
    const today = new Date();
    const diffTime = endDate - today;
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return Math.max(0, diffDays);
  };

  const getStatusColor = () => {
    const progress = calculateProgress();
    if (progress >= 80) return 'bg-green-500';
    if (progress >= 50) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  const isActive = () => {
    const today = new Date();
    const startDate = new Date(habit.start_date);
    const endDate = new Date(habit.end_date);
    return today >= startDate && today <= endDate;
  };

  return (
      <div className="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start mb-4">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <h3 className="text-lg font-semibold text-gray-900">
                {habit.title}
              </h3>
              {isActive() && (
                  <span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded-full">
                Active
              </span>
              )}
            </div>

            <div className="text-sm text-gray-500 mb-3">
              <div>📅 {formatDate(habit.start_date)} - {formatDate(habit.end_date)}</div>
              {habit.reminder_time && (
                  <div>⏰ Daily reminder: {habit.reminder_time}</div>
              )}
              <div>📊 {habit.completed_days || 0} / {habit.total_days || 0} days completed</div>
            </div>

            {/* Progress Bar */}
            <div className="mb-3">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm font-medium text-gray-700">Progress</span>
                <span className="text-sm text-gray-600">{calculateProgress()}%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                    className={`h-2 rounded-full transition-all duration-300 ${getStatusColor()}`}
                    style={{ width: `${calculateProgress()}%` }}
                />
              </div>
            </div>

            {/* Days Remaining */}
            <div className="text-sm text-gray-600">
              {getDaysRemaining() > 0 ? (
                  <span>🗓️ {getDaysRemaining()} days remaining</span>
              ) : (
                  <span className="text-red-600">⏰ Habit period ended</span>
              )}
            </div>
          </div>

          <div className="flex flex-col space-y-2 ml-4">
            {isActive() && onLog && (
                <Button
                    variant="success"
                    size="sm"
                    onClick={() => onLog(habit)}
                >
                  Log Today
                </Button>
            )}

            {onEdit && (
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => onEdit(habit)}
                >
                  Edit
                </Button>
            )}

            {onDelete && (
                <Button
                    variant="danger"
                    size="sm"
                    onClick={() => onDelete(habit.id)}
                >
                  Delete
                </Button>
            )}
          </div>
        </div>
      </div>
  );
};

export default HabitCard;