"use client"

import { useState, useEffect, Suspense } from "react"
import { useAuth } from "@/lib/contexts/AuthContext"
import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { itemsAPI } from "@/lib/api/items"
import { ItemAction } from "@/lib/types/items"
import { useRouter, useSearchParams } from "next/navigation"
import {
  Box,
  Typography,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  CircularProgress,
  Alert,
  Chip,
  IconButton,
} from "@mui/material"
import ArrowBackIcon from "@mui/icons-material/ArrowBack"
import DownloadIcon from "@mui/icons-material/Download"

// Компонент для содержимого страницы истории
function HistoryContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { hasRole } = useAuth()
  const { isLoading, isAuthenticated } = useAuthGuard()

  const [history, setHistory] = useState<ItemAction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Получаем ID товара из URL
  const itemId = searchParams.get("itemId")

  // Проверяем права доступа - только для Аудитора
  useEffect(() => {
    if (!isLoading && isAuthenticated && !hasRole(["Аудитор"])) {
      router.push("/")
    }
  }, [isLoading, isAuthenticated, hasRole, router])

  // Загружаем историю изменений
  useEffect(() => {
    const loadHistory = async () => {
      if (!itemId) {
        setError("ID товара не указан")
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        setError(null)
        
        const response = await itemsAPI.getItemHistory(parseInt(itemId))
        setHistory(response.history || [])
      } catch (err) {
        console.error("Failed to load history:", err)
        setError(err instanceof Error ? err.message : "Не удалось загрузить историю")
      } finally {
        setLoading(false)
      }
    }

    if (isAuthenticated && !isLoading && itemId && hasRole(["Аудитор"])) {
      loadHistory()
    }
  }, [isAuthenticated, isLoading, itemId, hasRole])

  // Форматирование даты
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("ru-RU")
  }

  // Форматирование типа действия
  const formatActionType = (actionType: string) => {
    switch (actionType) {
      case "create":
        return { label: "Создан", color: "success" as const }
      case "update":
        return { label: "Изменен", color: "warning" as const }
      case "delete":
        return { label: "Удален", color: "error" as const }
      default:
        return { label: actionType, color: "default" as const }
    }
  }

  // Форматирование изменений для отображения
  const formatChanges = (changes: string | undefined) => {
    if (!changes || changes === "null") return "Нет данных"
    
    try {
      const parsed = JSON.parse(changes)
      const parts = []
      
      // Обрабатываем разные типы изменений
      for (const [key, value] of Object.entries(parsed)) {
        if (typeof value === 'object' && value !== null && 'old' in value && 'new' in value) {
          // Это изменение (есть old и new)
          parts.push(
            <Box key={key} sx={{ mb: 1 }}>
              <Typography variant="body2" sx={{ fontWeight: 'medium', color: 'primary.main' }}>
                {key}:
              </Typography>
              <Box sx={{ ml: 2 }}>
                <Typography variant="body2" sx={{ color: 'error.main' }}>
                  Было: {String(value.old)}
                </Typography>
                <Typography variant="body2" sx={{ color: 'success.main' }}>
                  Стало: {String(value.new)}
                </Typography>
              </Box>
            </Box>
          )
        } else {
          // Это создание (просто значение)
          parts.push(
            <Typography key={key} variant="body2" sx={{ mb: 0.5 }}>
              <strong>{key}:</strong> {String(value)}
            </Typography>
          )
        }
      }
      
      return <Box>{parts}</Box>
    } catch (error) {
      return <Typography variant="body2" color="error">Ошибка парсинга данных</Typography>
    }
  }

  // Обработчик экспорта в CSV
  const handleExportCSV = async () => {
    if (!itemId) return
    
    try {
      setError(null)
      const blob = await itemsAPI.exportItemHistoryCSV(parseInt(itemId))
      
      // Создаем ссылку для скачивания
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement("a")
      link.href = url
      link.download = `item_${itemId}_history.csv`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (err) {
      console.error("Failed to export CSV:", err)
      setError(err instanceof Error ? err.message : "Не удалось экспортировать CSV")
    }
  }

  const handleBackToItems = () => {
    router.push("/")
  }

  // Показываем загрузку
  if (isLoading || loading) {
    return (
      <Box
        sx={{
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          minHeight: "50vh",
        }}
      >
        <CircularProgress />
      </Box>
    )
  }

  // Если не авторизован или не Аудитор, useAuthGuard перенаправит
  if (!isAuthenticated || !hasRole(["Аудитор"])) {
    return null
  }

  return (
    <Box sx={{ maxWidth: 1200, mx: "auto", p: 2 }}>
      {/* Заголовок и кнопки */}
      <Box sx={{ display: "flex", alignItems: "center", mb: 3, gap: 2 }}>
        <IconButton onClick={handleBackToItems} color="primary">
          <ArrowBackIcon />
        </IconButton>
        <Typography variant="h4" sx={{ flex: 1 }}>
          История изменений товара #{itemId}
        </Typography>
        <Button
          variant="outlined"
          startIcon={<DownloadIcon />}
          onClick={handleExportCSV}
          color="primary"
        >
          Экспорт в CSV
        </Button>
      </Box>

      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      {/* Таблица истории */}
      <TableContainer
        component={Paper}
        elevation={2}
        sx={{
          backgroundColor: "background.paper",
          border: "1px solid",
          borderColor: "divider",
        }}
      >
        <Table>
          <TableHead>
            <TableRow sx={{ backgroundColor: "action.hover" }}>
              <TableCell sx={{ fontWeight: "bold" }}>ID</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Дата изменения</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Действие</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Пользователь</TableCell>
              <TableCell sx={{ fontWeight: "bold" }}>Что изменилось</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {history.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} sx={{ textAlign: "center", py: 4 }}>
                  <Typography color="text.secondary">
                    История изменений не найдена
                  </Typography>
                </TableCell>
              </TableRow>
            ) : (
              history.map((action) => {
                const actionTypeInfo = formatActionType(action.action_type)
                return (
                  <TableRow
                    key={action.id}
                    hover
                    sx={{ "&:hover": { backgroundColor: "action.hover" } }}
                  >
                    <TableCell>{action.id}</TableCell>
                    <TableCell>{formatDate(action.created_at)}</TableCell>
                    <TableCell>
                      <Chip
                        label={actionTypeInfo.label}
                        color={actionTypeInfo.color}
                        size="small"
                      />
                    </TableCell>
                    <TableCell sx={{ fontWeight: "medium" }}>
                      {action.user_name}
                    </TableCell>
                    <TableCell>
                      <Box sx={{ maxWidth: 400 }}>
                        {formatChanges(action.changes)}
                      </Box>
                    </TableCell>
                  </TableRow>
                )
              })
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  )
}

// Главный компонент страницы с Suspense
export default function HistoryPage() {
  return (
    <Suspense fallback={
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "50vh" }}>
        <CircularProgress />
      </Box>
    }>
      <HistoryContent />
    </Suspense>
  )
}