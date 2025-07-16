import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { statsService } from '../../services/statsService';
import { aiService } from '../../services/aiService';
import { useToast } from '../../contexts/ToastContext';
import Card from '../../components/UI/Card';
import Button from '../../components/UI/Button';

const DashboardPage = () => {
  const { user } = useAuth();
  const [dashboardData, setDashboardData] = useState(null);
  const [dailySummary, setDailySummary] = useState('');
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const [statsResponse, summaryResponse] = await Promise.all([
          statsService.getDashboard(),
          aiService.getDailySummary()
        ]);

        setDashboardData(statsResponse.data);
        setDailySummary(summaryResponse.data?.summary_text || 'No summary available yet. Start adding activities to get AI insights!');
      } catch (error) {
        showToast('error', 'Failed to load dashboard data');
        // Set default data if API fails
        setDashboardData({
          total_activities: 0,
          total_hours: 0,
          total_expenses: 0,
          active_habits: 0
        });
        setDailySummary('Welcome to DailyTrackr! Start by adding your first activity or habit to get personalized AI insights.');
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, [showToast]);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">
          Welcome back, {user?.username}! 👋
        </h1>
        <p className="text-gray-600">Here's your productivity overview</p>
      </div>
      
      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <Card>
          <div className="text-center">
            <div className="text-3xl font-bold text-blue-600">
              {dashboardData?.total_activities || 0}
            </div>
            <div className="text-sm text-gray-600">Total Activities</div>
          </div>
        </Card>
        
        <Card>
          <div className="text-center">
            <div className="text-3xl font-bold text-green-600">
              {Math.round(dashboardData?.total_hours || 0)}h
            </div>
            <div className="text-sm text-gray-600">Total Hours</div>
          </div>
        </Card>
        
        <Card>
          <div className="text-center">
            <div className="text-3xl font-bold text-purple-600">
              {dashboardData?.active_habits || 0}
            </div>
            <div className="text-sm text-gray-600">Active Habits</div>
          </div>
        </Card>
        
        <Card>
          <div className="text-center">
            <div className="text-3xl font-bold text-orange-600">
              Rp{(dashboardData?.total_expenses || 0).toLocaleString()}
            </div>
            <div className="text-sm text-gray-600">Total Expenses</div>
          </div>
        </Card>
      </div>

      {/* AI Summary */}
      <Card title="🤖 AI Daily Summary">
        <div className="prose max-w-none">
          <p className="text-gray-700 leading-relaxed">{dailySummary}</p>
        </div>
      </Card>

      {/* Quick Actions */}
      <Card title="Quick Actions">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Button className="w-full">
            📝 Add Activity
          </Button>
          <Button variant="secondary" className="w-full">
            🎯 Log Habit
          </Button>
          <Button variant="secondary" className="w-full">
            📊 View Stats
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default DashboardPage;