import { User, RegisterResponse } from "../types/auth"
import apiClient from "./client"

interface UpdateUserRequest {
  name?: string
  user_role?: string
}

// Интерфейсы для API ответов
interface LoginResponse {
  message: string
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
      const response = await apiClient.post<LoginResponse>("/auth/login", {
        username,
        password,
      })

      return response
    } catch (error) {
      console.error("Login failed:", error)
      throw new Error("Не удалось выполнить вход. Проверьте логин и пароль.")
    }
  },

  /**
   * Регистрация нового пользователя
   * @param username Имя пользователя
   * @param password Пароль
   * @param name Отображаемое имя
   * @param role Роль пользователя
   * @returns Данные пользователя
   */
  async register(
    username: string,
    password: string,
    name: string,
    role: string
  ): Promise<RegisterResponse> {
    try {
      const response = await apiClient.post<RegisterResponse>(
        "/auth/register",
        {
          username,
          password,
          name,
          role,
        }
      )

      return response
    } catch (error) {
      console.error("Registration failed:", error)
      throw new Error(
        "Не удалось зарегистрироваться. Проверьте введенные данные."
      )
    }
  },

  /**
   * Выход из системы
   */
  async logout(): Promise<void> {
    try {
      await apiClient.post("/auth/logout")
      // Cookie будут удалены сервером
    } catch (error) {
      console.error("Logout failed:", error)
    }
  },

  /**
   * Получение данных текущего пользователя
   * @returns Данные пользователя
   */
  async getCurrentUser(): Promise<User> {
    try {
      const response = await apiClient.get<User>("/auth/me")
      return response
    } catch (error) {
      console.error("Get current user failed:", error)
      if (error instanceof Error) {
        throw error
      }
      throw new Error("Не удалось получить данные пользователя")
    }
  },

  /**
   * Проверка валидности токена
   * @returns true если токен валиден
   */
  async validateToken(): Promise<boolean> {
    try {
      const response = await apiClient.get<ValidateTokenResponse>(
        "/auth/validate"
      )
      return response.valid
    } catch (error) {
      console.error("Token validation failed:", error)
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
      console.error("Server status check failed:", error)
      return false
    }
  },

  /**
   * Обновление данных пользователя
   * @param userData Данные для обновления
   * @returns Обновленные данные пользователя
   */
  async updateUser(userData: UpdateUserRequest): Promise<User> {
    try {
      const response = await apiClient.put<User>("/auth/update", userData)
      return response
    } catch (error) {
      console.error("Update user failed:", error)
      throw new Error("Не удалось обновить данные пользователя")
    }
  },

  /**
   * Удаление пользователя
   */
  async deleteUser(): Promise<void> {
    try {
      await apiClient.delete("/auth/delete")
    } catch (error) {
      console.error("Delete user failed:", error)
      throw new Error("Не удалось удалить пользователя")
    }
  },
}

export type {
  LoginResponse,
  GetUserResponse,
  RefreshTokenResponse,
  ValidateTokenResponse,
}
