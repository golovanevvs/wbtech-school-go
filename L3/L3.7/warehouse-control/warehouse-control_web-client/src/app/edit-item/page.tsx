"use client"

import { useState, useEffect, Suspense } from "react"
import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { itemsAPI } from "@/lib/api/items"
import { Item } from "@/lib/types/items"
import ItemForm from "@/ui/ItemForm"
import { Box, Alert, CircularProgress } from "@mui/material"
import { useRouter, useSearchParams } from "next/navigation"

// Компонент для содержимого страницы редактирования
function EditItemContent() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { isLoading, isAuthenticated, hasRole } = useAuthGuard()
  
  const [item, setItem] = useState<Item | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Получаем ID товара из URL
  const itemId = searchParams.get("id")

  // Загружаем данные товара
  useEffect(() => {
    const loadItem = async () => {
      if (!itemId) {
        setError("ID товара не указан")
        setLoading(false)
        return
      }

      try {
        setLoading(true)
        setError(null)
        
        // Получаем один товар по ID
        const foundItem = await itemsAPI.getItem(parseInt(itemId))
        setItem(foundItem)
      } catch (err) {
        console.error("Failed to load item:", err)
        setError(err instanceof Error ? err.message : "Не удалось загрузить товар")
      } finally {
        setLoading(false)
      }
    }

    if (isAuthenticated && !isLoading && itemId) {
      loadItem()
    }
  }, [isAuthenticated, isLoading, itemId])

  // Проверяем права доступа
  useEffect(() => {
    if (!isLoading && isAuthenticated && !hasRole(["Кладовщик"])) {
      router.push("/")
    }
  }, [isLoading, isAuthenticated, hasRole, router])

  const handleSubmit = async (data: Omit<Item, 'id' | 'created_at' | 'updated_at'>) => {
    if (!item) return
    
    try {
      setError(null)
      await itemsAPI.updateItem(item.id, data)
      router.push("/")
    } catch (err) {
      console.error("Failed to update item:", err)
      setError(err instanceof Error ? err.message : "Не удалось обновить товар")
      throw err // Перебрасываем ошибку для обработки в форме
    }
  }

  const handleCancel = () => {
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
          minHeight: "50vh" 
        }}
      >
        <CircularProgress />
      </Box>
    )
  }

  // Если не авторизован или не Кладовщик, useAuthGuard перенаправит
  if (!isAuthenticated || !hasRole(["Кладовщик"])) {
    return null
  }

  // Если товар не найден
  if (!item) {
    return (
      <Box sx={{ maxWidth: 600, mx: "auto", p: 3 }}>
        <Alert severity="error">
          {error || "Товар не найден"}
        </Alert>
      </Box>
    )
  }

  return (
    <Box>
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}
      
      <ItemForm
        title={`Редактирование товара #${item.id}`}
        submitButtonText="Сохранить"
        initialData={item}
        onSubmit={handleSubmit}
        onCancel={handleCancel}
        loading={loading}
      />
    </Box>
  )
}

// Главный компонент страницы с Suspense
export default function EditItemPage() {
  return (
    <Suspense fallback={
      <Box sx={{ display: "flex", justifyContent: "center", alignItems: "center", minHeight: "50vh" }}>
        <CircularProgress />
      </Box>
    }>
      <EditItemContent />
    </Suspense>
  )
}