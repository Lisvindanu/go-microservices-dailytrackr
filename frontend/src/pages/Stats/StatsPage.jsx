import React, { useState, useEffect } from 'react';
import { statsService } from '../../services/statsService';
import { aiService } from '../../services/aiService';
import { useToast } from '../../contexts/ToastContext';
import Card from '../../components/UI/Card';
import Button from '../../components/UI/Button';

const StatsPage = () => {
  const [dashboardStats, setDashboardStats] = useState({
    avg_daily_hours: 0,
    streak_days: 0,
    hours_growth: 0,
    this_week_hours: 0,
    last_week_hours: 0
  });
  const [insights, setInsights] = useState('');
  const [habitRecommendation, setHabitRecommendation] = useState('');
  const [loading, setLoading] = useState(true);
  const [apiAvailable, setApiAvailable] = useState(false);
  const { showToast } = useToast();

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);

        // Try to fetch data with timeout
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), 5000);

        const [statsResponse, insightsResponse, recommendationResponse] = await Promise.all([
          statsService.getDashboard(),
          aiService.getInsights(),
          aiService.getHabitRecommendation()
        ]);

        clearTimeout(timeoutId);

        // If we get here, API is working
        setApiAvailable(true);
        setDashboardStats(statsResponse.data || dashboardStats);
        setInsights(insightsResponse.data?.ai_insights || 'No insights available yet.');
        setHabitRecommendation(recommendationResponse.data?.recommendation || 'No recommendations available yet.');

      } catch (error) {
        console.log('Stats API not available:', error.message);
        setApiAvailable(false);

        // Set demo data
        setDashboardStats({
          avg_daily_hours: 4.5,
          streak_days: 7,
          hours_growth: 15,
          this_week_hours: 32,
          last_week_hours: 28
        });
        setInsights('Demo Mode: Your productivity has been consistent this week. You\'ve maintained a good balance between different activities. Keep up the great work!');
        setHabitRecommendation('Demo Mode: Based on your activity patterns, consider adding a "Daily Reading" habit for 30 minutes before bed to enhance your learning routine.');

        // Only show toast if it's not a network error
        if (!error.message.includes('Network error')) {
          showToast('warning', 'Backend services not available. Showing demo data.');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

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
                      Showing demo statistics. Start backend services to see real data.
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
          Statistics & Insights {!apiAvailable ? '(Demo)' : ''}
        </h1>

        {/* Progress Overview */}
        <Card title="📊 Progress Overview">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600">
                {Math.round(dashboardStats?.avg_daily_hours || 0)}h
              </div>
              <div className="text-sm text-gray-600">Avg Daily Hours</div>
            </div>

            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">
                {dashboardStats?.streak_days || 0}
              </div>
              <div className="text-sm text-gray-600">Current Streak</div>
            </div>

            <div className="text-center">
              <div className="text-2xl font-bold text-purple-600">
                {Math.round(dashboardStats?.hours_growth || 0)}%
              </div>
              <div className="text-sm text-gray-600">Weekly Growth</div>
            </div>
          </div>
        </Card>

        {/* AI Insights */}
        <Card title="🧠 AI Insights">
          <div className="prose max-w-none">
            <p className="text-gray-700 leading-relaxed">{insights}</p>
          </div>
        </Card>

        {/* Habit Recommendations */}
        <Card title="💡 Smart Habit Recommendations">
          <div className="prose max-w-none">
            <p className="text-gray-700 leading-relaxed">{habitRecommendation}</p>
          </div>
        </Card>

        {/* Week Comparison */}
        <Card title="📈 This Week vs Last Week">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h4 className="font-semibold text-gray-700 mb-2">This Week</h4>
              <div className="text-2xl font-bold text-blue-600">
                {Math.round(dashboardStats?.this_week_hours || 0)}h
              </div>
            </div>

            <div>
              <h4 className="font-semibold text-gray-700 mb-2">Last Week</h4>
              <div className="text-2xl font-bold text-gray-600">
                {Math.round(dashboardStats?.last_week_hours || 0)}h
              </div>
            </div>
          </div>

          {dashboardStats?.hours_growth !== undefined && (
              <div className="mt-4 pt-4 border-t">
                <div className={`text-lg font-semibold ${dashboardStats.hours_growth >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                  {dashboardStats.hours_growth >= 0 ? '↗️' : '↘️'}
                  {Math.abs(Math.round(dashboardStats.hours_growth))}%
                  {dashboardStats.hours_growth >= 0 ? 'increase' : 'decrease'}
                </div>
              </div>
          )}
        </Card>
      </div>
  );
};

export default StatsPage;