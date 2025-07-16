import React, { useState, useEffect } from 'react';
import { activitiesService } from '../../services/activitiesService';
import { useToast } from '../../contexts/ToastContext';
import Button from '../../components/UI/Button';
import Card from '../../components/UI/Card';
import Modal from '../../components/UI/Modal';
import Input from '../../components/UI/Input';

const ActivitiesPage = () => {
  const [activities, setActivities] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    title: '',
    start_time: '',
    duration_mins: '',
    cost: '',
    note: ''
  });
  const { showToast } = useToast();

  useEffect(() => {
    fetchActivities();
  }, []);

  const fetchActivities = async () => {
    try {
      const response = await activitiesService.getAll();
      setActivities(response.data.activities || []);
    } catch (error) {
      showToast('error', 'Failed to load activities');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await activitiesService.create({
        ...formData,
        duration_mins: parseInt(formData.duration_mins),
        cost: formData.cost ? parseInt(formData.cost) : null
      });
      
      showToast('success', 'Activity created successfully!');
      setShowModal(false);
      setFormData({ title: '', start_time: '', duration_mins: '', cost: '', note: '' });
      fetchActivities();
    } catch (error) {
      showToast('error', error.message || 'Failed to create activity');
    }
  };

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
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
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-900">Activities</h1>
        <Button onClick={() => setShowModal(true)}>
          + Add Activity
        </Button>
      </div>

      <div className="grid gap-4">
        {activities.length === 0 ? (
          <Card>
            <div className="text-center py-8 text-gray-500">
              No activities yet. Create your first activity!
            </div>
          </Card>
        ) : (
          activities.map((activity) => (
            <Card key={activity.id}>
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{activity.title}</h3>
                  <p className="text-gray-600 mt-1">
                    {new Date(activity.start_time).toLocaleString()} • {activity.duration_mins} mins
                  </p>
                  {activity.note && (
                    <p className="text-gray-700 mt-2">{activity.note}</p>
                  )}
                </div>
                <div className="text-right">
                  {activity.cost && (
                    <div className="text-lg font-semibold text-green-600">
                      Rp{activity.cost.toLocaleString()}
                    </div>
                  )}
                </div>
              </div>
            </Card>
          ))
        )}
      </div>

      <Modal
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        title="Add New Activity"
      >
        <form onSubmit={handleSubmit}>
          <Input
            label="Title"
            name="title"
            value={formData.title}
            onChange={handleChange}
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
            required
          />
          
          <Input
            label="Cost (optional)"
            type="number"
            name="cost"
            value={formData.cost}
            onChange={handleChange}
          />
          
          <Input
            label="Note (optional)"
            name="note"
            value={formData.note}
            onChange={handleChange}
          />
          
          <div className="flex justify-end space-x-3 mt-6">
            <Button 
              type="button" 
              variant="secondary" 
              onClick={() => setShowModal(false)}
            >
              Cancel
            </Button>
            <Button type="submit">
              Create Activity
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default ActivitiesPage;