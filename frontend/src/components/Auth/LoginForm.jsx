// LoginForm component - Login form dengan email & password
// Props: onSubmit, loading, error
import React, { useState } from 'react';
import Button from '../UI/Button';
import Input from '../UI/Input';

const LoginForm = ({ onSubmit, loading = false, error = null }) => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <div className="flex items-center space-x-2">
                <span className="text-red-600">❌</span>
                <span className="text-red-700 text-sm">{error}</span>
              </div>
            </div>
        )}

        <Input
            label="Email Address"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="your@email.com"
            required
        />

        <Input
            label="Password"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Enter your password"
            required
        />

        <Button
            type="submit"
            loading={loading}
            disabled={loading}
            className="w-full"
        >
          {loading ? 'Signing in...' : 'Sign In'}
        </Button>

        <div className="text-center">
          <p className="text-xs text-gray-500">
            Make sure Backend Gateway is running on port 3000
          </p>
        </div>
      </form>
  );
};

export default LoginForm;