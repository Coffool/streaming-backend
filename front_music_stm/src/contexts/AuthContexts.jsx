// src/contexts/AuthContext.jsx
import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { authAPI } from '@/services/auth.js';

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth debe ser usado dentro de un AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [currentUser, setCurrentUser] = useState(null);
  const [authTokens, setAuthTokens] = useState(null);
  const [loading, setLoading] = useState(true);

  // Función para guardar tokens en localStorage
  const saveTokens = useCallback((tokens, user) => {
    localStorage.setItem('authTokens', JSON.stringify(tokens));
    localStorage.setItem('userData', JSON.stringify(user));
    setAuthTokens(tokens);
    setCurrentUser(user);
  }, []);

  // Función para limpiar tokens
  const clearTokens = useCallback(() => {
    localStorage.removeItem('authTokens');
    localStorage.removeItem('userData');
    setAuthTokens(null);
    setCurrentUser(null);
  }, []);

  // Verificar autenticación al cargar la app
  useEffect(() => {
    const initAuth = async () => {
      try {
        const storedTokens = JSON.parse(localStorage.getItem('authTokens'));
        const storedUser = JSON.parse(localStorage.getItem('userData'));
        
        if (storedTokens && storedUser) {
          setAuthTokens(storedTokens);
          setCurrentUser(storedUser);
        }
      } catch (error) {
        console.error('Error initializing auth:', error);
        clearTokens();
      } finally {
        setLoading(false);
      }
    };

    initAuth();
  }, [clearTokens]);

  const register = async (userData) => {
    try {
      const response = await authAPI.register(userData);
      return response;
    } catch (error) {
      console.error('Register error:', error);
      throw error;
    }
  };

  const login = async (credentials) => {
    try {
      const response = await authAPI.login(credentials);
      
      const tokens = {
        accessToken: response.access_token,
        refreshToken: response.refresh_token,
        expiresIn: response.expires_in
      };
      
      saveTokens(tokens, response.user);
      return response;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  };

  const logout = async () => {
    try {
      if (authTokens?.accessToken) {
        await authAPI.logout(authTokens.accessToken);
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      clearTokens();
    }
  };

  const value = {
    currentUser,
    authTokens,
    loading,
    login,
    logout,
    register
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};