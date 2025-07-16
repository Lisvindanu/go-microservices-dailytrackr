import React from 'react';
import { Link } from 'react-router-dom';

const Sidebar = ({ currentPath }) => {
  const menuItems = [
    { path: '/dashboard', label: 'Dashboard', icon: '📊' },
    { path: '/activities', label: 'Activities', icon: '📝' },
    { path: '/habits', label: 'Habits', icon: '🎯' },
    { path: '/stats', label: 'Statistics', icon: '📈' },
    { path: '/profile', label: 'Profile', icon: '👤' }
  ];

  return (
      <aside className="w-64 bg-white border-r border-gray-200 fixed left-0 top-16 bottom-0 z-40">
        <div className="p-4">
          <nav className="space-y-2">
            {menuItems.map((item) => (
                <Link
                    key={item.path}
                    to={item.path}
                    className={`flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors ${
                        currentPath === item.path
                            ? 'bg-blue-50 text-blue-700 border border-blue-200'
                            : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                    }`}
                >
                  <span className="text-lg">{item.icon}</span>
                  <span className="font-medium">{item.label}</span>
                </Link>
            ))}
          </nav>
        </div>
      </aside>
  );
};

export default Sidebar;