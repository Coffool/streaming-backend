// src/components/ProtectedRoute.jsx
import { useAuth } from '@/contexts/AuthContexts';
import { Navigate } from 'react-router-dom';

const ProtectedRoute = ({ children, requiredRole = null }) => {
  const { currentUser, loading } = useAuth();

  if (loading) {
    return <div>Cargando...</div>; // O un spinner
  }

  if (!currentUser) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && currentUser.role !== requiredRole) {
    return <Navigate to="/unauthorized" replace />;
  }

  return children;
};

export default ProtectedRoute;