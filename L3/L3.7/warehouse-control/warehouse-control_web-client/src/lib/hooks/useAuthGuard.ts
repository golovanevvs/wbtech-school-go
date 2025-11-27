"use client"

import { useRouter } from "next/navigation"
import { useAuth } from "../contexts/AuthContext"
import { UserRole } from "../types/auth"
import { useEffect } from "react"

export const useAuthGuard = (requiredRoles?: UserRole[]) => {
  const { user, isLoading, isAuthenticated, hasRole } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!isLoading) {
      if (!isAuthenticated) {
        router.push("/login")
        return
      }

      if (requiredRoles && !hasRole(requiredRoles)) {
        router.push("/login")
        return
      }
    }
  }, [isLoading, isAuthenticated, requiredRoles, hasRole, router])
  return { user, isLoading, isAuthenticated }
}
