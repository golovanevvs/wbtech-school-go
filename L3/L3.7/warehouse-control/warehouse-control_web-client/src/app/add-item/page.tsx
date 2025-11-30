"use client"

import { useAuthGuard } from "@/lib/hooks/useAuthGuard"
import { itemsAPI } from "@/lib/api/items"
import { Item } from "@/lib/types/items"
import ItemForm from "@/ui/ItemForm"
import { Box, Alert, CircularProgress } from "@mui/material"
import { useRouter } from "next/navigation"
import { useState, useEffect } from "react"
import { getFullPath } from "@/lib/utils/paths"

export default function AddItemPage() {
  const router = useRouter()
  const { isLoading, isAuthenticated, hasRole } = useAuthGuard()
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!isLoading && isAuthenticated && !hasRole(["Кладовщик"])) {
      router.push(getFullPath("/"))
    }
  }, [isLoading, isAuthenticated, hasRole, router])

  const handleSubmit = async (
    data: Omit<Item, "id" | "created_at" | "updated_at">
  ) => {
    try {
      setError(null)
      await itemsAPI.createItem(data)
      router.push(getFullPath("/"))
    } catch (err) {
      console.error("Failed to create item:", err)
      setError(err instanceof Error ? err.message : "Не удалось создать товар")
      throw err
    }
  }

  const handleCancel = () => {
    router.push(getFullPath("/"))
  }

  if (isLoading) {
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

  if (!isAuthenticated || !hasRole(["Кладовщик"])) {
    return null
  }

  return (
    <Box>
      {error && (
        <Alert severity="error" sx={{ mb: 3 }}>
          {error}
        </Alert>
      )}

      <ItemForm
        title="Добавление товара"
        submitButtonText="Сохранить"
        onSubmit={handleSubmit}
        onCancel={handleCancel}
        loading={isLoading}
      />
    </Box>
  )
}
