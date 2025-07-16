import React, { useState, useEffect } from 'react';
import { habitsService } from '../../services/habitsService';
import { useToast } from '../../contexts/ToastContext';
import Button from '../../components/UI/Button';
import Card from '../../components/UI/Card';
import Modal from '../../components/UI/Modal';
import Input from '../../components/UI/Input';

const HabitsPage = () => {
  const [habits, setHabits] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showModal, setShowModal] = useState(false);
  const [showLogModal, setShowLogModal] = useState(false);
  const [selectedHabit, setSelectedHabit] = useState(null);
  const [formData, setFormData] = useState({
    title: '',
    start_date: '',
    end_date: '',
    reminder_time: '09:00'
  });
  const [logData, setLogData] = useState({
    date: new Date().toISOString().split('T')[0],
    status: 'DONE',
    note: ''
  });
  const { showToast } = useToast();

  useEffect(() => {
    fetchHabits();
  }, []);

  const fetchHabits = async () => {
    try {
      const response = await habitsService.getAll(true);
      setHabits(response.data || []);
    } catch (error) {
      showToast('error', 'Failed to load habits');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await habitsService.create(formData);
      showToast('success', 'Habit created successfully!');
      setShowModal(false);
      setFormData({ title: '', start_date: '', end_date: '', reminder_time: '09:00' });
      fetchHabits();
    } catch (error) {
      showToast('error', error.message || 'Failed to create habit');
    }
  };

  const handleLogSubmit = async (e) => {
    e.preventDefault();
    try {
      await habitsService.createLog(selectedHabit.id, logData);
      showToast('success', 'Habit logged successfully!');
      setShowLogModal(false);
      setSelectedHabit(null);
      fetchHabits();
    } catch (error) {
      showToast('error', error.message || 'Failed to log habit');
   }
 };

 const handleChange = (e) => {
   setFormData(prev => ({
     ...prev,
     [e.target.name]: e.target.value
   }));
 };

 const handleLogChange = (e) => {
   setLogData(prev => ({
     ...prev,
     [e.target.name]: e.target.value
   }));
 };

 const openLogModal = (habit) => {
   setSelectedHabit(habit);
   setShowLogModal(true);
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
       <h1 className="text-2xl font-bold text-gray-900">Habits</h1>
       <Button onClick={() => setShowModal(true)}>
         + Add Habit
       </Button>
     </div>

     <div className="grid gap-4">
       {habits.length === 0 ? (
         <Card>
           <div className="text-center py-8 text-gray-500">
             No active habits. Create your first habit!
           </div>
         </Card>
       ) : (
         habits.map((habit) => (
           <Card key={habit.id}>
             <div className="flex justify-between items-start">
               <div>
                 <h3 className="font-semibold text-lg">{habit.title}</h3>
                 <p className="text-gray-600 mt-1">
                   {new Date(habit.start_date).toLocaleDateString()} - {new Date(habit.end_date).toLocaleDateString()}
                 </p>
                 {habit.reminder_time && (
                   <p className="text-gray-500 text-sm mt-1">
                     Reminder: {habit.reminder_time}
                   </p>
                 )}
               </div>
               <div className="flex space-x-2">
                 <Button
                   size="sm"
                   variant="success"
                   onClick={() => openLogModal(habit)}
                 >
                   Log Today
                 </Button>
               </div>
             </div>
           </Card>
         ))
       )}
     </div>

     {/* Create Habit Modal */}
     <Modal
       isOpen={showModal}
       onClose={() => setShowModal(false)}
       title="Create New Habit"
     >
       <form onSubmit={handleSubmit}>
         <Input
           label="Habit Title"
           name="title"
           value={formData.title}
           onChange={handleChange}
           required
         />
         
         <Input
           label="Start Date"
           type="date"
           name="start_date"
           value={formData.start_date}
           onChange={handleChange}
           required
         />
         
         <Input
           label="End Date"
           type="date"
           name="end_date"
           value={formData.end_date}
           onChange={handleChange}
           required
         />
         
         <Input
           label="Reminder Time"
           type="time"
           name="reminder_time"
           value={formData.reminder_time}
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
             Create Habit
           </Button>
         </div>
       </form>
     </Modal>

     {/* Log Habit Modal */}
     <Modal
       isOpen={showLogModal}
       onClose={() => setShowLogModal(false)}
       title={`Log: ${selectedHabit?.title}`}
     >
       <form onSubmit={handleLogSubmit}>
         <Input
           label="Date"
           type="date"
           name="date"
           value={logData.date}
           onChange={handleLogChange}
           required
         />
         
         <div className="mb-4">
           <label className="block text-sm font-medium text-gray-700 mb-1">
             Status
           </label>
           <select
             name="status"
             value={logData.status}
             onChange={handleLogChange}
             className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
           >
             <option value="DONE">✅ Done</option>
             <option value="SKIPPED">⏭️ Skipped</option>
             <option value="FAILED">❌ Failed</option>
           </select>
         </div>
         
         <Input
           label="Note (optional)"
           name="note"
           value={logData.note}
           onChange={handleLogChange}
         />
         
         <div className="flex justify-end space-x-3 mt-6">
           <Button 
             type="button" 
             variant="secondary" 
             onClick={() => setShowLogModal(false)}
           >
             Cancel
           </Button>
           <Button type="submit">
             Log Habit
           </Button>
         </div>
       </form>
     </Modal>
   </div>
 );
};

export default HabitsPage;