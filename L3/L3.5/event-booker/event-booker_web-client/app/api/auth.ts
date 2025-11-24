import {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  UpdateUserRequest,
} from "../lib/types"

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
  console.log(`Making request to: ${API_BASE_URL}${endpoint}`)

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include",
    ...options,
  })

  console.log(`Response status: ${response.status}`)

  if (!response.ok) {
    const errorData = await response.text()
    console.error(`API Error: ${response.status}`, errorData)
    throw new ApiError(errorData, response.status)
  }

  return response.json()
}

const removeCookie = (name: string): void => {
  document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;`
}

interface LoginApiResponse {
  message: string
}

interface RegisterApiResponse {
  message: string
  user: User
}

export const login = async (
  credentials: LoginRequest
): Promise<AuthResponse> => {
  console.log("Attempting login with:", credentials.email)
  await apiRequest<LoginApiResponse>("/auth/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  })

  console.log("Login successful, cookies after login:", document.cookie)

  const user = await getCurrentUser()
  console.log("User data retrieved:", user)

  return {
    user,
  }
}

export const register = async (
  userData: RegisterRequest
): Promise<AuthResponse> => {
  console.log("Attempting registration with:", userData.email)
  await apiRequest<RegisterApiResponse>("/auth/register", {
    method: "POST",
    body: JSON.stringify(userData),
  })
  console.log(
    "Registration successful, cookies after registration:",
    document.cookie
  )
  const user = await getCurrentUser()
  console.log("User data retrieved:", user)
  return {
    user,
  }
}

export const getCurrentUser = async (): Promise<User> => {
  return apiRequest<User>("/auth/me")
}

export const updateUser = async (
  userData: UpdateUserRequest
): Promise<User> => {
  console.log("Sending to backend:", userData)
  return apiRequest<User>("/auth/update", {
    method: "PUT",
    body: JSON.stringify(userData),
  })
}

export const deleteUser = async (): Promise<{ message: string }> => {
  console.log("Deleting user profile")
  return apiRequest<{ message: string }>("/auth/delete", {
    method: "DELETE",
  })
}

export const logout = async (): Promise<void> => {
  try {
    await apiRequest<{ message: string }>("/auth/logout", {
      method: "POST",
    })
    console.log("Server logout successful")
  } catch (error) {
    console.error("Server logout failed:", error)
    removeCookie("access_token")
    removeCookie("refresh_token")
    throw error
  }
}
