import React, { useState, useEffect } from 'react';
import { statsService } from '../../services/statsService';
import { aiService } from '../../services/aiService';
import { useToast } from '../../contexts/ToastContext';
import Card from '../../components/UI/Card';

const StatsPage = () => {
  const [dashboardStats, setDashboardStats] = useState(null);
  const [insights, setInsights] = useState('');
  const [habitRecommendation, setHabitRecommendation] = useState('');
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [statsResponse, insightsResponse, recommendationResponse] = await Promise.all([
          statsService.getDashboard(),
          aiService.getInsights(),
          aiService.getHabitRecommendation()
        ]);

        setDashboardStats(statsResponse.data);
        setInsights(insightsResponse.data?.ai_insights || 'No insights available yet.');
        setHabitRecommendation(recommendationResponse.data?.recommendation || 'No recommendations available yet.');
      } catch (error) {
        showToast('error', 'Failed to load statistics');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
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
      <h1 className="text-2xl font-bold text-gray-900">Statistics & Insights</h1>
      
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