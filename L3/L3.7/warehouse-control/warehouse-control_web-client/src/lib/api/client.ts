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
    options: RequestInit = {},
    hasRetried = false
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

    // const token = this.getAccessToken()
    // if (token) {
    //   config.headers = {
    //     ...config.headers,
    //     Authorization: `Bearer ${token}`,
    //   }
    // }

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        if (response.status === 401 && !hasRetried) {
          try {
            await this.handleUnauthorized()
            return this.request(endpoint, options, true)
          } catch {
            await this.redirectToLogin()
            throw new Error("Unauthorized")
          }
        }
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

    const cookies = document.cookie.split(";")
    for (const cookie of cookies) {
      const [name, ...valueParts] = cookie.trim().split("=")
      if (name === "access_token") {
        return valueParts.join("=") // Восстанавливаем значение, если в нем были знаки '='
      }
    }
    return null
  }

  private getRefreshToken(): string | null {
    if (typeof document === "undefined") return null

    const cookies = document.cookie.split(";")
    for (const cookie of cookies) {
      const [name, ...valueParts] = cookie.trim().split("=")
      if (name === "refresh_token") {
        return valueParts.join("=") // Восстанавливаем значение, если в нем были знаки '='
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

    throw error
  }

  private async handleUnauthorized(): Promise<void> {
    try {
      const refreshResponse = await fetch(`${this.baseURL}/auth/refresh`, {
        method: "POST",
        credentials: "include",
      })

      if (refreshResponse.ok) {
        return  // Успешно обновили токен
      }
      
      // Если ответ не ok, пробрасываем ошибку
      const errorData = await refreshResponse.json().catch(() => ({}))
      throw new Error(errorData.message || "Failed to refresh token")
    } catch (error) {
      console.error("Token refresh failed:", error)
      
      // Пробрасываем ошибку, чтобы request() знал, что обновление не удалось
      throw error
    }
  }

  private async redirectToLogin(): Promise<void> {
    if (typeof window !== "undefined") {
      // Сохраняем текущий путь для возврата после авторизации
      const currentPath = window.location.pathname + window.location.search
      if (currentPath !== "/auth") {
        sessionStorage.setItem("redirectAfterLogin", currentPath)
      }

      try {
        const { redirect } = await import("next/navigation")
        redirect("/auth")
      } catch {
        window.location.href = "/auth"
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
