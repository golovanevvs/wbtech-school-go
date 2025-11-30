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
    console.log("ApiClient.request:", endpoint, "hasRetried:", hasRetried)

    const url = `${this.baseURL}${endpoint}`

    const config: RequestInit = {
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    }

    try {
      const response = await fetch(url, config)
      console.log(
        "ApiClient.request response:",
        endpoint,
        "status:",
        response.status
      )

      if (!response.ok) {
        if (response.status === 401 && !hasRetried) {
          console.log(
            "ApiClient.request: 401 detected, calling handleUnauthorized"
          )
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
      console.log("ApiClient.request error:", endpoint, error)
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
        return
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
      // Сохранение текущего пути для возврата после авторизации
      const currentPath = window.location.pathname + window.location.search
      // Избегаем сохранения пути /auth, чтобы не создавать циклы переадресации
      if (currentPath !== "/auth") {
        sessionStorage.setItem("redirectAfterLogin", currentPath)
      }

      // Формируем полный путь с basePath
      const basePath = process.env.NEXT_PUBLIC_BASE_PATH || ""
      const authPath = basePath ? `${basePath}/auth` : "/auth"

      try {
        // Используем Next.js redirect с полным путем
        const { redirect } = await import("next/navigation")
        redirect(authPath)
      } catch (error) {
        console.error("Next.js redirect failed, falling back to window.location:", error)
        // Fallback: window.location с полным путем
        window.location.href = authPath
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
