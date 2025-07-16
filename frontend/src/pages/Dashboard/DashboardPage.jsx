import React, { useState, useEffect } from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { statsService } from '../../services/statsService';
import { aiService } from '../../services/aiService';
import { useToast } from '../../contexts/ToastContext';
import Card from '../../components/UI/Card';
import Button from '../../components/UI/Button';

const DashboardPage = () => {
  const { user } = useAuth();
  const [dashboardData, setDashboardData] = useState({
    total_activities: 0,
    total_hours: 0,
    total_expenses: 0,
    active_habits: 0
  });
  const [dailySummary, setDailySummary] = useState('');
  const [loading, setLoading] = useState(true);
  const [apiAvailable, setApiAvailable] = useState(false);
  const { showToast } = useToast();

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);

        // Try to fetch data with timeout
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 second timeout

        const statsPromise = statsService.getDashboard();
        const summaryPromise = aiService.getDailySummary();

        const [statsResponse, summaryResponse] = await Promise.all([
          statsPromise,
          summaryPromise
        ]);

        clearTimeout(timeoutId);

        // If we get here, API is working
        setApiAvailable(true);
        setDashboardData(statsResponse.data || dashboardData);
        setDailySummary(summaryResponse.data?.summary_text || 'No summary available yet. Start adding activities to get AI insights!');

      } catch (error) {
        console.log('Dashboard API not available:', error.message);
        setApiAvailable(false);

        // Set default data and summary for offline mode
        setDashboardData({
          total_activities: 0,
          total_hours: 0,
          total_expenses: 0,
          active_habits: 0
        });
        setDailySummary('Welcome to DailyTrackr! Start by adding your first activity or habit to get personalized AI insights.');

        // Only show toast if it's not a network error
        if (!error.message.includes('Network error')) {
          showToast('warning', 'Backend services not available. Using offline mode.');
        }
      } finally {
        setLoading(false);
      }
    };

    // Only fetch data once
    fetchDashboardData();
  }, []); // Remove showToast dependency to prevent re-fetching

  const handleRetryConnection = async () => {
    setLoading(true);

    try {
      const response = await fetch('http://localhost:3000/', {
        method: 'GET',
        signal: AbortSignal.timeout(3000) // 3 second timeout
      });

      if (response.ok) {
        setApiAvailable(true);
        showToast('success', 'Backend connection restored!');
        // Refresh the page to load data
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
                      Make sure Gateway is running on port 3000. You can still use the app in demo mode.
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

        <div>
          <h1 className="text-2xl font-bold text-gray-900">
            Welcome back, {user?.username || 'testuser'}! 👋
          </h1>
          <p className="text-gray-600">
            {apiAvailable ? "Here's your productivity overview" : "Demo mode - Add some activities to see your progress"}
          </p>
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
        <Card title={`🤖 AI Daily Summary ${!apiAvailable ? '(Demo)' : ''}`}>
          <div className="prose max-w-none">
            <p className="text-gray-700 leading-relaxed">{dailySummary}</p>
          </div>
        </Card>

        {/* Quick Actions */}
        <Card title="Quick Actions">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Button className="w-full" onClick={() => window.location.href = '/activities'}>
              📝 Add Activity
            </Button>
            <Button variant="secondary" className="w-full" onClick={() => window.location.href = '/habits'}>
              🎯 Log Habit
            </Button>
            <Button variant="secondary" className="w-full" onClick={() => window.location.href = '/stats'}>
              📊 View Stats
            </Button>
          </div>
        </Card>
      </div>
  );
};

export default DashboardPage;