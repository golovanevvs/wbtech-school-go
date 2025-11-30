"use client"

import { useRouter } from "next/navigation"
import { useAuth } from "../contexts/AuthContext"
import { UserRole } from "../types/auth"
import { useEffect } from "react"
import { getFullPath } from "../utils/paths"

export const useAuthGuard = (requiredRoles?: UserRole[]) => {
  const { user, isLoading, isAuthenticated, isChecking, hasRole } = useAuth()
  const router = useRouter()

  useEffect(() => {
    console.log("useAuthGuard effect:", {
      isLoading,
      isChecking,
      isAuthenticated,
      requiredRoles,
    })

    if (!isLoading && !isChecking) {
      console.log("AuthGuard: checking permissions...")
      if (!isAuthenticated) {
        console.log("AuthGuard: not authenticated, redirecting to /auth")
        router.push(getFullPath("/auth"))
        return
      }

      if (requiredRoles && !hasRole(requiredRoles)) {
        console.log("AuthGuard: insufficient role, redirecting to /auth")
        router.push(getFullPath("/auth"))
        return
      }

      console.log("AuthGuard: access granted")
    } else {
      console.log("AuthGuard: waiting for auth check...", {
        isLoading,
        isChecking,
      })
    }
  }, [isLoading, isChecking, isAuthenticated, requiredRoles, hasRole, router])
  return { user, isLoading, isAuthenticated, hasRole }
}
