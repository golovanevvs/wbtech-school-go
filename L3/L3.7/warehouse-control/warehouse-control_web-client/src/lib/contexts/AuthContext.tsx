"use client"

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react"
import { User, UserRole } from "../types/auth"
import { authAPI } from "../api/auth"

// Интерфейс для контекста авторизации
interface AuthContextType {
  // Состояние
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  error: string | null

  // Функции
  login: (username: string, password: string) => Promise<boolean>
  register: (username: string, password: string, name: string, role: string) => Promise<boolean>
  logout: () => void
  checkAuth: () => Promise<void>
  hasRole: (roles: UserRole[]) => boolean
  clearError: () => void
}

// Контекст
const AuthContext = createContext<AuthContextType | undefined>(undefined)

// Хук для использования контекста
export const useAuth = () => {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}

// Провайдер для контекста
interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  // Состояние
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Авторизован?
  const isAuthenticated = user !== null

  // Функция проверки авторизации при загрузке
  const checkAuth = async () => {
    try {
      setIsLoading(true)
      setError(null)
      const currentUser = await authAPI.getCurrentUser()
      setUser(currentUser)
    } catch (error) {
      console.error("Check auth failed:", error)
      setUser(null)
      if (error instanceof Error) {
        setError(error.message)
      }
    } finally {
      setIsLoading(false)
    }
  }

  // Функция входа в систему
  const login = async (
    username: string,
    password: string
  ): Promise<boolean> => {
    try {
      setIsLoading(true)
      setError(null)
      
      await authAPI.login(username, password)

      // Получение данных пользователя
      const currentUser = await authAPI.getCurrentUser()
      setUser(currentUser)

      return true
    } catch (error) {
      console.error("Login failed:", error)
      if (error instanceof Error) {
        setError(error.message)
      }
      return false
    } finally {
      setIsLoading(false)
    }
  }

  // Функция регистрации нового пользователя
  const register = async (
    username: string,
    password: string,
    name: string,
    role: string
  ): Promise<boolean> => {
    try {
      setIsLoading(true)
      setError(null)
      
      await authAPI.register(username, password, name, role)

      // Получение данных пользователя после регистрации
      const currentUser = await authAPI.getCurrentUser()
      setUser(currentUser)

      return true
    } catch (error) {
      console.error("Registration failed:", error)
      if (error instanceof Error) {
        setError(error.message)
      }
      return false
    } finally {
      setIsLoading(false)
    }
  }

  // Функция выхода из системы
  const logout = () => {
    document.cookie = "access_token=; path=/; max-age=0"
    document.cookie = "refresh_token=; path=/; max-age=0"

    // Очищение состояния
    setUser(null)

    authAPI.logout().catch(console.error)
  }

  // Функция проверки ролей
  const hasRole = (roles: UserRole[]): boolean => {
    if (!user) return false
    return roles.includes(user.user_role)
  }

  // Функция очистки ошибки
  const clearError = () => {
    setError(null)
  }

  // Проверка авторизации при монтировании компонента
  useEffect(() => {
    checkAuth()
  }, [])

  // Значение контекста
  const value: AuthContextType = {
    user,
    isLoading,
    isAuthenticated,
    error,
    login,
    register,
    logout,
    checkAuth,
    hasRole,
    clearError,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}