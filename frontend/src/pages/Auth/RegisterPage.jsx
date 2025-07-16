// RegisterPage - Registration page dengan register form
import React, { useState } from 'react';
import { Link, useNavigate, Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/AuthContext';
import { useToast } from '../../contexts/ToastContext';
import RegisterForm from '../../components/Auth/RegisterForm';

const RegisterPage = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const { register, isAuthenticated } = useAuth();
  const { showToast } = useToast();
  const navigate = useNavigate();

  // Redirect if already authenticated
  if (isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  const handleRegister = async (formData) => {
    setLoading(true);
    setError('');

    try {
      await register(formData.username, formData.email, formData.password);
      showToast('success', 'Registration successful! Welcome to DailyTrackr!');
      navigate('/');
    } catch (err) {
      const errorMessage = err.message || 'Registration failed. Please try again.';
      setError(errorMessage);
      showToast('error', errorMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4">
        <div className="max-w-md w-full space-y-8">
          {/* Header */}
          <div className="text-center">
            <div className="mx-auto w-16 h-16 bg-blue-600 rounded-xl flex items-center justify-center mb-4">
              <span className="text-white font-bold text-xl">DT</span>
            </div>
            <h2 className="text-3xl font-bold text-gray-900">
              Create Account
            </h2>
            <p className="mt-2 text-gray-600">
              Start tracking your daily activities and build better habits
            </p>
          </div>

          {/* Register Form */}
          <div className="bg-white py-8 px-6 shadow-sm rounded-lg border">
            <RegisterForm
                onSubmit={handleRegister}
                loading={loading}
                error={error}
            />
          </div>

          {/* Login Link */}
          <div className="text-center">
            <span className="text-gray-600">Already have an account? </span>
            <Link
                to="/login"
                className="text-blue-600 hover:text-blue-700 font-medium transition-colors"
            >
              Sign in
            </Link>
          </div>

          {/* Features Preview */}
          <div className="bg-white rounded-lg border p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">
              What you'll get:
            </h3>
            <div className="space-y-3">
              <div className="flex items-center space-x-3">
                <span className="text-blue-600">📝</span>
                <span className="text-gray-700">Track daily activities with photos</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="text-green-600">🎯</span>
                <span className="text-gray-700">Build and monitor habits</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="text-purple-600">📊</span>
                <span className="text-gray-700">View progress statistics</span>
              </div>
              <div className="flex items-center space-x-3">
                <span className="text-orange-600">🤖</span>
                <span className="text-gray-700">Get AI-powered insights</span>
              </div>
            </div>
          </div>

          {/* Footer */}
          <div className="text-center">
            <p className="text-xs text-gray-500">
              By creating an account, you agree to our{' '}
              <a href="#" className="text-blue-600 hover:text-blue-700">
                Terms of Service
              </a>{' '}
              and{' '}
              <a href="#" className="text-blue-600 hover:text-blue-700">
                Privacy Policy
              </a>
            </p>
          </div>
        </div>
      </div>
  );
};

export default RegisterPage;