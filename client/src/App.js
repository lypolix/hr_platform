import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

// ВАЖНО: именованный импорт для AuthProvider
import { AuthProvider } from './context/AuthContext';

// ВАЖНО: default-импорты для компонентов
import Navbar from './components/Navbar';
import PrivateRoute from './components/PrivateRoute';

// Pages
import Landing from './pages/Landing';

// University
import UniversityLogin from './pages/university/UniversityLogin';
import UniversityRegister from './pages/university/UniversityRegister';
import UniversityDashboard from './pages/university/UniversityDashboard';
import CreateApplication from './pages/university/CreateApplication';

// Company
import CompanyLogin from './pages/company/CompanyLogin';
import CompanyRegister from './pages/company/CompanyRegister';
import CompanyDashboard from './pages/company/CompanyDashboard';
import CreateVacancy from './pages/company/CreateVacancy';
import ViewResponses from './pages/company/ViewResponses';

// Admin
import AdminLogin from './pages/admin/AdminLogin';
import AdminDashboard from './pages/admin/AdminDashboard';

import './App.css';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Navbar />
        <Routes>
          {/* Public */}
          <Route path="/" element={<Landing />} />

          {/* University */}
          <Route path="/university/login" element={<UniversityLogin />} />
          <Route path="/university/register" element={<UniversityRegister />} />
          <Route
            path="/university/dashboard"
            element={
              <PrivateRoute allowedTypes={['university']}>
                <UniversityDashboard />
              </PrivateRoute>
            }
          />
          <Route
            path="/university/create-application"
            element={
              <PrivateRoute allowedTypes={['university']}>
                <CreateApplication />
              </PrivateRoute>
            }
          />

          {/* Company */}
          <Route path="/company/login" element={<CompanyLogin />} />
          <Route path="/company/register" element={<CompanyRegister />} />
          <Route
            path="/company/dashboard"
            element={
              <PrivateRoute allowedTypes={['company']}>
                <CompanyDashboard />
              </PrivateRoute>
            }
          />
          <Route
            path="/company/create-vacancy"
            element={
              <PrivateRoute allowedTypes={['company']}>
                <CreateVacancy />
              </PrivateRoute>
            }
          />
          <Route
            path="/company/responses/:vacancyId"
            element={
              <PrivateRoute allowedTypes={['company']}>
                <ViewResponses />
              </PrivateRoute>
            }
          />

          {/* Admin */}
          <Route path="/admin/login" element={<AdminLogin />} />
          <Route
            path="/admin/dashboard"
            element={
              <PrivateRoute allowedTypes={['admin']}>
                <AdminDashboard />
              </PrivateRoute>
            }
          />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
