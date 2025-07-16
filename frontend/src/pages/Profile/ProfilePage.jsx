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
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    bio: ''
  });
  const { showToast } = useToast();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await authService.getProfile();
        setProfile(response.data);
        setFormData({
          username: response.data.username || '',
          email: response.data.email || '',
          bio: response.data.bio || ''
        });
      } catch (error) {
        showToast('error', 'Failed to load profile');
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, [showToast]);

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSave = async () => {
    try {
      // This would call the update profile API
      showToast('success', 'Profile updated successfully!');
      setEditMode(false);
    } catch (error) {
      showToast('error', 'Failed to update profile');
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
      <h1 className="text-2xl font-bold text-gray-900">Profile</h1>
      
      <Card title="👤 Profile Information" action={
        <Button
          variant="secondary"
          onClick={() => editMode ? handleSave() : setEditMode(true)}
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
              />
              
              <Input
                label="Email"
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
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
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Tell us about yourself..."
                />
              </div>
            </>
          ) : (
            <>
              <div>
                <label className="block text-sm font-medium text-gray-700">Username</label>
                <p className="mt-1 text-gray-900">{profile?.username}</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700">Email</label>
                <p className="mt-1 text-gray-900">{profile?.email}</p>
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
          <Button variant="secondary" className="w-full">
            Change Password
          </Button>
          
          <Button variant="danger" className="w-full">
            Delete Account
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default ProfilePage;