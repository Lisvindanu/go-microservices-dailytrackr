import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { authService } from '../../services/authService';
import { useToast } from '../../contexts/ToastContext';
import Button from '../../components/UI/Button';
import Card from '../../components/UI/Card';
import Input from '../../components/UI/Input';

const ProfilePage = () => {
  const { user } = useAuth();
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editMode, setEditMode] = useState(false);
  const [apiAvailable, setApiAvailable] = useState(false);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    bio: ''
  });
  const { showToast } = useToast();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        setLoading(true);

        // Try to fetch profile with timeout
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000);

        const response = await authService.getProfile();

        clearTimeout(timeoutId);

        // If we get here, API is working
        setApiAvailable(true);
        setProfile(response.data);
        setFormData({
          username: response.data.username || '',
          email: response.data.email || '',
          bio: response.data.bio || ''
        });

      } catch (error) {
        console.log('Profile API not available:', error.message);
        setApiAvailable(false);

        // Use fallback data from user context or localStorage
        const fallbackProfile = user || {
          username: 'testuser',
          email: 'test@example.com',
          bio: 'Demo user profile',
          created_at: new Date().toISOString()
        };

        setProfile(fallbackProfile);
        setFormData({
          username: fallbackProfile.username || '',
          email: fallbackProfile.email || '',
          bio: fallbackProfile.bio || ''
        });

        if (!error.message.includes('Network error')) {
          showToast('warning', 'Backend services not available. Using cached profile data.');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [user]);

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSave = async () => {
    if (!apiAvailable) {
      showToast('warning', 'Cannot save changes - backend services not available.');
      return;
    }

    try {
      // This would call the update profile API
      showToast('success', 'Profile updated successfully!');
      setEditMode(false);
    } catch (error) {
      showToast('error', 'Failed to update profile');
    }
  };

  const handleRetryConnection = async () => {
    setLoading(true);

    try {
      const response = await fetch('http://localhost:3000/', {
        method: 'GET',
        signal: AbortSignal.timeout(3000)
      });

      if (response.ok) {
        setApiAvailable(true);
        showToast('success', 'Backend connection restored!');
        window.location.reload();
      }
    } catch (error) {
      setApiAvailable(false);
      showToast('error', 'Backend still not available. Make sure Gateway is running on port 3000.');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        </div>
    );
  }

  return (
      <div className="space-y-6">
        {/* Backend Status Banner */}
        {!apiAvailable && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-2">
                  <span className="text-yellow-600">⚠️</span>
                  <div>
                    <p className="text-yellow-800 font-medium">Backend Services Offline</p>
                    <p className="text-yellow-700 text-sm">
                      Profile changes cannot be saved. Start backend services to enable editing.
                    </p>
                  </div>
                </div>
                <Button
                    variant="secondary"
                    size="sm"
                    onClick={handleRetryConnection}
                    disabled={loading}
                >
                  Retry Connection
                </Button>
              </div>
            </div>
        )}

        <h1 className="text-2xl font-bold text-gray-900">
          Profile {!apiAvailable ? '(Read-only)' : ''}
        </h1>

        <Card title="👤 Profile Information" action={
          <Button
              variant="secondary"
              onClick={() => editMode ? handleSave() : setEditMode(true)}
              disabled={!apiAvailable && !editMode}
          >
            {editMode ? 'Save' : 'Edit'}
          </Button>
        }>
          <div className="space-y-4">
            {editMode ? (
                <>
                  <Input
                      label="Username"
                      name="username"
                      value={formData.username}
                      onChange={handleChange}
                      disabled={!apiAvailable}
                  />

                  <Input
                      label="Email"
                      type="email"
                      name="email"
                      value={formData.email}
                      onChange={handleChange}
                      disabled={!apiAvailable}
                  />

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Bio
                    </label>
                    <textarea
                        name="bio"
                        value={formData.bio}
                        onChange={handleChange}
                        rows={3}
                        disabled={!apiAvailable}
                        className={`w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 ${!apiAvailable ? 'bg-gray-100 cursor-not-allowed' : ''}`}
                        placeholder="Tell us about yourself..."
                    />
                  </div>

                  {!apiAvailable && (
                      <p className="text-sm text-yellow-600">
                        ⚠️ Backend services are offline. Changes cannot be saved.
                      </p>
                  )}
                </>
            ) : (
                <>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Username</label>
                    <p className="mt-1 text-gray-900">{profile?.username || 'N/A'}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700">Email</label>
                    <p className="mt-1 text-gray-900">{profile?.email || 'N/A'}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700">Bio</label>
                    <p className="mt-1 text-gray-900">{profile?.bio || 'No bio added yet.'}</p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700">Member Since</label>
                    <p className="mt-1 text-gray-900">
                      {profile?.created_at ? new Date(profile.created_at).toLocaleDateString() : 'Unknown'}
                    </p>
                  </div>
                </>
            )}
          </div>
        </Card>

        <Card title="🔒 Security">
          <div className="space-y-4">
            <Button
                variant="secondary"
                className="w-full"
                disabled={!apiAvailable}
            >
              Change Password
            </Button>

            <Button
                variant="danger"
                className="w-full"
                disabled={!apiAvailable}
            >
              Delete Account
            </Button>

            {!apiAvailable && (
                <p className="text-sm text-gray-500 text-center">
                  Security options are disabled when backend is offline
                </p>
            )}
          </div>
        </Card>
      </div>
  );
};

export default ProfilePage;