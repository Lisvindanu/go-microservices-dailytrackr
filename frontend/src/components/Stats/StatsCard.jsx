// StatsCard component - Display statistics dengan icon dan value
// Props: title, value, icon, trend, color
import React from 'react';

const StatsCard = ({ title, value, icon, trend, color = 'blue' }) => {
  const colorClasses = {
    blue: 'text-blue-600 bg-blue-50 border-blue-200',
    green: 'text-green-600 bg-green-50 border-green-200',
    purple: 'text-purple-600 bg-purple-50 border-purple-200',
    orange: 'text-orange-600 bg-orange-50 border-orange-200',
    red: 'text-red-600 bg-red-50 border-red-200'
  };

  const getTrendIcon = () => {
    if (!trend) return null;
    return trend > 0 ? '↗️' : trend < 0 ? '↘️' : '➡️';
  };

  const getTrendColor = () => {
    if (!trend) return '';
    return trend > 0 ? 'text-green-600' : trend < 0 ? 'text-red-600' : 'text-gray-600';
  };

  return (
      <div className={`p-6 rounded-lg border ${colorClasses[color]}`}>
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              {icon && <span className="text-2xl">{icon}</span>}
              <h3 className="text-sm font-medium text-gray-700">{title}</h3>
            </div>

            <div className="text-2xl font-bold text-gray-900 mb-1">
              {value}
            </div>

            {trend !== undefined && (
                <div className={`flex items-center space-x-1 text-sm ${getTrendColor()}`}>
                  <span>{getTrendIcon()}</span>
                  <span>{Math.abs(trend)}%</span>
                </div>
            )}
          </div>
        </div>
      </div>
  );
};

export default StatsCard;