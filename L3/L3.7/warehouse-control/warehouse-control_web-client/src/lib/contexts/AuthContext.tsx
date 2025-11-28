"use client"

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
  useCallback,
} from "react"
import { usePathname, useRouter } from "next/navigation"
import { User, UserRole } from "../types/auth"
import { authAPI } from "../api/auth"

// Интерфейс для контекста авторизации
interface AuthContextType {
  // Состояние
  user: User | null
  isLoading: boolean
  isChecking: boolean
  isAuthenticated: boolean
  error: string | null

  // Функции
  login: (username: string, password: string) => Promise<boolean>
  register: (
    username: string,
    password: string,
    name: string,
    role: string
  ) => Promise<boolean>
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
  const [isChecking, setIsChecking] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const pathname = usePathname()
  const router = useRouter()

  // Логирование изменений error (упрощенное)
  // console.log("AuthProvider render - error:", error)

  // Авторизован?
  const isAuthenticated = user !== null

  // Функция входа в систему
  const login = async (
    username: string,
    password: string
  ): Promise<boolean> => {
    try {
      setIsChecking(true)
      setIsLoading(true)
      setError(null)

      await authAPI.login(username, password)

      // Получение данных пользователя
      const currentUser = await authAPI.getCurrentUser()
      setUser(currentUser)

      // Проверяем, есть ли сохраненный путь для перенаправления
      const redirectPath = sessionStorage.getItem("redirectAfterLogin")
      if (redirectPath) {
        sessionStorage.removeItem("redirectAfterLogin")
        router.push(redirectPath)
      } else if (pathname === "/auth") {
        router.push("/")
      }

      return true
    } catch (error) {
      console.error("Login failed:", error)
      if (error instanceof Error) {
        setError(error.message)
      }
      return false
    } finally {
      setIsLoading(false)
      setIsChecking(false)
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
      setIsChecking(true)
      setIsLoading(true)
      setError(null)

      await authAPI.register(username, password, name, role)

      // Получение данных пользователя после регистрации
      const currentUser = await authAPI.getCurrentUser()
      setUser(currentUser)

      // Проверяем, есть ли сохраненный путь для перенаправления
      const redirectPath = sessionStorage.getItem("redirectAfterLogin")
      if (redirectPath) {
        sessionStorage.removeItem("redirectAfterLogin")
        router.push(redirectPath)
      } else if (pathname === "/auth") {
        router.push("/")
      }

      return true
    } catch (error) {
      console.error("Registration failed:", error)
      if (error instanceof Error) {
        setError(error.message)
      }
      return false
    } finally {
      setIsLoading(false)
      setIsChecking(false)
    }
  }

  // Функция выхода из системы
  const logout = () => {
    document.cookie = "access_token=; path=/; max-age=0"
    document.cookie = "refresh_token=; path=/; max-age=0"

    // Очищение состояния
    setUser(null)

    // Очищаем сохраненный путь для перенаправления
    sessionStorage.removeItem("redirectAfterLogin")

    authAPI.logout().catch(console.error)
  }

  // Функция проверки авторизации при загрузке
  const checkAuth = useCallback(async () => {
    console.log("checkAuth: starting...")
    try {
      setIsChecking(true)
      setIsLoading(true)
      setError(null)
      console.log("checkAuth: making request to get current user")
      const currentUser = await authAPI.getCurrentUser()
      console.log("checkAuth: success, user:", currentUser)
      console.log("checkAuth: user type:", typeof currentUser)
      console.log("checkAuth: user keys:", currentUser ? Object.keys(currentUser) : "null/undefined")
      setUser(currentUser)
    } catch (error) {
      console.error("Check auth failed:", error)
      setUser(null)
      if (error instanceof Error) {
        console.log("checkAuth: setting error:", error.message)
        setError(error.message)
      }
    } finally {
      setIsLoading(false)
      setIsChecking(false)
      console.log("checkAuth: finished")
    }
  }, [])

  // Функция проверки ролей
  const hasRole = (roles: UserRole[]): boolean => {
    if (!user) return false
    return roles.includes(user.user_role)
  }

  // Функция очистки ошибки
  const clearError = () => {
    console.log("clearError called")
    setError(null)
  }

  // Проверка авторизации при монтировании компонента
  useEffect(() => {
    console.log("AuthContext useEffect triggered, pathname:", pathname)
    // Не проверяем авторизацию на странице входа
    if (pathname !== "/auth") {
      console.log("AuthContext: starting checkAuth for pathname:", pathname)
      checkAuth()
    } else {
      console.log("AuthContext: skipping checkAuth for /auth page")
    }
  }, [pathname, checkAuth])

  // Значение контекста
  const value: AuthContextType = {
    user,
    isLoading,
    isChecking,
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