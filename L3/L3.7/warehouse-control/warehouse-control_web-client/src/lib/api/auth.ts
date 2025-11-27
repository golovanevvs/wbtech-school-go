import { User } from "../types/auth";
import apiClient from "./client";

// Интерфейсы для API ответов
interface LoginResponse {
  user: User
}

interface GetUserResponse {
  user: User
}

interface RefreshTokenResponse {
  access_token: string
  refresh_token?: string
}

interface ValidateTokenResponse {
  valid: boolean
}

// API для работы с авторизацией
export const authAPI = {
  /**
   * Вход в систему
   * @param username Имя пользователя
   * @param password Пароль
   * @returns Данные пользователя
   */
  async login(username: string, password: string): Promise<LoginResponse> {
    try {
      const response = await apiClient.post<LoginResponse>('/auth/login', {
        username,
        password,
      })
      
      // Сервер уже установил cookie с токенами
      return response
    } catch (error) {
      console.error('Login failed:', error)
      throw new Error('Не удалось выполнить вход. Проверьте логин и пароль.')
    }
  },

  /**
   * Выход из системы
   */
  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout')
      // Cookie будут удалены сервером
    } catch (error) {
      console.error('Logout failed:', error)
      // Очищение локальных данных, даже если сервер не ответил
    }
  },

  /**
   * Получение данных текущего пользователя
   * @returns Данные пользователя
   */
  async getCurrentUser(): Promise<User> {
    try {
      const response = await apiClient.get<GetUserResponse>('/auth/me')
      return response.user
    } catch (error) {
      console.error('Get current user failed:', error)
      throw new Error('Не удалось получить данные пользователя')
    }
  },

  /**
   * Обновление access токена
   * @returns Новый access токен
   */
  async refreshToken(): Promise<RefreshTokenResponse> {
    try {
      const response = await apiClient.post<RefreshTokenResponse>('/auth/refresh')
      return response
    } catch (error) {
      console.error('Token refresh failed:', error)
      throw new Error('Не удалось обновить токен')
    }
  },

  /**
   * Проверка валидности токена
   * @returns true если токен валиден
   */
  async validateToken(): Promise<boolean> {
    try {
      const response = await apiClient.get<ValidateTokenResponse>('/auth/validate')
      return response.valid
    } catch (error) {
      console.error('Token validation failed:', error)
      return false
    }
  },

  /**
   * Проверка доступности сервера авторизации
   * @returns true если сервер доступен
   */
  async checkServerStatus(): Promise<boolean> {
    try {
      return await apiClient.ping()
    } catch (error) {
      console.error('Server status check failed:', error)
      return false
    }
  },
}

export type {
  LoginResponse,
  GetUserResponse,
  RefreshTokenResponse,
  ValidateTokenResponse,
}