import { User } from '../lib/types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Интерфейс для обработки ошибок API
class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message);
    this.name = 'ApiError';
  }
}

// Функция для выполнения API запросов
const apiRequest = async <T>(endpoint: string, options: RequestInit = {}): Promise<T> => {
  const token = getToken();
  
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const errorData = await response.text();
    throw new ApiError(errorData, response.status);
  }

  return response.json();
};

// Функция для получения токена из localStorage
const getToken = (): string | null => {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
};

// Функция для получения текущего пользователя
export const getCurrentUser = async (): Promise<User> => {
  const token = getToken();
  
  if (!token) {
    throw new ApiError('No authentication token', 401);
  }

  return apiRequest<User>('/auth/me', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
};

// Функция для обновления пользователя
export const updateUser = async (userData: User): Promise<User> => {
  const token = getToken();
  
  if (!token) {
    throw new ApiError('No authentication token', 401);
  }

  return apiRequest<User>('/auth/update', {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
};

// Функция для выхода
export const logout = (): void => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('token');
  }
};