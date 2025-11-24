"use client"

import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from "react"
import { getCurrentUser, logout as apiLogout } from "../api/auth"
import { User } from "../lib/types"

type AuthContextType = {
  user: User | null
  loading: boolean
  login: () => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const logout = async () => {
    try {
      await apiLogout()
      console.log("Logout successful")
    } catch (error) {
      console.error("Logout failed:", error)
    } finally {
      setUser(null)
    }
  }

  const login = async (): Promise<void> => {
    try {
      const userData = await getCurrentUser()
      setUser(userData)
    } catch (error) {
      throw error
    }
  }

  useEffect(() => {
    const checkAuthStatus = async () => {
      try {
        console.log("Checking auth status...")
        console.log("Current cookies:", document.cookie)

        const userData = await getCurrentUser()
        console.log("Auth check successful:", userData)
        setUser(userData)
      } catch (error) {
        if (error instanceof Error && error.message.includes("401")) {
          console.log("User is not authenticated (expected behavior)")
        } else {
          console.error("Auth check failed:", error)
        }
        setUser(null)
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
