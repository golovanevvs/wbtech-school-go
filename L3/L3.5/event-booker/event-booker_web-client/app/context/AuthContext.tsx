"use client"

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react"
import { getCurrentUser } from "../api/auth"
import { User } from "../lib/types"

type AuthContextType = {
  user: User | null
  loading: boolean
  login: (accessToken: string, refreshToken?: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const logout = () => {
    localStorage.removeItem("token")
    localStorage.removeItem("refreshToken")
    setUser(null)
  }

  const login = async (
    accessToken: string,
    refreshToken?: string
  ): Promise<void> => {
    localStorage.setItem("token", accessToken)
    if (refreshToken) {
      localStorage.setItem("refreshToken", refreshToken)
    }

    try {
      const userData = await getCurrentUser()
      setUser(userData)
    } catch (error) {
      localStorage.removeItem("token")
      localStorage.removeItem("refreshToken")
      setUser(null)
      throw error
    }
  }

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        // Пробуем получить текущего пользователя
        // Если токен просрочен, api/auth.ts автоматически попытается обновить его
        const userData = await getCurrentUser()
        setUser(userData)
      } catch (error) {
        // Если обновление токена не помогло (или нет токена вообще)
        console.error("Auth check failed:", error)
        logout()
      } finally {
        setLoading(false)
      }
    }

    checkAuthStatus()
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}