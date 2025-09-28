// src/services/api.js
import { authAPI } from './auth';

const API_BASE_URL = 'http://localhost:8080';

// Función para hacer requests autenticados
export const authenticatedRequest = async (endpoint, options = {}, authContext) => {
  const url = `${API_BASE_URL}${endpoint}`;
  
  let token;
  if (authContext) {
    token = await authContext.getAccessToken();
  }
  
  const defaultOptions = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
  };
  
  const response = await fetch(url, { ...defaultOptions, ...options });
  
  // Si recibimos un 401 (Unauthorized), intentamos refrescar el token
  if (response.status === 401 && authContext) {
    try {
      const newToken = await authContext.refreshAuthToken();
      
      // Reintentar la request con el nuevo token
      const retryOptions = {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${newToken}`,
        },
      };
      
      return fetch(url, retryOptions);
    } catch (refreshError) {
      // Si el refresh falla, hacemos logout
      await authContext.logout();
      throw new Error('Session expired. Please login again.');
    }
  }
  
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `Error ${response.status}: ${response.statusText}`);
  }
  
  return response.json();
};

// Ejemplo de uso para endpoints específicos
export const userAPI = {
  getProfile: (authContext) => authenticatedRequest('/user/profile', {}, authContext),
  updateProfile: (data, authContext) => authenticatedRequest('/user/profile', {
    method: 'PUT',
    body: JSON.stringify(data),
  }, authContext),
};