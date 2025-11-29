"use client"

import { useState, useEffect } from "react"
import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { itemsAPI } from "@/lib/api/items"
import { Item } from "@/lib/types/items"
import ItemForm from "@/ui/ItemForm"
import { Box, Alert, CircularProgress } from "@mui/material"
import { useRouter, useSearchParams } from "next/navigation"

export default function EditItemPage() {
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
        
        // Получаем список товаров и ищем нужный
        const response = await itemsAPI.getItems()
        const foundItem = response.items.find(item => item.id === parseInt(itemId))
        
        if (!foundItem) {
          setError("Товар не найден")
          return
        }
        
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