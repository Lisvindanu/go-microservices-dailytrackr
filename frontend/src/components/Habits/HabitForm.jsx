// HabitForm component - Form untuk create/edit habit
// Props: habit, onSubmit, loading
import React, { useState, useEffect } from 'react';
import Button from '../UI/Button';
import Input from '../UI/Input';

const HabitForm = ({ habit = null, onSubmit, loading = false }) => {
  const [formData, setFormData] = useState({
    title: '',
    start_date: '',
    end_date: '',
    reminder_time: '09:00'
  });

  useEffect(() => {
    if (habit) {
      setFormData({
        title: habit.title || '',
        start_date: habit.start_date ? habit.start_date.split('T')[0] : '',
        end_date: habit.end_date ? habit.end_date.split('T')[0] : '',
        reminder_time: habit.reminder_time || '09:00'
      });
    }
  }, [habit]);

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    // Validate dates
    const startDate = new Date(formData.start_date);
    const endDate = new Date(formData.end_date);

    if (endDate <= startDate) {
      alert('End date must be after start date');
      return;
    }

    onSubmit(formData);
  };

  const getTodayDate = () => {
    return new Date().toISOString().split('T')[0];
  };

  const getDefaultEndDate = () => {
    const today = new Date();
    const endDate = new Date(today);
    endDate.setDate(today.getDate() + 30); // Default 30 days
    return endDate.toISOString().split('T')[0];
  };

  return (
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
            label="Habit Title"
            name="title"
            value={formData.title}
            onChange={handleChange}
            placeholder="e.g., Exercise 30 minutes daily"
            required
        />

        <Input
            label="Start Date"
            type="date"
            name="start_date"
            value={formData.start_date}
            onChange={handleChange}
            min={getTodayDate()}
            required
        />

        <Input
            label="End Date"
            type="date"
            name="end_date"
            value={formData.end_date}
            onChange={handleChange}
            min={formData.start_date || getTodayDate()}
            required
        />

        <Input
            label="Daily Reminder Time"
            type="time"
            name="reminder_time"
            value={formData.reminder_time}
            onChange={handleChange}
        />

        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex items-start space-x-2">
            <span className="text-blue-600">💡</span>
            <div className="text-sm text-blue-700">
              <p className="font-medium mb-1">Habit Tip:</p>
              <p>Start with small, achievable goals. Consistency is more important than intensity!</p>
            </div>
          </div>
        </div>

        <div className="flex justify-end space-x-3 pt-4">
          <Button
              type="submit"
              loading={loading}
              disabled={loading}
          >
            {habit ? 'Update Habit' : 'Create Habit'}
          </Button>
        </div>
      </form>
  );
};

export default HabitForm;