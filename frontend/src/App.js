import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ToastProvider } from './contexts/ToastContext';
import LoginPage from './pages/Auth/LoginPage';
import DashboardPage from './pages/Dashboard/DashboardPage';
import ActivitiesPage from './pages/Activities/ActivitiesPage';
import HabitsPage from './pages/Habits/HabitsPage';
import StatsPage from './pages/Stats/StatsPage';
import ProfilePage from './pages/Profile/ProfilePage';
import Layout from './components/Layout/Layout';
import ProtectedRoute from './components/Auth/ProtectedRoute';

function App() {
  return (
      <Router>
        <AuthProvider>
          <ToastProvider>
            <Routes>
              {/* Public Routes */}
              <Route path="/login" element={<LoginPage />} />

              {/* Protected Routes - NO nested Routes */}
              <Route path="/" element={
                <ProtectedRoute>
                  <Layout>
                    <DashboardPage />
                  </Layout>
                </ProtectedRoute>
              } />

              <Route path="/dashboard" element={
                <ProtectedRoute>
                  <Layout>
                    <DashboardPage />
                  </Layout>
                </ProtectedRoute>
              } />

              <Route path="/activities" element={
                <ProtectedRoute>
                  <Layout>
                    <ActivitiesPage />
                  </Layout>
                </ProtectedRoute>
              } />

              <Route path="/habits" element={
                <ProtectedRoute>
                  <Layout>
                    <HabitsPage />
                  </Layout>
                </ProtectedRoute>
              } />

              <Route path="/stats" element={
                <ProtectedRoute>
                  <Layout>
                    <StatsPage />
                  </Layout>
                </ProtectedRoute>
              } />

              <Route path="/profile" element={
                <ProtectedRoute>
                  <Layout>
                    <ProfilePage />
                  </Layout>
                </ProtectedRoute>
              } />

              {/* Catch all - redirect to home */}
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </ToastProvider>
        </AuthProvider>
      </Router>
  );
}

export default App;