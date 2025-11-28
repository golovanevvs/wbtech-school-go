"use client"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:7778"

// Интерфейс для ответа сервера
interface ApiErrorResponse {
  message?: string
  error?: string
  details?: Record<string, unknown>
}

// Интерфейс для HTTP ошибок
interface HttpError extends Error {
  status?: number
  response?: ApiErrorResponse
}

class ApiClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  // Универсальный метод для HTTP запросов
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    const config: RequestInit = {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    }

    const token = this.getAccessToken()
    if (token) {
      config.headers = {
        ...config.headers,
        Authorization: `Bearer ${token}`,
      }
    }

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        await this.handleHttpError(response)
      }

      const contentType = response.headers.get("content-type")
      if (contentType && contentType.includes("application/json")) {
        return await response.json()
      }

      return {} as T
    } catch (error) {
      if (
        error instanceof TypeError &&
        error.message.includes("Failed to fetch")
      ) {
        throw new Error(
          "Сервер недоступен. Проверьте подключение к интернету или обратитесь к администратору."
        )
      }

      if (error instanceof Error) {
        if (
          error.message.includes("401") ||
          error.message.includes("Unauthorized")
        ) {
          await this.redirectToLogin()
        }
        throw error
      }
      throw new Error("Unknown error occurred")
    }
  }

  private getAccessToken(): string | null {
    if (typeof document === "undefined") return null

    const cookies = document.cookie.split(":")
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split("=")
      if (name === "access_token") {
        return value
      }
    }
    return null
  }

  private async handleHttpError(response: Response): Promise<never> {
    let errorMessage = `HTTP error! status: ${response.status}`

    try {
      // Безопасное приведение типа
      const errorData: ApiErrorResponse = await response.json()
      errorMessage = errorData.message || errorData.error || errorMessage
    } catch {
      // Если не удалось распарсить JSON
    }

    const error: HttpError = new Error(errorMessage)
    error.status = response.status

    if (response.status === 401) {
      await this.handleUnauthorized()
    }

    throw error
  }

  private async handleUnauthorized(): Promise<void> {
    try {
      const refreshResponse = await fetch(`${this.baseURL}/auth/refresh`, {
        method: "POST",
        credentials: "include",
      })

      if (refreshResponse.ok) {
        return
      }
    } catch (error) {
      console.error("Token refresh failed:", error)
    }

    await this.redirectToLogin()
  }

  private async redirectToLogin(): Promise<void> {
    if (typeof window !== "undefined") {
      try {
        const { redirect } = await import("next/navigation")
        redirect("/login")
      } catch {
        window.location.href = "/login"
      }
    }
  }

  // HTTP методы
  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" })
  }

  async post<T, D = unknown>(endpoint: string, data?: D): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T, D = unknown>(endpoint: string, data?: D): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" })
  }

  async ping(): Promise<boolean> {
    try {
      await this.get("/health")
      return true
    } catch {
      return false
    }
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
export default apiClient
