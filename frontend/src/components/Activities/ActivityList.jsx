import React from 'react';
import Card from '../UI/Card';
import Button from '../UI/Button';

const ActivityList = ({ activities, loading, onEdit, onDelete }) => {
  if (loading) {
    return (
        <div className="space-y-4">
          {[1, 2, 3].map(i => (
              <div key={i} className="bg-white rounded-lg border p-6">
                <div className="animate-pulse">
                  <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                  <div className="h-3 bg-gray-200 rounded w-1/2 mb-2"></div>
                  <div className="h-3 bg-gray-200 rounded w-2/3"></div>
                </div>
              </div>
          ))}
        </div>
    );
  }

  if (activities.length === 0) {
    return (
        <div className="bg-white rounded-lg border p-8 text-center">
          <div className="text-gray-400 text-6xl mb-4">📝</div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No activities yet</h3>
          <p className="text-gray-500">Start tracking your daily activities to see them here.</p>
        </div>
    );
  }

  return (
      <div className="space-y-4">
        {activities.map((activity) => (
            <ActivityCard
                key={activity.id}
                activity={activity}
                onEdit={onEdit}
                onDelete={onDelete}
            />
        ))}

        {onLoadMore && (
            <div className="text-center pt-4">
              <button
                  onClick={onLoadMore}
                  className="text-blue-600 hover:text-blue-700 text-sm font-medium"
              >
                Load More Activities
              </button>
            </div>
        )}
      </div>
  );
};

export default ActivityList;