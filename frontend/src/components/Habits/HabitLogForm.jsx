// HabitLogForm component - Form untuk log daily habit (DONE/SKIPPED/FAILED)
// Props: habit, onSubmit, loading, defaultDate
import React, { useState, useEffect } from 'react';
import Button from '../UI/Button';
import Input from '../UI/Input';

const HabitLogForm = ({ habit, onSubmit, loading = false, defaultDate = null }) => {
  const [formData, setFormData] = useState({
    date: defaultDate || new Date().toISOString().split('T')[0],
    status: 'DONE',
    note: ''
  });

  useEffect(() => {
    if (defaultDate) {
      setFormData(prev => ({
        ...prev,
        date: defaultDate
      }));
    }
  }, [defaultDate]);

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

  const statusOptions = [
    { value: 'DONE', label: '✅ Done', color: 'text-green-600', description: 'Successfully completed' },
    { value: 'SKIPPED', label: '⏭️ Skipped', color: 'text-yellow-600', description: 'Intentionally skipped' },
    { value: 'FAILED', label: '❌ Failed', color: 'text-red-600', description: 'Unable to complete' }
  ];

  const getMinDate = () => {
    if (!habit?.start_date) return new Date().toISOString().split('T')[0];
    const startDate = new Date(habit.start_date).toISOString().split('T')[0];
    return startDate;
  };

  const getMaxDate = () => {
    const today = new Date().toISOString().split('T')[0];
    if (!habit?.end_date) return today;
    const endDate = new Date(habit.end_date).toISOString().split('T')[0];
    return endDate < today ? endDate : today;
  };

  return (
      <form onSubmit={handleSubmit} className="space-y-4">
        {habit && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
              <h3 className="font-medium text-blue-900 mb-1">
                🎯 {habit.title}
              </h3>
              <p className="text-sm text-blue-700">
                Log your progress for this habit
              </p>
            </div>
        )}

        <Input
            label="Date"
            type="date"
            name="date"
            value={formData.date}
            onChange={handleChange}
            min={getMinDate()}
            max={getMaxDate()}
            required
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-3">
            Status
          </label>
          <div className="space-y-2">
            {statusOptions.map((option) => (
                <label key={option.value} className="flex items-start space-x-3 cursor-pointer">
                  <input
                      type="radio"
                      name="status"
                      value={option.value}
                      checked={formData.status === option.value}
                      onChange={handleChange}
                      className="mt-1"
                  />
                  <div className="flex-1">
                    <div className={`font-medium ${option.color}`}>
                      {option.label}
                    </div>
                    <div className="text-sm text-gray-500">
                      {option.description}
                    </div>
                  </div>
                </label>
            ))}
          </div>
        </div>

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
              placeholder="Add any notes about your progress..."
          />
        </div>

        <div className="flex justify-end space-x-3 pt-4">
          <Button
              type="submit"
              loading={loading}
              disabled={loading}
          >
            Log Habit
          </Button>
        </div>
      </form>
  );
};

export default HabitLogForm;