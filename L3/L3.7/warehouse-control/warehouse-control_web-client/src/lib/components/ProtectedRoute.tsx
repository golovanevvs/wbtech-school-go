"use client"

import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { UserRole } from "@/lib/types/auth"
import { ReactNode } from "react"

interface ProtectedRouteProps {
  children: ReactNode
  requiredRoles?: UserRole[]
  fallback?: ReactNode
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({
  children,
  requiredRoles,
  fallback = <div>Загрузка...</div>,
}) => {
  const { isLoading } = useAuthGuard(requiredRoles)

  if (isLoading) {
    return fallback
  }

  return <>{children}</>
}

export default ProtectedRoute
