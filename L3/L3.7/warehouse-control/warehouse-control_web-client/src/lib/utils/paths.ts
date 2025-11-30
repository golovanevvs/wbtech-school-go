/**
 * Утилиты для работы с путями с учетом basePath
 */

const getBasePath = (): string => {
  return process.env.NEXT_PUBLIC_BASE_PATH || ""
}

/**
 * Формирует полный путь с учетом basePath
 */
export const getFullPath = (path: string): string => {
  const basePath = getBasePath()
  
  // Если path уже содержит basePath, возвращаем как есть
  if (basePath && path.startsWith(basePath)) {
    return path
  }
  
  // Убираем начальный слеш из path для корректного формирования
  const cleanPath = path.startsWith("/") ? path.slice(1) : path
  
  // Формируем полный путь
  return basePath ? `/${basePath}/${cleanPath}` : `/${cleanPath}`
}

/**
 * Формирует относительный путь для использования в router.push()
 * Если basePath пустой, возвращает исходный путь
 */
export const getRelativePath = (path: string): string => {
  const basePath = getBasePath()
  
  if (!basePath) {
    return path
  }
  
  // Если path начинается с basePath, убираем его
  if (path.startsWith(basePath)) {
    return path.slice(basePath.length) || "/"
  }
  
  return path
}

/**
 * Проверяет, является ли текущий путь страницей авторизации
 */
export const isAuthPage = (path?: string): boolean => {
  const currentPath = path || (typeof window !== "undefined" ? window.location.pathname : "")
  return currentPath === "/auth" || currentPath === getFullPath("/auth")
}