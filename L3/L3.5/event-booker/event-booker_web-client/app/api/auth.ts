import { User, LoginRequest, RegisterRequest, AuthResponse, UpdateUserRequest } from "../lib/types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL

class ApiError extends Error {
  constructor(message: string, public status: number) {
    super(message)
    this.name = "ApiError"
  }
}

const apiRequest = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  let token = getToken()

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
    ...options,
  })

  if (response.status === 401 && token) {
    try {
      const newTokens = await refreshTokens()
      token = newTokens.token

      const retryResponse = await fetch(`${API_BASE_URL}${endpoint}`, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
          ...options.headers,
        },
        ...options,
      })

      if (!retryResponse.ok) {
        const errorData = await retryResponse.text()
        throw new ApiError(errorData, retryResponse.status)
      }

      return retryResponse.json()
    } catch (refreshError) {
      logout()
      throw refreshError
    }
  }

  if (!response.ok) {
    const errorData = await response.text()
    throw new ApiError(errorData, response.status)
  }

  return response.json()
}

const getToken = (): string | null => {
  if (typeof window !== "undefined") {
    return localStorage.getItem("token")
  }
  return null
}

const setToken = (token: string): void => {
  if (typeof window !== "undefined") {
    localStorage.setItem("token", token)
  }
}

const removeToken = (): void => {
  if (typeof window !== "undefined") {
    localStorage.removeItem("token")
  }
}

export interface LoginApiResponse {
  token: string
  refreshToken: string
}

export interface RegisterApiResponse {
  token: string
  refreshToken: string
}

export const login = async (
  credentials: LoginRequest
): Promise<AuthResponse> => {
  const response = await apiRequest<LoginApiResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  })

  setToken(response.token)

  if (typeof window !== "undefined") {
    localStorage.setItem("refreshToken", response.refreshToken)
  }

  const user = await getCurrentUser()

  return {
    user,
    token: response.token,
    refreshToken: response.refreshToken,
  }
}

export const register = async (
  userData: RegisterRequest
): Promise<AuthResponse> => {
  const response = await apiRequest<RegisterApiResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(userData),
  })

  setToken(response.token)

  if (typeof window !== "undefined") {
    localStorage.setItem("refreshToken", response.refreshToken)
  }

  const user = await getCurrentUser()

  return {
    user,
    token: response.token,
  }
}

export const getCurrentUser = async (): Promise<User> => {
  const token = getToken()

  if (!token) {
    throw new ApiError("No authentication token", 401)
  }

  return apiRequest<User>("/auth/me", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
}

export const updateUser = async (userData: UpdateUserRequest): Promise<User> => {
  const token = getToken()

  if (!token) {
    throw new ApiError("No authentication token", 401)
  }

  console.log("Sending to backend:", userData)

  return apiRequest<User>("/auth/update", {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  })
}

export const logout = (): void => {
  removeToken()
}

export const refreshTokens = async (): Promise<{
  token: string
  refreshToken: string
}> => {
  const refreshToken =
    typeof window !== "undefined" ? localStorage.getItem("refreshToken") : null

  if (!refreshToken) {
    throw new ApiError("No refresh token", 401)
  }

  const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ refreshToken }),
  })

  if (!response.ok) {
    const errorData = await response.text()
    throw new ApiError(errorData, response.status)
  }

  const data = await response.json()

  setToken(data.token)
  if (typeof window !== "undefined") {
    localStorage.setItem("refreshToken", data.refreshToken)
  }

  return data
}
