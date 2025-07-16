import React from 'react';
import { useAuth } from '../../contexts/AuthContext';

const Navbar = ({ onSidebarToggle }) => {
  const { user, logout } = useAuth();

  return (
    <nav className="bg-white shadow-sm border-b border-gray-200 fixed w-full top-0 z-40">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <button
              onClick={onSidebarToggle}
              className="lg:hidden p-2 rounded-lg text-gray-600 hover:bg-gray-100"
            >
              ☰
            </button>
            <h1 className="ml-2 text-xl font-bold text-gray-900">
              🧠 DailyTrackr
            </h1>
          </div>
          
          <div className="flex items-center space-x-4">
            <span className="text-sm text-gray-700">
              Hi, {user?.username || 'User'}!
            </span>
            <button
              onClick={logout}
              className="text-sm text-gray-500 hover:text-gray-700 px-3 py-2 rounded-lg hover:bg-gray-100"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;