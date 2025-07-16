// Chart component - Simple chart component untuk activity trends
// Props: data, type, title
import React from 'react';

const Chart = ({ data = [], type = 'bar', title }) => {
  if (!data || data.length === 0) {
    return (
        <div className="bg-white rounded-lg border p-6">
          {title && <h3 className="text-lg font-semibold mb-4">{title}</h3>}
          <div className="text-center py-8 text-gray-500">
            <div className="text-4xl mb-2">📊</div>
            <p>No data available for chart</p>
          </div>
        </div>
    );
  }

  const maxValue = Math.max(...data.map(item => item.value));

  const renderBarChart = () => (
      <div className="space-y-3">
        {data.map((item, index) => (
            <div key={index} className="flex items-center space-x-3">
              <div className="w-16 text-sm text-gray-600 text-right">
                {item.label}
              </div>
              <div className="flex-1 bg-gray-200 rounded-full h-4 relative">
                <div
                    className="bg-blue-500 h-4 rounded-full transition-all duration-500"
                    style={{ width: `${(item.value / maxValue) * 100}%` }}
                />
              </div>
              <div className="w-12 text-sm text-gray-700 text-right">
                {item.value}
              </div>
            </div>
        ))}
      </div>
  );

  const renderLineChart = () => (
      <div className="flex items-end space-x-2 h-32">
        {data.map((item, index) => (
            <div key={index} className="flex-1 flex flex-col items-center">
              <div
                  className="bg-blue-500 w-full rounded-t transition-all duration-500"
                  style={{ height: `${(item.value / maxValue) * 100}%` }}
              />
              <div className="text-xs text-gray-600 mt-2 text-center">
                {item.label}
              </div>
            </div>
        ))}
      </div>
  );

  return (
      <div className="bg-white rounded-lg border p-6">
        {title && <h3 className="text-lg font-semibold mb-4">{title}</h3>}
        {type === 'bar' ? renderBarChart() : renderLineChart()}
      </div>
  );
};

export default Chart;