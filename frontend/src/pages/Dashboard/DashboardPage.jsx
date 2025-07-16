import React, { useState, useEffect } from 'react';
import {
  Activity,
  Target,
  TrendingUp,
  Clock,
  DollarSign,
  Zap,
  Plus,
  Calendar,
  CheckCircle,
  BarChart3,
  Sparkles,
  RefreshCw,
  Wifi,
  WifiOff
} from 'lucide-react';

const DashboardPage = () => {
  const [stats, setStats] = useState({
    totalActivities: 0,
    totalHours: 0,
    activeHabits: 0,
    totalExpenses: 0
  });
  const [apiAvailable, setApiAvailable] = useState(false);
  const [loading, setLoading] = useState(true);
  const [recentActivities, setRecentActivities] = useState([]);
  const [recentHabits, setRecentHabits] = useState([]);

  useEffect(() => {
    checkBackendStatus();
    loadDashboardData();
  }, []);

  const checkBackendStatus = async () => {
    try {
      const response = await fetch('http://localhost:3000/', {
        method: 'GET',
        signal: AbortSignal.timeout(3000)
      });
      setApiAvailable(response.ok);
    } catch (error) {
      setApiAvailable(false);
    }
  };

  const loadDashboardData = async () => {
    try {
      // Simulate loading demo data
      setStats({
        totalActivities: 12,
        totalHours: 45.5,
        activeHabits: 3,
        totalExpenses: 125000
      });

      setRecentActivities([
        { id: 1, title: 'Belajar React.js', duration: 120, cost: 0, date: '2025-07-16' },
        { id: 2, title: 'Workout Session', duration: 60, cost: 15000, date: '2025-07-15' },
        { id: 3, title: 'Reading Time', duration: 45, cost: 0, date: '2025-07-15' }
      ]);

      setRecentHabits([
        { id: 1, title: 'Morning Exercise', streak: 7, status: 'active' },
        { id: 2, title: 'Daily Reading', streak: 12, status: 'active' },
        { id: 3, title: 'Meditation', streak: 3, status: 'active' }
      ]);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const StatCard = ({ icon: Icon, title, value, subtitle, color = "blue", trend }) => (
      <div className={`bg-gradient-to-br from-${color}-50 to-${color}-100 border border-${color}-200 rounded-xl p-6 hover:shadow-lg transition-all duration-300 hover:scale-105`}>
        <div className="flex items-center justify-between">
          <div>
            <div className="flex items-center space-x-2 mb-2">
              <Icon className={`h-5 w-5 text-${color}-600`} />
              <p className={`text-${color}-700 text-sm font-medium`}>{title}</p>
            </div>
            <p className="text-2xl font-bold text-gray-900">{value}</p>
            {subtitle && <p className="text-sm text-gray-600 mt-1">{subtitle}</p>}
          </div>
          {trend && (
              <div className={`flex items-center space-x-1 text-${color}-600`}>
                <TrendingUp className="h-4 w-4" />
                <span className="text-sm font-medium">{trend}</span>
              </div>
          )}
        </div>
      </div>
  );

  const ActivityCard = ({ activity }) => (
      <div className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <h4 className="font-medium text-gray-900">{activity.title}</h4>
            <div className="flex items-center space-x-4 mt-2 text-sm text-gray-500">
              <div className="flex items-center space-x-1">
                <Clock className="h-4 w-4" />
                <span>{activity.duration}m</span>
              </div>
              {activity.cost > 0 && (
                  <div className="flex items-center space-x-1">
                    <DollarSign className="h-4 w-4" />
                    <span>Rp{activity.cost.toLocaleString()}</span>
                  </div>
              )}
              <div className="flex items-center space-x-1">
                <Calendar className="h-4 w-4" />
                <span>{activity.date}</span>
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-2">
          <span className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">
            Completed
          </span>
          </div>
        </div>
      </div>
  );

  const HabitCard = ({ habit }) => (
      <div className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-purple-100 rounded-lg">
              <Target className="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <h4 className="font-medium text-gray-900">{habit.title}</h4>
              <p className="text-sm text-gray-500">{habit.streak} day streak</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <CheckCircle className="h-5 w-5 text-green-500" />
            <span className="px-2 py-1 bg-green-100 text-green-800 text-xs rounded-full">
            Active
          </span>
          </div>
        </div>
      </div>
  );

  const QuickActionButton = ({ icon: Icon, label, color, onClick }) => (
      <button
          onClick={onClick}
          className={`flex items-center space-x-2 px-4 py-3 bg-gradient-to-r from-${color}-500 to-${color}-600 text-white rounded-lg hover:from-${color}-600 hover:to-${color}-700 transition-all duration-200 hover:scale-105 shadow-lg`}
      >
        <Icon className="h-5 w-5" />
        <span className="font-medium">{label}</span>
      </button>
  );

  if (loading) {
    return (
        <div className="flex items-center justify-center h-64">
          <div className="flex items-center space-x-2">
            <RefreshCw className="h-6 w-6 animate-spin text-blue-600" />
            <span className="text-gray-600">Loading dashboard...</span>
          </div>
        </div>
    );
  }

  return (
      <div className="space-y-8 max-w-7xl mx-auto">
        {/* Header */}
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl p-8 text-white">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold mb-2">Welcome back, Anaphygon! 👋</h1>
              <p className="text-blue-100 text-lg">
                {apiAvailable ?
                    "Ready to track your progress today?" :
                    "Running in demo mode - Start your backend services to sync data"
                }
              </p>
            </div>
            <div className="flex items-center space-x-4">
              <div className={`flex items-center space-x-2 px-3 py-2 rounded-lg ${
                  apiAvailable ? 'bg-green-500/20 text-green-100' : 'bg-red-500/20 text-red-100'
              }`}>
                {apiAvailable ? <Wifi className="h-4 w-4" /> : <WifiOff className="h-4 w-4" />}
                <span className="text-sm font-medium">
                {apiAvailable ? 'Connected' : 'Offline'}
              </span>
              </div>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatCard
              icon={Activity}
              title="Total Activities"
              value={stats.totalActivities}
              color="blue"
              trend="+12%"
          />
          <StatCard
              icon={Clock}
              title="Total Hours"
              value={`${stats.totalHours}h`}
              color="green"
              trend="+8%"
          />
          <StatCard
              icon={Target}
              title="Active Habits"
              value={stats.activeHabits}
              color="purple"
              trend="+2"
          />
          <StatCard
              icon={DollarSign}
              title="Total Expenses"
              value={`Rp${stats.totalExpenses.toLocaleString()}`}
              color="orange"
              trend="-5%"
          />
        </div>

        {/* Quick Actions */}
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-xl font-bold text-gray-900 mb-4 flex items-center space-x-2">
            <Zap className="h-6 w-6 text-yellow-500" />
            <span>Quick Actions</span>
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <QuickActionButton
                icon={Plus}
                label="Add Activity"
                color="blue"
                onClick={() => console.log('Add Activity')}
            />
            <QuickActionButton
                icon={Target}
                label="Log Habit"
                color="purple"
                onClick={() => console.log('Log Habit')}
            />
            <QuickActionButton
                icon={BarChart3}
                label="View Stats"
                color="green"
                onClick={() => console.log('View Stats')}
            />
            <QuickActionButton
                icon={Sparkles}
                label="AI Insights"
                color="pink"
                onClick={() => console.log('AI Insights')}
            />
          </div>
        </div>

        {/* Recent Activities & Habits */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Recent Activities */}
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900 flex items-center space-x-2">
                <Activity className="h-6 w-6 text-blue-500" />
                <span>Recent Activities</span>
              </h2>
              <button className="text-blue-600 hover:text-blue-700 text-sm font-medium">
                View All
              </button>
            </div>
            <div className="space-y-4">
              {recentActivities.map(activity => (
                  <ActivityCard key={activity.id} activity={activity} />
              ))}
            </div>
          </div>

          {/* Active Habits */}
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-bold text-gray-900 flex items-center space-x-2">
                <Target className="h-6 w-6 text-purple-500" />
                <span>Active Habits</span>
              </h2>
              <button className="text-purple-600 hover:text-purple-700 text-sm font-medium">
                View All
              </button>
            </div>
            <div className="space-y-4">
              {recentHabits.map(habit => (
                  <HabitCard key={habit.id} habit={habit} />
              ))}
            </div>
          </div>
        </div>

        {/* AI Summary Section */}
        <div className="bg-gradient-to-r from-purple-50 to-pink-50 border border-purple-200 rounded-xl p-6">
          <div className="flex items-center space-x-2 mb-4">
            <Sparkles className="h-6 w-6 text-purple-600" />
            <h2 className="text-xl font-bold text-gray-900">AI Daily Summary</h2>
            <span className="px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-full">Demo</span>
          </div>
          <div className="bg-white rounded-lg p-4 border border-purple-100">
            <p className="text-gray-700 leading-relaxed">
              Welcome to DailyTrackr! Start by adding your first activity or habit to get personalized AI insights.
              Track your progress, build consistent habits, and let our AI help you optimize your productivity journey.
            </p>
          </div>
        </div>
      </div>
  );
};

export default DashboardPage;