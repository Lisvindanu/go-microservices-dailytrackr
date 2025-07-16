// ActivityForm component - Form untuk create/edit activity
// Props: activity, onSubmit, loading
import React, { useState, useEffect } from 'react';
import Button from '../UI/Button';
import Input from '../UI/Input';

const ActivityForm = ({ activity = null, onSubmit, loading = false }) => {
  const [formData, setFormData] = useState({
    title: '',
    start_time: '',
    duration_mins: '',
    cost: '',
    note: ''
  });

  useEffect(() => {
    if (activity) {
      // Format datetime untuk input datetime-local
      const startTime = activity.start_time
          ? new Date(activity.start_time).toISOString().slice(0, 16)
          : '';

      setFormData({
        title: activity.title || '',
        start_time: startTime,
        duration_mins: activity.duration_mins || '',
        cost: activity.cost || '',
        note: activity.note || ''
      });
    }
  }, [activity]);

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const submitData = {
      ...formData,
      duration_mins: parseInt(formData.duration_mins),
      cost: formData.cost ? parseInt(formData.cost) : null
    };

    onSubmit(submitData);
  };

  return (
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
            label="Activity Title"
            name="title"
            value={formData.title}
            onChange={handleChange}
            placeholder="e.g., Morning workout"
            required
        />

        <Input
            label="Start Time"
            type="datetime-local"
            name="start_time"
            value={formData.start_time}
            onChange={handleChange}
            required
        />

        <Input
            label="Duration (minutes)"
            type="number"
            name="duration_mins"
            value={formData.duration_mins}
            onChange={handleChange}
            placeholder="e.g., 60"
            min="1"
            required
        />

        <Input
            label="Cost (optional)"
            type="number"
            name="cost"
            value={formData.cost}
            onChange={handleChange}
            placeholder="e.g., 50000"
            min="0"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Note (optional)
          </label>
          <textarea
              name="note"
              value={formData.note}
              onChange={handleChange}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Add any additional notes..."
          />
        </div>

        <div className="flex justify-end space-x-3 pt-4">
          <Button
              type="submit"
              loading={loading}
              disabled={loading}
          >
            {activity ? 'Update Activity' : 'Create Activity'}
          </Button>
        </div>
      </form>
  );
};

export default ActivityForm;