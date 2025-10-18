import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const PrivateRoute = ({ children, allowedTypes }) => {
  const { user, userType, loading } = useAuth();

  if (loading) {
    return <div className="loading">Загрузка...</div>;
  }

  if (!user || !allowedTypes.includes(userType)) {
    return <Navigate to="/" />;
  }

  return children;
};

export default PrivateRoute;
